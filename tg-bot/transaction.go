package tg_bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

type TransactionCommand struct {
	Options
}

func (c *TransactionCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] bot received the text '%s'", m.Text)
	_, err := c.Bot.Send(m.Sender, "answer to your text")

	return err
}
