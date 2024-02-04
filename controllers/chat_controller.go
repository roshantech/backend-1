package controllers

import (
	model "backend/models"
	"backend/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllConversations(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	conversations, err := services.GetAllConversations(strconv.Itoa(int(user.ID)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.JSON(conversations)
}
