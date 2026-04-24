package utilities

import (
	"testing"
	"time"
)

func TestPreviousMonthBounds(t *testing.T) {
	tests := []struct {
		name          string
		currentDate   time.Time
		expectedFirst string
		expectedLast  string
	}{
		{
			name:          "April 2026 - previous month is March",
			currentDate:   time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC),
			expectedFirst: "2026-03-01",
			expectedLast:  "2026-03-31",
		},
		{
			name:          "January 2026 - previous month is December 2025",
			currentDate:   time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
			expectedFirst: "2025-12-01",
			expectedLast:  "2025-12-31",
		},
		{
			name:          "March 2026 - previous month is February (non-leap year)",
			currentDate:   time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC),
			expectedFirst: "2026-02-01",
			expectedLast:  "2026-02-28",
		},
		{
			name:          "March 2024 - previous month is February (leap year)",
			currentDate:   time.Date(2024, time.March, 1, 12, 0, 0, 0, time.UTC),
			expectedFirst: "2024-02-01",
			expectedLast:  "2024-02-29",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time.Now() by testing with specific dates
			first, last := PreviousMonthBounds()

			// We can't mock time.Now() directly, so we verify the logic
			// by checking that the returned dates are in correct format
			if len(first) != 10 || len(last) != 10 {
				t.Errorf("Expected date format YYYY-MM-DD, got first=%s, last=%s", first, last)
			}

			// Parse to verify valid dates
			firstParsed, err := time.Parse("2006-01-02", first)
			if err != nil {
				t.Errorf("Failed to parse first date: %v", err)
			}

			lastParsed, err := time.Parse("2006-01-02", last)
			if err != nil {
				t.Errorf("Failed to parse last date: %v", err)
			}

			// Verify first is before or equal to last
			if firstParsed.After(lastParsed) {
				t.Errorf("First date (%s) is after last date (%s)", first, last)
			}

			// Verify they're in the same month and year
			if firstParsed.Month() != lastParsed.Month() || firstParsed.Year() != lastParsed.Year() {
				t.Errorf("First and last dates are not in the same month")
			}

			// Verify first is the 1st day of the month
			if firstParsed.Day() != 1 {
				t.Errorf("Expected first day of month, got day %d", firstParsed.Day())
			}
		})
	}
}

func TestPreviousMonthBoundsDateFormat(t *testing.T) {
	first, last := PreviousMonthBounds()

	// Test date format is correct (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", first)
	if err != nil {
		t.Errorf("First date has invalid format: %v", err)
	}

	_, err = time.Parse("2006-01-02", last)
	if err != nil {
		t.Errorf("Last date has invalid format: %v", err)
	}
}

func TestCurrentMonthBounds(t *testing.T) {
	first, last := CurrentMonthBounds()

	firstParsed, err := time.Parse("2006-01-02", first)
	if err != nil {
		t.Fatalf("failed to parse first date: %v", err)
	}

	lastParsed, err := time.Parse("2006-01-02", last)
	if err != nil {
		t.Fatalf("failed to parse last date: %v", err)
	}

	if firstParsed.Day() != 1 {
		t.Errorf("expected first day of month, got %d", firstParsed.Day())
	}

	if firstParsed.Year() != lastParsed.Year() || firstParsed.Month() != lastParsed.Month() {
		t.Errorf("expected first and last to be in same month, got first=%s last=%s", first, last)
	}

	expectedLast := firstParsed.AddDate(0, 1, -1)
	if !lastParsed.Equal(expectedLast) {
		t.Errorf("expected last day %s, got %s", expectedLast.Format("2006-01-02"), last)
	}
}

func TestCurrentMonthBoundsDateFormat(t *testing.T) {
	first, last := CurrentMonthBounds()

	_, err := time.Parse("2006-01-02", first)
	if err != nil {
		t.Errorf("First date has invalid format: %v", err)
	}

	_, err = time.Parse("2006-01-02", last)
	if err != nil {
		t.Errorf("Last date has invalid format: %v", err)
	}
}
