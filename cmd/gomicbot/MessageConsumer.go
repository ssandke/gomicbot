package main

import (
	"github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type MessageConsumer interface {
	Initialize(*Configuration, StateStore, *telebot.Bot) error
	ConsumeMessage(telebot.Message) (consumed bool, err error)
}
