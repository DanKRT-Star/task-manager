package logger

func LabelCreated(labelID, projectID, userID uint) {
	Info("label_created", "label created successfully",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelCreateMissingName(userID, projectID uint) {
	Warn("label_create_missing_name", "label creation failed: name is required",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}

func LabelCreateFailed(projectID, userID uint, name string, err error) {
	Error("label_create_failed", "failed to create label",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "name", Value: name},
		Field{Key: "error", Value: err},
	)
}

func LabelFetchFailed(projectID uint, err error) {
	Error("label_fetch_failed", "failed to fetch labels for project",
		Field{Key: "project_id", Value: projectID},
		Field{Key: "error", Value: err},
	)
}

func LabelNotFound(labelID, userID uint) {
	Warn("label_not_found", "label not found",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelOwnerRequired(userID, projectID, labelID uint) {
	Warn("label_owner_required", "only the project owner can delete a label",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
		Field{Key: "label_id", Value: labelID},
	)
}

func LabelDeleteFailed(labelID, userID uint, err error) {
	Error("label_delete_failed", "failed to delete label",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func LabelDeleted(labelID, userID uint) {
	Info("label_deleted", "label deleted successfully",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelAttachFailed(labelID, taskID, userID uint, err error) {
	Error("label_attach_failed", "failed to attach label",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func LabelAttached(taskID, labelID, userID uint) {
	Info("label_attached", "label attached to task successfully",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "label_id", Value: labelID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelTaskNotFound(taskID, userID uint) {
	Warn("label_task_not_found", "task not found",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelTaskNotInProject(taskID, userID uint) {
	Warn("label_task_not_in_project", "labels can only be attached to tasks within a project",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelTaskNoLabels(taskID, userID uint) {
	Warn("label_task_no_labels", "task has no labels",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
	)
}

func LabelProjectMismatch(taskID, labelID, projectID uint) {
	Warn("label_project_mismatch", "label does not belong to the same project as the task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "label_id", Value: labelID},
		Field{Key: "project_id", Value: projectID},
	)
}

func LabelDetachFailed(labelID, taskID, userID uint, err error) {
	Error("label_detach_failed", "failed to detach label",
		Field{Key: "label_id", Value: labelID},
		Field{Key: "task_id", Value: taskID},
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func LabelDetached(taskID, labelID, userID uint) {
	Info("label_detached", "label detached from task successfully",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "label_id", Value: labelID},
		Field{Key: "user_id", Value: userID},
	)
}

func TaskLabelsFetchFailed(taskID uint, err error) {
	Error("task_labels_fetch_failed", "failed to fetch labels for task",
		Field{Key: "task_id", Value: taskID},
		Field{Key: "error", Value: err},
	)
}

func LabelProjectAccessDenied(userID, projectID uint) {
	Warn("label_project_access_denied", "user is not a member of this project",
		Field{Key: "user_id", Value: userID},
		Field{Key: "project_id", Value: projectID},
	)
}