package http

import (
	"be-service-auth/auth/delivery/http/handler"
	// "be-service-auth/delivery/http/handler"
	"be-service-auth/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

// RouterAPI is the main router for this Service Insurance REST API
func RouterAPI(app *fiber.App, Auth domain.AuthUseCase, OAuth domain.OAuthUseCase) {
	handlerAuth := &handler.AuthHandler{AuthUseCase: Auth, OAuthUseCase: OAuth}

	basePath := viper.GetString("server.base_path")

	cms := app.Group(basePath)

	cms.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// Authentication & Authorization Route
	cms.Post("/login", handlerAuth.Login)
	cms.Get("/auth", handlerAuth.AuthorizationAuth(), handlerAuth.Auth)
	cms.Get("/oauth2", handlerAuth.TokenOauth(), handlerAuth.OAuth2)
	cms.Post("b2b/token", adaptor.HTTPHandlerFunc(handlerAuth.PostTokenOAuth2))
	cms.Post("/logout", handlerAuth.AuthorizationAuth(), handlerAuth.Logout)

}
