package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"golearning/internal/adapter/http/dto"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/webhook"
)

type UserHandler struct {
	userService     port.UserService
	orderService    port.OrderService
	checkoutService port.CheckoutService
}

func NewUserHandler(userService port.UserService, orderService port.OrderService, checkoutService port.CheckoutService) *UserHandler {
	return &UserHandler{userService: userService, orderService: orderService, checkoutService: checkoutService}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return c.Context().Err()
	}
	token, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		log.Printf("error log:%v", err)
		return c.SendStatus(400)
	}
	return c.JSON(token)
}
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	info := domain.User{Email: req.Email, Username: req.Username, Password: req.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = h.userService.Register(ctx, info)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	err = registerNotification(req.Email)
	if err != nil {
		fmt.Println("Notification function err:", err)
	}

	return c.SendStatus(200)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := c.Locals("userID").(string)
	data, err := h.userService.GetUser(ctx, userId)
	if err != nil {
		c.SendStatus(500)
	}
	return c.JSON(data)
}
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := getRefresh(c)
	if refreshToken == "No-Token" {
		return c.Status(401).SendString("please Login")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	token, err := h.userService.CheckRefresh(ctx, refreshToken)
	if err != nil {
		log.Println(err)
		return c.SendString("Token invalid, Login again.")
	} else {
		if token.AccessToken != "" {
			return c.JSON(token)
		}
	}
	return c.JSON(token)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {

	var req dto.UpdateUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(500)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	info := domain.User{FirstName: req.FirstName,
		LastName:       req.LastName,
		AddressDetails: req.AddressDetails}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.UpdateUser(ctx, info, userId)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {

	var req dto.ChangePasswordRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	newPassword := req.NewPassword
	oldPassword := req.OldPassword
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.ChangePassword(ctx, userId, oldPassword, newPassword)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Password Invalid")
	}
	return c.SendStatus(200)
}
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.EmailRequest
	err := c.BodyParser(&req)
	if err != nil {
		c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newPassword, err := h.userService.ResetPassword(ctx, req.Email)
	if err != nil {
		return c.Status(400).SendString("email doesn't exist")
	}

	err = resetPasswordEmail(req.Email, newPassword)
	if err != nil {
		log.Fatal(err)
	}
	return c.SendStatus(200)

}

func (h *UserHandler) AddToCart(c *fiber.Ctx) error {

	// have to handle amount, req not more than actual amount of products
	var req dto.Product
	err := c.BodyParser(&req)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	product := domain.Cart{
		ProductId:     req.ProductId,
		ProductName:   req.ProductName,
		Amount:        req.Amount,
		Images:        req.Images,
		PricePerPiece: req.PricePerPiece,
		PriceId:       req.PriceId,
		Details:       req.Details,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.AddtoCart(ctx, product, userId)
	if err != nil {
		return c.Status(400).SendString("Not Enough Product In Stock")
	}
	return c.SendStatus(200)
}
func (h *UserHandler) DeleteItemInCart(c *fiber.Ctx) error {

	var req dto.DeleteItemInCartRequest
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:,%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.DeleteItemInCart(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) IncreaseItemInCart(c *fiber.Ctx) error {

	var req dto.IncreaseItemInCartRequest
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.IncreaseCartProduct(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) DecreaseItemInCart(c *fiber.Ctx) error {

	var req dto.DecreaseItemInCartRequest
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	err = h.userService.DecreaseCartProduct(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) GetCart(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	result, err := h.userService.GetCart(ctx, userId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.JSON(result)
}
func (h *UserHandler) Checkout(c *fiber.Ctx) error {
	var list dto.CheckoutRequest
	err := c.BodyParser(&list)
	if err != nil {
		return err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(list)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	var newList domain.ProductList
	for _, item := range list.Product {
		Product := domain.StripeProduct{ProductId: item.ProductId,
			ProductName: item.ProductName, Images: item.Images, Details: item.Details, Amount: item.Amount, PricePerPiece: item.PricePerPiece, PriceId: item.PriceId}
		newList.ProductList = append(newList.ProductList, Product)
	}
	userId := c.Locals("userID").(string)
	result, err := h.checkoutService.Checkout(c.Context(), newList, userId, time.Now())
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *UserHandler) StripeWebHook(c *fiber.Ctx) error {
	//env := godotenv.Load()
	//if env != nil {
	//	fmt.Println("fail to load env")
	//}
	secret := os.Getenv("WEBHOOKENDPOINTSECRET")
	endpointSecret := secret
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		fmt.Println("err:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Webhook verification failed")
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse webhook JSON")
		}
		fmt.Println("data", session.Created)
		// Update the order status in database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = h.orderService.UpdateOrderStatus(ctx, session.ID, string(session.PaymentStatus))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Payment was successful for session: %s\n", session.ID)
	}

	return c.SendStatus(fiber.StatusOK)
}
func (h *UserHandler) GetOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := c.Locals("userID").(string)
	getOrder, err := h.orderService.GetOrder(ctx, userId)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	return c.JSON(getOrder)
}
