package service

import (
	"context"
	"errors"
	"fmt"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo       port.UserRepo
	ProductService port.ProductService
	TokenGenerator port.TokenGenerator
	PasswordHasher port.PasswordHasher
}

func NewUserService(userRepo port.UserRepo, productService port.ProductService, tokenGenerator port.TokenGenerator, passwordHasher port.PasswordHasher) *UserService {
	return &UserService{userRepo: userRepo, ProductService: productService, TokenGenerator: tokenGenerator, PasswordHasher: passwordHasher}
}

func (s *UserService) Login(ctx context.Context, username, password string) (*domain.Token, error) {
	result, err := s.userRepo.CheckUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	dbHashedPassword := result.Password
	err = s.PasswordHasher.ComparePassword(ctx, password, dbHashedPassword)
	if err != nil {
		return nil, err
	}
	token, err := s.TokenGenerator.GenerateToken(ctx, result.Username, result.UserId)
	if err != nil {
		return nil, err
	}
	update := domain.User{UserId: result.UserId, Token: token.RefreshToken}
	err = s.UpdateUser(ctx, update)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserService) Register(ctx context.Context, user domain.User) error {
	username := user.Username
	password := user.Password
	email := user.Email
	result, err := s.userRepo.CheckUsername(ctx, username)
	if err != nil {
		return err
	}
	if result.Username == username {
		return fmt.Errorf("user already exist")
	}
	emailExist, err := s.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return err
	}
	if emailExist {
		return fmt.Errorf("email already exist")
	}

	hash, err := s.PasswordHasher.HashPassword(ctx, password)
	if err != nil {
		return err
	}
	// have to check username and email
	return s.userRepo.Register(ctx, username, email, hash)
}

func (s *UserService) GetUser(ctx context.Context, userId string) (*domain.User, error) {

	return s.userRepo.GetUser(ctx, userId)
}

func (s *UserService) UpdateUser(ctx context.Context, info domain.User) error {

	return s.userRepo.UpdateUser(ctx, info)
}
func (s *UserService) ChangePassword(ctx context.Context, userId, oldPass, newPass string) error {

	return s.userRepo.ChangePassword(ctx, userId, oldPass, newPass)
}

func (s *UserService) ResetPassword(ctx context.Context, email string) (string, error) {
	_, err := s.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("err:%w", err)
	}
	return s.userRepo.ResetPassword(ctx, email)

}
func (s *UserService) RefreshToken(ctx context.Context, refreshtoken string) (*domain.Token, error) {

	result, err := s.userRepo.CheckRefresh(ctx, refreshtoken)
	if err != nil {
		return nil, err
	}
	tokens, err := s.TokenGenerator.GenerateToken(ctx, result.Username, result.UserId)
	if err != nil {
		return nil, err
	}
	//user := domain.User{UserId: result.UserId, Token: tokens.RefreshToken}
	err = s.userRepo.UpdateUser(ctx, domain.User{UserId: result.UserId, Token: tokens.RefreshToken})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *UserService) AddtoCart(ctx context.Context, Product domain.Cart, userId string) error {
	productId := Product.ProductId
	reqAmount := Product.Amount
	actualAmount, err := r.ProductService.CheckAmount(ctx, productId)
	if err != nil {
		return errors.New("no product found")
	}
	if reqAmount > actualAmount {
		return errors.New("not enough product in stock")
	}
	return r.userRepo.AddtoCart(ctx, Product, userId)

}
func (r *UserService) GetCart(ctx context.Context, userId string) ([]domain.Cart, error) {
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
func (r *UserService) IncreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {

	return r.userRepo.IncreaseCartProduct(ctx, userId, productId)
}
func (r *UserService) DecreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {

	return r.userRepo.DecreaseCartProduct(ctx, userId, productId)
}
func (r *UserService) DeleteItemInCart(ctx context.Context, userId string, productId primitive.ObjectID) error {

	return r.userRepo.DeleteItemInCart(ctx, userId, productId)
}
func (r *UserService) ClearCart(ctx context.Context, userId string) error {

	return r.userRepo.ClearCart(ctx, userId)
}
