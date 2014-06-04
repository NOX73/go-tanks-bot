package bot

import (
	"math"

	client "github.com/NOX73/go-tanks-client"
)

const (
	TURN_SPEED_DELIMETER = 90
	MOVE_SPEED_DELIMETER = 200
	ANGLE_DIFF           = 5
	DISTANCE_DIFF        = 100
)

type RingStrategy struct {
	world   client.Message
	tank    client.Tank
	command client.Message

	//Move
	direction int64
}

func NewRingStrategy() *RingStrategy {
	return &RingStrategy{direction: 0}
}

func (s *RingStrategy) Perform(world client.Message, tank client.Tank) (command client.Message) {
	s.command = client.NewMessage()
	s.world = world
	s.tank = tank

	s.PerformMove()

	return s.command
}

func (s *RingStrategy) PerformMove() {
	var x, y float64

	switch s.direction {
	case 0:
		//left top
		x = 0
		y = 0
	case 1:
		//rigth top
		x = float64(s.world.Map.Width)
		y = 0
	case 2:
		//right bottom
		x = float64(s.world.Map.Width)
		y = float64(s.world.Map.Height)
	case 3:
		//left bottom
		x = 0
		y = float64(s.world.Map.Height)

	}

	s.GoTo(x, y)
}

func (s *RingStrategy) GoTo(x, y float64) {
	tankX := s.tank.Coords.X
	tankY := s.tank.Coords.Y

	alpha := getAlpha(x, y, tankX, tankY)
	diff := degreeDiff(alpha, s.tank.Direction)

	dist := distanceBetween(x, y, tankX, tankY)

	if dist < DISTANCE_DIFF {
		s.nextDirection()
	}

	if math.Abs(diff) < ANGLE_DIFF {
		speed := dist / MOVE_SPEED_DELIMETER
		s.Forward(speed)
	} else {
		s.TurnTo(alpha)
	}
}

func distanceBetween(x, y, x2, y2 float64) float64 {
	a := x - x2
	b := y - y2

	return math.Hypot(a, b)
}

func degreeDiff(a, b float64) float64 {
	diff := a - b

	if diff > 180 {
		diff = diff - 360
	}

	return diff
}

func getAlpha(x, y, x2, y2 float64) float64 {

	a := x - x2
	b := y - y2

	if a < 0 && b < 0 {
		return 270 - (math.Atan(a/b) * 180 / math.Pi)
	}

	if a > 0 && b < 0 {
		return 270 - (math.Atan(a/b) * 180 / math.Pi)
	}

	if a > 0 && b > 0 {
		return 90 - (math.Atan(a/b) * 180 / math.Pi)
	}

	if a < 0 && b > 0 {
		return 90 - (math.Atan(a/b) * 180 / math.Pi)
	}

	return 0
}

func (s *RingStrategy) Forward(speed float64) {
	if speed > 1 {
		speed = 1
	}

	s.command = s.command.Motors(speed, speed)
}

func (s *RingStrategy) TurnTo(angle float64) {
	diff := degreeDiff(angle, s.tank.Direction)

	speed := math.Abs(diff / TURN_SPEED_DELIMETER)

	if diff < 0 {
		s.TurnLeft(speed)
	} else {
		s.TurnRight(speed)
	}
}

func (s *RingStrategy) TurnLeft(speed float64) {
	if speed > 1 {
		speed = 1
	}
	s.command = s.command.Motors(-speed, speed)
}

func (s *RingStrategy) TurnRight(speed float64) {
	if speed > 1 {
		speed = 1
	}
	s.command = s.command.Motors(speed, -speed)
}

func (s *RingStrategy) nextDirection() {
	s.direction++
	if s.direction > 3 {
		s.direction = 0
	}
}
