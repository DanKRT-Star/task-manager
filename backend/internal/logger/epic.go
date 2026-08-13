package logger

func EpicCreated(epicID, projectID, userID uint) {
	Info("epic_created", "epic created successfully",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func EpicCreateMissingTitle(userID, projectID uint) {
	Warn("epic_create_missing_title", "epic creation failed: title is required",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func EpicCreateAccessDenied(userID, projectID uint) {
	Warn("epic_create_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func EpicCreateFailed(projectID, userID uint, title string, err error) {
	Error("epic_create_failed", "failed to create epic",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "title", Value: title},
		Field{Key: "error", Value: err},
	)
}

func EpicFetchAccessDenied(userID, projectID uint) {
	Warn("epic_fetch_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func EpicFetchFailed(projectID uint, err error) {
	Error("epic_fetch_failed", "failed to fetch epics for project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "error", Value: err},
	)
}

func EpicNotFound(epicID, userID uint) {
	Warn("epic_not_found", "epic not found",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
	)
}

func EpicUpdateAccessDenied(epicID, userID, projectID uint) {
	Warn("epic_update_access_denied", "user is not a member of this project",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func EpicUpdateFailed(epicID, userID uint, err error) {
	Error("epic_update_failed", "failed to update epic",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func EpicUpdated(epicID, userID uint) {
	Info("epic_updated", "epic updated successfully",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
	)
}

func EpicDeleteAccessDenied(epicID, userID, projectID uint) {
	Warn("epic_delete_access_denied", "user is not a member of this project",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func EpicOwnerRequired(epicID, projectID, userID uint) {
	Warn("epic_owner_required", "only the project owner can delete an epic",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func EpicDeleteFailed(epicID, userID uint, err error) {
	Error("epic_delete_failed", "failed to delete epic",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func EpicDeleted(epicID, userID uint) {
	Info("epic_deleted", "epic deleted successfully",
		Field{Key: "epic_id", Value: epicID},
		Field{Key: "user_id", Value: userID},
	)
}