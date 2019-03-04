package store

import (
	"shopaholic/utils"
	"time"
)

type User struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"time" bson:"time"`
	Balance   utils.Money `json:"balance"`
}
