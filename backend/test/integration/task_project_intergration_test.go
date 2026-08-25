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

// ---------------------------------------------------------------------
// CreateTask: epic/milestone/sprint phải thuộc đúng project
// ---------------------------------------------------------------------

func TestIntegration_Task_Create_EpicMilestoneSprint_Validation(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "taskepicowner@example.com", "12345678")

	projectA := createTestProject(t, ownerToken, "Project A")
	projectAStr := strconv.Itoa(int(projectA))
	projectB := createTestProject(t, ownerToken, "Project B")

	epicBody, _ := json.Marshal(map[string]string{"title": "Epic A"})
	reqEpic := httptest.NewRequest("POST", "/api/v1/projects/"+projectAStr+"/epics", bytes.NewReader(epicBody))
	reqEpic.Header.Set("Content-Type", "application/json")
	reqEpic.Header.Set("Authorization", "Bearer "+ownerToken)
	respEpic, err := app.Test(reqEpic)
	require.NoError(t, err)
	require.Equal(t, 201, respEpic.StatusCode)
	var epic struct {
		EpicID uint `json:"epicId"`
	}
	require.NoError(t, json.NewDecoder(respEpic.Body).Decode(&epic))

	milestoneBody, _ := json.Marshal(map[string]string{"title": "Milestone A"})
	reqMilestone := httptest.NewRequest("POST", "/api/v1/projects/"+projectAStr+"/milestones", bytes.NewReader(milestoneBody))
	reqMilestone.Header.Set("Content-Type", "application/json")
	reqMilestone.Header.Set("Authorization", "Bearer "+ownerToken)
	respMilestone, err := app.Test(reqMilestone)
	require.NoError(t, err)
	require.Equal(t, 201, respMilestone.StatusCode)
	var milestone struct {
		MilestoneID uint `json:"milestoneId"`
	}
	require.NoError(t, json.NewDecoder(respMilestone.Body).Decode(&milestone))

	sprintBody, _ := json.Marshal(map[string]string{"name": "Sprint A"})
	reqSprint := httptest.NewRequest("POST", "/api/v1/projects/"+projectAStr+"/sprints", bytes.NewReader(sprintBody))
	reqSprint.Header.Set("Content-Type", "application/json")
	reqSprint.Header.Set("Authorization", "Bearer "+ownerToken)
	respSprint, err := app.Test(reqSprint)
	require.NoError(t, err)
	require.Equal(t, 201, respSprint.StatusCode)
	var sprint struct {
		SprintID uint `json:"sprintId"`
	}
	require.NoError(t, json.NewDecoder(respSprint.Body).Decode(&sprint))

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "epic from another project rejected",
			body: map[string]interface{}{
				"title": "Task X", "deadline": "2026-08-15T17:00:00Z",
				"projectId": projectB, "epicId": epic.EpicID,
			},
			expectedStatus: 400,
		},
		{
			name: "milestone from another project rejected",
			body: map[string]interface{}{
				"title": "Task Y", "deadline": "2026-08-15T17:00:00Z",
				"projectId": projectB, "milestoneId": milestone.MilestoneID,
			},
			expectedStatus: 400,
		},
		{
			name: "sprint from another project rejected",
			body: map[string]interface{}{
				"title": "Task Z", "deadline": "2026-08-15T17:00:00Z",
				"projectId": projectB, "sprintId": sprint.SprintID,
			},
			expectedStatus: 400,
		},
		{
			name: "personal task cannot have epic",
			body: map[string]interface{}{
				"title": "Personal task", "deadline": "2026-08-15T17:00:00Z",
				"epicId": epic.EpicID,
			},
			expectedStatus: 400,
		},
		{
			name: "epic/milestone/sprint from same project accepted",
			body: map[string]interface{}{
				"title": "Task with valid links", "deadline": "2026-08-15T17:00:00Z",
				"projectId": projectA, "epicId": epic.EpicID,
				"milestoneId": milestone.MilestoneID, "sprintId": sprint.SprintID,
			},
			expectedStatus: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+ownerToken)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// ---------------------------------------------------------------------
// UpdateTask: gán assignee/epic/milestone/sprint
// ---------------------------------------------------------------------

