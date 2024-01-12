package controllers

import "github.com/gofiber/fiber/v2"

func GetJobFile(c *fiber.Ctx) error {
	idStr := c.Query("loc")
	return c.SendFile(idStr)
}