package utilities

import (
	"testing"
	"time"
)

func TestPreviousMonthBoundsAt(t *testing.T) {
	tests := []struct {
		name          string
		currentDate   time.Time
		wantFirst     string
		wantLast      string
	}{
		{
			name:        "April 2026 — previous month is March",
			currentDate: time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC),
			wantFirst:   "2026-03-01",
			wantLast:    "2026-03-31",
		},
		{
			name:        "January 2026 — previous month is December 2025",
			currentDate: time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
			wantFirst:   "2025-12-01",
			wantLast:    "2025-12-31",
		},
		{
			name:        "March 2026 — previous month is February (non-leap)",
			currentDate: time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC),
			wantFirst:   "2026-02-01",
			wantLast:    "2026-02-28",
		},
		{
			name:        "March 2024 — previous month is February (leap year)",
			currentDate: time.Date(2024, time.March, 1, 12, 0, 0, 0, time.UTC),
			wantFirst:   "2024-02-01",
			wantLast:    "2024-02-29",
		},
		{
			name:        "first day of month — previous month is full prior month",
			currentDate: time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
			wantFirst:   "2025-05-01",
			wantLast:    "2025-05-31",
		},
		{
			name:        "non-UTC location — calendar month still uses UTC (matches implementation)",
			currentDate: time.Date(2026, time.April, 2, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			wantFirst:   "2026-03-01",
			wantLast:    "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirst, gotLast := previousMonthBoundsAt(tt.currentDate)
			if gotFirst != tt.wantFirst || gotLast != tt.wantLast {
				t.Fatalf("previousMonthBoundsAt(%v) = (%q, %q), want (%q, %q)",
					tt.currentDate, gotFirst, gotLast, tt.wantFirst, tt.wantLast)
			}
		})
	}
}

func TestRequiredWorkingMinutesPreviousMonthAt(t *testing.T) {
	tests := []struct {
		name         string
		currentDate  time.Time
		minHours     int
		wantMinutes  int32
		wantFirstStr string
		wantLastStr  string
	}{
		{
			name:         "April 2026 — March 2026 has 22 weekdays × 8h",
			currentDate:  time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC),
			minHours:     8,
			wantMinutes:  22 * 8 * 60,
			wantFirstStr: "2026-03-01",
			wantLastStr:  "2026-03-31",
		},
		{
			name:         "January 2026 — December 2025 has 23 weekdays × 8h",
			currentDate:  time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
			minHours:     8,
			wantMinutes:  23 * 8 * 60,
			wantFirstStr: "2025-12-01",
			wantLastStr:  "2025-12-31",
		},
		{
			name:         "March 2024 — February 2024 (leap) has 21 weekdays × 8h",
			currentDate:  time.Date(2024, time.March, 1, 12, 0, 0, 0, time.UTC),
			minHours:     8,
			wantMinutes:  21 * 8 * 60,
			wantFirstStr: "2024-02-01",
			wantLastStr:  "2024-02-29",
		},
		{
			name:         "March 2026 — February 2026 has 20 weekdays × 8h",
			currentDate:  time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC),
			minHours:     8,
			wantMinutes:  20 * 8 * 60,
			wantFirstStr: "2026-02-01",
			wantLastStr:  "2026-02-28",
		},
		{
			name:         "custom minHours scales linearly",
			currentDate:  time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC),
			minHours:     6,
			wantMinutes:  22 * 6 * 60,
			wantFirstStr: "2026-03-01",
			wantLastStr:  "2026-03-31",
		},
		{
			name:         "minHours 0 yields zero minutes but correct bounds",
			currentDate:  time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC),
			minHours:     0,
			wantMinutes:  0,
			wantFirstStr: "2026-03-01",
			wantLastStr:  "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotFirst, gotLast := requiredWorkingMinutesPreviousMonthAt(tt.currentDate, tt.minHours)
			if gotMin != tt.wantMinutes {
				t.Errorf("minutes = %d, want %d", gotMin, tt.wantMinutes)
			}
			if gotFirst.Format("2006-01-02") != tt.wantFirstStr {
				t.Errorf("first bound = %s, want %s", gotFirst.Format("2006-01-02"), tt.wantFirstStr)
			}
			if gotLast.Format("2006-01-02") != tt.wantLastStr {
				t.Errorf("last bound = %s, want %s", gotLast.Format("2006-01-02"), tt.wantLastStr)
			}

			pFirst, pLast := previousMonthBoundsAt(tt.currentDate)
			if pFirst != tt.wantFirstStr || pLast != tt.wantLastStr {
				t.Errorf("previousMonthBoundsAt inconsistent with required bounds: got (%s,%s) want (%s,%s)",
					pFirst, pLast, tt.wantFirstStr, tt.wantLastStr)
			}
		})
	}
}

func TestPreviousMonthBounds_smoke(t *testing.T) {
	first, last := PreviousMonthBounds()
	if len(first) != 10 || len(last) != 10 {
		t.Fatalf("expected YYYY-MM-DD length, got first=%q last=%q", first, last)
	}
	firstParsed, err := time.Parse("2006-01-02", first)
	if err != nil {
		t.Fatalf("parse first: %v", err)
	}
	lastParsed, err := time.Parse("2006-01-02", last)
	if err != nil {
		t.Fatalf("parse last: %v", err)
	}
	if firstParsed.After(lastParsed) {
		t.Fatalf("first after last: %s > %s", first, last)
	}
	if firstParsed.Day() != 1 {
		t.Fatalf("first day should be 1st, got %v", firstParsed)
	}
	if firstParsed.Month() != lastParsed.Month() || firstParsed.Year() != lastParsed.Year() {
		t.Fatalf("first and last not same month: %s, %s", first, last)
	}
}

func TestRequiredWorkingMinutesPreviousMonth_smoke(t *testing.T) {
	minutes, first, last := RequiredWorkingMinutesPreviousMonth(8)
	if minutes <= 0 {
		t.Fatalf("expected positive minutes, got %d", minutes)
	}
	if minutes%(8*60) != 0 {
		t.Fatalf("expected multiple of 8h in minutes, got %d", minutes)
	}
	pf, pl := previousMonthBoundsAt(time.Now())
	if first.Format("2006-01-02") != pf || last.Format("2006-01-02") != pl {
		t.Fatalf("time bounds mismatch: got (%s,%s) vs PreviousMonthBounds (%s,%s)",
			first.Format("2006-01-02"), last.Format("2006-01-02"), pf, pl)
	}
}
