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

func TestIntegration_CreateTask(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "creator@example.com", "12345678")

	tests := []struct {
		name           string
		body           interface{}
		withToken      bool
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]string{
				"title":    "Học Golang",
				"deadline": "2026-08-15T17:00:00Z",
			},
			withToken:      true,
			expectedStatus: 201,
		},
		{
			name: "missing title",
			body: map[string]string{
				"deadline": "2026-08-15T17:00:00Z",
			},
			withToken:      true,
			expectedStatus: 400,
		},
		{
			name: "invalid deadline format",
			body: map[string]string{
				"title":    "Task X",
				"deadline": "15/08/2026",
			},
			withToken:      true,
			expectedStatus: 400,
		},
		{
			name: "no token",
			body: map[string]string{
				"title": "Task Y",
			},
			withToken:      false,
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.withToken {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_GetTasks_FilterSortPaginate(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "lister@example.com", "12345678")

	// Seed 3 task với status khác nhau
	seedTasks := []map[string]string{
		{"title": "Task pending", "status": "pending", "deadline": "2026-08-10T00:00:00Z"},
		{"title": "Task in progress", "status": "in_progress", "deadline": "2026-08-20T00:00:00Z"},
		{"title": "Task done", "status": "done", "deadline": "2026-08-05T00:00:00Z"},
	}
	for _, task := range seedTasks {
		jsonBody, _ := json.Marshal(task)
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		app.Test(req)
	}

	tests := []struct {
		name           string
		query          string
		expectedLength int
	}{
		{
			name:           "get all",
			query:          "/api/v1/tasks?page=1&limit=10",
			expectedLength: 3,
		},
		{
			name:           "filter by status pending",
			query:          "/api/v1/tasks?status=pending&page=1&limit=10",
			expectedLength: 1,
		},
		{
			name:           "pagination limit 2",
			query:          "/api/v1/tasks?page=1&limit=2",
			expectedLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.query, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			var body struct {
				Data []map[string]interface{} `json:"data"`
			}
			json.NewDecoder(resp.Body).Decode(&body)
			assert.Len(t, body.Data, tt.expectedLength)
		})
	}
}

