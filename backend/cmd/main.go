package main

import (
	"log"
	"os"
	"time"

	apperror "github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/config"
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	v1 "github.com/DanKRT-Star/task-manager/internal/route/v1"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/joho/godotenv"

	_ "github.com/DanKRT-Star/task-manager/docs"
	"github.com/gofiber/contrib/v3/swaggo"
)

// @title           Task Manager API
// @version         1.0
// @description     REST API for a task management application with JWT authentication.
// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token.

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.ConnectDatabase()
	if err := config.InitSchema(); err != nil {
		log.Fatal("Failed to initialize database schema: ", err)
	}

	userRepo := repository.NewUserRepository(config.DB)
	taskRepo := repository.NewTaskRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	memberRepo := repository.NewProjectMemberRepository(config.DB)
	commentRepo := repository.NewCommentRepository(config.DB)
	labelRepo := repository.NewLabelRepository(config.DB)
	activityRepo := repository.NewActivityLogRepository(config.DB)
	epicRepo := repository.NewEpicRepository(config.DB)
	milestoneRepo := repository.NewMilestoneRepository(config.DB)
	sprintRepo := repository.NewSprintRepository(config.DB)

	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo, memberRepo, activityRepo)
	projectService := service.NewProjectService(projectRepo, memberRepo, userRepo)
	epicService := service.NewEpicService(epicRepo, memberRepo)
	milestoneService := service.NewMilestoneService(milestoneRepo, memberRepo)
	sprintService := service.NewSprintService(sprintRepo, memberRepo)
	commentService := service.NewCommentService(commentRepo, taskRepo, memberRepo)
	labelService := service.NewLabelService(labelRepo, taskRepo, memberRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	projectHandler := handler.NewProjectHandler(projectService)
	epicHandler := handler.NewEpicHandler(epicService)
	milestoneHandler := handler.NewMilestoneHandler(milestoneService)
	sprintHandler := handler.NewSprintHandler(sprintService)
	commentHandler := handler.NewCommentHandler(commentService)
	labelHandler := handler.NewLabelHandler(labelService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3001", "https://task-manager-phi-one-73.vercel.app"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := config.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(503).JSON(fiber.Map{"status": "error", "database": "unreachable"})
		}
		return c.JSON(fiber.Map{"status": "ok", "message": "Server is running"})
	})

	app.Get("/swagger/*", swaggo.New(swaggo.Config{}))

	v1.SetupRoutes(app, authHandler, taskHandler, projectHandler, epicHandler, milestoneHandler, sprintHandler, commentHandler, labelHandler, true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
