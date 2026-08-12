package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type EpicHandler struct {
	EpicService *service.EpicService
}

func NewEpicHandler(epicService *service.EpicService) *EpicHandler {
	return &EpicHandler{EpicService: epicService}
}

// CreateEpic godoc
// @Summary      Create an epic
// @Description  Create a new epic within a project; requester must be a project member
// @Tags         epics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Project ID"
// @Param        request body dto.CreateEpicRequest true "Epic payload"
// @Success      201 {object} model.Epic
// @Failure      400 {object} map[string]string "validation error or access denied"
// @Router       /projects/{id}/epics [post]
func (h *EpicHandler) CreateEpic(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.CreateEpicRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	epic, err := h.EpicService.CreateEpic(userID, uint(projectID), req.Title, req.Description)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(epic)
}

// GetProjectEpics godoc
// @Summary      List epics in a project
// @Tags         epics
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} map[string]interface{} "data: array of epics"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Router       /projects/{id}/epics [get]
func (h *EpicHandler) GetProjectEpics(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	epics, err := h.EpicService.GetProjectEpics(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": epics})
}

// UpdateEpic godoc
// @Summary      Update an epic
// @Tags         epics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        epicId  path int true "Epic ID"
// @Param        request body dto.UpdateEpicRequest true "Fields to update"
// @Success      200 {object} model.Epic
// @Failure      400 {object} map[string]string "validation error, not found, or access denied"
// @Router       /epics/{epicId} [put]
func (h *EpicHandler) UpdateEpic(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	epicID, err := strconv.Atoi(c.Params("epicId"))
	if err != nil {
		return apperror.BadRequest("invalid epic id")
	}

	var req dto.UpdateEpicRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	epic, err := h.EpicService.UpdateEpic(userID, uint(epicID), req.Title, req.Description)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(epic)
}

// DeleteEpic godoc
// @Summary      Delete an epic
// @Description  Only the project owner can delete an epic
// @Tags         epics
// @Produce      json
// @Security     BearerAuth
// @Param        epicId path int true "Epic ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "not found, or not the owner"
// @Router       /epics/{epicId} [delete]
func (h *EpicHandler) DeleteEpic(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	epicID, err := strconv.Atoi(c.Params("epicId"))
	if err != nil {
		return apperror.BadRequest("invalid epic id")
	}

	if err := h.EpicService.DeleteEpic(userID, uint(epicID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "epic deleted successfully"})
}