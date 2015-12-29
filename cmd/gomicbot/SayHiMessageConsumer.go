package main

import (
	"log"

	"github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type SayHiMessageConsumer struct {
	bot    *telebot.Bot
	config *Configuration
	store  StateStore
}

func (c *SayHiMessageConsumer) Initialize(config *Configuration, store StateStore, bot *telebot.Bot) error {
	c.bot = bot
	c.config = config
	c.store = store

	return nil
}

func (c *SayHiMessageConsumer) ConsumeMessage(message telebot.Message) (consumed bool, err error) {

	err = nil
	consumed = false

	log.Printf("SayHi got message '%s' Pers: %v", message.Text, message.IsPersonal())

	matched, _ := checkCommandMatch(message.Text, "/hi", message.IsPersonal(), false, c.config)
	if matched {
		c.bot.SendMessage(message.Chat,
			"Hello, "+message.Sender.FirstName+"!", &telebot.SendOptions{ReplyTo: message})
		consumed = true
	}

	return
}
