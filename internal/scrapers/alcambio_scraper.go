package scrapers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type alcambioRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type alcambioResponse struct {
	Data struct {
		GetCountryConversions *struct {
			DateBcv         int64 `json:"dateBcv"`
			ConversionRates []struct {
				Official     bool    `json:"official"`
				Type         string  `json:"type"`
				RateValue    float64 `json:"rateValue"`
				BaseValue    float64 `json:"baseValue"`
				UsesRateValue bool   `json:"usesRateValue"`
				RateCurrency struct {
					Code string `json:"code"`
				} `json:"rateCurrency"`
			} `json:"conversionRates"`
		} `json:"getCountryConversions"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

const alcambioQuery = `
query getCountryConversions($countryCode: String!, $dateSearch: DateSearchInput) {
  getCountryConversions(
    payload: { countryCode: $countryCode }
    dateSearch: $dateSearch
  ) {
    dateBcv
    conversionRates {
      official
      type
      rateValue
      baseValue
      usesRateValue
      rateCurrency { code }
    }
  }
}
`

// scrapeAlCambio obtiene USD/EUR oficiales vía GraphQL de api.alcambio.app
func scrapeAlCambio() (*BCVRates, error) {
	payload := alcambioRequest{
		Query:     alcambioQuery,
		Variables: map[string]interface{}{"countryCode": "VE"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("alcambio: error armando JSON: %w", err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("POST", "https://api.alcambio.app/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("alcambio: error creando petición: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://alcambio.app")
	req.Header.Set("Referer", "https://alcambio.app/")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alcambio: error conectando: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("alcambio: error leyendo respuesta: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alcambio: HTTP %d: %s", res.StatusCode, string(raw))
	}

	var parsed alcambioResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("alcambio: JSON inválido: %w", err)
	}
	if len(parsed.Errors) > 0 {
		return nil, fmt.Errorf("alcambio: graphql: %s", parsed.Errors[0].Message)
	}
	if parsed.Data.GetCountryConversions == nil {
		return nil, fmt.Errorf("alcambio: sin datos getCountryConversions")
	}

	rates := &BCVRates{}
	for _, cr := range parsed.Data.GetCountryConversions.ConversionRates {
		if !cr.Official {
			continue
		}
		value := cr.BaseValue
		if cr.UsesRateValue || cr.RateValue > 0 {
			if cr.RateValue > 0 {
				value = cr.RateValue
			}
		}
		if value <= 0 {
			continue
		}
		switch cr.RateCurrency.Code {
		case "USD":
			// Preferir SECONDARY (vista home BCV); si ya hay valor, no pisar con OTHER
			if rates.USD == 0 || cr.Type == "SECONDARY" {
				rates.USD = value
			}
		case "EUR":
			if rates.EUR == 0 {
				rates.EUR = value
			}
		}
	}

	if rates.USD == 0 || rates.EUR == 0 {
		return nil, fmt.Errorf("alcambio: no se pudieron extraer USD/EUR oficiales")
	}

	if ts := parsed.Data.GetCountryConversions.DateBcv; ts > 0 {
		log.Printf("📡 Al Cambio: dateBcv=%s USD=%.4f EUR=%.4f",
			time.UnixMilli(ts).UTC().Format("2006-01-02"), rates.USD, rates.EUR)
	}

	return rates, nil
}
