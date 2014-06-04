package bot

import (
	"math"

	client "github.com/NOX73/go-tanks-client"
)

const (
	TURN_SPEED_DELIMETER        = 90
	MOVE_SPEED_DELIMETER        = 200
	ANGLE_DIFF                  = 2
	DISTANCE_DIFF               = 20
	TURN_AND_MOVE_ANGLE_DIFF    = 50
	TURN_AND_MOVE_DISTANCE_DIFF = 50
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
		x = DISTANCE_DIFF
		y = DISTANCE_DIFF
	case 1:
		//rigth top
		x = float64(s.world.Map.Width) - DISTANCE_DIFF
		y = DISTANCE_DIFF
	case 2:
		//right bottom
		x = float64(s.world.Map.Width) - DISTANCE_DIFF
		y = float64(s.world.Map.Height) - DISTANCE_DIFF
	case 3:
		//left bottom
		x = DISTANCE_DIFF
		y = float64(s.world.Map.Height) - DISTANCE_DIFF

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
		if math.Abs(diff) < TURN_AND_MOVE_ANGLE_DIFF && dist > TURN_AND_MOVE_DISTANCE_DIFF {
			speed := dist / MOVE_SPEED_DELIMETER
			s.TurnToAndForward(alpha, speed)
		} else {
			s.TurnTo(alpha)
		}
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

	if diff < -180 {
		diff = diff + 360
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

func (s *RingStrategy) TurnToAndForward(angle, moveSpeed float64) {
	diff := degreeDiff(angle, s.tank.Direction)
	turnSpeed := math.Abs(diff / TURN_SPEED_DELIMETER)

	if moveSpeed > 1 {
		moveSpeed = 1
	}

	if turnSpeed > 1 {
		turnSpeed = 1
	}

	left := moveSpeed
	right := moveSpeed

	if diff < 0 {
		left = left - turnSpeed
		right = right + turnSpeed
	} else {
		left = left + turnSpeed
		right = right - turnSpeed
	}

	if left > 1 {
		left = 1
	}

	if right > 1 {
		right = 1
	}

	s.command = s.command.Motors(left, right)
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
