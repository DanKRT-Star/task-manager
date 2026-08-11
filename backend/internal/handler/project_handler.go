package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type ProjectHandler struct {
	ProjectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{ProjectService: projectService}
}

// CreateProject godoc
// @Summary      Create a project
// @Description  Create a new project; the creator automatically becomes the project owner
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateProjectRequest true "Project payload"
// @Success      201 {object} model.Project
// @Failure      400 {object} map[string]string "validation error"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects [post]
func (h *ProjectHandler) CreateProject(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.CreateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	project, err := h.ProjectService.CreateProject(userID, req.Name, req.Description, req.Deadline)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(project)
}

// GetProjects godoc
// @Summary      List my projects
// @Description  Get all projects the authenticated user is a member or owner of
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "data: array of projects"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects [get]
func (h *ProjectHandler) GetProjects(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projects, err := h.ProjectService.GetUserProjects(userID)
	if err != nil {
		return apperror.Internal("failed to fetch projects")
	}

	return c.JSON(fiber.Map{"data": projects})
}

// GetProject godoc
// @Summary      Get a project
// @Description  Get project details; the requester must be a member of the project
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} model.Project
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id} [get]
func (h *ProjectHandler) GetProject(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	project, err := h.ProjectService.GetProject(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(project)
}

// UpdateProject godoc
// @Summary      Update a project
// @Description  Update project details; only the project owner can perform this action
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Project ID"
// @Param        request body dto.UpdateProjectRequest true "Fields to update"
// @Success      200 {object} model.Project
// @Failure      400 {object} map[string]string "validation error, invalid id, or not the owner"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.UpdateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	project, err := h.ProjectService.UpdateProject(userID, uint(projectID), req.Name, req.Description, req.Deadline)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(project)
}

// DeleteProject godoc
// @Summary      Delete a project
// @Description  Delete a project; only the project owner can perform this action. Tasks in the project are unlinked, not deleted.
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "invalid id or not the owner"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	if err := h.ProjectService.DeleteProject(userID, uint(projectID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "project deleted successfully"})
}

// AddMember godoc
// @Summary      Add a project member
// @Description  Add a user to the project by email; only the project owner can perform this action
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int true "Project ID"
// @Param        request body dto.AddMemberRequest true "Member email"
// @Success      201 {object} model.ProjectMember
// @Failure      400 {object} map[string]string "validation error, not the owner, user not found, or already a member"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id}/members [post]
func (h *ProjectHandler) AddMember(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	var req dto.AddMemberRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	member, err := h.ProjectService.AddMember(userID, uint(projectID), req.Email)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(member)
}

// RemoveMember godoc
// @Summary      Remove a project member
// @Description  Remove a member from the project; only the project owner can perform this action. The owner cannot be removed.
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id     path int true "Project ID"
// @Param        userId path int true "User ID to remove"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "invalid id, not the owner, or attempting to remove the owner"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id}/members/{userId} [delete]
func (h *ProjectHandler) RemoveMember(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}
	targetUserID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	if err := h.ProjectService.RemoveMember(userID, uint(projectID), uint(targetUserID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "member removed successfully"})
}

// GetMembers godoc
// @Summary      List project members
// @Description  Get all members of a project; the requester must be a member of the project
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} map[string]interface{} "data: array of members"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /projects/{id}/members [get]
func (h *ProjectHandler) GetMembers(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	projectID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid project id")
	}

	members, err := h.ProjectService.GetMembers(userID, uint(projectID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": members})
}