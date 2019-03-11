package store

import (
	"time"
)

type Category struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"time" bson:"time"`
	Type      Type      `json:"type"`
	IsDefault bool      `json:"is_default"`
	UserID    string    `json:"user_id"`
}
