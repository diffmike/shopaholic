package store

import (
	"shopaholic/utils"
	"time"
)

type Type int

const (
	Income  Type = 0
	Expense Type = 1
)

type Transaction struct {
	ID         string      `json:"id" bson:"_id"`
	User       User        `json:"user"`
	Amount     utils.Money `json:"amount"`
	CreatedAt  time.Time   `json:"time" bson:"time"`
	Type       Type        `json:"type"`
	BalanceWas utils.Money `json:"balance_was"`
	BalanceNow utils.Money `json:"balance_now"`
	Category   Category    `json:"category"`
}
