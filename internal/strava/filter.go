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

func FilterByZone2(activities []Activity, lowHR, highHR float64) []Activity {
	var filtered []Activity
	for _, a := range activities {
		if a.HasHeartrate && a.AverageHeartrate >= lowHR && a.AverageHeartrate <= highHR {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
