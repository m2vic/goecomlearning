package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"golearning/internal/adapter/http/dto"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"
	errs "golearning/internal/error"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService     port.UserService
	orderService    port.OrderService
	checkoutService port.CheckoutService
	EmailService    port.EmailService
	cryptoService   port.CryptoService
}

func NewUserHandler(userService port.UserService, orderService port.OrderService, checkoutService port.CheckoutService, emailservice port.EmailService, cryptoService port.CryptoService) *UserHandler {
	return &UserHandler{userService: userService, orderService: orderService, checkoutService: checkoutService, EmailService: emailservice, cryptoService: cryptoService}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := dto.LoginRequest{}
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
	req := dto.RegisterRequest{}
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
	id, _ := primitive.ObjectIDFromHex(userId)
	data, err := h.userService.GetUser(ctx, id)
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
	token, err := h.userService.RefreshToken(ctx, refreshToken)
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
	req := dto.UpdateUserRequest{}
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
	userId := c.Locals("userID").(string)
	id, _ := primitive.ObjectIDFromHex(userId)
	info := domain.User{
		UserId:         id,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		AddressDetails: req.AddressDetails}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = h.userService.UpdateUser(ctx, info)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {

	req := dto.ChangePasswordRequest{}
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
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	err = h.userService.ChangePassword(ctx, userId, oldPassword, newPassword)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Password Invalid")
	}
	return c.SendStatus(200)
}
func (h *UserHandler) ResetPasswordLink(c *fiber.Ctx) error {
	req := dto.EmailRequest{}
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
	_, err = h.userService.CheckEmail(ctx, req.Email)
	if err != nil {
		return errs.EmailNotFound
	}
	//have to set this link to an email but have to encrypt
	encodeEmail, err := h.cryptoService.Encrypt(req.Email)
	if err != nil {
		return err
	}
	err = h.EmailService.SetResetPasswordLink(req.Email, encodeEmail)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	email := c.Params("email")
	if email != "" {
		c.SendStatus(400)
	}
	decodeEmail, err := h.cryptoService.Decrypt(email)
	if err != nil {
		return err
	}
	//where to placing crypto service? // have to decrypt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newPassword, err := h.userService.ResetPassword(ctx, decodeEmail)
	if err != nil {
		return c.Status(400).SendString("email doesn't exist")
	}

	err = h.EmailService.NewPasswordNotify(decodeEmail, newPassword)
	if err != nil {
		log.Fatal(err)
	}
	return c.Status(200).SendString("Password Reset!")

}

func (h *UserHandler) AddToCart(c *fiber.Ctx) error {
	req := dto.Product{}
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
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	err = h.userService.AddtoCart(ctx, product, userId)
	if err != nil {
		return c.Status(400).SendString("Not Enough Product In Stock")
	}
	return c.SendStatus(200)
}
func (h *UserHandler) DeleteItemInCart(c *fiber.Ctx) error {
	req := dto.DeleteItemInCartRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:,%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	err = h.userService.DeleteItemInCart(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) IncreaseItemInCart(c *fiber.Ctx) error {
	req := dto.IncreaseItemInCartRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	err = h.userService.IncreaseCartProduct(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) DecreaseItemInCart(c *fiber.Ctx) error {
	req := dto.DecreaseItemInCartRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	productId := req.Product.ProductId
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	err = h.userService.DecreaseCartProduct(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.SendStatus(200)
}
func (h *UserHandler) GetCart(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	result, err := h.userService.GetCart(ctx, userId)
	if err != nil {
		return fmt.Errorf("err:%w", err)
	}
	return c.JSON(result)
}
func (h *UserHandler) Checkout(c *fiber.Ctx) error {
	list := dto.CheckoutRequest{}
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
	newList := domain.ProductList{}
	for _, item := range list.Product {
		Product := domain.StripeProduct{ProductId: item.ProductId,
			ProductName: item.ProductName, Images: item.Images, Details: item.Details, Amount: item.Amount, PricePerPiece: item.PricePerPiece, PriceId: item.PriceId}
		newList.ProductList = append(newList.ProductList, Product)
	}
	id := c.Locals("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(id)
	result, err := h.checkoutService.Checkout(c.Context(), newList, userId, time.Now())
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *UserHandler) StripeWebHook(c *fiber.Ctx) error {
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
