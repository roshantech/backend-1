package utils

import (
	model "backend/models"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(email string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = email
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix()
	t, err := token.SignedString([]byte("cd626409f4a042f353b6f532aa7da9ed803fcfc081e41599f8d97813250c06e0daaa30688159ad791f5dccfee7f6b2f9f801b5758f4256a752d9e5aee7f2c6bb"))
	if err != nil {
		return t, t, err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = email
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := refreshToken.SignedString([]byte("0011356368826a9af0ce158050c9dcb10e7292256c40dfba6ffa7f9b8a0c06ddfd667683d43105841a33136aa836a70700899309d7871a4ffe920b73a10ae0f4"))
	if err != nil {
		return t, t, err
	}
	//return token and refresh token
	return t, rt, nil
}

// RandomString generates a secure random string of length n.
func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandomInt returns a cryptographically secure random integer
func RandomInt() int {
	return rand.Intn(9999999)
}

func JWTConfig() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("cd626409f4a042f353b6f532aa7da9ed803fcfc081e41599f8d97813250c06e0daaa30688159ad791f5dccfee7f6b2f9f801b5758f4256a752d9e5aee7f2c6bb"),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT"})
	} else {
		return c.Status(fiber.StatusForbidden).JSON(err.Error())
	}
}

func JWTFilter(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	userEmail := token.Claims.(jwt.MapClaims)["sub"].(string)
	var user model.User
	err := db.Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		return err
	}
	var address model.Address
	err = db.Where("user_id = ?", user.ID).First(&address).Error
	if err != nil {
		return err
	}
	user.Address = address
	c.Locals("user", user)
	return c.Next()
}
