package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

type AuthService struct {
	UserRepo         repository.UserRepositoryInterface
	RefreshTokenRepo repository.RefreshTokenRepositoryInterface
}

func NewAuthService(userRepo repository.UserRepositoryInterface, refreshTokenRepo repository.RefreshTokenRepositoryInterface) *AuthService {
	return &AuthService{UserRepo: userRepo, RefreshTokenRepo: refreshTokenRepo}
}

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

// Login trả về (accessToken, refreshToken, error)
func (s *AuthService) Login(email, password string) (string, string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		logger.AuthLoginUserNotFound(email)
		return "", "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		logger.AuthLoginInvalidPassword(user.UserID, email)
		return "", "", errors.New("invalid email or password")
	}

	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		logger.AuthLoginTokenGenerationFailed(user.UserID, err)
		return "", "", err
	}

	refreshToken, err := s.issueRefreshToken(user.UserID)
	if err != nil {
		logger.AuthLoginTokenGenerationFailed(user.UserID, err)
		return "", "", err
	}

	logger.AuthLoginSuccess(user.UserID, email)
	return accessToken, refreshToken, nil
}

// RefreshAccessToken xoay vòng refresh token: thu hồi token cũ, cấp cặp token mới
func (s *AuthService) RefreshAccessToken(rawToken string) (string, string, error) {
	record, err := s.RefreshTokenRepo.FindByHash(hashToken(rawToken))
	if err != nil {
		logger.AuthRefreshTokenInvalid()
		return "", "", errors.New("invalid or expired refresh token")
	}

	if record.RevokedAt != nil || time.Now().After(record.ExpiresAt) {
		logger.AuthRefreshTokenExpiredOrRevoked(record.UserID)
		return "", "", errors.New("invalid or expired refresh token")
	}

	if err := s.RefreshTokenRepo.Revoke(record.RefreshTokenID); err != nil {
		return "", "", err
	}

	newAccessToken, err := generateAccessToken(record.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.issueRefreshToken(record.UserID)
	if err != nil {
		return "", "", err
	}

	logger.AuthRefreshTokenRotated(record.UserID)
	return newAccessToken, newRefreshToken, nil
}

// Logout thu hồi refresh token — access token cũ vẫn còn hiệu lực tới khi hết hạn (15 phút), đây là đánh đổi chấp nhận được cho JWT stateless
func (s *AuthService) Logout(rawToken string) error {
	record, err := s.RefreshTokenRepo.FindByHash(hashToken(rawToken))
	if err != nil {
		return nil // token không tồn tại/đã hết hạn -> coi như đã logout, idempotent
	}
	if err := s.RefreshTokenRepo.Revoke(record.RefreshTokenID); err != nil {
		return err
	}
	logger.AuthLogoutSuccess(record.UserID)
	return nil
}

func (s *AuthService) GetMe(userID uint) (*model.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		logger.AuthGetMeFailed(userID, err)
		return nil, err
	}
	return user, nil
}

func (s *AuthService) issueRefreshToken(userID uint) (string, error) {
	raw, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	record := &model.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(raw),
		ExpiresAt: time.Now().Add(RefreshTokenExpiry),
	}
	if err := s.RefreshTokenRepo.Create(record); err != nil {
		return "", err
	}
	return raw, nil
}

func generateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(AccessTokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func generateRandomToken(numBytes int) (string, error) {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}