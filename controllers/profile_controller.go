package controllers

import (
	"backend/dto"
	"backend/services"
	"log"

	"github.com/gofiber/fiber/v2"
)


func UpdateProfile(c *fiber.Ctx) error {
	ID := c.Query("ID")

	var userInfo dto.UserInfo
	if err := c.BodyParser(&userInfo); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	user ,err :=services.GetUserByID(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	user.About = userInfo.About
	user.Company = userInfo.Company
	user.Education = userInfo.School
	user.Email = user.Email 
	user.Address.Country = userInfo.Country
	user.PhoneNumber = userInfo.PhoneNo

	_,err = services.UpdateUserById(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.Status(fiber.StatusOK).SendString("Successfully Updated User")
}