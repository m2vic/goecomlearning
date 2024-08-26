package service

import (
	"fmt"
	"golearning/internal/core/domain"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

type StripeService struct {
	StripeKey string
}

func NewStripeService(stripeKey string) *StripeService {
	return &StripeService{StripeKey: stripeKey}
}

func (s StripeService) CreateSession(req domain.ProductList) (*stripe.CheckoutSession, error) {
	url := os.Getenv("SUCCESSURL")
	stripe.Key = s.StripeKey
	list := req.ProductList
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for i := 0; i < len(list); i++ {
		itemDetail := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("thb"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(list[i].ProductName),
					Description: stripe.String(list[i].Details),
					Images:      stripe.StringSlice(list[i].Images),
				},
				UnitAmount: stripe.Int64(int64(list[i].PricePerPiece * 100)),
			},
			Quantity: stripe.Int64(int64(list[i].Amount)),
		}
		lineItems = append(lineItems, itemDetail)
	}

	Sessionparams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(url),
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	makeSession, err := session.New(Sessionparams)
	if err != nil {
		return nil, fmt.Errorf("session err:%w", err)
	}
	return makeSession, nil
}
