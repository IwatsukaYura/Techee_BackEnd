package service

import (
	"context"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// UserRepository はユーザーデータへのアクセスインターフェースです。
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*model.User, error)
	UpdateUserTags(ctx context.Context, userID string, tags []string) error
}

// UserService はユーザー関連のビジネスロジックを扱います。
type UserService struct {
	repo UserRepository
}

// NewUserService はUserServiceの新しいインスタンスを作成します。
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUser はユーザー情報を取得します。
func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	// リポジトリからユーザー情報を取得するロジックを実装
	return s.repo.GetUser(ctx, userID)
}

// UpdateUserTags はユーザーのタグを更新します。
func (s *UserService) UpdateUserTags(ctx context.Context, userID string, tags []string) error {
	// リポジトリを使ってユーザーのタグを更新するロジックを実装
	return s.repo.UpdateUserTags(ctx, userID, tags)
}
