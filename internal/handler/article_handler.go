package handler

import (
	"net/http"
	"strconv"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

// 記事一覧取得
func GetArticles(c echo.Context) error {
	tag := c.QueryParam("tag")
	limitStr := c.QueryParam("limit")
	limit := 15
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 30 {
			limit = l
		}
	}
	articles, err := repository.GetCachedArticles(tag, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch articles"})
	}
	return c.JSON(http.StatusOK, articles)
}
