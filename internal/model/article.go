package model

type Article struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	Likes       int      `json:"likes"`
	PublishedAt string   `json:"publishedAt"`
	Source      string   `json:"source"`
	FetchedAt   string   `json:"fetchedAt"`
}
