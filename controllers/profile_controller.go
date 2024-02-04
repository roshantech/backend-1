package controllers

import (
	"backend/dto"
	model "backend/models"
	"backend/services"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfile(c *fiber.Ctx) error {
	ID := c.Query("ID")

	var userInfo dto.UserInfo
	if err := c.BodyParser(&userInfo); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	user, err := services.GetUserByID(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	user.About = userInfo.About
	user.Company = userInfo.Company
	user.Education = userInfo.School
	user.Email = user.Email
	user.Address.Country = userInfo.Country
	user.PhoneNumber = userInfo.PhoneNo

	_, err = services.UpdateUserById(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.Status(fiber.StatusOK).SendString("Successfully Updated User")
}

func GetUserFollowers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	followers, err := services.GetUserFollowers(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.JSON(followers)
}

func AddFollowerToUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	ID := c.Query("FID")
	FID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	err = services.AddFollowerToUser(user.ID, uint(FID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.SendString("Follower Added Successfully!")
}
