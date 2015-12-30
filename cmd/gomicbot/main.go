package main

import (
	"log"
	"time"

	"github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const BotName = "@gomicbot"

func main() {

	config, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	log.Printf("Configuration Loaded:\n%s", config)

	stateStore, err := makeStateStore(config)
	if err != nil {
		panic(err)
	}

	bot, err := telebot.NewBot(config.token)
	if err != nil {
		panic(err)
	}

	consumers, err := loadConsumers()
	if err != nil {
		panic(err)
	}

	err = initConsumers(consumers, config, stateStore, bot)
	if err != nil {
		panic(err)
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		for _, consumer := range consumers {
			consumed, err := consumer.ConsumeMessage(message)

			if err != nil {
				// TODO: what if there's an error? I say we ignore it...should log it
			}

			if consumed {
				break
			}
		}
	}
}

func loadConsumers() (consumers []MessageConsumer, err error) {
	err = nil
	consumers = []MessageConsumer{new(CommenterMessageConsumer), new(SayHiMessageConsumer)}

	return
}

func initConsumers(consumers []MessageConsumer, config *Configuration, store StateStore, bot *telebot.Bot) error {
	for _, consumer := range consumers {
		err := consumer.Initialize(config, store, bot)
		if err != nil {
			// TODO: build a compostion of all the errors to return rather than just the first
			return err
		}
	}
	return nil
}

func makeStateStore(config *Configuration) (store StateStore, err error) {
	store = new(InMemoryStateStore)
	err = store.Initialize(config)
	store.StoreSaying("MIC™")
	store.StoreSaying("Tiny bubbles!")
	return
}
