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

func (h *TaskHandler) CreateTask(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.CreateTaskRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	task, err := h.TaskService.CreateTask(userID, req.Title, req.Description, req.Status, req.Deadline)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(task)
}

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

	task, err := h.TaskService.UpdateTask(uint(taskID), userID, req.Title, req.Description, req.Status, req.Deadline)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(task)
}

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