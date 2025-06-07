package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"context"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// ArticleService は記事サービス層へのインターフェースです。
type ArticleService interface {
	GetPopularArticles(ctx context.Context, tag string) ([]model.Article, error)
	// 必要に応じて他のメソッドを追加
}

// ArticleHandler は記事関連のリクエストを処理するハンドラーです。
type ArticleHandler struct {
	service ArticleService
}

// NewArticleHandler はArticleHandlerの新しいインスタンスを作成します。
func NewArticleHandler(service ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// 記事一覧取得ハンドラー
func (h *ArticleHandler) GetArticles(c echo.Context) error {
	tag := c.QueryParam("tag")
	// limitStr := c.QueryParam("limit") // 現在はlimitを使っていない
	// limit := 15
	// if limitStr != "" {
	// 	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 30 {
	// 		limit = l
	// 	}
	// }

	ctx := c.Request().Context()

	// サービス層を介して記事を取得
	articles, err := h.service.GetPopularArticles(ctx, tag)
	if err != nil {
		// エラーハンドリングを適切に行う
		c.Logger().Errorf("failed to get articles from service: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "記事の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, articles)
}
