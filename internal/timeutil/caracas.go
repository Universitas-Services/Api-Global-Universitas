package timeutil

import "time"

// CaracasLocation returns America/Caracas, or a fixed UTC-4 zone as fallback.
func CaracasLocation() *time.Location {
	loc, err := time.LoadLocation("America/Caracas")
	if err != nil {
		return time.FixedZone("VET", -4*60*60)
	}
	return loc
}

// HoyCaracas returns today's calendar date in Caracas as UTC midnight
// (suitable for PostgreSQL date columns and seed CSV dates).
func HoyCaracas() time.Time {
	now := time.Now().In(CaracasLocation())
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// MananaCaracas returns tomorrow's calendar date in Caracas as UTC midnight.
func MananaCaracas() time.Time {
	return HoyCaracas().AddDate(0, 0, 1)
}

const bcvCutoffHour = 17

// FechaEfectivaBCV returns the calendar date that a BCV scrape should be stored under.
// Before 17:00 Caracas the rate applies to today; at 17:00 or later it applies to tomorrow
// (the bank publishes the next day's official rate in the afternoon).
func FechaEfectivaBCV() time.Time {
	return FechaEfectivaBCVAt(time.Now())
}

// FechaEfectivaBCVAt is like FechaEfectivaBCV but uses the provided instant (for tests).
func FechaEfectivaBCVAt(t time.Time) time.Time {
	now := t.In(CaracasLocation())
	hoy := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if now.Hour() >= bcvCutoffHour {
		return hoy.AddDate(0, 0, 1)
	}
	return hoy
}
