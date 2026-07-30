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
