package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type MilestoneHandler struct {
	MilestoneService *service.MilestoneService
}

func NewMilestoneHandler(milestoneService *service.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{MilestoneService: milestoneService}
}

// CreateMilestone godoc
// @Summary      Create a milestone
// @Description  Create a new milestone within a project; requester must be a project member
// @Tags         milestones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Project ID"
// @Param        request body dto.CreateMilestoneRequest true "Milestone payload"
// @Success      201 {object} model.Milestone
// @Failure      400 {object} map[string]string "validation error or access denied"
// @Router       /projects/{id}/milestones [post]
func (h *MilestoneHandler) CreateMilestone(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.CreateMilestoneRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	milestone, err := h.MilestoneService.CreateMilestone(userID, uint(projectID), req.Title, req.Description, req.DueDate)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(milestone)
}

// GetProjectMilestones godoc
// @Summary      List milestones in a project
// @Tags         milestones
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} map[string]interface{} "data: array of milestones"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Router       /projects/{id}/milestones [get]
func (h *MilestoneHandler) GetProjectMilestones(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	milestones, err := h.MilestoneService.GetProjectMilestones(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": milestones})
}

// UpdateMilestone godoc
// @Summary      Update a milestone
// @Tags         milestones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        milestoneId path int true "Milestone ID"
// @Param        request     body dto.UpdateMilestoneRequest true "Fields to update"
// @Success      200 {object} model.Milestone
// @Failure      400 {object} map[string]string "validation error, not found, or access denied"
// @Router       /milestones/{milestoneId} [put]
func (h *MilestoneHandler) UpdateMilestone(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	milestoneID, err := strconv.Atoi(c.Params("milestoneId"))
	if err != nil {
		return apperror.BadRequest("invalid milestone id")
	}

	var req dto.UpdateMilestoneRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	milestone, err := h.MilestoneService.UpdateMilestone(userID, uint(milestoneID), req.Title, req.Description, req.DueDate)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(milestone)
}

// DeleteMilestone godoc
// @Summary      Delete a milestone
// @Description  Only the project owner can delete a milestone
// @Tags         milestones
// @Produce      json
// @Security     BearerAuth
// @Param        milestoneId path int true "Milestone ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "not found, or not the owner"
// @Router       /milestones/{milestoneId} [delete]
func (h *MilestoneHandler) DeleteMilestone(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	milestoneID, err := strconv.Atoi(c.Params("milestoneId"))
	if err != nil {
		return apperror.BadRequest("invalid milestone id")
	}

	if err := h.MilestoneService.DeleteMilestone(userID, uint(milestoneID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "milestone deleted successfully"})
}