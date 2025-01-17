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
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/markbates/pkger"
)

func main() {

	utils.GetDB()
	defer utils.CloseDB()

	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: "Fiber",
		AppName:      "FakeBook 1.0",
		BodyLimit:    1024 * 1024 * 1024,
	})

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			MaxAge:       3600,
		},
	))

	app.Use("/Files", filesystem.New(filesystem.Config{
		Root:   pkger.Dir("/Files"),
		Browse: true,
	}))
	public := app.Group("/core")
	public.Post("/signup", controllers.Signup)
	public.Post("/login", controllers.Login)
	v1 := app.Group("/v1", utils.JWTConfig())
	private := v1.Group("/core", utils.JWTFilter)
	private.Get("/getLoggedInUser", controllers.GetLoggedInUser)
	private.Get("/getUserByID", controllers.GetUserByID)
	private.Get("/getAllUsers", controllers.GetAllUsers)

	private.Get("/getProfilePic", controllers.GetJobFile)
	private.Post("/updateProfile", controllers.UpdateProfile)

	private.Get("/getPosts", controllers.GetPosts)
	private.Post("/createPost", controllers.CreatePost)
	private.Post("/createPostComments", controllers.CreatePostComments)
	private.Get("/getPostComments", controllers.GetPostComments)
	private.Get("/getUserFollowers", controllers.GetUserFollowers)
	private.Get("/addFollowerToUser", controllers.AddFollowerToUser)
	private.Get("/getAllConversations", controllers.GetAllConversations)

	private.Get("/likePost", controllers.LikePost)
	public.Use("/ws", controllers.NotificationWs)
	private.Get("/getAllNotifications", controllers.GetAllNotifications)

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
