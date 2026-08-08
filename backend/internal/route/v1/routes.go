package v1

import (
	"time"
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, enableRateLimit bool) {
	api := app.Group("/api/v1")

	// Auth routes - không cần token
	auth := api.Group("/auth")
	if enableRateLimit {
		auth.Use(limiter.New(limiter.Config{
			Max:        10,
			Expiration: 1 * time.Minute,
		}))
	}
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", middleware.AuthRequired, authHandler.GetMe)

	// Task routes - bắt buộc phải có JWT hợp lệ
	tasks := api.Group("/tasks", middleware.AuthRequired)
	tasks.Post("/", taskHandler.CreateTask)
	tasks.Get("/", taskHandler.GetTasks)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Delete("/:id", taskHandler.DeleteTask)
}
