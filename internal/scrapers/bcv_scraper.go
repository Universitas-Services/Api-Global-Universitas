package scrapers

import (
	"crypto/tls"
	"fmt"
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

// ScrapeBCV se conecta a bcv.org.ve y extrae las tasas del día
func ScrapeBCV() (*BCVRates, error) {
	// Crear un transporte que ignore los errores de certificado (InsecureSkipVerify)
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Aplicar el transporte al cliente HTTP
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: customTransport,
	}

	req, err := http.NewRequest("GET", "https://www.bcv.org.ve/", nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %v", err)
	}

	// 2. Falsificar el User-Agent para que el firewall del BCV no nos bloquee
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error conectando al BCV: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("el BCV respondió con código de error: %d", res.StatusCode)
	}

	// 3. Cargar el HTML en goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo el HTML: %v", err)
	}

	rates := &BCVRates{}

	// 4. Extraer y limpiar el Dólar (El BCV usa un div con ID "dolar" y dentro un strong)
	usdText := doc.Find("#dolar strong").Text()
	rates.USD = cleanNumber(usdText)

	// 5. Extraer y limpiar el Euro
	eurText := doc.Find("#euro strong").Text()
	rates.EUR = cleanNumber(eurText)

	// Validamos que haya extraído algo coherente
	if rates.USD == 0 || rates.EUR == 0 {
		return nil, fmt.Errorf("no se pudieron extraer las tasas, es posible que el BCV haya cambiado su diseño")
	}

	return rates, nil
}

// cleanNumber quita espacios, cambia la coma por punto y convierte a float64
func cleanNumber(text string) float64 {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", ".")
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	return val
}
