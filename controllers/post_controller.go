package controllers

import (
	model "backend/models"
	"backend/services"
	"backend/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func CreatePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(utils.PARSE_FORM)
	}


	file := form.File["file"]
	if len(file) > 0 {
		uploadedFile, err := file[0].Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_OPEN_FILE)
		}
		defer uploadedFile.Close()
		uploadDir := filepath.Join("Files", strconv.Itoa(int(user.ID)))

		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating directory")
		}
		createfile := filepath.Join(uploadDir, file[0].Filename)
		
		deviceFile, err := os.Create(createfile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_CREATE_FILE)
		}
		defer deviceFile.Close()

		_, err = io.Copy(deviceFile, uploadedFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(utils.ERR_COPY_FILE)
		}
	}

	post := &model.Post{
		UserID   :user.ID,
		Caption   : form.Value["caption"][0],
		MediaURL  : "Files/"+strconv.Itoa(int(user.ID))+"/" + file[0].Filename,
		MediaType: form.Value["media_type"][0],
		Comments : []model.Comment{},
		Likes     : []model.Like{},
	}
	_ ,err = services.CreatePost(*post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	
	return c.Status(fiber.StatusOK).SendString("Successfully Updated User")
}


func GetPosts(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	
	posts ,err := services.GetPostUsingId(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}


	return c.JSON(posts)
}