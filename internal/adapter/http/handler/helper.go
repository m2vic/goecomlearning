package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	err := godotenv.Load("/Users/m2vic/Desktop/gosystem/cmd/.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func getRefresh(c *fiber.Ctx) string {
	auth := c.Get("Refresh", "No-Token")
	if auth == "No-Token" {
		return "No-Token"
	}
	return auth
}
func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	parsedToken, err := Authenticate(token)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	userID := claims["userid"].(string)
	role := claims["role"].(string)
	fmt.Println(role)
	// Inject user info into the context
	c.Locals("userID", userID)
	c.Locals("role", role)

	// Attach user to the context
	return c.Next() // Proceed to the next handler
}
func Authenticate(Token interface{}) (*jwt.Token, error) {
	tokenstring, ok := Token.(string)
	if !ok {
		return nil, fmt.Errorf("token is not a string")
	}
	tokenstring = strings.ReplaceAll(tokenstring, "Bearer ", "")
	pass := os.Getenv("PASS")
	if pass == "" {
		return nil, errors.New("missing JWT signing key in environment")
	}
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
			return nil, errors.New("token expired")
		}
	} else {
		return nil, errors.New("invalid token")
	}
	return token, nil

}

func registerNotification(receiverEmail string) error {
	env := godotenv.Load()
	if env != nil {
		fmt.Println("fail to load env")
	}
	sender := os.Getenv("EMAILSENDER")
	smtp := os.Getenv("SMTP")
	smtpport := os.Getenv("SMTPPORT")
	port, err := strconv.Atoi(smtpport)
	if err != nil {
		return fmt.Errorf("fail to convert smtpport to int")
	}
	emailPassword := os.Getenv("EMAILPASSWORD")
	subject := "Registration Successful!"
	text := "Welcome to our website,<br><p>we hope you have great experience on our shopping sotre<p>Best regards"

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiverEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(smtp, port, sender, emailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("fail to dial and send:%w", err)
	}
	return nil
}

func resetPasswordEmail(receiverEmail, newPassword string) error {
	env := godotenv.Load()
	if env != nil {
		fmt.Println("fail to load env")
	}
	sender := os.Getenv("EMAILSENDER")
	smtp := os.Getenv("SMTP")
	smtpport := os.Getenv("SMTPPORT")
	port, err := strconv.Atoi(smtpport)
	if err != nil {
		return fmt.Errorf("fail to convert smtpport to int")
	}
	emailPassword := os.Getenv("EMAILPASSWORD")
	subject := "Reset Password"

	text := fmt.Sprintf("Hello, System generate a temporary password for you, Sign in to change your password as you wish<br><p>%s<p>", newPassword)

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiverEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(smtp, port, sender, emailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("fail to dial and send:%w", err)
	}
	return nil
}
