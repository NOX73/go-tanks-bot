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

	//GUN

	FIRE_MIN_X float64 = 150 // z
	FIRE_MIN_Y float64 = 1   // n
	FIRE_MAX_X float64 = 800 // l
	FIRE_MAX_Y float64 = 20  // m

	FIRE_DISTANCE float64 = 800
)

var FIRE_THRESHOLD_C = FIRE_MAX_Y
var FIRE_THRESHOLD_A = (FIRE_THRESHOLD_C - FIRE_MIN_Y) / math.Pow(FIRE_MIN_X, 2)
var FIRE_THRESHOLD_B = -2 * FIRE_THRESHOLD_A * FIRE_MIN_X

var FIRE_THRESHOLD_C_BIG = (-FIRE_MIN_Y*math.Pow(FIRE_MAX_X, 2) + 2*FIRE_MAX_X*FIRE_MIN_X*FIRE_MIN_Y - FIRE_MAX_Y*math.Pow(FIRE_MIN_X, 2)) /
	(math.Pow(FIRE_MIN_X, 2) - math.Pow(FIRE_MAX_X, 2) + 2*FIRE_MAX_X*FIRE_MIN_X)

var FIRE_THRESHOLD_A_BIG = (FIRE_THRESHOLD_C_BIG - FIRE_MIN_Y) / math.Pow(FIRE_MIN_X, 2)
var FIRE_THRESHOLD_B_BIG = -2 * FIRE_THRESHOLD_A_BIG * FIRE_MIN_X

type RingStrategy struct {
	world   client.Message
	tank    client.Tank
	command client.Message

	//Move
	direction int64

	//Target
	target *client.Tank
}

func NewRingStrategy() *RingStrategy {
	return &RingStrategy{direction: 0}
}

func (s *RingStrategy) Perform(world client.Message, tank client.Tank) (command client.Message) {
	s.command = client.NewMessage()
	s.world = world
	s.tank = tank

	s.PerformMove()
	s.PerformGun()

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

	alpha := getAngle(x, y, tankX, tankY)
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

func (s *RingStrategy) selectTarget() *client.Tank {

	if s.target != nil {
		tank := s.world.GetTankById(s.target.Id)
		if tank != nil && s.isFireDistanceTo(tank) {
			s.target = tank
			return tank
		}
	}

	for _, tank := range s.world.Tanks {
		if tank.Id != s.tank.Id && s.isFireDistanceTo(&tank) {
			s.target = &tank
			return &tank
		}
	}

	return nil
}

func (s *RingStrategy) PerformGun() {
	var target = s.selectTarget()

	if target == nil {
		return
	}

	x := target.Coords.X
	y := target.Coords.Y

	if s.isFireDistanceTo(target) {
		s.FireTo(x, y)
	}
}

func (s *RingStrategy) FireTo(x, y float64) {
	tankX := s.tank.Coords.X
	tankY := s.tank.Coords.Y
	gunAngle := s.tank.Gun.Direction + s.tank.Direction

	if gunAngle > 360 {
		gunAngle = gunAngle - 360
	}

	alpha := getAngle(x, y, tankX, tankY)
	diff := degreeDiff(alpha, gunAngle)

	dist := distanceBetween(x, y, tankX, tankY)

	s.command = s.command.TurnGun(diff)

	var a, b, c float64

	if dist < FIRE_MIN_X {
		a = FIRE_THRESHOLD_A
		b = FIRE_THRESHOLD_B
		c = FIRE_THRESHOLD_C
	} else {
		a = FIRE_THRESHOLD_A_BIG
		b = FIRE_THRESHOLD_B_BIG
		c = FIRE_THRESHOLD_C_BIG
	}

	threshold := a*math.Pow(dist, 2) + b*dist + c

	//log.Println(math.Abs(diff), "\t", threshold, "\t", dist, "\t", math.Abs(diff) < threshold)

	if math.Abs(diff) < threshold {
		s.command = s.command.SetFire()
	}

}

func (s *RingStrategy) isFireDistanceTo(target *client.Tank) bool {
	x := target.Coords.X
	y := target.Coords.Y
	tankX := s.tank.Coords.X
	tankY := s.tank.Coords.Y

	dist := distanceBetween(x, y, tankX, tankY)

	return dist < FIRE_DISTANCE
}
