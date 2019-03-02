package store

import "time"

type Purchase struct {
	ID        string    `json:"id" bson:"_id"`
	User      User      `json:"user"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"time" bson:"time"`
}
