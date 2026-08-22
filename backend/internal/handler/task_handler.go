package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type TaskHandler struct {
	TaskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{TaskService: taskService}
}

// CreateTask godoc
// @Summary      Create a task
// @Description  Create a new task for the authenticated user, optionally attached to a project, epic, milestone, and assigned to a project member
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateTaskRequest true "Task payload"
// @Success      201 {object} model.Task
// @Failure      400 {object} map[string]string "validation error"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.CreateTaskRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	task, err := h.TaskService.CreateTask(userID, req.Title, req.Description, req.Status, req.Deadline, req.ProjectID, req.AssigneeID, req.EpicID, req.MilestoneID, req.SprintID)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(task)
}

// GetTasks godoc
// @Summary      List tasks
// @Description  Get a paginated list of tasks for the authenticated user, with optional filtering and sorting
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        status query string false "Filter by status" Enums(pending, in_progress, done)
// @Param        sort   query string false "Sort by deadline" Enums(deadline_asc, deadline_desc)
// @Param        page   query int    false "Page number" default(1)
// @Param        limit  query int    false "Items per page (max 100)" default(10)
// @Success      200 {object} dto.TaskListResponse
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /tasks [get]
func (h *TaskHandler) GetTasks(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	status := c.Query("status")
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tasks, total, err := h.TaskService.GetTasks(userID, status, sort, page, limit)
	if err != nil {
		return apperror.Internal("failed to fetch tasks")
	}

	return c.JSON(dto.TaskListResponse{
		Data:  tasks,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update a task owned by the authenticated user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Task ID"
// @Param        request body dto.UpdateTaskRequest true "Fields to update"
// @Success      200 {object} model.Task
// @Failure      400 {object} map[string]string "validation error, invalid id, or not found/not owned"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	var req dto.UpdateTaskRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	task, err := h.TaskService.UpdateTask(uint(taskID), userID, req.Title, req.Description, req.Status, req.Deadline, req.AssigneeID, req.EpicID, req.MilestoneID, req.SprintID)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(task)
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task owned by the authenticated user
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Task ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "invalid id or not found/not owned"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	if err := h.TaskService.DeleteTask(uint(taskID), userID); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "task deleted successfully"})
}

// GetProjectTasks godoc
// @Summary      List tasks in a project
// @Description  Get a paginated list of tasks belonging to a project; the requester must be a member of the project
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id     path  int    true  "Project ID"
// @Param        status query string false "Filter by status" Enums(pending, in_progress, done)
// @Param        sort   query string false "Sort by deadline" Enums(deadline_asc, deadline_desc)
// @Param        page   query int    false "Page number" default(1)
// @Param        limit  query int    false "Items per page (max 100)" default(10)
// @Success      200 {object} dto.TaskListResponse
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id}/tasks [get]
func (h *TaskHandler) GetProjectTasks(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	status := c.Query("status")
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tasks, total, err := h.TaskService.GetProjectTasks(userID, uint(projectID), status, sort, page, limit)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(dto.TaskListResponse{
		Data:  tasks,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}