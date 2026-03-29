package stats

import "github.com/z2-cli/internal/strava"

// EfficiencyFactor calculates speed (m/s) divided by average heart rate.
// A higher value means better aerobic efficiency — running faster at the same effort.
func EfficiencyFactor(a strava.Activity) float64 {
	if a.MovingTime == 0 || !a.HasHeartrate || a.AverageHeartrate == 0 {
		return 0
	}
	speed := a.Distance / float64(a.MovingTime)
	return speed / a.AverageHeartrate
}
