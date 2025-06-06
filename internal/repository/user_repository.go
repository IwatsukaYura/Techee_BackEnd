package repository

import (
	"context"
	"time"

	firestore "cloud.google.com/go/firestore"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

const userCollection = "users"

// Firestoreからユーザー情報を取得
func GetUserByID(uid string) (*model.User, error) {
	ctx := context.Background()
	dsnap, err := config.FirestoreClient.Collection(userCollection).Doc(uid).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := dsnap.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Firestoreでユーザーのタグを更新
func UpdateUserTags(uid string, tags []string) error {
	ctx := context.Background()
	_, err := config.FirestoreClient.Collection(userCollection).Doc(uid).Set(ctx, map[string]interface{}{
		"tags":       tags,
		"updated_at": time.Now(),
	}, firestore.MergeAll)
	return err
}
