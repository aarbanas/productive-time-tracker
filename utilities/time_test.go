package utilities

import (
	"testing"
	"time"
)

func isolateHolidayCache(t *testing.T) {
	t.Helper()
	holidayMu.Lock()
	holidayCache = nil
	holidayMu.Unlock()
	t.Cleanup(func() {
		holidayMu.Lock()
		holidayCache = nil
		holidayMu.Unlock()
	})
}

func seedHolidayCache(t *testing.T, year int, dates map[string]bool) {
	t.Helper()
	holidayMu.Lock()
	if holidayCache == nil {
		holidayCache = make(map[int]map[string]bool)
	}
	cp := make(map[string]bool, len(dates))
	for k, v := range dates {
		cp[k] = v
	}
	holidayCache[year] = cp
	holidayMu.Unlock()
}

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

func TestParseExecutionTime(t *testing.T) {
	tests := []struct {
		in           string
		wantH, wantM int
		wantErr      bool
	}{
		{in: "09:30", wantH: 9, wantM: 30},
		{in: "00:00", wantH: 0, wantM: 0},
		{in: "23:59", wantH: 23, wantM: 59},
		{in: "not-a-time", wantErr: true},
		{in: "25:00", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			h, m, err := parseExecutionTime(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h != tt.wantH || m != tt.wantM {
				t.Fatalf("got hour=%d min=%d, want hour=%d min=%d", h, m, tt.wantH, tt.wantM)
			}
		})
	}
}

func TestSameScheduleMinute(t *testing.T) {
	t.Run("same local minute, different seconds and nanoseconds", func(t *testing.T) {
		a := time.Date(2026, 1, 2, 10, 30, 45, 0, time.Local)
		b := time.Date(2026, 1, 2, 10, 30, 0, 999999999, time.Local)
		if !SameScheduleMinute(a, b) {
			t.Fatal("expected times to match at minute precision")
		}
	})

	t.Run("different minute", func(t *testing.T) {
		a := time.Date(2026, 1, 2, 10, 30, 0, 0, time.Local)
		b := time.Date(2026, 1, 2, 10, 31, 0, 0, time.Local)
		if SameScheduleMinute(a, b) {
			t.Fatal("expected times not to match")
		}
	})

	t.Run("same instant in different zones normalizes to local minute", func(t *testing.T) {
		plus2 := time.FixedZone("Plus2", 2*3600)
		utc := time.Date(2026, 1, 2, 12, 0, 30, 0, time.UTC)
		plus2Wall := time.Date(2026, 1, 2, 14, 0, 45, 0, plus2)
		if !SameScheduleMinute(utc, plus2Wall) {
			t.Fatal("expected same schedule minute after conversion to local")
		}
	})
}

func TestLastWorkingDay(t *testing.T) {
	t.Setenv("TZ", "UTC")

	t.Run("April 2026 — last calendar day is a weekday", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{})

		got, err := lastWorkingDay(2026, time.April)
		if err != nil {
			t.Fatalf("lastWorkingDay: %v", err)
		}
		want := time.Date(2026, time.April, 30, 0, 0, 0, 0, time.Local)
		if !got.Equal(want) {
			t.Fatalf("got %v, want %v", got.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	})

	t.Run("May 2026 — month ends on weekend, walks backward", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{})

		got, err := lastWorkingDay(2026, time.May)
		if err != nil {
			t.Fatalf("lastWorkingDay: %v", err)
		}
		want := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.Local)
		if !got.Equal(want) {
			t.Fatalf("got %v, want %v", got.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	})

	t.Run("April 2026 — public holiday on last weekday shifts earlier", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{"2026-04-30": true})

		got, err := lastWorkingDay(2026, time.April)
		if err != nil {
			t.Fatalf("lastWorkingDay: %v", err)
		}
		want := time.Date(2026, time.April, 29, 0, 0, 0, 0, time.Local)
		if !got.Equal(want) {
			t.Fatalf("got %v, want %v", got.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	})
}

func TestNextRunLastWorkingDayUTC(t *testing.T) {
	t.Setenv("TZ", "UTC")

	t.Run("last working day of April at 12:00 UTC", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{})

		now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC)
		got, err := NextRunLastWorkingDayUTC(now, 12, 0)
		if err != nil {
			t.Fatalf("NextRunLastWorkingDayUTC: %v", err)
		}
		want := time.Date(2026, time.April, 30, 12, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Fatalf("got %v, want %v", got.UTC().Format(time.RFC3339), want.UTC().Format(time.RFC3339))
		}
	})
}

func TestNextRunTime(t *testing.T) {
	t.Setenv("TZ", "UTC")

	t.Run("invalid execution time", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{})

		_, err := NextRunTime(time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC), "not-a-time")
		if err == nil {
			t.Fatal("expected error for bad execution time")
		}
	})

	t.Run("combines last working day of month with parsed clock", func(t *testing.T) {
		isolateHolidayCache(t)
		seedHolidayCache(t, 2026, map[string]bool{})

		now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC)
		got, err := NextRunTime(now, "14:05")
		if err != nil {
			t.Fatalf("NextRunTime: %v", err)
		}
		want := time.Date(2026, time.April, 30, 14, 5, 0, 0, time.Local)
		if !got.Equal(want) {
			t.Fatalf("got %v, want %v", got.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	})
}
