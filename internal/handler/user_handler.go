package handler

import (
	"net/http"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/middleware"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

// ユーザー情報取得
func GetUser(c echo.Context) error {
	uid, ok := c.Get(middleware.ContextUIDKey).(string)
	if !ok || uid == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	user, err := repository.GetUserByID(uid)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, user)
}

// ユーザーのタグ更新
func UpdateUserTags(c echo.Context) error {
	type reqBody struct {
		Tags []string `json:"tags"`
	}
	var req reqBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	uid, ok := c.Get(middleware.ContextUIDKey).(string)
	if !ok || uid == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if err := repository.UpdateUserTags(uid, req.Tags); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update tags"})
	}
	return c.JSON(http.StatusOK, echo.Map{"tags": req.Tags})
}
