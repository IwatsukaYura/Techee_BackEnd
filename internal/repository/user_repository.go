package repository

import (
	"context"
	"fmt"
	"time"

	firestore "cloud.google.com/go/firestore"
	// configは直接Clientを受け取るため不要
	// "github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const userCollection = "users"

// UserRepository はユーザーデータへのアクセスを抽象化するインターフェースです。
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*model.User, error)
	UpdateUserTags(ctx context.Context, userID string, tags []string) error
	// 必要に応じて他のメソッドを追加
}

// firestoreUserRepository はFirestoreをデータストアとして使用するUserRepositoryの実装です。
type firestoreUserRepository struct {
	client *firestore.Client
}

// NewUserRepository はfirestoreUserRepositoryの新しいインスタンスを作成します。
func NewUserRepository(client *firestore.Client) UserRepository {
	return &firestoreUserRepository{client: client}
}

// Firestoreからユーザー情報を取得
func (r *firestoreUserRepository) GetUser(ctx context.Context, userID string) (*model.User, error) {
	dsnap, err := r.client.Collection(userCollection).Doc(userID).Get(ctx)
	if err != nil {
		// ドキュメントが存在しない場合はnilとnilエラーを返す（または専用エラーを返す）
		if status.Code(err) == codes.NotFound {
			return nil, nil // ユーザーが存在しない場合はエラーではなくnilユーザーを返す
		}
		return nil, fmt.Errorf("failed to get user from firestore: %w", err)
	}
	var user model.User
	if err := dsnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to map firestore data to user model: %w", err)
	}
	return &user, nil
}

// Firestoreでユーザーのタグを更新
func (r *firestoreUserRepository) UpdateUserTags(ctx context.Context, userID string, tags []string) error {
	_, err := r.client.Collection(userCollection).Doc(userID).Set(ctx, map[string]interface{}{
		"tags":       tags,
		"updated_at": time.Now().Format(time.RFC3339), // stringにフォーマット
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update user tags in firestore: %w", err)
	}
	return nil
}
