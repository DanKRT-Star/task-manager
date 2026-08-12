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
