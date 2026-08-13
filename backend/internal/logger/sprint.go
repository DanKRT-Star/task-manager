package logger

func SprintCreated(sprintID, projectID, userID uint) {
	Info("sprint_created", "sprint created successfully",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func SprintProjectAccessDenied(userID, projectID uint) {
	Warn("sprint_project_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func SprintDeleteFailed(sprintID, userID uint, err error) {
	Error("sprint_delete_failed", "failed to delete sprint",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func SprintUpdated(sprintID, userID uint) {
	Info("sprint_updated", "sprint updated successfully",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "user_id", Value: userID},
	)
}

func SprintUpdateFailed(sprintID, userID uint, err error) {
	Error("sprint_update_failed", "failed to update sprint",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func SprintCreateMissingName(userID, projectID uint) {
	Warn("sprint_create_missing_name", "sprint name is required",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func SprintCreateFailed(projectID, userID uint, name string, err error) {
	Error("sprint_create_failed", "failed to create sprint",
		Field{Key: "name", Value: name},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func SprintFetchFailed(projectID uint, err error) {
	Error("sprint_fetch_failed", "failed to fetch sprints for project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "error", Value: err},
	)
}

func SprintNotFound(sprintID, userID uint) {
	Warn("sprint_not_found", "sprint can not found or access denied",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "user_id", Value: userID},
	)
}

func SprintInvalidStatus(sprintID uint, status string) {
	Warn("sprint_invalid_status", "invalid status value",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "status", Value: status},
	)
}

func SprintOwnerRequired(sprintID, projectID, userID uint) {
	Warn("sprint_owner_required", "only the project owner can delete a sprint",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func SprintDeleted(sprintID, userID uint) {
	Info("sprint_deleted", "sprint deleted successfully",
		Field{Key: "sprint_id", Value: sprintID},
		Field{Key: "user_id", Value: userID},
	)
}