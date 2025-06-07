package service

import (
	"context"
	"fmt"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/fetcher"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// ArticleRepository は記事データへのアクセスインターフェースです。
type ArticleRepository interface {
	SaveArticles(ctx context.Context, articles []model.Article) error
	GetArticlesByTag(ctx context.Context, tag string) ([]model.Article, error)
	// 必要に応じて他のメソッドを追加
}

// ArticleService は記事関連のビジネスロジックを扱います。
type ArticleService struct {
	repo ArticleRepository
}

// NewArticleService はArticleServiceの新しいインスタンスを作成します。
func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

// GetPopularArticles は人気記事を取得します。
// 必要に応じて、定期実行処理からこの関数を呼び出し、
// Fetcherで最新記事を取得してRepositoryで保存・更新する処理を実装します。
// 現在はキャッシュからの取得のみを行います。
func (s *ArticleService) GetPopularArticles(ctx context.Context, tag string) ([]model.Article, error) {
	// TODO: 定期実行処理でFetcherを呼び出し、Repository.SaveArticlesを呼び出す

	// 現在はキャッシュから記事を取得して返すのみ
	articles, err := s.repo.GetArticlesByTag(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles from repository: %w", err)
	}

	// TODO: 必要に応じて、キャッシュが古い場合はFetcherを呼び出して更新するロジックを追加

	return articles, nil
}

// FetchAndSaveArticles はQiitaとZennから記事を取得し、リポジトリに保存します。
// この関数はバッチ処理や定期実行される関数から呼び出されることを想定しています。
func (s *ArticleService) FetchAndSaveArticles(ctx context.Context, tags []string) error {
	var allArticles []model.Article

	// 各タグごとに記事を取得
	for _, tag := range tags {
		qiitaArticles, err := fetcher.FetchQiitaArticles(tag)
		if err != nil {
			// エラーをログに出力するなどして、処理を続行
			fmt.Printf("Error fetching Qiita articles for tag %s: %v\n", tag, err)
		}
		allArticles = append(allArticles, qiitaArticles...)

		zennArticles, err := fetcher.FetchZennArticles(tag)
		if err != nil {
			// エラーをログに出力するなどして、処理を続行
			fmt.Printf("Error fetching Zenn articles for tag %s: %v\n", tag, err)
		}
		allArticles = append(allArticles, zennArticles...)
	}

	// TODO: 取得した記事の重複排除やソーティングを行う

	// リポジトリに保存
	err := s.repo.SaveArticles(ctx, allArticles)
	if err != nil {
		return fmt.Errorf("failed to save articles to repository: %w", err)
	}

	fmt.Printf("Successfully fetched and saved %d articles.\n", len(allArticles))

	return nil
}
