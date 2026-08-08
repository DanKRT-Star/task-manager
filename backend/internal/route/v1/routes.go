package v1

import (
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) {
	api := app.Group("/api/v1")

	// Auth routes - không cần token
	auth := api.Group("/auth")
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
