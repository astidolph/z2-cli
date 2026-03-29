package strava

import "time"

func FilterByWeekday(activities []Activity, day time.Weekday) []Activity {
	var filtered []Activity
	for _, a := range activities {
		t, err := a.StartTime()
		if err != nil {
			continue
		}
		if t.Weekday() == day {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func FilterByMaxHR(activities []Activity, maxHR float64) []Activity {
	var filtered []Activity
	for _, a := range activities {
		if a.HasHeartrate && a.AverageHeartrate <= maxHR {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func FilterByMinDistance(activities []Activity, minKm float64) []Activity {
	var filtered []Activity
	minMeters := minKm * 1000.0
	for _, a := range activities {
		if a.Distance >= minMeters {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
