package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type CommenterMessageConsumer struct {
	bot    *telebot.Bot
	config *Configuration
	store  StateStore
	rng    *rand.Rand
}

func (c *CommenterMessageConsumer) Initialize(config *Configuration, store StateStore, bot *telebot.Bot) error {
	c.bot = bot
	c.config = config
	c.store = store

	c.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	return nil
}

func (c *CommenterMessageConsumer) ConsumeMessage(message telebot.Message) (consumed bool, err error) {

	err = nil
	consumed = true

	writer := c.config.isWriter(message.Sender.ID)
	reader := c.config.isReader(message.Sender.ID)

	if writer && c.commandMatch(message, "/list") {
		c.ListSayings(message)
	} else if args := c.commandMatchWithArgs(message, "/add"); writer && args != "" {
		c.AddSaying(args, message)
	} else if args := c.commandMatchWithArgs(message, "/remove"); writer && args != "" {
		c.RemoveSaying(args, message)
	} else if reader && strings.Contains(message.Text, BotName) {
		c.EmitSaying(message, true)
	} else {
		roll := c.rng.Float64()
		if reader && roll < c.config.sayingChance {
			c.EmitSaying(message, false)
			consumed = true
		}

		consumed = false
	}

	return
}

func (c *CommenterMessageConsumer) EmitSaying(message telebot.Message, reply bool) {
	sayings, _ := c.store.LoadSayings()

	if len(sayings) > 0 {
		var options *telebot.SendOptions
		if reply {
			options = &telebot.SendOptions{ReplyTo: message}
		}

		c.bot.SendMessage(message.Chat, sayings[c.rng.Intn(len(sayings))], options)
	}
}

func (c *CommenterMessageConsumer) ListSayings(message telebot.Message) {
	sayings, err := c.store.LoadSayings()

	if err != nil {
		c.bot.SendMessage(message.Chat, fmt.Sprintf("Error loading sayings: %s", err), &telebot.SendOptions{ReplyTo: message})
		return
	}

	sort.Strings(sayings)

	buffer := bytes.NewBufferString(fmt.Sprintf("%d Sayings:\n", len(sayings)))

	for _, saying := range sayings {
		buffer.WriteString(saying)
		buffer.WriteString("\n")

		// Telegram doesn't like messages >4k it seems, so if we're getting close, send what we have
		// and start building up another buffer to send.
		if buffer.Len() > 4000 {
			c.bot.SendMessage(message.Chat, buffer.String(), nil)
			buffer.Reset()
		}
	}

	if buffer.Len() > 0 {
		c.bot.SendMessage(message.Chat, buffer.String(), nil)
	}
}

func (c *CommenterMessageConsumer) AddSaying(saying string, message telebot.Message) {
	c.store.StoreSaying(saying)
	c.bot.SendMessage(message.Chat, "Added", &telebot.SendOptions{ReplyTo: message})
}

func (c *CommenterMessageConsumer) RemoveSaying(saying string, message telebot.Message) {
	present, err := c.store.RemoveSaying(saying)
	var result string
	switch {
	case err != nil:
		result = fmt.Sprintf("Error removing saying: %s", err)
	case !present:
		result = "No such saying in my database."
	default:
		result = "Removed."
	}

	c.bot.SendMessage(message.Chat, result, &telebot.SendOptions{ReplyTo: message})
}

func (c *CommenterMessageConsumer) commandMatch(message telebot.Message, command string) (matched bool) {
	matched, _ = checkCommandMatch(message.Text, command, message.IsPersonal(), false, c.config)
	return
}

func (c *CommenterMessageConsumer) commandMatchWithArgs(message telebot.Message, command string) (args string) {
	_, args = checkCommandMatch(message.Text, command, message.IsPersonal(), true, c.config)
	return
}