func TestIntegration_Task_Update_EpicMilestoneSprintAssignee(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "updateepicowner@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Update Epic Project")
	projectIDStr := strconv.Itoa(int(projectID))

	addBody, _ := json.Marshal(map[string]string{"email": "updateepicmember@example.com"})
	registerAndGetToken(t, "updateepicmember@example.com", "12345678")
	reqAdd := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+ownerToken)
	respAdd, err := app.Test(reqAdd)
	require.NoError(t, err)
	require.Equal(t, 201, respAdd.StatusCode)
	var addedMember struct {
		UserID uint `json:"userId"`
	}
	require.NoError(t, json.NewDecoder(respAdd.Body).Decode(&addedMember))

	epicBody, _ := json.Marshal(map[string]string{"title": "Epic"})
	reqEpic := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/epics", bytes.NewReader(epicBody))
	reqEpic.Header.Set("Content-Type", "application/json")
	reqEpic.Header.Set("Authorization", "Bearer "+ownerToken)
	respEpic, err := app.Test(reqEpic)
	require.NoError(t, err)
	require.Equal(t, 201, respEpic.StatusCode)
	var epic struct {
		EpicID uint `json:"epicId"`
	}
	require.NoError(t, json.NewDecoder(respEpic.Body).Decode(&epic))

	taskBody, _ := json.Marshal(map[string]interface{}{
		"title": "Task to update", "deadline": "2026-08-15T17:00:00Z", "projectId": projectID,
	})
	reqCreate := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(taskBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+ownerToken)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)
	require.Equal(t, 201, respCreate.StatusCode)
	var createdTask struct {
		TaskID uint `json:"taskId"`
	}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&createdTask))
	taskIDStr := strconv.Itoa(int(createdTask.TaskID))

	personalBody, _ := json.Marshal(map[string]string{"title": "Personal task", "deadline": "2026-08-15T17:00:00Z"})
	reqPersonal := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(personalBody))
	reqPersonal.Header.Set("Content-Type", "application/json")
	reqPersonal.Header.Set("Authorization", "Bearer "+ownerToken)
	respPersonal, err := app.Test(reqPersonal)
	require.NoError(t, err)
	require.Equal(t, 201, respPersonal.StatusCode)
	var personalTask struct {
		TaskID uint `json:"taskId"`
	}
	require.NoError(t, json.NewDecoder(respPersonal.Body).Decode(&personalTask))
	personalTaskIDStr := strconv.Itoa(int(personalTask.TaskID))

	t.Run("assign member and epic to project task succeeds", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]interface{}{
			"assigneeId": addedMember.UserID, "epicId": epic.EpicID,
		})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("cannot assign epic to a personal task", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]interface{}{"epicId": epic.EpicID})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+personalTaskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("cannot assign assignee to a personal task", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]interface{}{"assigneeId": addedMember.UserID})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+personalTaskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("assignee not a project member rejected", func(t *testing.T) {
		nonMemberID := uint(999999)
		updateBody, _ := json.Marshal(map[string]interface{}{"assigneeId": nonMemberID})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

// ---------------------------------------------------------------------
// Permission: owner của project override quyền sửa/xoá task của member khác
// ---------------------------------------------------------------------

func TestIntegration_Task_ProjectOwnerCanModifyAnyMemberTask(t *testing.T) {
	cleanTables()
	ownerToken := registerAndGetToken(t, "taskpermowner@example.com", "12345678")
	memberAToken := registerAndGetToken(t, "taskpermmemberA@example.com", "12345678")
	memberBToken := registerAndGetToken(t, "taskpermmemberB@example.com", "12345678")

	projectID := createTestProject(t, ownerToken, "Permission Task Project")
	projectIDStr := strconv.Itoa(int(projectID))

	for _, email := range []string{"taskpermmemberA@example.com", "taskpermmemberB@example.com"} {
		addBody, _ := json.Marshal(map[string]string{"email": email})
		req := httptest.NewRequest("POST", "/api/v1/projects/"+projectIDStr+"/members", bytes.NewReader(addBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)
	}

	taskBody, _ := json.Marshal(map[string]interface{}{
		"title": "Member A's task", "deadline": "2026-08-15T17:00:00Z", "projectId": projectID,
	})
	reqCreate := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(taskBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+memberAToken)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)
	require.Equal(t, 201, respCreate.StatusCode)
	var createdTask struct {
		TaskID uint `json:"taskId"`
	}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&createdTask))
	taskIDStr := strconv.Itoa(int(createdTask.TaskID))

	t.Run("member B (not creator, not assignee, not owner) cannot update", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]string{"status": "done"})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberBToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("owner can update member A's task", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]string{"status": "in_progress"})
		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskIDStr, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("owner can delete member A's task", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+taskIDStr, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// ---------------------------------------------------------------------
// CreateTask: status không hợp lệ
// ---------------------------------------------------------------------

func TestIntegration_Task_Create_InvalidStatus(t *testing.T) {
	cleanTables()
	token := registerAndGetToken(t, "invalidstatus@example.com", "12345678")

	body, _ := json.Marshal(map[string]string{
		"title": "Bad status task", "deadline": "2026-08-15T17:00:00Z", "status": "not_a_real_status",
	})
	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}