package controllers

import (
	"backend/dto"
	"backend/services"
	"backend/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	// Parse the request body into a UserLogin DTO
	var userLogin dto.UserLogin
	if err := c.BodyParser(&userLogin); err != nil {
		// Return a 400 Bad Request response if there's an error parsing the request body
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Retrieve the user from the database based on the provided username
	user, err := services.GetUserByUsername(userLogin.Username)
	if err != nil {
		// Handle user retrieval errors
		log.Println("Error retrieving user:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// Verify the user's password
	if err := user.VerifyPassword(userLogin.Password); err != nil {
		// Handle password verification errors
		log.Println("Error verifying password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateToken(user.Email)
	if err != nil {
		// Handle token generation errors
		log.Println("Error generating tokens:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	// Return a successful response with access and refresh tokens, and user information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func Signup(c *fiber.Ctx) error {
	var user dto.UserCreateUpdate
	err := c.BodyParser(&user)
	if err != nil {
		fmt.Println(err)
	}
	_, err = services.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.SendString("User Created Successfully")
}
