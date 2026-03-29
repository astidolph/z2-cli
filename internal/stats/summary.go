package stats

import "github.com/strava-cli/internal/strava"

type Summary struct {
	Count    int
	AvgEF    float64
	AvgHR    float64
	AvgPace  float64 // seconds per km
	TotalKm  float64
}

func Summarise(runs []strava.Activity) Summary {
	if len(runs) == 0 {
		return Summary{}
	}

	var totalEF, totalHR, totalPaceSeconds float64
	var efCount, hrCount, paceCount int

	for _, r := range runs {
		if ef := EfficiencyFactor(r); ef > 0 {
			totalEF += ef
			efCount++
		}
		if r.HasHeartrate && r.AverageHeartrate > 0 {
			totalHR += r.AverageHeartrate
			hrCount++
		}
		if r.Distance > 0 && r.MovingTime > 0 {
			paceSeconds := float64(r.MovingTime) / (r.Distance / 1000.0)
			totalPaceSeconds += paceSeconds
			paceCount++
		}
	}

	s := Summary{Count: len(runs)}

	if efCount > 0 {
		s.AvgEF = totalEF / float64(efCount)
	}
	if hrCount > 0 {
		s.AvgHR = totalHR / float64(hrCount)
	}
	if paceCount > 0 {
		s.AvgPace = totalPaceSeconds / float64(paceCount)
	}

	for _, r := range runs {
		s.TotalKm += r.Distance / 1000.0
	}

	return s
}

// TrendPercent calculates the percentage change in EF between two summaries.
// Positive means improvement (current period is faster at same effort).
func TrendPercent(current, previous Summary) float64 {
	if previous.AvgEF == 0 {
		return 0
	}
	return ((current.AvgEF - previous.AvgEF) / previous.AvgEF) * 100
}
