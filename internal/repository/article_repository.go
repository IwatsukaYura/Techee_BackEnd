package repository

import (
	"context"
	"fmt"

	firestore "cloud.google.com/go/firestore"
	// "github.com/iwatsukayugaku/my-tech-articles-app/backend/config" // 直接Clientを受け取るため不要
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

const articleCollection = "articles"

// ArticleRepository は記事データへのアクセスを抽象化するインターフェースです。
type ArticleRepository interface {
	SaveArticles(ctx context.Context, articles []model.Article) error
	GetArticlesByTag(ctx context.Context, tag string) ([]model.Article, error)
	// 必要に応じて他のメソッドを追加
}

// firestoreArticleRepository はFirestoreをデータストアとして使用するArticleRepositoryの実装です。
type firestoreArticleRepository struct {
	client *firestore.Client
}

// NewArticleRepository はfirestoreArticleRepositoryの新しいインスタンスを作成します。
func NewArticleRepository(client *firestore.Client) ArticleRepository {
	return &firestoreArticleRepository{client: client}
}

// Firestoreに記事をキャッシュ保存
func (r *firestoreArticleRepository) SaveArticles(ctx context.Context, articles []model.Article) error {
	batch := r.client.Batch()
	for _, a := range articles {
		// ドキュメントIDとして記事のIDまたはURLを使用
		docID := a.ID
		if docID == "" {
			docID = a.URL // IDがない場合はURLを使用
		}
		ref := r.client.Collection(articleCollection).Doc(docID)

		// model.Article構造体をmap[string]interface{}に変換してからSetに渡す
		data := map[string]interface{}{
			"id":          a.ID,
			"title":       a.Title,
			"url":         a.URL,
			"tags":        a.Tags,
			"likes":       a.Likes,
			"publishedAt": a.PublishedAt,
			"source":      a.Source,
			// "fetchedAt": time.Now().Format(time.RFC3339), // 記事取得日時を追加する場合はここで設定
		}

		batch.Set(ref, data, firestore.MergeAll)
	}

	// バッチが大きすぎる場合に分割してコミットする処理を追加するべきだが、簡易実装として省略

	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit articles batch: %w", err)
	}

	return nil
}

// Firestoreから記事キャッシュを取得（タグでフィルタ）
func (r *firestoreArticleRepository) GetArticlesByTag(ctx context.Context, tag string) ([]model.Article, error) {
	q := r.client.Collection(articleCollection).OrderBy("likes", firestore.Desc).Limit(50) // 例として最大50件取得
	if tag != "" {
		// タグによる絞り込み。tagsフィールドがstring[]なのでarray-containsを使用
		q = q.Where("tags", "array-contains", tag)
	}
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents from firestore: %w", err)
	}
	var articles []model.Article
	for _, doc := range docs {
		var a model.Article
		// PublishedAtがtime.Time型としてFirestoreに保存されている場合、ここで変換が必要になる可能性
		// model.ArticleのPublishedAtをstring型にしているので、Firestoreへの保存時に適切に変換されている前提
		if err := doc.DataTo(&a); err == nil {
			articles = append(articles, a)
		}
	}
	return articles, nil
}

// TODO: User関連のリポジトリ関数（GetUser, UpdateUserTagsなど）もこのファイルにまとめるか、別途user_repository.goに実装する
// 今回はuser_repository.goに実装することにする
