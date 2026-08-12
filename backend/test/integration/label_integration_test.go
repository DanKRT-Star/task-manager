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

func TestIntegration_Label_CreateAttachAndList(t *testing.T) {
	cleanTables()

	ownerToken := registerAndGetToken(t, "labelowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "labelmember@example.com", "12345678")
	outsiderToken := registerAndGetToken(t, "labeloutsider@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Label Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addMemberBody, _ := json.Marshal(map[string]string{"email": "labelmember@example.com"})
	addMemberReq := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addMemberBody))
	addMemberReq.Header.Set("Content-Type", "application/json")
	addMemberReq.Header.Set("Authorization", "Bearer "+ownerToken)
	addMemberResp, err := app.Test(addMemberReq)
	require.NoError(t, err)
	require.Equal(t, 201, addMemberResp.StatusCode)

	createTaskBody, _ := json.Marshal(map[string]interface{}{
		"title":     "Task with labels",
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

	labelBody, _ := json.Marshal(map[string]string{"name": "frontend", "color": "#ff0000"})
	createLabelReq := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/labels", bytes.NewReader(labelBody))
	createLabelReq.Header.Set("Content-Type", "application/json")
	createLabelReq.Header.Set("Authorization", "Bearer "+ownerToken)
	createLabelResp, err := app.Test(createLabelReq)
	require.NoError(t, err)
	require.Equal(t, 201, createLabelResp.StatusCode)

	var createdLabel struct {
		LabelID uint `json:"labelId"`
	}
	require.NoError(t, json.NewDecoder(createLabelResp.Body).Decode(&createdLabel))
	labelIDStr := strconv.Itoa(int(createdLabel.LabelID))

	tests := []struct {
		name           string
		method         string
		path           string
		authorization  string
		expectedStatus int
		assertBody     func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "member can attach label to task",
			method:         http.MethodPost,
			path:           "/api/v1/tasks/" + taskIDStr + "/labels/" + labelIDStr,
			authorization:  "Bearer " + memberToken,
			expectedStatus: 200,
		},
		{
			name:           "outsider cannot attach label",
			method:         http.MethodPost,
			path:           "/api/v1/tasks/" + taskIDStr + "/labels/" + labelIDStr,
			authorization:  "Bearer " + outsiderToken,
			expectedStatus: 400,
		},
		{
			name:           "member can list task labels",
			method:         http.MethodGet,
			path:           "/api/v1/tasks/" + taskIDStr + "/labels",
			authorization:  "Bearer " + memberToken,
			expectedStatus: 200,
			assertBody: func(t *testing.T, resp *http.Response) {
				var body struct {
					Data []map[string]interface{} `json:"data"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				assert.Len(t, body.Data, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Authorization", tt.authorization)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			if tt.assertBody != nil {
				tt.assertBody(t, resp)
			}
		})
	}
}
