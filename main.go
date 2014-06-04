package main

import (
	"flag"
	"log"
	"sync"

	"./bot"
)

var count = flag.Int("bots-count", 1, "Count to running bots.")
var login = flag.String("login", "login", "Login")
var pass = flag.String("pass", "pass", "Password")
var host = flag.String("host", "nox73.ru:9292", "Host and port to connect. Ex: nox73.ru:9292")

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	for i := 0; i < *count; i++ {

		wg.Add(1)
		go func() {

			strategy := bot.NewRingStrategy()
			bot := bot.NewBot(*login, *pass, *host, strategy)
			err := bot.Go()

			if err != nil {
				log.Panic(err)
			}

			wg.Done()

			log.Println("Bot finished")
		}()
	}

	wg.Wait()

	log.Println("All bots finished")
}
