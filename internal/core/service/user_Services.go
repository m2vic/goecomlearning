package service

import (
	"context"
	"fmt"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"
	errs "golearning/internal/error"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	userRepo          port.UserRepo
	ProductService    port.ProductService
	TokenGenerator    port.TokenGenerator
	PasswordHasher    port.PasswordHasher
	PasswordGenerator port.PasswordGenerator
}

func NewUserService(userRepo port.UserRepo, productService port.ProductService, tokenGenerator port.TokenGenerator, passwordHasher port.PasswordHasher, passwordGenerator port.PasswordGenerator) *UserService {
	return &UserService{userRepo: userRepo, ProductService: productService, TokenGenerator: tokenGenerator, PasswordHasher: passwordHasher, PasswordGenerator: passwordGenerator}
}

func (s *UserService) Login(ctx context.Context, username, password string) (*domain.Token, error) {
	result, err := s.userRepo.CheckUsername(ctx, username)
	if err != nil {
		return nil, errs.UsernameInvalid
	}
	dbHashedPassword := result.Password
	err = s.PasswordHasher.ComparePassword(ctx, password, dbHashedPassword)
	if err != nil {
		return nil, errs.PasswordInvalid
	}
	token, err := s.TokenGenerator.GenerateToken(ctx, result.Username, result.UserId)
	if err != nil {
		return nil, err
	}
	update := domain.User{UserId: result.UserId, Token: token.RefreshToken}
	err = s.UpdateUser(ctx, update)
	if err != nil {
		fmt.Println(err)
		return nil, errs.UpdateUserFail
	}
	return token, nil
}

func (s *UserService) Register(ctx context.Context, user domain.User) error {
	username := user.Username
	password := user.Password
	email := user.Email
	result, err := s.userRepo.CheckUsername(ctx, username)

	if result != nil && result.Username == username {
		return errs.UserAlreadyExist
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	emailExist, err := s.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return err
	}
	if emailExist {
		return errs.EmailAlreadyExist
	}

	hash, err := s.PasswordHasher.HashPassword(ctx, password)
	if err != nil {
		return errs.HashPasswordFail
	}
	// have to check username and email
	return s.userRepo.Register(ctx, username, email, hash)
}

func (s *UserService) GetUser(ctx context.Context, userId primitive.ObjectID) (*domain.User, error) {

	return s.userRepo.GetUser(ctx, userId)
}

func (s *UserService) UpdateUser(ctx context.Context, info domain.User) error {

	return s.userRepo.UpdateUser(ctx, info)
}
func (s *UserService) ChangePassword(ctx context.Context, userId primitive.ObjectID, oldPass, newPass string) error {
	hash, err := s.PasswordHasher.HashPassword(ctx, oldPass)
	if err != nil {
		return errs.HashPasswordFail
	}
	err = s.userRepo.CheckPassword(ctx, string(hash))
	if err != nil {
		return errs.PasswordInvalid
	}
	newHash, err := s.PasswordHasher.HashPassword(ctx, newPass)
	if err != nil {
		return errs.HashPasswordFail
	}
	return s.userRepo.ChangePassword(ctx, userId, string(newHash))
}
func (s *UserService) CheckEmail(ctx context.Context, email string) (string, error) {
	_, err := s.userRepo.CheckEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return "", errs.EmailNotFound
	}
	return "linkinemailthatresetpassword", nil
}
func (s *UserService) ResetPassword(ctx context.Context, email string) (string, error) {
	_, err := s.userRepo.CheckEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return "", errs.EmailNotFound
	}
	//problem now is if somebody know our password , they can reset it

	// password generator
	randomPass, err := s.PasswordGenerator.RandomPassword()
	if err != nil {
		return "", err
	}
	hash, err := s.PasswordHasher.HashPassword(ctx, randomPass)
	if err != nil {
		return "", errs.GenerateTokenFail
	}
	return s.userRepo.ResetPassword(ctx, email, string(hash))

}
func (s *UserService) RefreshToken(ctx context.Context, refreshtoken string) (*domain.Token, error) {

	result, err := s.userRepo.CheckRefresh(ctx, refreshtoken)
	if err != nil {
		return nil, errs.TokenNotFound
	}
	tokens, err := s.TokenGenerator.GenerateToken(ctx, result.Username, result.UserId)
	if err != nil {
		return nil, errs.GenerateTokenFail
	}
	//user := domain.User{UserId: result.UserId, Token: tokens.RefreshToken}
	err = s.userRepo.UpdateUser(ctx, domain.User{UserId: result.UserId, Token: tokens.RefreshToken})
	if err != nil {
		return nil, errs.UpdateUserFail
	}
	return tokens, nil
}

func (r *UserService) AddtoCart(ctx context.Context, Product domain.Cart, userId primitive.ObjectID) error {
	productId := Product.ProductId
	reqAmount := Product.Amount
	actualAmount, err := r.ProductService.CheckAmount(ctx, productId)
	if err != nil {
		return errs.ProductNotFound
	}
	if reqAmount > actualAmount {
		return errs.NotEnoughProduct
	}
	return r.userRepo.AddtoCart(ctx, Product, userId)

}
func (r *UserService) GetCart(ctx context.Context, userId primitive.ObjectID) ([]domain.Cart, error) {
	// have to check if sufficient and then compare
	newCart := []domain.Cart{}
	userCart, err := r.userRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, item := range userCart {
		amount, err := r.ProductService.CheckAmount(ctx, item.ProductId)
		if amount >= item.Amount {
			newCart = append(newCart, item)
		}
		if err != nil {
			return nil, err
		}
	}

	return newCart, nil
}
func (r *UserService) IncreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error {

	return r.userRepo.IncreaseCartProduct(ctx, userId, productId)
}
func (r *UserService) DecreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error {

	return r.userRepo.DecreaseCartProduct(ctx, userId, productId)
}
func (r *UserService) DeleteItemInCart(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error {

	return r.userRepo.DeleteItemInCart(ctx, userId, productId)
}
func (r *UserService) ClearCart(ctx context.Context, userId primitive.ObjectID) error {

	return r.userRepo.ClearCart(ctx, userId)
}
