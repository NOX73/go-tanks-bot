package bot

import (
	client "github.com/NOX73/go-tanks-client"
)

type Strategy interface {
	Perform(world client.Message) (command client.Message)
}
