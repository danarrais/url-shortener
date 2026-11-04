package models

import (
	"time"
)

type URL struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShortCode string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"short_code"`
	LongURL   string    `gorm:"type:text;not null" json:"long_url"`
	Clicks    int       `gorm:"type:int;default:0" json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
