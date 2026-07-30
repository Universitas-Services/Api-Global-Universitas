package jobs

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"api-global/internal/repositories"
	"api-global/internal/scrapers"
	"api-global/internal/timeutil"

	"gorm.io/gorm"
)

// StartBCVDailyJob scrapea el BCV a horas fijas America/Caracas (default 07:00 y 17:00).
// En cada disparo hace upsert: inserta si no hay fila del día, o actualiza el valor si cambió.
func StartBCVDailyJob(db *gorm.DB, hoursCSV string) {
	hours := parseHours(hoursCSV)
	if len(hours) == 0 {
		hours = []int{7, 17}
	}

	loc := timeutil.CaracasLocation()

	go func() {
		log.Printf("⏰ Job BCV iniciado (America/Caracas horas: %v)", hours)

		now := time.Now().In(loc)
		if shouldCatchUp(now, hours) {
			runBCVCapture(db)
		}

		for {
			now = time.Now().In(loc)
			next := nextRun(now, hours, loc)
			log.Printf("⏰ BCV: próxima captura a las %s (Caracas)", next.Format("2006-01-02 15:04"))
			time.Sleep(time.Until(next))
			runBCVCapture(db)
		}
	}()
}

func parseHours(csv string) []int {
	seen := map[int]bool{}
	var hours []int
	for _, part := range strings.Split(csv, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		h, err := strconv.Atoi(part)
		if err != nil || h < 0 || h > 23 || seen[h] {
			continue
		}
		seen[h] = true
		hours = append(hours, h)
	}
	sort.Ints(hours)
	return hours
}

func shouldCatchUp(now time.Time, hours []int) bool {
	for _, h := range hours {
		if now.Hour() > h || (now.Hour() == h && now.Minute() >= 0) {
			return true
		}
	}
	return false
}

func nextRun(now time.Time, hours []int, loc *time.Location) time.Time {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	for _, h := range hours {
		candidate := today.Add(time.Duration(h) * time.Hour)
		if candidate.After(now) {
			return candidate
		}
	}
	// Mañana a la primera hora
	tomorrow := today.Add(24 * time.Hour)
	return tomorrow.Add(time.Duration(hours[0]) * time.Hour)
}

func runBCVCapture(db *gorm.DB) {
	prevUSD, errUSD := repositories.GetIndicadorDeHoy(db, "USD_BCV")
	prevEUR, errEUR := repositories.GetIndicadorDeHoy(db, "EUR_BCV")

	rates, err := scrapers.ScrapeBCV()
	if err != nil {
		log.Printf("❌ BCV job: error scrapeando: %v", err)
		return
	}

	if err := repositories.SaveIndicador(db, "USD_BCV", rates.USD); err != nil {
		log.Printf("❌ BCV job: error guardando USD: %v", err)
		return
	}
	if err := repositories.SaveIndicador(db, "EUR_BCV", rates.EUR); err != nil {
		log.Printf("❌ BCV job: error guardando EUR: %v", err)
		return
	}

	hoy := timeutil.HoyCaracas().Format("2006-01-02")
	usdSame := errUSD == nil && almostEqual(prevUSD.Valor, rates.USD)
	eurSame := errEUR == nil && almostEqual(prevEUR.Valor, rates.EUR)

	switch {
	case errUSD != nil || errEUR != nil:
		log.Printf("✅ BCV job: tasas creadas para %s USD=%.4f EUR=%.4f", hoy, rates.USD, rates.EUR)
	case !usdSame || !eurSame:
		prevU, prevE := 0.0, 0.0
		if errUSD == nil {
			prevU = prevUSD.Valor
		}
		if errEUR == nil {
			prevE = prevEUR.Valor
		}
		log.Printf("✅ BCV job: tasas actualizadas para %s USD=%.4f→%.4f EUR=%.4f→%.4f",
			hoy, prevU, rates.USD, prevE, rates.EUR)
	default:
		log.Printf("⚡ BCV job: sin cambios para %s USD=%.4f EUR=%.4f", hoy, rates.USD, rates.EUR)
	}
}

func almostEqual(a, b float64) bool {
	const eps = 0.00005 // precisión ~4 decimales
	if a > b {
		return a-b < eps
	}
	return b-a < eps
}
