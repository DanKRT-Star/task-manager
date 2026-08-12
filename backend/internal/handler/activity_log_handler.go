package handler

import (
	"strconv"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/gofiber/fiber/v3"
)

type ActivityLogHandler struct {
	ActivityLogService *service.ActivityLogService
}

func NewActivityLogHandler(activityLogService *service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{ActivityLogService: activityLogService}
}

// GetTaskActivity godoc
// @Summary      Get task activity history
// @Description  Get the change history of a task (created, status changes, etc.); requester must have access to the task
// @Tags         activity
// @Produce      json
// @Security     BearerAuth
// @Param        taskId path int true "Task ID"
// @Success      200 {object} map[string]interface{} "data: array of activity logs"
// @Failure      400 {object} map[string]string "invalid id or access denied"
// @Router       /tasks/{taskId}/activity [get]
func (h *ActivityLogHandler) GetTaskActivity(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return apperror.BadRequest("invalid task id")
	}

	logs, err := h.ActivityLogService.GetTaskActivity(userID, uint(taskID))
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.JSON(fiber.Map{"data": logs})
}