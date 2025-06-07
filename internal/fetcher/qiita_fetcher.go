package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

// Qiita API レスポンスの構造体を定義
type qiitaArticle struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	LikesCount int    `json:"likes_count"`
	CreatedAt  string `json:"created_at"` // ISO 8601 format
	Tags       []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

// Qiita APIから記事を取得する関数
func FetchQiitaArticles(tag string) ([]model.Article, error) {
	baseURL := "https://qiita.com/api/v2/items"

	// クエリパラメータを設定
	params := url.Values{}
	params.Add("sort", "likes") // いいね数でソート
	if tag != "" {
		params.Add("query", fmt.Sprintf("tag:%s", tag)) // タグで絞り込み
	}

	// リクエストURLを構築
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	reqURL.RawQuery = params.Encode()

	// HTTP GET リクエストを実行
	res, err := http.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Qiita articles: %w", err)
	}
	defer res.Body.Close()

	// ステータスコードを確認
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Qiita API returned non-200 status: %d, body: %s", res.StatusCode, bodyBytes)
	}

	// レスポンスボディを読み込み
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Qiita API response body: %w", err)
	}

	// JSONをパース
	var qiitaArticles []qiitaArticle
	err = json.Unmarshal(body, &qiitaArticles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Qiita API response: %w", err)
	}

	// 内部モデルにマッピング
	articles := make([]model.Article, 0, len(qiitaArticles))
	for _, qa := range qiitaArticles {
		var tags []string
		for _, t := range qa.Tags {
			tags = append(tags, t.Name)
		}

		publishedAt, err := time.Parse(time.RFC3339, qa.CreatedAt)
		if err != nil {
			// パースエラーの場合はスキップまたはエラー処理を検討
			fmt.Printf("Warning: failed to parse Qiita article created_at '%s': %v\n", qa.CreatedAt, err)
			// エラーで中断せず、次の記事に進む
			// publishedAtがパースできなかった場合はゼロ値または適切なデフォルト値を使用
			publishedAt = time.Time{}
		}

		articles = append(articles, model.Article{
			ID:          qa.ID,
			Title:       qa.Title,
			URL:         qa.URL,
			Tags:        tags,
			Likes:       qa.LikesCount,
			PublishedAt: publishedAt.Format(time.RFC3339), // time.Timeをstringにフォーマット（二重定義を解消）
			Source:      "Qiita",
			// FetchedAtは記事取得時に設定（後で実装）
		})
	}

	return articles, nil
}
