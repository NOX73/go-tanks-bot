package main

import (
	"log"

	"./bot"
)

func main() {
	strategy := bot.NewRingStrategy()
	bot := bot.NewBot("login", "pass", "nox73.ru:9292", strategy)

	bot.ShowSendingCommand = true

	err := bot.Go()

	if err != nil {
		log.Panic(err)
	}

}
