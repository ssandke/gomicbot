package main

import "github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/tucnak/telebot"

type SayHiMessageConsumer struct {
	bot   *telebot.Bot
	store StateStore
}

func (c *SayHiMessageConsumer) Initialize(config *Configuration, store StateStore, bot *telebot.Bot) error {
	c.bot = bot
	c.store = store

	return nil
}

func (c *SayHiMessageConsumer) ConsumeMessage(message telebot.Message) (consumed bool, err error) {

	err = nil
	consumed = false

	if message.Text == "/hi" {
		c.bot.SendMessage(message.Chat,
			"Hello, "+message.Sender.FirstName+"!", &telebot.SendOptions{ReplyTo: message})
		consumed = true
	}

	return
}
