package product_routes

import (
	"golearning/internal/adapter/http/handler"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	productGroup := app.Group("/private/product", handler.AuthMiddleware)
	productGroup.Get("/:productid", productHandler.GetProductById)
	productGroup.Post("/update", productHandler.UpdateProduct)
	productGroup.Post("/new", productHandler.AddNewProduct)
	productGroup.Post("/delete", productHandler.DeleteProduct)
	app.Get("/product/all", productHandler.GetAllProduct)

	//app.Get("/checkamount", productHandler.CheckAmount)

}
