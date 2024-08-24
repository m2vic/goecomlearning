package service_test

import (
	"context"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v79"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCheckout(t *testing.T) {
	mockOrderService := new(service.MockOrderService)
	mockStripeService := new(service.MockStripeService)
	mockProductService := new(service.MockProductService)
	mockUserService := new(service.MockUserService)

	// Create the checkout service with the mocks
	checkoutService := service.NewCheckoutService(mockOrderService, mockStripeService, mockProductService, mockUserService)

	// Define expected results
	successfulURL := "https://example.com/success"
	expectedCheckoutURL := "checkoutURL"
	ctx := context.Background()

	// Mock product and user details
	img := []string{"1", "2"}
	productId, _ := primitive.ObjectIDFromHex("1")
	productName := "test"
	amount := 2
	userId := "1"
	sessionId := "1"
	details := "test"
	price := float64(122)
	totalPrice := amount * int(price)
	mockTime := time.Now()

	// Define the product list and order details
	productList := domain.ProductList{
		ProductList: []domain.StripeProduct{
			{
				ProductId:     productId,
				ProductName:   productName,
				Amount:        amount,
				PricePerPiece: price,
				Details:       details,
				Images:        img,
			},
		},
	}

	// Mock session and order details
	lineItems := []domain.ProductDetails{
		{
			ProductName: productName,
			UnitPrice:   int(price),
			Quantity:    amount,
			Images:      img,
			Description: details,
		},
	}
	order := domain.Order{
		UserId:         userId,
		OrderID:        sessionId,
		Ordered_At:     mockTime, // Replace this with a fixed time for testing
		TotalPrice:     totalPrice,
		Payment_Method: "card",
		LineItems:      lineItems,
		Status:         "unpaid",
	}

	// Mock return values from services
	url := &stripe.CheckoutSession{URL: expectedCheckoutURL, ID: sessionId, AmountTotal: int64(amount) * int64(price) * 100}

	// Mock the services
	mockStripeService.On("CreateSession", productList, successfulURL).Return(url, nil)
	mockProductService.On("UpdateStock", ctx, productId, amount).Return(nil)
	mockOrderService.On("NewOrder", ctx, order).Return(nil)
	mockUserService.On("ClearCart", ctx, userId).Return(nil)

	// Call the Checkout method
	actual, err := checkoutService.Checkout(ctx, productList, userId, mockTime)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCheckoutURL, actual)

	// Ensure the NewOrder method was called with the correct order
	mockOrderService.AssertCalled(t, "NewOrder", ctx, mock.MatchedBy(func(o domain.Order) bool {
		return o.UserId == userId &&
			o.OrderID == sessionId &&
			o.TotalPrice == totalPrice &&
			o.Payment_Method == "card" &&
			len(o.LineItems) == len(lineItems)
	}))
}
