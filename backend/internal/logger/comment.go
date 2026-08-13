package logger

func CommentCreated(commentID, taskID, userID uint) {
	Info("comment_created", "comment created successfully",
		Field{Key: "comment_id", Value: commentID},
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func CommentCreateFailed(taskID, userID uint, err error) {
	Error("comment_create_failed", "failed to create comment",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func CommentTaskAccessDenied(userID, taskID uint) {
	Warn("comment_task_access_denied", "user does not have access to this task",
		Field{Key: "user_id", Value: userID},
		Field{Key: "task_id", Value: taskID},
	)
}

func CommentNotFound(commentID, userID uint) {
	Warn("comment_not_found", "comment not found",
		Field{Key: "comment_id", Value: commentID},
		Field{Key: "user_id", Value: userID},
	)
}

func CommentDeleteForbidden(commentID, userID, commentOwnerID uint) {
	Warn("comment_delete_forbidden", "comment delete forbidden",
		Field{Key: "comment_id", Value: commentID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "comment_owner_id", Value: commentOwnerID},
	)
}

func CommentDeleteFailed(commentID, userID uint, err error) {
	Error("comment_delete_failed", "failed to delete comment",
		Field{Key: "comment_id", Value: commentID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func CommentTaskNotFound(taskID, userID uint) {
	Warn("comment_task_not_found", "task not found for comment",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func CommentEmptyContent(userID, taskID uint) {
	Warn("comment_empty_content", "comment content is empty",
		Field{Key: "user_id", Value: userID},
		Field{Key: "task_id", Value: taskID},
	)
}

func CommentFetchFailed(taskID uint, err error) {
	Error("comment_fetch_failed", "failed to fetch comments for task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "error", Value: err},
	)
}

func CommentDeleted(commentID, userID uint) {
	Info("comment_deleted", "comment deleted successfully",
		Field{Key: "comment_id", Value: commentID},
		Field{Key: "user_id", Value: userID},
	)
}