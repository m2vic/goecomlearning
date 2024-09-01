package service

import (
	"context"
	"golearning/internal/core/domain"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenGenerator struct {
}

func (g *TokenGenerator) GenerateToken(ctx context.Context, username string, id primitive.ObjectID) (*domain.Token, error) {
	accessKey := os.Getenv("PASS")
	refreshKey := os.Getenv("REFRESHTOKEN")
	exp := time.Now().Add(time.Hour * 1).Unix()
	expRT := time.Now().Add(time.Hour * 1).Unix()
	role := ""
	if username == "admin" {
		role = "admin"
	}
	setAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": username, "userid": id, "role": role, "exp": exp})
	accessTokenString, err := setAccessToken.SignedString([]byte(accessKey))
	if err != nil {
		log.Fatal(err)
	}
	setRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": username, "userid": id, "role": role, "exp": expRT})
	refreshTokenString, err := setRefreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		log.Fatal("from refresh:", err)
	}
	token := domain.Token{AccessToken: accessTokenString, RefreshToken: refreshTokenString}
	return &token, err
}
