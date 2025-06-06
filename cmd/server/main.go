package main

import (
	"log"
	"net/http"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/handler"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitFirebase()

	e := echo.New()
	// CORSミドルウェアを追加（このブロックを追加）
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // フロントエンドのオリジンを許可
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// 記事一覧API（仮実装）
	e.GET("/api/articles", handler.GetArticles)
	e.GET("/api/user", handler.GetUser, middleware.FirebaseAuth)
	e.PUT("/api/user/tags", handler.UpdateUserTags, middleware.FirebaseAuth)

	log.Println("Server started at :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
