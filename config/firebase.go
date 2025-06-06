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

	// サービスアカウントJSONのパスを指定
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")

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
