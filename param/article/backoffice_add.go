package articleparam

import "time"

type AddArticleRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Body            string   `json:"body"`
	IsPublished     bool     `json:"is_published"`
	Tags            []string `json:"tags"`
	MetaDescription string   `json:"meta_description"`
	MediaID         string   `json:"media_id"`
}

type ArticleInfo struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Body            string    `json:"body"`
	IsPublished     bool      `json:"is_published"`
	Tags            []string  `json:"tags"`
	MetaDescription string    `json:"meta_description"`
	BackOfficeID    string    `json:"back_office_id"`
	BackOfficeName  string    `json:"back_office_name"`
	MediaUrl        string    `json:"media_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AddArticleResponse struct {
	Article ArticleInfo `json:"article"`
}
