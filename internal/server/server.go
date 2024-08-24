package server

import (
	routes "golearning/internal/adapter/http"
	"golearning/internal/adapter/http/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start(userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {
	app := fiber.New()
	app.Use(cors.New())
	routes.SetupRoutes(app, userHandler, productHandler)
	app.Listen(":8080")
}
