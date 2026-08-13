package service

import (
	"errors"
	"os"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repository.UserRepositoryInterface
}

func NewAuthService(userRepo repository.UserRepositoryInterface) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Register tạo user mới, hash password trước khi lưu
func (s *AuthService) Register(userName, email, password string) (*model.User, error) {
	existing, _ := s.UserRepo.FindByEmail(email)
	if existing != nil {
		logger.AuthRegisterEmailExists(email)
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.AuthRegisterPasswordHashFailed(email, err)
		return nil, err
	}

	user := &model.User{
		UserName:       userName,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	if err := s.UserRepo.Create(user); err != nil {
		logger.AuthRegisterFailed(email, err)
		return nil, err
	}

	logger.AuthRegisterSuccess(user.UserID, email)
	return user, nil
}

// Login kiểm tra email/password, trả về JWT token nếu đúng
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		logger.AuthLoginUserNotFound(email)
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		logger.AuthLoginInvalidPassword(user.UserID, email)
		return "", errors.New("invalid email or password")
	}

	token, err := generateJWT(user.UserID)
	if err != nil {
		logger.AuthLoginTokenGenerationFailed(user.UserID, err)
		return "", err
	}

	logger.AuthLoginSuccess(user.UserID, email)
	return token, nil
}

func (s *AuthService) GetMe(userID uint) (*model.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		logger.AuthGetMeFailed(userID, err)
		return nil, err
	}

	return user, nil
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}
