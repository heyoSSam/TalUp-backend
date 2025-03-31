package rating

import "math"

func CalculateXP(completionTimeSec float64) int {
	optimalTime := 8.0
	maxXP := 5
	minXP := -5

	if completionTimeSec > optimalTime*2 {
		return minXP
	}

	xp := maxXP - int(math.Round(completionTimeSec/optimalTime*5))

	return int(math.Max(float64(minXP), math.Min(float64(maxXP), float64(xp))))
}
