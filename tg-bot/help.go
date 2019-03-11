package tg_bot

import tb "gopkg.in/tucnak/telebot.v2"

type HelpCommand struct {
	Options
}

func (c *HelpCommand) Execute(m *tb.Message) error {
	helpText := "Thanks for asking!\n❔Available commands:\n" +
		"/history - showing latest transactions 📜\n" +
		"/category {title} - create new category 🧳\n" +
		"/info - show your details 🤓"
	_, err := c.Bot.Send(m.Sender, helpText)

	return err
}
