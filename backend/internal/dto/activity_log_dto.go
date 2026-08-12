package dto

type ActivityLogResponse struct {
	ActivityLogID uint   `json:"activityLogId"`
	TaskID        uint   `json:"taskId"`
	UserID        uint   `json:"userId"`
	Action        string `json:"action"`
	Detail        string `json:"detail"`
	CreatedAt     string `json:"createdAt"`
}