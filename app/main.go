package main

import (
	"be-service-auth/config"
	"be-service-auth/helper"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"strconv"

	_RepoGRPCAuthObject "be-service-auth/auth/delivery/grpc"
	_RepoGRPCAuthServer "be-service-auth/auth/delivery/grpc/authorization"
	_DeliveryHTTP "be-service-auth/auth/delivery/http"
	_RepoMySQLAuth "be-service-auth/auth/repository/mysql"
	_RepoRedisAuth "be-service-auth/auth/repository/redis"
	_UsecaseAuth "be-service-auth/auth/usecase"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	errorrs "gopkg.in/oauth2.v3/errors"
	manage "gopkg.in/oauth2.v3/manage"
	modelsOAuth "gopkg.in/oauth2.v3/models"
	serverOAuth "gopkg.in/oauth2.v3/server"
	store "gopkg.in/oauth2.v3/store"
)

const (
	dbFlag = "mysql"
)

func main() {
	// CLI options parse
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	// Config file
	config.ReadConfig(*configFile)

	// Set log level
	switch viper.GetString("server.log_level") {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Initialize database connection
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.database"))

	dsn := fmt.Sprintf("%s?%s", connection, "") // MySQL does not require additional parameters in DSN

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatal(err)
	}

	// Initial OAuth2
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()

	manager.MapClientStorage(clientStore)

	srv := serverOAuth.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(serverOAuth.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errorrs.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errorrs.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	// clientId := viper.GetStringSlice("oauth2_credential.client_id")
	// clientSecret := viper.GetStringSlice("oauth2_credential.client_secret")
	// domainStore := viper.GetStringSlice("oauth2_credential.client_domain")
	// log.Println(domainStore)

	// Migrate database if any new schema
	driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err == nil {
		mig, err := migrate.NewWithDatabaseInstance(viper.GetString("database.path_migrate"), viper.GetString("mysql.database"), driver)
		log.Info(viper.GetString("database.path_migrate"))
		if err == nil {
			err = mig.Up()
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Debug("No database migration")
				} else {
					log.Error(err)
				}
			} else {
				log.Info("Migrate database success")
			}
			version, dirty, err := mig.Version()
			if err != nil && err != migrate.ErrNilVersion {
				log.Error(err)
			}
			log.Debug("Current DB version: " + strconv.FormatUint(uint64(version), 10) + "; Dirty: " + strconv.FormatBool(dirty))
		} else {
			log.Warn(err)
		}
	} else {
		log.Warn(err)
	}

	// Initialize Redis
	ctx := context.Background()
	dbRedis := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.database"),
		PoolSize: viper.GetInt("redis.max_connection"),
	})

	_, err = dbRedis.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Redis connection established")

	// Register repository & usecase auth

	repoMySQLAuth := _RepoMySQLAuth.NewMySQLAuthRepository(dbConn)
	repoMySQLOAuth := _RepoMySQLAuth.NewMySQLOAuthRepository(dbConn)
	repoRedisAuth := _RepoRedisAuth.NewRedisAuthRepository(dbRedis)

	recachingData, err := helper.RecachingB2BData(repoMySQLOAuth)
	if err != nil {
		log.Error(err)
		return
	}

	for _, i := range recachingData {
		store := clientStore.Set(i.ClientID, &modelsOAuth.Client{
			ID:     i.ClientID,
			Secret: i.ClientSecret,
			Domain: i.Domain,
		})

		fmt.Println("Client Store :", store)
	}
	oAuthHttpInit := srv
	usecaseAuth := _UsecaseAuth.NewAuthUsecase(repoMySQLAuth, repoRedisAuth)
	usecaseOAuth := _UsecaseAuth.NewOAuthUsecase(repoMySQLOAuth, repoRedisAuth, oAuthHttpInit)
	serverAuth := _RepoGRPCAuthObject.NewGRPCAuth(usecaseAuth)

	// Initialize gRPC server
	go func() {
		listen, err := net.Listen("tcp", ":"+viper.GetString("server.grpc_port"))
		if err != nil {
			log.Fatalf("[ERROR] Failed to listen tcp: %v", err)
		}

		grpcServer := grpc.NewServer()
		_RepoGRPCAuthServer.RegisterAuthorizationServiceServer(grpcServer, serverAuth)
		log.Println("gRPC server is running in port", viper.GetString("server.grpc_port"))
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Initialize HTTP web framework
	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.prefork"),
		StrictRouting: viper.GetBool("server.strict_routing"),
		CaseSensitive: viper.GetBool("server.case_sensitive"),
		BodyLimit:     viper.GetInt("server.body_limit"),
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// HTTP routing
	app.Get(viper.GetString("server.base_path")+"/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	_DeliveryHTTP.RouterAPI(app, usecaseAuth, usecaseOAuth)

	// Start Fiber HTTP server
	if err := app.Listen(":" + viper.GetString("server.port")); err != nil {
		log.Fatal(err)
	}
}
