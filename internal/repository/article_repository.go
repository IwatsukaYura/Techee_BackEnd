package repository

import (
	"context"

	firestore "cloud.google.com/go/firestore"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/config"
	"github.com/iwatsukayugaku/my-tech-articles-app/backend/internal/model"
)

const articleCollection = "articles"

// Firestoreに記事をキャッシュ保存
func SaveArticles(articles []model.Article) error {
	ctx := context.Background()
	batch := config.FirestoreClient.Batch()
	for _, a := range articles {
		ref := config.FirestoreClient.Collection(articleCollection).Doc(a.ID)
		batch.Set(ref, a, firestore.MergeAll)
	}
	_, err := batch.Commit(ctx)
	return err
}

// Firestoreから記事キャッシュを取得（タグでフィルタ）
func GetCachedArticles(tag string, limit int) ([]model.Article, error) {
	ctx := context.Background()
	q := config.FirestoreClient.Collection(articleCollection).OrderBy("likes", firestore.Desc).Limit(limit)
	if tag != "" {
		q = q.Where("tags", "array-contains", tag)
	}
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var articles []model.Article
	for _, doc := range docs {
		var a model.Article
		if err := doc.DataTo(&a); err == nil {
			articles = append(articles, a)
		}
	}
	return articles, nil
}
