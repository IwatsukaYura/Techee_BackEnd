package config

import (
	"context"
	"log"

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
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("failed to initialize firebase app: %v", err)
	}
	FirebaseApp = app

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to initialize firestore: %v", err)
	}
	FirestoreClient = client
}
