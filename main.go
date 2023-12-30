package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&User{})
}

func main() {
	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			MaxAge:       3600,
		},
	))
	public := app.Group("/core")

	public.Post("/login", func(c *fiber.Ctx) error {
		var userd UserDto
		err := c.BodyParser(&userd)
		if err != nil {
			fmt.Println(err)
		}
		var user User

		err = db.Table("users").Where("username = ?", userd.Username).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User not found",
				})
			}

			// Handle other errors
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		t, rt, err := GenerateToken(user.Email)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"accessToken":        t,
			"refreshToken": rt,
			"user":         user,
		})
	})

	public.Post("/signup", func(c *fiber.Ctx) error {
		var user User
		err := c.BodyParser(&user)
		if err != nil {
			fmt.Println(err)
		}

		err = db.Table("users").Create(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		return c.SendString("User Created Successfully")
	})

	// v1 := app.Group("/v1", JWTConfig())
	// private := v1.Group("/core", JWTFilter)

	app.Listen(":3001")

}
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
	userEmail := token.Claims.(jwt.MapClaims)["sub"].(float64)
	var user User
	err := db.Where("ID = ?", userEmail).First(&user).Error
	if err != nil {
		return err
	}
	c.Locals("user", user)
	return c.Next()
}
