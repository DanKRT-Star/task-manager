package logger

func ProjectCreated(projectID, userID uint) {
	Info("project_created", "project created successfully",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectCreateFailed(userID uint, name string, err error) {
	Error("project_create_failed", "failed to create project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_name", Value: name},
		Field{Key: "error", Value: err},
	)
}

func ProjectCreateMissingName(userID uint) {
	Warn("project_create_missing_name", "project creation failed: name is required",
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectAccessDenied(userID, projectID uint) {
	Warn("project_access_denied", "project not found or access denied",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func ProjectUpdated(projectID, userID uint) {
	Info("project_updated", "project updated successfully",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectUpdateFailed(projectID, userID uint, err error) {
	Error("project_update_failed", "failed to update project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func ProjectDeleted(projectID, userID uint) {
	Info("project_deleted", "project deleted successfully",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectDeleteFailed(projectID, userID uint, err error) {
	Error("project_delete_failed", "failed to delete project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func ProjectOwnerRequired(projectID, userID uint) {
	Warn("project_owner_required", "Only project owner can do this action",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectMemberNotFound(projectID uint, email string) {
	Warn("project_member_not_found", "user with this email does not exist",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "email", Value: email},
	)
}

func ProjectMemberAlreadyExists(projectID, userID uint) {
	Warn("project_member_already_exists", "user is already a member of this project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectAddMemberFailed(projectID, userID uint, err error) {
	Error("project_add_member_failed", "failed to add project member",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func ProjectMemberAdded(projectID, userID uint) {
	Info("project_member_added", "project member added successfully",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectOwnerCanNotBeRemoved(projectID, userID uint) {
	Warn("project_owner_can_not_be_removed", "cannot remove the project owner",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectRemoveMemberFailed(projectID, userID uint, err error) {
	Error("project_remove_member_failed", "failed to remove member",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func ProjectMemberRemoved(projectID, userID uint) {
	Info("project_remove_member_successfully", "project member removed successfully",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func ProjectFetchMembersFailed(projectID, userID uint, err error) {
	Error("project_fetch_members_failed", "failed to fetch project members",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}