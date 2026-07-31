package scrapers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// BCVRates contendrá los valores extraídos
type BCVRates struct {
	USD float64
	EUR float64
}

// ScrapeBCV intenta bcv.org.ve (timeout 20s) y, si falla, Al Cambio como fallback.
func ScrapeBCV() (*BCVRates, error) {
	rates, err := scrapeBCVPrimary()
	if err == nil {
		log.Printf("📡 Tasas obtenidas desde BCV.org.ve USD=%.4f EUR=%.4f", rates.USD, rates.EUR)
		return rates, nil
	}

	log.Printf("⚠️ BCV primario falló (%v). Intentando fallback Al Cambio...", err)

	fallback, fbErr := scrapeAlCambio()
	if fbErr != nil {
		return nil, fmt.Errorf("bcv: %v; alcambio: %v", err, fbErr)
	}

	log.Printf("📡 Tasas obtenidas desde Al Cambio (fallback) USD=%.4f EUR=%.4f", fallback.USD, fallback.EUR)
	return fallback, nil
}

func scrapeBCVPrimary() (*BCVRates, error) {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: customTransport,
	}

	req, err := http.NewRequest("GET", "https://www.bcv.org.ve/", nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error conectando al BCV: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("el BCV respondió con código de error: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo el HTML: %v", err)
	}

	rates := &BCVRates{}
	rates.USD = cleanNumber(doc.Find("#dolar strong").Text())
	rates.EUR = cleanNumber(doc.Find("#euro strong").Text())

	if rates.USD == 0 || rates.EUR == 0 {
		return nil, fmt.Errorf("no se pudieron extraer las tasas, es posible que el BCV haya cambiado su diseño")
	}

	return rates, nil
}

func cleanNumber(text string) float64 {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", ".")
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	return val
}
