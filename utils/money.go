package utils

import "github.com/rhymond/go-money"

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

func ToMoney(m Money) *money.Money {
	return money.New(m.Amount, m.Currency)
}
