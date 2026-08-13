package logger

func TaskMissingTitle(userID uint) {
	Warn("task_create_missing_title", "task creation failed: title is required",
		Field{Key: "user_id", Value: userID},
	)
}

func TaskCreateInvalidStatus(userID uint, status string) {
	Warn("task_create_invalid_status", "invalid status value",
		Field{Key: "user_id", Value: userID},
		Field{Key: "status", Value: status},
	)
}

func TaskCreateAccessDenied(userID, projectID uint) {
	Warn("task_create_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func TaskCreateAssigneeAccessDenied(userID, projectID, assigneeID uint) {
	Warn("task_create_assignee_access_denied", "assignee is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "assignee_id", Value: assigneeID},
	)
}

func TaskCreateAssigneeWithoutProject(userID, assigneeID uint) {
	Warn("task_create_assignee_without_project", "cannot assign a task that does not belong to a project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "assignee_id", Value: assigneeID},
	)
}

func TaskCreateInvalidDeadline(userID uint, deadline string) {
	Warn("task_create_invalid_deadline", "invalid deadline format",
		Field{Key: "user_id", Value: userID},
		Field{Key: "deadline", Value: deadline},
	)
}

func TaskCreateFailed(userID uint, title string, err error) {
	Error("task_create_failed", "failed to create task",
		Field{Key: "user_id", Value: userID},
		Field{Key: "title", Value: title},
		Field{Key: "error", Value: err},
	)
}

func TaskActivityLogCreateFailed(taskID, userID uint, err error) {
	Error("task_activity_log_create_failed", "failed to create task activity log",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func TaskCreated(taskID, userID uint, title string) {
	Info("task_created", "task created successfully",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "title", Value: title},
	)
}

func TaskListFetchFailed(userID uint, status string, err error) {
	Error("task_list_fetch_failed", "failed to fetch tasks",
		Field{Key: "user_id", Value: userID},
		Field{Key: "status", Value: status},
		Field{Key: "error", Value: err},
	)
}

func ProjectTaskListAccessDenied(userID, projectID uint) {
	Warn("project_task_list_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func ProjectTaskListFetchFailed(projectID uint, err error) {
	Error("project_task_list_fetch_failed", "failed to fetch project tasks",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "error", Value: err},
	)
}

func TaskUpdateNotFound(taskID, userID uint) {
	Warn("task_update_not_found", "task not found",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskUpdateAccessDenied(taskID, userID uint) {
	Warn("task_update_access_denied", "user does not have permission to modify this task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskUpdateInvalidStatus(taskID uint, status string) {
	Warn("task_update_invalid_status", "invalid status value",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "status", Value: status},
	)
}

func TaskUpdateInvalidDeadline(taskID uint, deadline string) {
	Warn("task_update_invalid_deadline", "invalid deadline format",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "deadline", Value: deadline},
	)
}

func TaskUpdateAssigneeWithoutProject(taskID, assigneeID uint) {
	Warn("task_update_assignee_without_project", "cannot assign a task that does not belong to a project",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "assignee_id", Value: assigneeID},
	)
}

func TaskUpdateAssigneeNotProjectMember(taskID, projectID, assigneeID uint) {
	Warn("task_update_assignee_not_project_member", "assignee is not a member of this project",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "assignee_id", Value: assigneeID},
	)
}

func TaskUpdateFailed(taskID, userID uint, err error) {
	Error("task_update_failed", "failed to update task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func TaskStatusActivityLogFailed(taskID, userID uint, err error) {
	Error("task_status_activity_log_failed", "failed to create status change activity log",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func TaskUpdated(taskID, userID uint) {
	Info("task_updated", "task updated successfully",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskDeleteNotFound(taskID, userID uint) {
	Warn("task_delete_not_found", "task not found",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskDeleteAccessDenied(taskID, userID uint) {
	Warn("task_delete_access_denied", "user does not have permission to delete this task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskDeleteFailed(taskID, userID uint, err error) {
	Error("task_delete_failed", "failed to delete task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func TaskDeleted(taskID, userID uint) {
	Info("task_deleted", "task deleted successfully",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}