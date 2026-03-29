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
