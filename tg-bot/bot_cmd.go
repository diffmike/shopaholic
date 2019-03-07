package tg_bot

import (
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"shopaholic/store/service"
)

type Commander interface {
	Execute(m *tb.Message) error
	SetOptions(options Options)
}

type Options struct {
	Currency string
	Store    service.DataStore
	Bot      *telebot.Bot
}

func (o *Options) SetOptions(options Options) {
	o.Currency = options.Currency
	o.Store = options.Store
	o.Bot = options.Bot
}
