package v1

import (
	"time"

	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	taskHandler *handler.TaskHandler,
	projectHandler *handler.ProjectHandler,
	epicHandler *handler.EpicHandler,
	milestoneHandler *handler.MilestoneHandler,
	sprintHandler *handler.SprintHandler,
	commentHandler *handler.CommentHandler,
	labelHandler *handler.LabelHandler,
	enableRateLimit bool,
) {
	api := app.Group("/api/v1")

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

	tasks := api.Group("/tasks", middleware.AuthRequired)
	tasks.Post("/", taskHandler.CreateTask)
	tasks.Get("/", taskHandler.GetTasks)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Delete("/:id", taskHandler.DeleteTask)
	tasks.Post("/:taskId/comments", commentHandler.CreateComment)
	tasks.Get("/:taskId/comments", commentHandler.GetTaskComments)
	tasks.Delete("/:taskId/comments/:commentId", commentHandler.DeleteComment)
	tasks.Post("/:taskId/labels/:labelId", labelHandler.AttachLabel)
	tasks.Delete("/:taskId/labels/:labelId", labelHandler.DetachLabel)
	tasks.Get("/:taskId/labels", labelHandler.GetTaskLabels)

	epics := api.Group("/epics", middleware.AuthRequired)
	epics.Put("/:epicId", epicHandler.UpdateEpic)
	epics.Delete("/:epicId", epicHandler.DeleteEpic)

	milestones := api.Group("/milestones", middleware.AuthRequired)
	milestones.Put("/:milestoneId", milestoneHandler.UpdateMilestone)
	milestones.Delete("/:milestoneId", milestoneHandler.DeleteMilestone)

	sprints := api.Group("/sprints", middleware.AuthRequired)
	sprints.Put("/:sprintId", sprintHandler.UpdateSprint)
	sprints.Delete("/:sprintId", sprintHandler.DeleteSprint)

	projects := api.Group("/projects", middleware.AuthRequired)
	projects.Post("/", projectHandler.CreateProject)
	projects.Get("/", projectHandler.GetProjects)
	projects.Get("/:id", projectHandler.GetProject)
	projects.Put("/:id", projectHandler.UpdateProject)
	projects.Delete("/:id", projectHandler.DeleteProject)
	projects.Get("/:id/tasks", taskHandler.GetProjectTasks)
	projects.Get("/:id/members", projectHandler.GetMembers)
	projects.Post("/:id/members", projectHandler.AddMember)
	projects.Delete("/:id/members/:userId", projectHandler.RemoveMember)
	projects.Post("/:id/sprints", sprintHandler.CreateSprint)
	projects.Get("/:id/sprints", sprintHandler.GetProjectSprints)
	projects.Post("/:id/labels", labelHandler.CreateLabel)
	projects.Get("/:id/labels", labelHandler.GetProjectLabels)
	projects.Delete("/:id/labels/:labelId", labelHandler.DeleteLabel)
	projects.Post("/:id/epics", epicHandler.CreateEpic)
	projects.Get("/:id/epics", epicHandler.GetProjectEpics)
	projects.Post("/:id/milestones", milestoneHandler.CreateMilestone)
	projects.Get("/:id/milestones", milestoneHandler.GetProjectMilestones)
}
