package tg_bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"shopaholic/store"
	"shopaholic/utils"
	"strconv"
	"strings"
)

type TransactionCommand struct {
	Options
}

func (c *TransactionCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] is storing transaction about '%s'", m.Text)

	user, err := c.Store.Details(strconv.Itoa(m.Sender.ID))
	if err != nil {
		return err
	}

	amount, cat, err := parseTransactionMessage(m.Text)
	if err != nil {
		return err
	}

	category, err := c.Store.FindCategoryByTitle(cat)
	if err != nil {
		return err
	}

	transaction := store.Transaction{
		User:     user,
		Category: category,
		Amount:   utils.Money{amount, c.Currency},
	}

	transactionID, err := c.Store.StoreTransaction(transaction)
	if err != nil {
		return err
	}

	result := fmt.Sprintf("Transaction for %.2f$ was created in the category %s.\nCurrent balance is %.2f$",
		float64(amount/100), category.Title, float64(transaction.BalanceNow.Amount/100))
	log.Printf("[INFO] %s. ID: %s", result, transactionID)
	_, err = c.Bot.Send(m.Sender, result)
	return err
}

func parseTransactionMessage(message string) (amount int64, category string, err error) {
	s := strings.Split(message, " ")
	log.Printf("[INFO] amount: %+v", s)
	if len(s) < 2 {
		s = append(s, "")
	}
	a, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0, "", err
	}
	amount, category = int64(a*100), s[1]

	return amount, category, err
}
