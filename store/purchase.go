package store

import (
	"shopaholic/utils"
	"time"
)

type TransactionType int

const (
	Income  TransactionType = 0
	Expense TransactionType = 1
)

type Transaction struct {
	ID         string          `json:"id" bson:"_id"`
	User       User            `json:"user"`
	Amount     utils.Money     `json:"amount"`
	CreatedAt  time.Time       `json:"time" bson:"time"`
	Type       TransactionType `json:"type"`
	BalanceWas utils.Money     `json:"balance_was"`
	BalanceNow utils.Money     `json:"balance_now"`
}
