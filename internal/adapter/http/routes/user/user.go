package user_routes

import (
	"golearning/internal/adapter/http/handler"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/login", userHandler.Login)
	app.Post("/register", userHandler.Register)
	app.Post("/webhook", userHandler.StripeWebHook)
	app.Post("/forgetpassword", userHandler.ResetPasswordLink)
	app.Post("/resetpassword/:email", userHandler.ResetPassword)

	userGroup := app.Group("/user", handler.AuthMiddleware)
	userGroup.Get("/getuser", userHandler.GetUser)
	userGroup.Get("/refresh", userHandler.RefreshToken)
	userGroup.Post("/update", userHandler.UpdateUser)
	userGroup.Post("/changepassword", userHandler.ChangePassword)
	userGroup.Post("/addtocart", userHandler.AddToCart)
	userGroup.Post("/cart/deleteproduct", userHandler.DeleteItemInCart)
	userGroup.Post("/cart/increase", userHandler.IncreaseItemInCart)
	userGroup.Post("/cart/decrease", userHandler.DecreaseItemInCart)
	userGroup.Get("/cart/getcart", userHandler.GetCart)
	userGroup.Post("/checkout", userHandler.Checkout)
	userGroup.Get("/getorder", userHandler.GetOrder)

}
