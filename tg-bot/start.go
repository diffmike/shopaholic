package tg_bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"shopaholic/store"
	"shopaholic/utils"
	"strconv"
)

type StartCommand struct {
	Options
}

func (c *StartCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] user who start chat is %+v", m.Sender)

	user := store.User{
		Name:    m.Sender.Username,
		ID:      strconv.Itoa(m.Sender.ID),
		Balance: utils.Money{0, c.Currency},
	}

	userID, err := c.Store.Register(user)
	if err != nil {
		return err
	}

	log.Printf("[INFO] user %s was created with ID: %s", user.Name, userID)

	helloText := "Hi man! Let's try to add a bit of structure to your finances 🧐"
	_, err = c.Bot.Send(m.Sender, helloText)

	return err
}
