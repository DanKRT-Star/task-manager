package logger

func MilestoneCreated(milestoneID, projectID, userID uint) {
	Info("milestone_created", "milestone created successfully",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func MilestoneCreateMissingTitle(userID, projectID uint) {
	Warn("milestone_create_missing_title", "milestone creation failed: title is required",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func MilestoneCreateAccessDenied(userID, projectID uint) {
	Warn("milestone_create_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func MilestoneCreateFailed(projectID, userID uint, title string, err error) {
	Error("milestone_create_failed", "failed to create milestone",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "title", Value: title},
		Field{Key: "error", Value: err},
	)
}

func MilestoneFetchAccessDenied(userID, projectID uint) {
	Warn("milestone_fetch_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func MilestoneFetchFailed(projectID uint, err error) {
	Error("milestone_fetch_failed", "failed to fetch milestones for project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "error", Value: err},
	)
}

func MilestoneNotFound(milestoneID, userID uint) {
	Warn("milestone_not_found", "milestone not found",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
	)
}

func MilestoneUpdateAccessDenied(milestoneID, userID, projectID uint) {
	Warn("milestone_update_access_denied", "user is not a member of this project",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func MilestoneUpdateFailed(milestoneID, userID uint, err error) {
	Error("milestone_update_failed", "failed to update milestone",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func MilestoneUpdated(milestoneID, userID uint) {
	Info("milestone_updated", "milestone updated successfully",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
	)
}

func MilestoneDeleteAccessDenied(milestoneID, userID, projectID uint) {
	Warn("milestone_delete_access_denied", "user is not a member of this project",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func MilestoneOwnerRequired(milestoneID, projectID, userID uint) {
	Warn("milestone_owner_required", "only the project owner can delete a milestone",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func MilestoneDeleteFailed(milestoneID, userID uint, err error) {
	Error("milestone_delete_failed", "failed to delete milestone",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func MilestoneDeleted(milestoneID, userID uint) {
	Info("milestone_deleted", "milestone deleted successfully",
		Field{Key: "milestone_id", Value: milestoneID},
		Field{Key: "user_id", Value: userID},
	)
}