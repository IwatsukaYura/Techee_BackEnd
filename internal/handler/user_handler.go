package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"context"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/middleware"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// UserService はユーザーサービス層へのインターフェースです。
type UserService interface {
	GetUser(ctx context.Context, userID string) (*model.User, error)
	UpdateUserTags(ctx context.Context, userID string, tags []string) error
	// 必要に応じて他のメソッドを追加
}

// UserHandler はユーザー関連のリクエストを処理するハンドラーです。
type UserHandler struct {
	service UserService
}

// NewUserHandler はUserHandlerの新しいインスタンスを作成します。
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ユーザー情報取得ハンドラー
func (h *UserHandler) GetUser(c echo.Context) error {
	uid, ok := c.Get(middleware.ContextUIDKey).(string)
	if !ok || uid == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	ctx := c.Request().Context()

	// サービス層を介してユーザー情報を取得
	user, err := h.service.GetUser(ctx, uid)
	if err != nil {
		// エラーハンドリングを適切に行う
		c.Logger().Errorf("failed to get user from service: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ユーザー情報取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, user)
}

// ユーザーのタグ更新ハンドラー
func (h *UserHandler) UpdateUserTags(c echo.Context) error {
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

	ctx := c.Request().Context()

	// サービス層を介してユーザーのタグを更新
	if err := h.service.UpdateUserTags(ctx, uid, req.Tags); err != nil {
		// エラーハンドリングを適切に行う
		c.Logger().Errorf("failed to update user tags via service: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "タグ更新に失敗しました"})
	}

	return c.JSON(http.StatusOK, echo.Map{"tags": req.Tags})
}
