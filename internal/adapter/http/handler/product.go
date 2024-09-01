package handler

import (
	"context"
	"fmt"
	"golearning/internal/adapter/http/dto"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) AddNewProduct(c *fiber.Ctx) error {
	productreq := dto.ProductRequest{}
	if err := c.BodyParser(&productreq); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(productreq)
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
	role := c.Locals("role").(string)
	product := domain.Product{
		ProductName: productreq.ProductName,
		Details:     productreq.Details,
		Stock:       productreq.Stock,
		Category:    productreq.Category,
		Images:      productreq.Image,
		Price:       productreq.Price,
	}

	err = h.ProductService.AddNewProduct(ctx, role, product)
	if err != nil {
		return fmt.Errorf("handler:%w", err)
	}

	return c.SendString("product add successful!")

}

func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return c.Context().Err()
	}
	product, err := h.ProductService.GetAllProduct(ctx)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(product)
}

func (h *ProductHandler) GetProductById(c *fiber.Ctx) error {
	productId := c.Params("productid")
	if productId == "" {
		fmt.Println("querynotfound")
		return c.SendStatus(400)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return c.Context().Err()
	}
	result, err := h.ProductService.GetProductById(ctx, productId)
	if err != nil {
		return c.Status(400).SendString("product not found")
	}
	return c.JSON(result)

}
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	req := dto.ProductRequest{}
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
	role := c.Locals("role").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return c.Context().Err()
	}
	product := domain.Product{
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		Price:           req.Price,
		Details:         req.Details,
		Images:          req.Image,
		Category:        req.Category,
		Location:        req.Location,
		StripeProductId: req.StripeProductId,
	}
	err = h.ProductService.EditProduct(ctx, role, product)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	return c.Status(200).SendString("Update Successful!")
}
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	product := domain.Product{}
	err := c.BodyParser(&product)
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return c.Context().Err()
	}
	err = h.ProductService.DeleteProduct(ctx, role, product)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.Status(200).SendString("Product Deleted!")
}

func (h *ProductHandler) CheckAmount(c *fiber.Ctx) error {
	product := dto.Product{}
	err := c.BodyParser(&product)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(product)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"field": err.Field(), "error": fmt.Sprintf("Validation failed on '%s' tag", err.Tag())})
		}
	}
	result, err := h.ProductService.CheckAmount(c.Context(), product.ProductId)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}
	return c.JSON(result)
}
