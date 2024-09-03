package port

import (
	"context"
	"golearning/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	GetUser(ctx context.Context, userId primitive.ObjectID) (*domain.User, error)
	Register(ctx context.Context, username, email string, hash []byte) error
	UpdateUser(ctx context.Context, info domain.User) error
	CheckUsername(ctx context.Context, username string) (*domain.User, error)
	CheckPassword(ctx context.Context, password string) error
	CheckEmail(ctx context.Context, email string) (bool, error)
	ResetPassword(ctx context.Context, email, randomPassword string) (string, error)
	ChangePassword(ctx context.Context, userId primitive.ObjectID, newPassword string) error
	CheckRefresh(ctx context.Context, token string) (*domain.User, error)
	GetCart(ctx context.Context, userId primitive.ObjectID) ([]domain.Cart, error)
	AddtoCart(ctx context.Context, userProduct domain.Cart, userId primitive.ObjectID) error
	IncreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error
	DecreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error
	DeleteItemInCart(ctx context.Context, userId primitive.ObjectID, ProductId primitive.ObjectID) error
	DeleteItemFromSystem(ctx context.Context, productId primitive.ObjectID) error
	EditItemFromSystem(ctx context.Context, product domain.Product) error
	ClearCart(ctx context.Context, userId primitive.ObjectID) error
}

type UserService interface {
	GetUser(ctx context.Context, userId primitive.ObjectID) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.Token, error)
	Register(ctx context.Context, info domain.User) error
	UpdateUser(ctx context.Context, user domain.User) error
	ChangePassword(ctx context.Context, userId primitive.ObjectID, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) (string, error)
	RefreshToken(ctx context.Context, token string) (*domain.Token, error)
	GetCart(ctx context.Context, userId primitive.ObjectID) ([]domain.Cart, error)
	AddtoCart(ctx context.Context, userProduct domain.Cart, userId primitive.ObjectID) error
	IncreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error
	DecreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error
	DeleteItemInCart(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error
	ClearCart(ctx context.Context, userId primitive.ObjectID) error
	CheckEmail(ctx context.Context, email string) (string, error)
}

type TokenGenerator interface {
	GenerateToken(context.Context, string, primitive.ObjectID) (*domain.Token, error)
}
type PasswordHasher interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	ComparePassword(ctx context.Context, password, hash string) error
}
type PasswordGenerator interface {
	RandomPassword() (string, error)
}
