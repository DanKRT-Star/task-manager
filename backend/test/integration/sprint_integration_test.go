package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Sprint_CreateAndList(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "sprintowner@example.com", "12345678")
	projectID := createTestProject(t, token, "Sprint Project")
	projectIDStr := strconv.Itoa(int(projectID))

	tests := []struct {
		name         string
		method       string
		path         string
		body         map[string]string
		authToken    string
		expectedCode int
		assertFn     func(t *testing.T, resp *http.Response)
	}{
		{
			name:         "create sprint success",
			method:       "POST",
			path:         "/api/v1/projects/" + projectIDStr + "/sprints",
			body:         map[string]string{"name": "Sprint 1"},
			authToken:    token,
			expectedCode: 201,
			assertFn: func(t *testing.T, resp *http.Response) {
				var created struct {
					Status string `json:"status"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
				assert.Equal(t, "planned", created.Status)
			},
		},
		{
			name:         "create sprint missing name",
			method:       "POST",
			path:         "/api/v1/projects/" + projectIDStr + "/sprints",
			body:         map[string]string{"name": ""},
			authToken:    token,
			expectedCode: 400,
		},
		{
			name:         "list sprints returns created sprint",
			method:       "GET",
			path:         "/api/v1/projects/" + projectIDStr + "/sprints",
			authToken:    token,
			expectedCode: 200,
			assertFn: func(t *testing.T, resp *http.Response) {
				var body struct {
					Data []map[string]interface{} `json:"data"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				assert.Len(t, body.Data, 1)
			},
		},
		{
			name:         "non-member cannot list sprints",
			method:       "GET",
			path:         "/api/v1/projects/" + projectIDStr + "/sprints",
			authToken:    registerAndGetToken(t, "sprintoutsider@example.com", "12345678"),
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody io.Reader
			if tt.body != nil {
				payload, err := json.Marshal(tt.body)
				require.NoError(t, err)
				reqBody = bytes.NewReader(payload)
			}

			req := httptest.NewRequest(tt.method, tt.path, reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			req.Header.Set("Authorization", "Bearer "+tt.authToken)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.assertFn != nil {
				tt.assertFn(t, resp)
			}
		})
	}
}

func TestIntegration_Sprint_UpdateStatusAndDelete_PermissionRules(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "sprintrealowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "sprintregularmember@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Sprint Permission Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addBody, _ := json.Marshal(map[string]string{"email": "sprintregularmember@example.com"})
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	_, err := app.Test(reqAdd)
	require.NoError(t, err)

	sprintBody, _ := json.Marshal(map[string]string{"name": "Original Sprint"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/sprints", bytes.NewReader(sprintBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)

	var created struct {
		SprintID uint `json:"sprintId"`
	}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&created))
	sprintIDStr := strconv.Itoa(int(created.SprintID))

	tests := []struct {
		name         string
		method       string
		path         string
		body         map[string]string
		authToken    string
		expectedCode int
		assertFn     func(t *testing.T, resp *http.Response)
	}{
		{
			name:         "member can start the sprint",
			method:       "PUT",
			path:         "/api/v1/sprints/" + sprintIDStr,
			body:         map[string]string{"status": "active"},
			authToken:    memberToken,
			expectedCode: 200,
			assertFn: func(t *testing.T, resp *http.Response) {
				var body struct {
					Status string `json:"status"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				assert.Equal(t, "active", body.Status)
			},
		},
		{
			name:         "invalid status value rejected",
			method:       "PUT",
			path:         "/api/v1/sprints/" + sprintIDStr,
			body:         map[string]string{"status": "not_a_real_status"},
			authToken:    ownerToken,
			expectedCode: 400,
		},
		{
			name:         "member cannot delete sprint",
			method:       "DELETE",
			path:         "/api/v1/sprints/" + sprintIDStr,
			authToken:    memberToken,
			expectedCode: 400,
		},
		{
			name:         "owner can delete sprint",
			method:       "DELETE",
			path:         "/api/v1/sprints/" + sprintIDStr,
			authToken:    ownerToken,
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody io.Reader
			if tt.body != nil {
				payload, err := json.Marshal(tt.body)
				require.NoError(t, err)
				reqBody = bytes.NewReader(payload)
			}

			req := httptest.NewRequest(tt.method, tt.path, reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			req.Header.Set("Authorization", "Bearer "+tt.authToken)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.assertFn != nil {
				tt.assertFn(t, resp)
			}
		})
	}
}
