package bot

import (
	client "github.com/NOX73/go-tanks-client"
)

type RingStrategy struct {
}

func NewRingStrategy() *RingStrategy {
	return &RingStrategy{}
}

func (s *RingStrategy) Perform(world client.Message, tank client.Tank) (command client.Message) {
	command = client.NewMessage().Motors(1.0, 0.2)
	return command
}
