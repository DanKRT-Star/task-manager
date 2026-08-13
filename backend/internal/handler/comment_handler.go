package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type CommentHandler struct {
	CommentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}


// CreateComment godoc
// @Summary      Create a comment
// @Description  Add a comment to a task; requester must have access to the task
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskId  path int true "Task ID"
// @Param        request body dto.CreateCommentRequest true "Comment payload"
// @Success      201 {object} model.Comment
// @Failure      400 {object} map[string]string "validation error or access denied"
// @Router       /tasks/{taskId}/comments [post]
func (h *CommentHandler) CreateComment(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	var req dto.CreateCommentRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	comment, err := h.CommentService.CreateComment(userID, uint(taskID), req.Content)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(comment)
}

// GetTaskComments godoc
// @Summary      List task comments
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        taskId path int true "Task ID"
// @Success      200 {object} map[string]interface{} "data: array of comments"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Router       /tasks/{taskId}/comments [get]
func (h *CommentHandler) GetTaskComments(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	comments, err := h.CommentService.GetTaskComments(userID, uint(taskID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": comments})
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Only the comment author can delete it
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        commentId path int true "Comment ID"
// @Success      200 {object} map[string]string "message"
// @Failure      400 {object} map[string]string "not found or not the author"
// @Router       /comments/{commentId} [delete]
func (h *CommentHandler) DeleteComment(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	commentID, err := strconv.Atoi(c.Params("commentId"))
	if err != nil {
		return apperror.BadRequest("invalid comment id")
	}

	if err := h.CommentService.DeleteComment(userID, uint(commentID)); err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"message": "comment deleted successfully"})
}
