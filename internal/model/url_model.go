package model

import "time"

type URLModel struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	CreatedAt time.Time `json:"created_at"`
	ClickCount int `json:"click_count"`
}