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

func TestIntegration_Milestone_CreateAndList(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "milestoneowner@example.com", "12345678")
	projectID := createTestProject(t, token, "Milestone Project")
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
			name:           "create milestone success",
			method:         http.MethodPost,
			path:           "/api/v1/projects/" + projectIDStr + "/milestones",
			authorization:  "Bearer " + token,
			body:           map[string]string{"title": "v1.0 Release"},
			expectedStatus: 201,
		},
		{
			name:           "list milestones returns created milestone",
			method:         http.MethodGet,
			path:           "/api/v1/projects/" + projectIDStr + "/milestones",
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

func TestIntegration_Milestone_UpdateAndDelete_PermissionRules(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "milestonerealowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "milestoneregularmember@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Milestone Permission Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addBody, _ := json.Marshal(map[string]string{"email": "milestoneregularmember@example.com"})
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	app.Test(reqAdd)

	milestoneBody, _ := json.Marshal(map[string]string{"title": "Original Milestone"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/milestones", bytes.NewReader(milestoneBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		MilestoneID uint `json:"milestoneId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	milestoneIDStr := strconv.Itoa(int(created.MilestoneID))

	t.Run("member can update milestone", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]string{"title": "Updated by member"})
		req := httptest.NewRequest("PUT", "/api/v1/milestones/"+milestoneIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("member cannot delete milestone", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/milestones/"+milestoneIDStr, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("owner can delete milestone", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/milestones/"+milestoneIDStr, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
