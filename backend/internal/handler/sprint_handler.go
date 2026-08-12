package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type SprintHandler struct {
	SprintService *service.SprintService
}

func NewSprintHandler(sprintService *service.SprintService) *SprintHandler {
	return &SprintHandler{SprintService: sprintService}
}

// CreateSprint godoc
// @Summary      Create a sprint
// @Description  Create a new sprint within a project; requester must be a project member
// @Tags         sprints
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Project ID"
// @Param        request body dto.CreateSprintRequest true "Sprint payload"
// @Success      201 {object} model.Sprint
// @Failure      400 {object} map[string]string "validation error or access denied"
// @Router       /projects/{id}/sprints [post]
func (h *SprintHandler) CreateSprint(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.CreateSprintRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	sprint, err := h.SprintService.CreateSprint(userID, uint(projectID), req.Name, req.StartDate, req.EndDate)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(sprint)
}

// GetProjectSprints godoc
// @Summary      List sprints in a project
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} map[string]interface{} "data: array of sprints"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Router       /projects/{id}/sprints [get]
func (h *SprintHandler) GetProjectSprints(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	sprints, err := h.SprintService.GetProjectSprints(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": sprints})
}

// UpdateSprint godoc
// @Summary      Update a sprint
// @Description  Update sprint details including status transitions (planned/active/completed)
// @Tags         sprints
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        sprintId path int true "Sprint ID"
// @Param        request  body dto.UpdateSprintRequest true "Fields to update"
// @Success      200 {object} model.Sprint
// @Failure      400 {object} map[string]string "validation error, not found, or access denied"
// @Router       /sprints/{sprintId} [put]
func (h *SprintHandler) UpdateSprint(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	sprintID, err := strconv.Atoi(c.Params("sprintId"))
	if err != nil {
		return apperror.BadRequest("invalid sprint id")
	}

	var req dto.UpdateSprintRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	sprint, err := h.SprintService.UpdateSprint(userID, uint(sprintID), req.Name, req.Status, req.StartDate, req.EndDate)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(sprint)
}

// DeleteSprint godoc
// @Summary      Delete a sprint
// @Description  Only the project owner can delete a sprint
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        sprintId path int true "Sprint ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "not found, or not the owner"
// @Router       /sprints/{sprintId} [delete]
func (h *SprintHandler) DeleteSprint(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	sprintID, err := strconv.Atoi(c.Params("sprintId"))
	if err != nil {
		return apperror.BadRequest("invalid sprint id")
	}

	if err := h.SprintService.DeleteSprint(userID, uint(sprintID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "sprint deleted successfully"})
}