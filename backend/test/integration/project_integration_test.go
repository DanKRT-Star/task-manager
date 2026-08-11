package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_CreateProject(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "projectcreator@example.com", "12345678")

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "success",
			body:           map[string]string{"name": "My Project", "description": "A test project"},
			expectedStatus: 201,
		},
		{
			name:           "missing name",
			body:           map[string]string{"description": "No name here"},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_GetProjects_OnlyShowsUsersOwnProjects(t *testing.T) {
	cleanTables()
	tokenA := registerAndGetToken(t, "projA@example.com", "12345678")
	tokenB := registerAndGetToken(t, "projB@example.com", "12345678")

	// User A tạo 1 project
	body, _ := json.Marshal(map[string]string{"name": "A's Project"})
	req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenA)
	app.Test(req)

	// User B gọi GetProjects -> không thấy project của A
	reqB := httptest.NewRequest("GET", "/api/v1/projects", nil)
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	respB, err := app.Test(reqB)
	require.NoError(t, err)
	assert.Equal(t, 200, respB.StatusCode)

	var bodyB struct {
		Data []map[string]interface{} `json:"data"`
	}
	json.NewDecoder(respB.Body).Decode(&bodyB)
	assert.Empty(t, bodyB.Data)
}

func TestIntegration_Project_MemberManagement(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "owner@example.com", "12345678")
	memberEmail := "member@example.com"
	registerAndGetToken(t, memberEmail, "12345678") // chỉ cần user tồn tại

	// Owner tạo project
	body, _ := json.Marshal(map[string]string{"name": "Team Project"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)
	require.Equal(t, 201, respCreate.StatusCode)

	var created struct {
		ProjectID uint `json:"projectId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	projectIDStr := strconv.Itoa(int(created.ProjectID))

	t.Run("owner can add member by email", func(t *testing.T) {
		addBody, _ := json.Marshal(map[string]string{"email": memberEmail})
		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("cannot add the same member twice", func(t *testing.T) {
		addBody, _ := json.Marshal(map[string]string{"email": memberEmail})
		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("cannot add a non-existent user", func(t *testing.T) {
		addBody, _ := json.Marshal(map[string]string{"email": "ghost@example.com"})
		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("get members list includes both owner and member", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectIDStr+"/members", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var body struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Len(t, body.Data, 2)
	})
}

func TestIntegration_Project_OnlyOwnerCanManage(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "realowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "regularmember@example.com", "12345678")

	// Owner tạo project
	body, _ := json.Marshal(map[string]string{"name": "Restricted Project"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		ProjectID uint `json:"projectId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	projectIDStr := strconv.Itoa(int(created.ProjectID))

	// Thêm member vào project
	addBody, _ := json.Marshal(map[string]string{"email": "regularmember@example.com"})
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	app.Test(reqAdd)

	t.Run("member cannot update project", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]string{"name": "Hacked Name"})
		req := httptest.NewRequest("PUT", "/api/v1/projects/"+projectIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("member cannot delete project", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/projects/"+projectIDStr, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("member cannot add other members", func(t *testing.T) {
		addBody, _ := json.Marshal(map[string]string{"email": "someone@example.com"})
		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("owner can update project", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]string{"name": "Updated By Owner"})
		req := httptest.NewRequest("PUT", "/api/v1/projects/"+projectIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestIntegration_CreateTask_InProject_MemberAssignment(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "taskowner@example.com", "12345678")
	memberToken := registerAndGetToken(t, "taskmember@example.com", "12345678")

	// Tạo project + thêm member
	body, _ := json.Marshal(map[string]string{"name": "Task Project"})
	reqCreate := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		ProjectID uint `json:"projectId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	projectID := created.ProjectID
	projectIDStr := strconv.Itoa(int(projectID))

	addBody, _ := json.Marshal(map[string]string{"email": "taskmember@example.com"})
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	respAdd, _ := app.Test(reqAdd)

	var addedMember struct {
		UserID uint `json:"userId"`
	}
	json.NewDecoder(respAdd.Body).Decode(&addedMember)

	t.Run("owner creates task in project and assigns to member", func(t *testing.T) {
		taskBody, _ := json.Marshal(map[string]interface{}{
			"title":      "Project task",
			"deadline":   "2026-08-15T17:00:00Z",
			"projectId":  projectID,
			"assigneeId": addedMember.UserID,
		})
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(taskBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("non-member cannot be assigned to project task", func(t *testing.T) {
		outsiderToken := registerAndGetToken(t, "outsider@example.com", "12345678")
		_ = outsiderToken

		fakeOutsiderID := uint(9999)
		taskBody, _ := json.Marshal(map[string]interface{}{
			"title":      "Should fail",
			"deadline":   "2026-08-15T17:00:00Z",
			"projectId":  projectID,
			"assigneeId": fakeOutsiderID,
		})
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(taskBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("member can view project tasks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectIDStr+"/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+memberToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("non-member cannot view project tasks", func(t *testing.T) {
		outsiderToken := registerAndGetToken(t, "viewer_outsider@example.com", "12345678")

		req := httptest.NewRequest("GET", "/api/v1/projects/"+projectIDStr+"/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+outsiderToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}