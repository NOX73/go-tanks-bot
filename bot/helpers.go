package bot

import "math"

func distanceBetween(xTo, yTo, xFrom, yFrom float64) float64 {
	a := xTo - xFrom
	b := yTo - yFrom

	return math.Hypot(a, b)
}

func degreeDiff(to, from float64) float64 {
	diff := to - from

	if diff > 180 {
		diff = diff - 360
	}

	if diff < -180 {
		diff = diff + 360
	}

	return diff
}

func getAngle(xTo, yTo, xFrom, yFrom float64) float64 {

	a := xTo - xFrom
	b := yTo - yFrom

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
