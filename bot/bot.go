package bot

import (
	"encoding/json"
	"log"

	client "github.com/NOX73/go-tanks-client"
)

type Bot struct {
	strategy Strategy
	client   client.Client
	addr     string
	worldCh  chan client.Message

	login string
	pass  string

	ShowSendingCommand bool
	WorldFrequency     int64
	Reconnect          bool

	SelfTank client.Tank
}

func NewBot(login, pass, addr string, strategy Strategy) *Bot {
	return &Bot{
		addr:           addr,
		worldCh:        make(chan client.Message),
		strategy:       strategy,
		login:          login,
		pass:           pass,
		WorldFrequency: 10,
		Reconnect:      true,
	}
}

func (b *Bot) Connect() error {
	var err error

	if b.client != nil {
		return nil
	}

	b.client, err = client.ConnectTo(b.addr)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) Go() error {
	var err error

	err = b.Connect()
	if err != nil {
		return err
	}

	err = b.auth()
	if err != nil {
		return err
	}

	go b.runReadMessages()

	b.loop()

	return nil
}

func (b *Bot) auth() error {
	err := b.client.Auth(b.login, b.pass)
	if err != nil {
		return err
	}

	err = b.readSelfTank()

	if err != nil {
		return err
	}

	b.setup()

	return nil
}

func (b *Bot) readSelfTank() error {

	for {

		message, err := b.client.ReadMessage()
		if err != nil {
			return err
		}

		if message.IsTank() {
			b.SelfTank = *message.Tank
			break
		}
	}

	return nil
}

func (b *Bot) setup() {
	b.client.WorldFrequency(b.WorldFrequency)
}

func (b *Bot) runReadMessages() {
loop:
	for {
		message, err := b.client.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		switch message.Type {
		case "World":

			select {
			case b.worldCh <- *message:
				// ok
			default:
				// drop message
			}

		case "Hit":
			if b.Reconnect {
				b.auth()
			} else {
				close(b.worldCh)
				break loop
			}

		default:
			log.Println("Tank #", b.SelfTank.Id, "Message received:", message.Type, "/", message.Message)
		}

	}

	log.Println("Tank #", b.SelfTank.Id, "Finish read messages.")
}

func (b *Bot) loop() {
loop:
	for {
		message, ok := <-b.worldCh

		if !ok {
			break loop
		}

		command := b.strategy.Perform(message, b.SelfTank)

		if b.ShowSendingCommand {
			jsonStr, _ := json.Marshal(command)
			log.Println("Tank #", b.SelfTank.Id, "Sending command:", string(jsonStr))
		}

		b.client.SendTankCommand(command)
	}

	log.Println("Tank #", b.SelfTank.Id, "Finish command loop.")
}
