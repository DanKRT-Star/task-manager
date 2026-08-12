package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Auth_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success",
			body: map[string]string{
				"userName": "Alice",
				"email":    "alice@example.com",
				"password": "12345678",
			},
			expectedStatus: 201,
			expectError:    false,
		},
		{
			name: "missing email",
			body: map[string]string{
				"userName": "Bob",
				"password": "12345678",
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name: "password too short",
			body: map[string]string{
				"userName": "Carol",
				"email":    "carol@example.com",
				"password": "123",
			},
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "invalid json",
			body:           "{invalid-json",
			expectedStatus: 400,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanTables()

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_Auth_RegisterDuplicateEmail(t *testing.T) {
	cleanTables()

	body := map[string]string{
		"userName": "Test",
		"email":    "dup@example.com",
		"password": "12345678",
	}
	jsonBody, _ := json.Marshal(body)

	req1 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, 201, resp1.StatusCode)

	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 400, resp2.StatusCode)
}

func TestIntegration_Auth_Login(t *testing.T) {
	cleanTables()

	registerBody := map[string]string{
		"userName": "Test",
		"email":    "login@example.com",
		"password": "correctpassword",
	}
	jsonRegister, _ := json.Marshal(registerBody)
	reqRegister := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonRegister))
	reqRegister.Header.Set("Content-Type", "application/json")
	app.Test(reqRegister)

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name: "correct credentials",
			body: map[string]string{
				"email":    "login@example.com",
				"password": "correctpassword",
			},
			expectedStatus: 200,
		},
		{
			name: "wrong password",
			body: map[string]string{
				"email":    "login@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: 401,
		},
		{
			name: "email not found",
			body: map[string]string{
				"email":    "notexist@example.com",
				"password": "whatever",
			},
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_Auth_LoginMissingFields(t *testing.T) {
	cleanTables()

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name: "missing password",
			body: map[string]string{
				"email": "someone@example.com",
			},
			expectedStatus: 400,
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": "12345678",
			},
			expectedStatus: 400,
		},
		{
			name:           "empty body",
			body:           map[string]string{},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIntegration_Auth_GetMe(t *testing.T) {
	cleanTables()

	token := registerAndGetToken(t, "me@example.com", "correctpassword")

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "Test User", result["userName"])
	assert.Equal(t, "me@example.com", result["email"])
}
