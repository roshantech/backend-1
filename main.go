package main

import (
	"backend/controllers"
	"backend/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	utils.GetDB()
	defer utils.CloseDB()

	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			MaxAge:       3600,
		},
	))
	public := app.Group("/core")

	public.Post("/signup", controllers.Signup)
	public.Post("/login", controllers.Login)

	v1 := app.Group("/v1", utils.JWTConfig())
	private := v1.Group("/core", utils.JWTFilter)
	private.Get("/all", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})
	private.Get("/getProfilePic" ,controllers.GetJobFile)
	go func() {
		if err := app.Listen(":3001"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown the server gracefully
	shutdownTimeout := 10 * time.Second
	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

}
