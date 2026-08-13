package logger

func ActivityFetchTaskNotFound(taskID, userID uint) {
	Warn("activity_fetch_task_not_found", "task not found while fetching activity",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func ActivityFetchAccessDenied(taskID, userID, projectID uint) {
	Warn("activity_fetch_access_denied", "user does not have access to this task activity",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func ActivityFetchFailed(taskID uint, err error) {
	Error("activity_fetch_failed", "failed to fetch activity logs",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "error", Value: err},
	)
}
