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

func TestIntegration_Comment_CreateAndList(t *testing.T) {
	cleanTables()

	ownerToken := registerAndGetToken(t, "commentowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "commentmember@example.com", "12345678")
	outsiderToken := registerAndGetToken(t, "commentoutsider@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Comment Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addMemberBody, _ := json.Marshal(map[string]string{"email": "commentmember@example.com"})
	addMemberReq := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addMemberBody))
	addMemberReq.Header.Set("Content-Type", "application/json")
	addMemberReq.Header.Set("Authorization", "Bearer "+ownerToken)
	addMemberResp, err := app.Test(addMemberReq)
	require.NoError(t, err)
	require.Equal(t, 201, addMemberResp.StatusCode)

	createTaskBody, _ := json.Marshal(map[string]interface{}{
		"title":     "Task with comments",
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

	tests := []struct {
		name           string
		method         string
		path           string
		authorization  string
		body           interface{}
		expectedStatus int
		assertBody     func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "member can create comment",
			method:         http.MethodPost,
			path:           "/api/v1/tasks/" + taskIDStr + "/comments",
			authorization:  "Bearer " + memberToken,
			body:           map[string]string{"content": "I reviewed this task"},
			expectedStatus: 201,
			assertBody: func(t *testing.T, resp *http.Response) {
				var body map[string]interface{}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				assert.Equal(t, "I reviewed this task", body["content"])
			},
		},
		{
			name:           "outsider cannot create comment",
			method:         http.MethodPost,
			path:           "/api/v1/tasks/" + taskIDStr + "/comments",
			authorization:  "Bearer " + outsiderToken,
			body:           map[string]string{"content": "This should fail"},
			expectedStatus: 400,
		},
		{
			name:           "member can list comments",
			method:         http.MethodGet,
			path:           "/api/v1/tasks/" + taskIDStr + "/comments",
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
			var requestBody []byte
			if tt.body != nil {
				requestBody, _ = json.Marshal(tt.body)
			}
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(requestBody))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
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
