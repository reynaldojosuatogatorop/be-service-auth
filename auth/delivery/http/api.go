package http

import (
	"be-service-auth/auth/delivery/http/handler"
	// "be-service-auth/delivery/http/handler"
	"be-service-auth/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

// RouterAPI is the main router for this Service Insurance REST API
func RouterAPI(app *fiber.App, Auth domain.AuthUseCase) {
	handlerAuth := &handler.AuthHandler{AuthUseCase: Auth}

	basePath := viper.GetString("server.base_path")

	cms := app.Group(basePath)

	cms.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// Authentication & Authorization Route
	cms.Post("/login", handlerAuth.Login)
	cms.Get("/auth", handlerAuth.AuthorizationAuth(), handlerAuth.Auth)
	cms.Post("/logout", handlerAuth.AuthorizationAuth(), handlerAuth.Logout)

}