func TestIntegration_Task_CannotAccessOtherUsersTask(t *testing.T) {
	cleanTables()

	tokenA := registerAndGetToken(t, "userA@example.com", "12345678")
	tokenB := registerAndGetToken(t, "userB@example.com", "12345678")

	// User A tạo task
	createBody := map[string]string{
		"title":    "User A's private task",
		"deadline": "2026-08-15T17:00:00Z",
	}
	jsonBody, _ := json.Marshal(createBody)
	reqCreate := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+tokenA)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)
	require.Equal(t, 201, respCreate.StatusCode)

	var created struct {
		TaskID uint `json:"taskId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	taskIDStr := strconv.Itoa(int(created.TaskID))

	tests := []struct {
		name           string
		method         string
		path           string
		token          string
		expectedStatus int
	}{
		{
			name:           "user B cannot update user A's task",
			method:         "PUT",
			path:           "/api/v1/tasks/" + taskIDStr,
			token:          tokenB,
			expectedStatus: 400,
		},
		{
			name:           "user B cannot delete user A's task",
			method:         "DELETE",
			path:           "/api/v1/tasks/" + taskIDStr,
			token:          tokenB,
			expectedStatus: 400,
		},
		{
			name:           "user A can update own task",
			method:         "PUT",
			path:           "/api/v1/tasks/" + taskIDStr,
			token:          tokenA,
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateBody := map[string]string{"status": "done"}
			jsonUpdate, _ := json.Marshal(updateBody)

			var req *httptest.ResponseRecorder
			_ = req

			r := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(jsonUpdate))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", "Bearer "+tt.token)

			resp, err := app.Test(r)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_UpdateTask(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "updater@example.com", "12345678")

	// Tạo task trước để update
	createBody := map[string]string{
		"title":    "Original title",
		"deadline": "2026-08-15T17:00:00Z",
	}
	jsonBody, _ := json.Marshal(createBody)
	reqCreate := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		TaskID uint `json:"taskId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	taskIDStr := strconv.Itoa(int(created.TaskID))

	tests := []struct {
		name           string
		taskID         string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:   "update title success",
			taskID: taskIDStr,
			body: map[string]string{
				"title": "Updated title",
			},
			expectedStatus: 200,
		},
		{
			name:   "update status success",
			taskID: taskIDStr,
			body: map[string]string{
				"status": "in_progress",
			},
			expectedStatus: 200,
		},
		{
			name:   "invalid status value",
			taskID: taskIDStr,
			body: map[string]string{
				"status": "not_a_real_status",
			},
			expectedStatus: 400,
		},
		{
			name:   "invalid deadline format",
			taskID: taskIDStr,
			body: map[string]string{
				"deadline": "not-a-date",
			},
			expectedStatus: 400,
		},
		{
			name:   "task not found",
			taskID: "99999",
			body: map[string]string{
				"title": "Doesn't matter",
			},
			expectedStatus: 400,
		},
		{
			name:   "invalid task id format",
			taskID: "abc",
			body: map[string]string{
				"title": "Doesn't matter",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonUpdate, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/api/v1/tasks/"+tt.taskID, bytes.NewReader(jsonUpdate))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_DeleteTask(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "deleter@example.com", "12345678")

	createBody := map[string]string{
		"title":    "Task to delete",
		"deadline": "2026-08-15T17:00:00Z",
	}
	jsonBody, _ := json.Marshal(createBody)
	reqCreate := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	respCreate, _ := app.Test(reqCreate)

	var created struct {
		TaskID uint `json:"taskId"`
	}
	json.NewDecoder(respCreate.Body).Decode(&created)
	taskIDStr := strconv.Itoa(int(created.TaskID))

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "delete success",
			taskID:         taskIDStr,
			expectedStatus: 200,
		},
		{
			name:           "delete already deleted task",
			taskID:         taskIDStr,
			expectedStatus: 400,
		},
		{
			name:           "delete non-existent task",
			taskID:         "99999",
			expectedStatus: 400,
		},
		{
			name:           "invalid task id format",
			taskID:         "abc",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+tt.taskID, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_GetTasks_Sort(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "sorter@example.com", "12345678")

	seedTasks := []map[string]string{
		{"title": "Deadline sớm nhất", "deadline": "2026-08-05T00:00:00Z"},
		{"title": "Deadline trễ nhất", "deadline": "2026-08-25T00:00:00Z"},
		{"title": "Deadline giữa", "deadline": "2026-08-15T00:00:00Z"},
	}
	for _, task := range seedTasks {
		jsonBody, _ := json.Marshal(task)
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		app.Test(req)
	}

	t.Run("sort ascending by default", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tasks?page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var body struct {
			Data []struct {
				Title string `json:"title"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		require.Len(t, body.Data, 3)
		assert.Equal(t, "Deadline sớm nhất", body.Data[0].Title)
	})

	t.Run("sort descending", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tasks?sort=deadline_desc&page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var body struct {
			Data []struct {
				Title string `json:"title"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		require.Len(t, body.Data, 3)
		assert.Equal(t, "Deadline trễ nhất", body.Data[0].Title)
	})
}

func TestIntegration_GetTasks_NoToken(t *testing.T) {
	cleanTables()

	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestIntegration_GetTasks_MalformedToken(t *testing.T) {
	cleanTables()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "no Bearer prefix",
			authHeader:     "sometoken123",
			expectedStatus: 401,
		},
		{
			name:           "invalid token signature",
			authHeader:     "Bearer invalid.token.string",
			expectedStatus: 401,
		},
		{
			name:           "empty Bearer",
			authHeader:     "Bearer ",
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
			req.Header.Set("Authorization", tt.authHeader)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_GetTasks_UserIsolation(t *testing.T) {
	cleanTables()

	tokenA := registerAndGetToken(t, "isolationA@example.com", "12345678")
	tokenB := registerAndGetToken(t, "isolationB@example.com", "12345678")

	// User A tạo 2 task
	for _, title := range []string{"A's task 1", "A's task 2"} {
		body := map[string]string{"title": title, "deadline": "2026-08-15T17:00:00Z"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenA)
		app.Test(req)
	}

	// User B tạo 1 task
	bodyB := map[string]string{"title": "B's task", "deadline": "2026-08-15T17:00:00Z"}
	jsonBodyB, _ := json.Marshal(bodyB)
	reqB := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBodyB))
	reqB.Header.Set("Content-Type", "application/json")
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	app.Test(reqB)

	// User B gọi GetTasks -> chỉ được thấy đúng 1 task của mình, không lẫn task của A
	req := httptest.NewRequest("GET", "/api/v1/tasks?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body struct {
		Data []struct {
			Title string `json:"title"`
		} `json:"data"`
		Total int64 `json:"total"`
	}
	json.NewDecoder(resp.Body).Decode(&body)

	assert.Equal(t, int64(1), body.Total)
	require.Len(t, body.Data, 1)
	assert.Equal(t, "B's task", body.Data[0].Title)
}

func TestIntegration_GetTasks_EmptyList(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "emptylist@example.com", "12345678")

	req := httptest.NewRequest("GET", "/api/v1/tasks?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body struct {
		Data  []interface{} `json:"data"`
		Total int64         `json:"total"`
	}
	json.NewDecoder(resp.Body).Decode(&body)

	assert.Equal(t, int64(0), body.Total)
	assert.Empty(t, body.Data)
}

func TestIntegration_GetTasks_PaginationEdgeCases(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "paginationedge@example.com", "12345678")

	// Seed 3 task
	for i := 0; i < 3; i++ {
		body := map[string]string{"title": "Task", "deadline": "2026-08-15T17:00:00Z"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		app.Test(req)
	}

	tests := []struct {
		name           string
		query          string
		expectedLength int
	}{
		{
			name:           "limit 0 falls back to default 10",
			query:          "/api/v1/tasks?page=1&limit=0",
			expectedLength: 3,
		},
		{
			name:           "limit exceeds max 100 falls back to default 10",
			query:          "/api/v1/tasks?page=1&limit=999",
			expectedLength: 3,
		},
		{
			name:           "page 0 falls back to page 1",
			query:          "/api/v1/tasks?page=0&limit=10",
			expectedLength: 3,
		},
		{
			name:           "negative page falls back to page 1",
			query:          "/api/v1/tasks?page=-1&limit=10",
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.query, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			var body struct {
				Data []interface{} `json:"data"`
			}
			json.NewDecoder(resp.Body).Decode(&body)
			assert.Len(t, body.Data, tt.expectedLength)
		})
	}
}