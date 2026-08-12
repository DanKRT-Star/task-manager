package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type LabelHandler struct {
	LabelService *service.LabelService
}

func NewLabelHandler(labelService *service.LabelService) *LabelHandler {
	return &LabelHandler{LabelService: labelService}
}

func (h *LabelHandler) CreateLabel(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.CreateLabelRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	label, err := h.LabelService.CreateLabel(userID, uint(projectID), req.Name, req.Color)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(label)
}

func (h *LabelHandler) GetProjectLabels(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	labels, err := h.LabelService.GetProjectLabels(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": labels})
}

func (h *LabelHandler) DeleteLabel(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	labelID, err := strconv.Atoi(c.Params("labelId"))
	if err != nil {
		return apperror.BadRequest("invalid label id")
	}

	if err := h.LabelService.DeleteLabel(userID, uint(labelID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "label deleted successfully"})
}

func (h *LabelHandler) AttachLabel(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}
	labelID, err := strconv.Atoi(c.Params("labelId"))
	if err != nil {
		return apperror.BadRequest("invalid label id")
	}

	if err := h.LabelService.AttachLabel(userID, uint(taskID), uint(labelID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "label attached successfully"})
}

func (h *LabelHandler) DetachLabel(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}
	labelID, err := strconv.Atoi(c.Params("labelId"))
	if err != nil {
		return apperror.BadRequest("invalid label id")
	}

	if err := h.LabelService.DetachLabel(userID, uint(taskID), uint(labelID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "label detached successfully"})
}

func (h *LabelHandler) GetTaskLabels(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	labels, err := h.LabelService.GetTaskLabels(userID, uint(taskID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": labels})
}
