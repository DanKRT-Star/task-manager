package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		email         string
		password      string
		existingUser  *model.User
		findByEmailErr error
		expectError   bool
	}{
		{
			name:            "success - email not taken",
			userName:        "Alice",
			email:           "alice@example.com",
			password:        "12345678",
			existingUser:    nil,
			findByEmailErr:  assert.AnError,
			expectError:     false,
		},
		{
			name:           "fail - email already taken",
			userName:       "Bob",
			email:          "bob@example.com",
			password:       "12345678",
			existingUser:   &model.User{UserID: 1, Email: "bob@example.com"},
			findByEmailErr: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("FindByEmail", tt.email).Return(tt.existingUser, tt.findByEmailErr)

			if !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			}

			authService := NewAuthService(mockRepo)
			user, err := authService.Register(tt.userName, tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				// Password phải được hash, không lưu plaintext
				assert.NotEqual(t, tt.password, user.HashedPassword)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Tạo sẵn 1 user với password đã hash để test compare
	tempService := &AuthService{}
	_ = tempService

	validUser := &model.User{
		UserID:   1,
		Email:    "test@example.com",
	}
	// Hash password "correctpassword" thật để test bcrypt compare hoạt động đúng
	hashed, _ := hashPasswordForTest("correctpassword")
	validUser.HashedPassword = hashed

	tests := []struct {
		name        string
		email       string
		password    string
		foundUser   *model.User
		findErr     error
		expectError bool
	}{
		{
			name:        "success - correct password",
			email:       "test@example.com",
			password:    "correctpassword",
			foundUser:   validUser,
			findErr:     nil,
			expectError: false,
		},
		{
			name:        "fail - wrong password",
			email:       "test@example.com",
			password:    "wrongpassword",
			foundUser:   validUser,
			findErr:     nil,
			expectError: true,
		},
		{
			name:        "fail - email not found",
			email:       "notexist@example.com",
			password:    "whatever",
			foundUser:   nil,
			findErr:     assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("FindByEmail", tt.email).Return(tt.foundUser, tt.findErr)

			authService := NewAuthService(mockRepo)
			token, err := authService.Login(tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}