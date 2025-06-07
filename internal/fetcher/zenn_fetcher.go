package fetcher

import (
	"fmt"
	"time"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// Zennから記事を取得する関数 (今回は簡易的に固定データを返す)
func FetchZennArticles(tag string) ([]model.Article, error) {
	// 実際にはZennのRSSフィードなどをパースして記事を取得
	fmt.Printf("Fetching Zenn articles for tag: %s (using fixed data)\n", tag)

	// 固定のサンプルデータを返す
	articles := []model.Article{
		{
			ID:          "zenn-sample-1",
			Title:       "Zennのサンプル記事1",
			URL:         "https://zenn.dev/sample/1",
			Tags:        []string{"Go", "Zenn"},
			Likes:       100,
			PublishedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339), // stringにフォーマット
			Source:      "Zenn",
		},
		{
			ID:          "zenn-sample-2",
			Title:       "Zennのサンプル記事2 (Python)",
			URL:         "https://zenn.dev/sample/2",
			Tags:        []string{"Python", "Zenn"},
			Likes:       50,
			PublishedAt: time.Now().Add(-48 * time.Hour).Format(time.RFC3339), // stringにフォーマット
			Source:      "Zenn",
		},
	}

	// タグでフィルタリング（固定データなので簡易フィルタリング）
	if tag != "" {
		filteredArticles := []model.Article{}
		for _, article := range articles {
			for _, articleTag := range article.Tags {
				if articleTag == tag {
					filteredArticles = append(filteredArticles, article)
					break
				}
			}
		}
		articles = filteredArticles
	}

	return articles, nil
}
