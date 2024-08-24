package service_test

import (
	"context"
	"golearning/internal/adapter/repo"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUser(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	expected := &domain.User{FirstName: "BOB", LastName: "CALLAWAY"}
	ctx := context.Background()
	userId := "1"
	userRepo.On("GetUser", ctx, userId).Return(expected)

	actual, err := userService.GetUser(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
func TestLogin(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	username := "123"
	password := "123"
	userRepo.On("Login", ctx, username, password).Return(expected)

	actual, err := userService.Login(ctx, username, password)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestRegister(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	ctx := context.Background()
	username := "123"
	password := "123"
	email := "123"
	data := domain.User{Email: email, Username: username, Password: password}
	userRepo.On("CheckUsername", ctx, username).Return(false)
	userRepo.On("CheckEmail", ctx, email).Return(false)
	userRepo.On("Register", ctx, username, password, email).Return(nil)

	actual := userService.Register(ctx, data)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestUpdateUser(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()

	data := domain.User{FirstName: "BOB"}
	userId := "1"
	userRepo.On("UpdateUser", ctx, data, userId).Return(nil)
	actual := userService.UpdateUser(ctx, data, userId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestChangePassword(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId := "1"
	oldPass := "1"
	newPass := "2"
	userRepo.On("ChangePassword", ctx, userId, oldPass, newPass).Return(nil)
	actual := userService.ChangePassword(ctx, userId, oldPass, newPass)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestResetPassword(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	ctx := context.Background()
	email := "123"
	userRepo.On("CheckEmail", ctx, email).Return(false)
	userRepo.On("ResetPassword", ctx, email).Return("OK")

	actual, err := userService.ResetPassword(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, "OK", actual)
}
func TestCheckRefresh(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)
	ctx := context.Background()
	expect := &domain.Token{AccessToken: "token", RefreshToken: "tokenbutred"}
	refreshToken := "aeiou"
	userRepo.On("CheckRefresh", ctx, refreshToken).Return(expect)
	actual, err := userService.CheckRefresh(ctx, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestGetCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	//productCached := new(repo.MockProductCached)
	//productRepo := new(repo.MockProductRepo)
	mockProductService := new(service.MockProductService)
	//productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, mockProductService)

	productId := "123456789012345678901234"
	id, _ := primitive.ObjectIDFromHex(productId)
	var carts []domain.Cart
	cart := domain.Cart{ProductId: id}
	carts = append(carts, cart)
	expected := carts
	ctx := context.Background()
	userId := "1"

	mockProductService.On("CheckAmount", ctx, id).Return(int(12), nil)
	userRepo.On("GetCart", ctx, userId).Return(expected)

	actual, err := userService.GetCart(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddToCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productService := new(service.MockProductService)
	userService := service.NewUserService(userRepo, productService)
	ctx := context.Background()
	productId := "1"
	primitiveProductId, _ := primitive.ObjectIDFromHex(productId)
	product := domain.Cart{ProductName: "Cupid", Amount: 1}
	userId := "1"

	userRepo.On("AddtoCart", ctx, product, userId).Return(nil)
	productService.On("CheckAmount", ctx, primitiveProductId).Return(1)

	actual := userService.AddtoCart(ctx, product, userId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestInCreaseCartProduct(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)
	ctx := context.Background()
	userId := "1"
	productId := primitive.NewObjectID()
	userRepo.On("IncreaseCartProduct", ctx, userId, productId).Return(nil)

	actual := userService.IncreaseCartProduct(ctx, userId, productId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
func TestDeCreaseCartProduct(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	ctx := context.Background()
	userId := "1"
	productId := primitive.NewObjectID()
	userRepo.On("DecreaseCartProduct", ctx, userId, productId).Return(nil)

	actual := userService.DecreaseCartProduct(ctx, userId, productId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
func TestDeleteItemInCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)
	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId := "1"
	productId := primitive.NewObjectID()
	userRepo.On("DeleteItemInCart", ctx, userId, productId).Return(nil)

	actual := userService.DeleteItemInCart(ctx, userId, productId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
func TestClearCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	userService := service.NewUserService(userRepo, productService)

	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId := "1"
	userRepo.On("ClearCart", ctx, userId).Return(nil)

	actual := userService.ClearCart(ctx, userId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
