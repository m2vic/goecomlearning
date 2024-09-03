package service

import (
	"context"
	"fmt"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"
	"time"

	"github.com/stripe/stripe-go/v79"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckoutService struct {
	OrderService   port.OrderService
	StripeService  port.StripeService
	ProductService port.ProductService
	UserService    port.UserService
}

func NewCheckoutService(OrderService port.OrderService, StripeService port.StripeService, productService port.ProductService, UserService port.UserService) *CheckoutService {
	return &CheckoutService{OrderService: OrderService, StripeService: StripeService, ProductService: productService, UserService: UserService}
}

func (s *CheckoutService) Checkout(ctx context.Context, list domain.ProductList, userId primitive.ObjectID, now time.Time) (string, error) {

	session, err := s.StripeService.CreateSession(list)
	if err != nil {
		return "", err
	}

	err = s.ProductService.UpdateStock(ctx, list.ProductList)
	if err != nil {
		return "", err
	}

	/// where to get userId
	err = s.UserService.ClearCart(ctx, userId)
	if err != nil {
		fmt.Println("fail to clear cart")
		return "", err
	}
	order, err := MapToOrders(list.ProductList, userId, session, now)
	if err != nil {
		return "", fmt.Errorf("fail to map line items:%w", err)
	}

	//save order
	err = s.OrderService.NewOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("fail to create new Order:%w", err)
	}
	return session.URL, nil
}

func MapToOrders(list []domain.StripeProduct, userId primitive.ObjectID, session *stripe.CheckoutSession, now time.Time) (domain.Order, error) {

	arrOfProductDetails := []domain.ProductDetails{}
	ProductDetails := domain.ProductDetails{}

	for _, item := range list {
		ProductDetails.Quantity = item.Amount
		ProductDetails.UnitPrice = int(item.PricePerPiece)
		ProductDetails.Description = item.Details
		ProductDetails.ProductName = item.ProductName
		ProductDetails.Images = item.Images
		arrOfProductDetails = append(arrOfProductDetails, ProductDetails)
	}

	order := domain.Order{}
	order.UserId = userId
	order.OrderID = session.ID
	order.Ordered_At = now
	order.TotalPrice = int(session.AmountTotal / 100)
	if len(session.PaymentMethodTypes) > 0 {
		order.Payment_Method = session.PaymentMethodTypes[0]
	} else {
		order.Payment_Method = "card"
	}
	order.Status = "unpaid"
	order.LineItems = arrOfProductDetails
	return order, nil
}
