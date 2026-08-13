package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c fiber.Ctx) error {
	start := time.Now()
	path := c.Route().Path
	if path == "" {
		path = c.Path()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		logger.AuthRequestRejected("missing_authorization_header", c.Method(), path, 401, time.Since(start), "missing authorization header")
		return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		logger.AuthRequestRejected("invalid_authorization_format", c.Method(), path, 401, time.Since(start), "invalid authorization format")
		return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		logger.AuthRequestRejected("invalid_or_expired_token", c.Method(), path, 401, time.Since(start), "invalid or expired token")
		return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.AuthRequestRejected("invalid_token_claims", c.Method(), path, 401, time.Since(start), "invalid token claims")
		return c.Status(401).JSON(fiber.Map{"error": "invalid token claims"})
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		logger.AuthRequestRejected("invalid_token_payload", c.Method(), path, 401, time.Since(start), "invalid token payload")
		return c.Status(401).JSON(fiber.Map{"error": "invalid token payload"})
	}

	c.Locals("userID", uint(userID))
	logger.RequestHandled(c.Method(), path, 200, time.Since(start), "auth middleware ok")
	return c.Next()
}