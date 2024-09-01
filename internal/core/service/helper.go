package service

import (
	"log"

	"github.com/golang-jwt/jwt"
)

func GetIdByToken(Token *jwt.Token) string {
	claims := Token.Claims.(jwt.MapClaims)
	id, ok := claims["userid"]
	if !ok {
		log.Fatal("field userid not found")
	}

	userIdString := id.(string)
	return userIdString
}
func GetRoleByToken(Token *jwt.Token) string {
	claims := Token.Claims.(jwt.MapClaims)
	role, ok := claims["role"]
	if !ok {
		log.Fatal("field userid not found")
	}

	userIdString := role.(string)
	return userIdString
}
