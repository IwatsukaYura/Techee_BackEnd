package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/labstack/echo/v4"
)

// Contextキー
const ContextUIDKey = "firebase_uid"

// Firebase認証ミドルウェア（本実装）
func FirebaseAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing or invalid token"})
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		ctx := context.Background()
		auth, err := config.FirebaseApp.Auth(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "firebase auth init failed"})
		}
		token, err := auth.VerifyIDToken(ctx, tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}
		// UIDをContextに格納
		c.Set(ContextUIDKey, token.UID)
		return next(c)
	}
}
