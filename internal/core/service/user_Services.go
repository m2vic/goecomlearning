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
}

func NewUserService(userRepo port.UserRepo, productService port.ProductService) *UserService {
	return &UserService{userRepo: userRepo, ProductService: productService}
}

func (s *UserService) Login(ctx context.Context, username, password string) (*domain.Token, error) {
	login, err := s.userRepo.Login(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("from service:%w", err)
	}
	return login, nil
}

func (s *UserService) Register(ctx context.Context, user domain.User) error {
	username := user.Username
	password := user.Password
	email := user.Email
	userExist, err := s.userRepo.CheckUsername(ctx, username)
	if err != nil {
		return err
	}
	if userExist {
		return fmt.Errorf("user already exist")
	}
	emailExist, err := s.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return err
	}
	if emailExist {
		return fmt.Errorf("email already exist")
	}
	// have to check username and email
	return s.userRepo.Register(ctx, username, password, email)
}

func (s *UserService) GetUser(ctx context.Context, userId string) (*domain.User, error) {

	return s.userRepo.GetUser(ctx, userId)
}

func (s *UserService) UpdateUser(ctx context.Context, info domain.User, userId string) error {

	return s.userRepo.UpdateUser(ctx, info, userId)
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
func (s *UserService) CheckRefresh(ctx context.Context, refreshtoken string) (*domain.Token, error) {
	refresh, err := s.userRepo.CheckRefresh(ctx, refreshtoken)
	if err != nil {
		return nil, err
	}
	return refresh, nil
}
func (r *UserService) AddtoCart(ctx context.Context, Product domain.Cart, userId string) error {
	productId := Product.ProductId
	reqAmount := Product.Amount
	actualAmount, err := r.ProductService.CheckAmount(ctx, productId)
	if err != nil {
		return errors.New("no product found")
	}
	if reqAmount > *actualAmount {
		return errors.New("not enough product in stock")
	}
	return r.userRepo.AddtoCart(ctx, Product, userId)

}
func (r *UserService) GetCart(ctx context.Context, userId string) ([]domain.Cart, error) {
	// have to check if sufficient and then compare
	var newCart []domain.Cart
	userCart, err := r.userRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, item := range userCart {
		amount, err := r.ProductService.CheckAmount(ctx, item.ProductId)
		if *amount >= item.Amount {
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
