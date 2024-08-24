package port

import (
	"context"
	"golearning/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	GetUser(ctx context.Context, userId string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.Token, error)
	Register(ctx context.Context, username, password, email string) error
	UpdateUser(ctx context.Context, info domain.User, id string) error
	CheckUsername(ctx context.Context, username string) (bool, error)
	CheckEmail(ctx context.Context, email string) (bool, error)
	ResetPassword(ctx context.Context, email string) (string, error)
	ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error
	CheckRefresh(ctx context.Context, token string) (*domain.Token, error)
	GetCart(ctx context.Context, userId string) ([]domain.Cart, error)
	AddtoCart(ctx context.Context, userProduct domain.Cart, userId string) error
	IncreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error
	DecreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error
	DeleteItemInCart(ctx context.Context, userId string, ProductId primitive.ObjectID) error
	DeleteItemFromSystem(ctx context.Context, productId primitive.ObjectID) error
	EditItemFromSystem(ctx context.Context, product domain.Product) error
	ClearCart(ctx context.Context, userId string) error
}

type UserService interface {
	GetUser(ctx context.Context, userId string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.Token, error)
	Register(ctx context.Context, info domain.User) error
	UpdateUser(ctx context.Context, user domain.User, token string) error
	ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) (string, error)
	CheckRefresh(ctx context.Context, token string) (*domain.Token, error)
	GetCart(ctx context.Context, userId string) ([]domain.Cart, error)
	AddtoCart(ctx context.Context, userProduct domain.Cart, userId string) error
	IncreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error
	DecreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error
	DeleteItemInCart(ctx context.Context, userId string, productId primitive.ObjectID) error
	ClearCart(ctx context.Context, userId string) error
}
