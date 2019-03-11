package cmd

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"shopaholic/tg-bot"
	"time"
)

type BotPollerCommand struct {
	CommonOpts
}

func (bpc *BotPollerCommand) Execute(args []string) error {
	log.Printf("[INFO] bot poller command is started")

	b, err := tb.NewBot(tb.Settings{
		Token:  bpc.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	opts := tg_bot.Options{
		Bot:      b,
		Currency: bpc.Currency,
		Store:    bpc.Store,
	}
	b.Handle("/help", func(m *tb.Message) {
		_ = processCommand(&tg_bot.HelpCommand{}, opts, m)
	})
	b.Handle("/start", func(m *tb.Message) {
		_ = processCommand(&tg_bot.StartCommand{}, opts, m)
	})
	b.Handle("/history", func(m *tb.Message) {
		_ = processCommand(&tg_bot.HistoryCommand{}, opts, m)
	})
	b.Handle("/info", func(m *tb.Message) {
		_ = processCommand(&tg_bot.InfoCommand{}, opts, m)
	})
	b.Handle("/category", func(m *tb.Message) {
		_ = processCommand(&tg_bot.CategoryCommand{}, opts, m)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		_ = processCommand(&tg_bot.TransactionCommand{}, opts, m)
	})

	b.Start()

	return nil
}

func processCommand(botCmd tg_bot.Commander, opts tg_bot.Options, m *tb.Message) (err error) {
	botCmd.SetOptions(opts)
	if err := botCmd.Execute(m); err != nil {
		log.Printf("[ERROR] error is '%s' in processing command of text '%s'", err.Error(), m.Text)
	}

	return err
}
