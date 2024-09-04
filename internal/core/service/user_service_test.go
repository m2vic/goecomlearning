package service_test

import (
	"context"
	"golearning/internal/adapter/repo"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	errs "golearning/internal/error"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUser(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	expected := &domain.User{FirstName: "BOB", LastName: "CALLAWAY"}
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
	userRepo.On("GetUser", ctx, userId).Return(expected)

	actual, err := userService.GetUser(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
func TestLogin(t *testing.T) {

	type testcase struct {
		username    string
		password    string
		expected    *domain.Token
		expectedErr string
	}
	testcases := []testcase{{username: "123", password: "123", expected: &domain.Token{AccessToken: "access", RefreshToken: "refresh"}, expectedErr: ""},
		{username: "wrong", password: "right", expected: nil, expectedErr: "username"},
		{username: "right", password: "wrong", expected: nil, expectedErr: "password"},
	}
	for _, tc := range testcases {
		t.Run(tc.username, func(t *testing.T) {

			userRepo := new(repo.MockUserRepo)
			productCached := new(repo.MockProductCached)
			productRepo := new(repo.MockProductRepo)
			productService := service.NewProductService(productRepo, productCached, userRepo)
			tokenGenerator := new(service.MockTokenGenerator)
			passwordHasher := new(service.MockPasswordHasher)
			passwordGenerator := new(service.MockPasswordGenerator)
			userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
			ctx := context.Background()
			id, _ := primitive.ObjectIDFromHex("1")
			user := &domain.User{Username: tc.username, Password: tc.password, UserId: id}

			switch tc.expectedErr {
			case "username":
				userRepo.On("CheckUsername", ctx, tc.username).Return(nil, errs.UsernameInvalid)
			case "password":
				userRepo.On("CheckUsername", ctx, tc.username).Return(user, nil)
				passwordHasher.On("ComparePassword", ctx, tc.password, tc.password).Return(errs.PasswordInvalid)
			default:
				userRepo.On("CheckUsername", ctx, tc.username).Return(user, nil)
				tokenGenerator.On("GenerateToken", ctx, tc.username, user.UserId).Return(tc.expected)
				passwordHasher.On("ComparePassword", ctx, tc.password, tc.password).Return(nil)
				userRepo.On("UpdateUser", ctx, domain.User{UserId: id, Token: tc.expected.RefreshToken}).Return(nil)
			}

			actual, err := userService.Login(ctx, tc.username, tc.password)

			switch tc.expectedErr {
			case "username":
				assert.ErrorIs(t, err, errs.UsernameInvalid)
				assert.Nil(t, actual)
			case "password":
				assert.ErrorIs(t, err, errs.PasswordInvalid)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}

}

func TestRegister(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)
	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	ctx := context.Background()

	username := "123"
	password := "123"
	email := "123"
	hash := []byte("123")
	data := &domain.User{Email: email, Username: username, Password: password}
	userRepo.On("CheckUsername", ctx, username).Return(&domain.User{})
	passwordHasher.On("HashPassword", ctx, password).Return(hash)
	userRepo.On("CheckEmail", ctx, email).Return(false)
	userRepo.On("Register", ctx, username, email, hash).Return(nil)

	actual := userService.Register(ctx, *data)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestUpdateUser(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)
	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()

	userId := "1"
	id, _ := primitive.ObjectIDFromHex(userId)
	data := domain.User{FirstName: "BOB", UserId: id}
	userRepo.On("UpdateUser", ctx, data).Return(nil)
	actual := userService.UpdateUser(ctx, data)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestChangePassword(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
	oldPass := "12345678"
	newPass := "12345678"
	passwordHasher.On("HashPassword", ctx, oldPass).Return([]byte(oldPass))
	userRepo.On("CheckPassword", ctx, oldPass).Return(nil)
	passwordHasher.On("HashPassword", ctx, newPass).Return([]byte(newPass))
	userRepo.On("ChangePassword", ctx, userId, newPass).Return(nil)
	actual := userService.ChangePassword(ctx, userId, oldPass, newPass)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestResetPassword(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)
	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	ctx := context.Background()
	email := "123"
	userRepo.On("CheckEmail", ctx, email).Return(false)
	passwordGenerator.On("RandomPassword").Return("123")
	passwordHasher.On("HashPassword", ctx, "123").Return([]byte("hash"))
	userRepo.On("ResetPassword", ctx, email, "hash").Return("OK")

	actual, err := userService.ResetPassword(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, "OK", actual)
}
func TestRefreshToken(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	productRepo := new(repo.MockProductRepo)
	productService := service.NewProductService(productRepo, productCached, userRepo)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex("1")
	refreshToken := "aeiou"
	expect := &domain.Token{AccessToken: "token", RefreshToken: refreshToken}
	user := domain.User{UserId: id, Token: refreshToken}
	userRepo.On("CheckRefresh", ctx, refreshToken).Return(&domain.User{Username: "somename"})
	tokenGenerator.On("GenerateToken", ctx, "somename", user.UserId).Return(expect)
	userRepo.On("UpdateUser", ctx, user).Return(nil)
	actual, err := userService.RefreshToken(ctx, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestGetCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	//productCached := new(repo.MockProductCached)
	//productRepo := new(repo.MockProductRepo)
	productService := new(service.MockProductService)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	productId := "123456789012345678901234"
	id, _ := primitive.ObjectIDFromHex(productId)
	var carts []domain.Cart
	cart := domain.Cart{ProductId: id}
	carts = append(carts, cart)
	expected := carts
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")

	productService.On("CheckAmount", ctx, id).Return(int(12), nil)
	userRepo.On("GetCart", ctx, userId).Return(expected)

	actual, err := userService.GetCart(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddToCart(t *testing.T) {
	userRepo := new(repo.MockUserRepo)
	productService := new(service.MockProductService)
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
	ctx := context.Background()
	productId := "1"
	primitiveProductId, _ := primitive.ObjectIDFromHex(productId)
	product := domain.Cart{ProductName: "Cupid", Amount: 1}
	userId, _ := primitive.ObjectIDFromHex("1")

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
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
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
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
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
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)
	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
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
	tokenGenerator := new(service.MockTokenGenerator)
	passwordHasher := new(service.MockPasswordHasher)

	passwordGenerator := new(service.MockPasswordGenerator)
	userService := service.NewUserService(userRepo, productService, tokenGenerator, passwordHasher, passwordGenerator)

	//expected := &domain.Token{AccessToken: "access", RefreshToken: "refresh"}
	ctx := context.Background()
	userId, _ := primitive.ObjectIDFromHex("1")
	userRepo.On("ClearCart", ctx, userId).Return(nil)

	actual := userService.ClearCart(ctx, userId)

	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
