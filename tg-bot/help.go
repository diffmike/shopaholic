package tg_bot

import tb "gopkg.in/tucnak/telebot.v2"

type HelpCommand struct {
	Options
}

func (c *HelpCommand) Execute(m *tb.Message) error {
	helpText := "Thanks for asking!\nAvailable commands:\n" +
		"/history - showing latest transactions\n" +
		"/info - show your current details: balance, transactions amount, creation date, etc"
	_, err := c.Bot.Send(m.Sender, helpText)

	return err
}
