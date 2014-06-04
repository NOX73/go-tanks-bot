package bot

import (
	client "github.com/NOX73/go-tanks-client"
)

type Strategy interface {
	Perform(world client.Message, tank client.Tank) (command client.Message)
}
