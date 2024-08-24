package routes

import (
	"golearning/internal/adapter/http/handler"
	product_routes "golearning/internal/adapter/http/routes/product"
	user_routes "golearning/internal/adapter/http/routes/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {
	user_routes.UserRoutes(app, userHandler)
	product_routes.ProductRoutes(app, productHandler)
}
