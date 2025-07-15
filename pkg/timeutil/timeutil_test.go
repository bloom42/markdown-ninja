package timeutil

import (
	"testing"
	"time"
)

func TestGetFirstDayOfTheMonth(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected time.Time
	}{
		{
			time.Date(2024, 1, 1, 0, 0, 1, 0, time.Local),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			time.Date(2025, 1, 3, 0, 3, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2026, 1, 31, 0, 0, 0, 1, time.UTC),
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2027, 1, 31, 1, 0, 0, 1, time.UTC),
			time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := GetFirstDayOfTheMonth(test.Date)
		if !result.Equal(test.Expected) {
			t.Errorf("Invalid result. Got: %v | Expected: %v", result, test.Expected)
		}
	}
}

func TestGetFirstDayOfNextMonth(t *testing.T) {
	tests := []struct {
		Date     time.Time
		Expected time.Time
	}{
		{
			time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2025, 1, 3, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2026, 1, 31, 0, 0, 1, 0, time.UTC),
			time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2027, 12, 1, 0, 0, 1, 0, time.UTC),
			time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := GetFirstDayOfNextMonth(test.Date)
		if !result.Equal(test.Expected) {
			t.Errorf("Invalid result. Got: %v | Expected: %v", result, test.Expected)
		}
	}
}
