package store

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"time" bson:"time"`
	Balance   int       `json:"balance"`
}
