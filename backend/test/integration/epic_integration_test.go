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

func createTestProject(t *testing.T, token, name string) uint {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"name": name})
	req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var created struct {
		ProjectID uint `json:"projectId"`
	}
	json.NewDecoder(resp.Body).Decode(&created)
	return created.ProjectID
}

func TestIntegration_Epic_CreateAndList(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "epicowner@example.com", "12345678")
	projectID := createTestProject(t, token, "Epic Project")
	projectIDStr := strconv.Itoa(int(projectID))

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
			name:           "create epic success",
			method:         http.MethodPost,
			path:           "/api/v1/projects/" + projectIDStr + "/epics",
			authorization:  "Bearer " + token,
			body:           map[string]string{"title": "Authentication System"},
			expectedStatus: 201,
		},
		{
			name:           "create epic missing title",
			method:         http.MethodPost,
			path:           "/api/v1/projects/" + projectIDStr + "/epics",
			authorization:  "Bearer " + token,
			body:           map[string]string{"title": ""},
			expectedStatus: 400,
		},
		{
			name:           "list epics returns created epic",
			method:         http.MethodGet,
			path:           "/api/v1/projects/" + projectIDStr + "/epics",
			authorization:  "Bearer " + token,
			expectedStatus: 200,
			assertBody: func(t *testing.T, resp *http.Response) {
				var body struct {
					Data []map[string]interface{} `json:"data"`
				}
				json.NewDecoder(resp.Body).Decode(&body)
				assert.Len(t, body.Data, 1)
			},
		},
		{
			name:           "non-member cannot list epics",
			method:         http.MethodGet,
			path:           "/api/v1/projects/" + projectIDStr + "/epics",
			authorization:  "Bearer " + registerAndGetToken(t, "epicoutsider@example.com", "12345678"),
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
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

func TestIntegration_Epic_UpdateAndDelete_PermissionRules(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "epicrealowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "epicregularmember@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Permission Test Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addBody, _ := json.Marshal(map[string]string{"email": "epicregularmember@example.com"})
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	app.Test(reqAdd)

	epicBody, _ := json.Marshal(map[string]string{"title": "Original Epic"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/epics", bytes.NewReader(epicBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		EpicID uint `json:"epicId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	epicIDStr := strconv.Itoa(int(created.EpicID))

	tests := []struct {
		name           string
		method         string
		path           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "member can update epic",
			method:         http.MethodPut,
			path:           "/api/v1/epics/" + epicIDStr,
			authorization:  "Bearer " + memberToken,
			expectedStatus: 200,
		},
		{
			name:           "member cannot delete epic",
			method:         http.MethodDelete,
			path:           "/api/v1/epics/" + epicIDStr,
			authorization:  "Bearer " + memberToken,
			expectedStatus: 400,
		},
		{
			name:           "owner can delete epic",
			method:         http.MethodDelete,
			path:           "/api/v1/epics/" + epicIDStr,
			authorization:  "Bearer " + ownerToken,
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.method == http.MethodPut {
				bodyBytes, _ = json.Marshal(map[string]string{"title": "Updated by member"})
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.authorization)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
