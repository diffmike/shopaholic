package tg_bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
)

type InfoCommand struct {
	Options
}

func (c *InfoCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] user who ask for info %+v", m.Sender)

	user, err := c.Store.Details(strconv.Itoa(m.Sender.ID))
	if err != nil {
		return err
	}

	balance := float64(user.Balance.Amount / 100)
	time := user.CreatedAt.Format("02.01.2006 15:04")
	result := fmt.Sprintf("%s 🤓 Balance is %.2f💲 Registion 🕰 %s", user.Name, balance, time)
	_, err = c.Bot.Send(m.Sender, result)

	return err
}
