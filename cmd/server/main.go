package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/handler"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/middleware"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/repository"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	// Firebaseの初期化
	config.InitFirebase() // 引数なしで呼び出す

	// Firestoreクライアントをグローバル変数から取得
	firestoreClient := config.FirestoreClient
	if firestoreClient == nil {
		log.Fatalf("Firestore client is not initialized") // 初期化失敗時のチェック
	}

	// リポジトリ層の初期化
	articleRepo := repository.NewArticleRepository(firestoreClient)
	userRepo := repository.NewUserRepository(firestoreClient) // userRepoも初期化

	// サービス層の初期化
	articleService := service.NewArticleService(articleRepo)
	userService := service.NewUserService(userRepo)

	// ★ サーバー起動時に一度だけ記事の取得と保存を実行 ★
	// 注意: これは簡易実装です。本来は定期実行されるバッチ処理などで実行すべきです。
	// サーバー起動をブロックする可能性があるため、非同期化やエラーハンドリングを適切に行う必要があります。

	// 取得対象のタグリスト (例)
	tagsToFetch := []string{"Go", "Python", "JavaScript", "TypeScript", "React", "Vue", "Docker", "Kubernetes", "AWS", "Firebase"}

	log.Println("Fetching and saving initial articles...")
	fetchCtx, cancel := context.WithTimeout(ctx, 60*time.Second) // タイムアウトを設定
	defer cancel()

	err := articleService.FetchAndSaveArticles(fetchCtx, tagsToFetch)
	if err != nil {
		log.Printf("Warning: failed to fetch and save initial articles: %v", err) // エラーでもサーバーは起動させる
	} else {
		log.Println("Initial articles fetched and saved successfully.")
	}

	// Echoサーバーの初期化
	e := echo.New()

	// CORSミドルウェアを追加
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // フロントエンドのオリジンを許可
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// ハンドラー層の初期化とルーティング設定
	articleHandler := handler.NewArticleHandler(articleService) // ArticleServiceをハンドラーに渡す
	userHandler := handler.NewUserHandler(userService)          // UserServiceをハンドラーに渡す

	// 記事一覧API
	e.GET("/api/articles", articleHandler.GetArticles)

	// ユーザー関連API（認証ミドルウェア適用）
	e.GET("/api/user", userHandler.GetUser, middleware.FirebaseAuth)
	e.PUT("/api/user/tags", userHandler.UpdateUserTags, middleware.FirebaseAuth)

	log.Println("Server started at :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
