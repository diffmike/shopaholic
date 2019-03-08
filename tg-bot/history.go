package tg_bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
)

type HistoryCommand struct {
	Options
}

func (c *HistoryCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] bot is preparing history")

	user, err := c.Store.Details(strconv.Itoa(m.Sender.ID))
	if err != nil {
		return err
	}

	transactions, err := c.Store.List(user)
	if err != nil {
		return err
	}

	results := []string{}
	for _, transaction := range transactions {
		amount := float64(transaction.Amount.Amount / 100)
		balanceWas := float64(transaction.BalanceWas.Amount / 100)
		balanceNow := float64(transaction.BalanceNow.Amount / 100)
		time := transaction.CreatedAt.Format("02.01.2006 15:04:05")
		result := fmt.Sprintf("Balance: was %.2f, now %.2f. Created: %s. Amount %.2f. Category %s",
			balanceWas, balanceNow, time, amount, transaction.Category.Title)
		results = append(results, result)
	}

	_, err = c.Bot.Send(m.Sender, strings.Join(results, "\n"))
	return err
}
