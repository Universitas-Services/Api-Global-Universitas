package jobs

import (
	"fmt"
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

// CaptureResult describe el resultado de un scrape+upsert BCV por fecha efectiva.
type CaptureResult struct {
	Fecha   string  `json:"fecha"`
	USD     float64 `json:"usd"`
	EUR     float64 `json:"eur"`
	Status  string  `json:"status"` // created | updated | unchanged
	Message string  `json:"message"`
}

// CaptureBCV scrapea el BCV y hace upsert de USD/EUR según la hora en Caracas:
// antes de las 17:00 → día actual; a las 17:00 o después → día siguiente.
func CaptureBCV(db *gorm.DB) (*CaptureResult, error) {
	fecha := timeutil.FechaEfectivaBCV()
	fechaStr := fecha.Format("2006-01-02")

	prevUSD, errUSD := repositories.GetIndicadorPorFecha(db, "USD_BCV", fecha)
	prevEUR, errEUR := repositories.GetIndicadorPorFecha(db, "EUR_BCV", fecha)

	rates, err := scrapers.ScrapeBCV()
	if err != nil {
		return nil, fmt.Errorf("error scrapeando BCV: %w", err)
	}

	if err := repositories.SaveIndicador(db, "USD_BCV", rates.USD, fecha); err != nil {
		return nil, fmt.Errorf("error guardando USD: %w", err)
	}
	if err := repositories.SaveIndicador(db, "EUR_BCV", rates.EUR, fecha); err != nil {
		return nil, fmt.Errorf("error guardando EUR: %w", err)
	}

	usdSame := errUSD == nil && almostEqual(prevUSD.Valor, rates.USD)
	eurSame := errEUR == nil && almostEqual(prevEUR.Valor, rates.EUR)

	result := &CaptureResult{
		Fecha: fechaStr,
		USD:   rates.USD,
		EUR:   rates.EUR,
	}

	switch {
	case errUSD != nil || errEUR != nil:
		result.Status = "created"
		result.Message = "Tasas BCV creadas para " + fechaStr
	case !usdSame || !eurSame:
		result.Status = "updated"
		result.Message = "Tasas BCV actualizadas para " + fechaStr
	default:
		result.Status = "unchanged"
		result.Message = "Tasas BCV sin cambios para " + fechaStr
	}

	return result, nil
}

// StartBCVDailyJob scrapea el BCV a horas fijas America/Caracas (default 07:00 y 17:00).
// En Cloud Run con min=0 preferir Cloud Scheduler + POST /bcv/capturar.
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
	tomorrow := today.Add(24 * time.Hour)
	return tomorrow.Add(time.Duration(hours[0]) * time.Hour)
}

func runBCVCapture(db *gorm.DB) {
	result, err := CaptureBCV(db)
	if err != nil {
		log.Printf("❌ BCV job: %v", err)
		return
	}
	log.Printf("✅ BCV job [%s]: %s USD=%.4f EUR=%.4f", result.Status, result.Message, result.USD, result.EUR)
}

func almostEqual(a, b float64) bool {
	const eps = 0.00005
	if a > b {
		return a-b < eps
	}
	return b-a < eps
}
