package controllers

import (
	"backend/dto"
	model "backend/models"
	"backend/services"
	"backend/utils"
	"io"
	"log"
	"os"
	"time"

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
		user, err = services.GetUserByEmail(userLogin.Username)
		if err != nil {
			log.Println("Error retrieving user:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
		}
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
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(utils.PARSE_FORM)
	}

	isEmailPresent, err1 := services.ValidateEmailId(form.Value["email"][0])
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err1.Error())
	}

	if isEmailPresent {
		return c.Status(fiber.StatusInternalServerError).SendString("Error:Email_Id already exists. Kindly provide different email id.")
	}

	file := form.File["file"]
	if len(file) > 0 {
		uploadedFile, err := file[0].Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_OPEN_FILE)
		}
		defer uploadedFile.Close()

		deviceFile, err := os.Create("Files/" + file[0].Filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_CREATE_FILE)
		}
		defer deviceFile.Close()

		_, err = io.Copy(deviceFile, uploadedFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_COPY_FILE)
		}
	}

	user := &model.User{
		Username:   form.Value["username"][0],
		Password:   form.Value["password"][0],
		Email:      form.Value["email"][0],
		ProfilePic: "Files/" + form.File["file"][0].Filename,
		Active:     true,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err = services.CreateUser(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.SendString("User Created Successfully")
}
