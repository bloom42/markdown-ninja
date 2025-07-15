package timeutil

import "time"

// GetFirstDayOfNextMonth returns the first day of the next month for a given date.
// If the current day is already the first day of the month then it returns date.
func GetFirstDayOfNextMonth(date time.Time) time.Time {
	return GetFirstDayOfTheMonth(date).AddDate(0, 1, 0)
}

func GetFirstDayOfTheMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func GetFirstDayOfLastMonth(date time.Time) time.Time {
	month := date.Month() - 1
	year := date.Year()

	// if the current month is January, we need to go back to December of the previous year
	if date.Month() == 1 {
		year -= 1
		month = 12
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
}

func GetLastHourOfLastMonth(date time.Time) time.Time {
	return GetFirstDayOfTheMonth(date).Add(-time.Hour)
}

func GetFirstMondayOfTheMonth(date time.Time) time.Time {
	day := GetFirstDayOfTheMonth(date)

	for day.Weekday() != time.Monday {
		day = day.AddDate(0, 0, 1) // Move to the next day
	}

	return day
}
