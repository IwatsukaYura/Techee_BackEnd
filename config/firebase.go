package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
)

// 初期化処理（mainで呼び出す想定）
func InitFirebase() {
	ctx := context.Background()

	// 環境変数からプロジェクトIDを取得
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Fatalf("FIREBASE_PROJECT_ID environment variable is not set")
	}

	// Firebase設定オブジェクトを作成（プロジェクトIDを指定）
	conf := &firebase.Config{
		ProjectID: projectID,
	}

	// 環境変数からサービスアカウントキーファイルのパスを取得し、optionとして使用
	serviceAccountKeyPath := os.Getenv("SERVICE_ACCOUNT_KEY_JSON") // Renderで設定したKey名
	if serviceAccountKeyPath == "" {
		log.Fatalf("SERVICE_ACCOUNT_KEY_JSON environment variable is not set")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyPath) // 環境変数で指定されたパスを使用

	// Firebase Appの初期化
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("failed to initialize firebase app: %v", err)
	}
	FirebaseApp = app

	// Firestoreクライアントの初期化
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to initialize firestore: %v", err)
	}
	FirestoreClient = client
}
