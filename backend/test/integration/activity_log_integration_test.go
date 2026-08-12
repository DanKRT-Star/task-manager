package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ActivityLog_TracksTaskChanges(t *testing.T) {
	cleanTables()

	ownerToken := registerAndGetToken(t, "activityowner@example.com", "12345678")
	outsiderToken := registerAndGetToken(t, "activityoutsider@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Activity Project")
	projectIDStr := strconv.Itoa(int(projectID))

	// Tạo task trong project
	createTaskBody, _ := json.Marshal(map[string]interface{}{
		"title":     "Tracked task",
		"deadline":  "2026-08-15T17:00:00Z",
		"projectId": projectID,
	})
	createTaskReq := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(createTaskBody))
	createTaskReq.Header.Set("Content-Type", "application/json")
	createTaskReq.Header.Set("Authorization", "Bearer "+ownerToken)
	createTaskResp, err := app.Test(createTaskReq)
	require.NoError(t, err)
	require.Equal(t, 201, createTaskResp.StatusCode)

	var createdTask struct {
		TaskID uint `json:"taskId"`
	}
	require.NoError(t, json.NewDecoder(createTaskResp.Body).Decode(&createdTask))
	taskIDStr := strconv.Itoa(int(createdTask.TaskID))

	// Đổi status để tạo thêm 1 activity log
	updateBody, _ := json.Marshal(map[string]string{"status": "in_progress"})
	updateReq := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskIDStr, bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+ownerToken)
	updateResp, err := app.Test(updateReq)
	require.NoError(t, err)
	require.Equal(t, 200, updateResp.StatusCode)

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
		assertBody     func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "owner can view activity history with both entries",
			authorization:  "Bearer " + ownerToken,
			expectedStatus: 200,
			assertBody: func(t *testing.T, resp *http.Response) {
				var body struct {
					Data []map[string]interface{} `json:"data"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				assert.Len(t, body.Data, 2)
			},
		},
		{
			name:           "outsider cannot view activity history",
			authorization:  "Bearer " + outsiderToken,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskIDStr+"/activity", nil)
			req.Header.Set("Authorization", tt.authorization)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			if tt.assertBody != nil {
				tt.assertBody(t, resp)
			}
		})
	}

	_ = projectIDStr
}