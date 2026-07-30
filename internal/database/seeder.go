package database

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"api-global/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedTerritories lee el JSON e inserta los datos si la tabla está vacía
func SeedTerritories(db *gorm.DB) {
	var count int64
	db.Model(&models.Estado{}).Count(&count)

	// Si ya hay registros, no hacemos nada para no duplicar
	if count > 0 {
		log.Println("⚡ Base de datos ya poblada. Omitiendo Seeder territorial.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando datos territoriales de Venezuela...")

	// Leemos el archivo JSON (La ruta relativa funcionará porque la configuramos en el Dockerfile)
	bytes, err := os.ReadFile("internal/database/seeds/venezuela.json")
	if err != nil {
		log.Printf("❌ Error leyendo archivo venezuela.json: %v\n", err)
		return
	}

	var estados []models.Estado
	if err := json.Unmarshal(bytes, &estados); err != nil {
		log.Printf("❌ Error decodificando el JSON: %v\n", err)
		return
	}

	// Magia de GORM: Inserta el estado, obtiene el ID, inserta los municipios con ese ID,
	// obtiene sus IDs e inserta las parroquias correspondientes. Todo en una transacción.
	if err := db.Create(&estados).Error; err != nil {
		log.Printf("❌ Error insertando datos en PostgreSQL: %v\n", err)
		return
	}

	log.Println("✅ Seeder completado: Venezuela cargada exitosamente en la base de datos.")
}

// ==========================================
// SEEDER DE TRIBUNALES
// ==========================================

// Tipos auxiliares para deserializar el JSON de tribunales
type tribunalSeedEntry struct {
	Nombre    string `json:"nombre"`
	Categoria string `json:"categoria"`
}

type tribunalMunicipalGroup struct {
	Municipios []string `json:"municipios"`
	Tribunales []struct {
		Nombre string `json:"nombre"`
	} `json:"tribunales"`
}

type tribunalSeedEstado struct {
	Estado                string                   `json:"estado"`
	TribunalesEstadales   []tribunalSeedEntry      `json:"tribunales_estadales"`
	TribunalesMunicipales []tribunalMunicipalGroup `json:"tribunales_municipales"`
}

// SeedTribunales lee los JSON de tribunales e inserta los datos si la tabla está vacía
func SeedTribunales(db *gorm.DB) {
	var count int64
	db.Model(&models.Tribunal{}).Count(&count)

	if count > 0 {
		log.Println("⚡ Base de datos ya poblada. Omitiendo Seeder de tribunales.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando tribunales judiciales...")

	files, err := filepath.Glob("internal/database/seeds/tribunales_*.json")
	if err != nil {
		log.Printf("❌ Error buscando archivos de tribunales: %v\n", err)
		return
	}

	if len(files) == 0 {
		log.Println("⚠️ No se encontraron archivos tribunales_*.json en internal/database/seeds/")
		return
	}

	for _, file := range files {
		bytes, err := os.ReadFile(file)
		if err != nil {
			log.Printf("❌ Error leyendo archivo %s: %v\n", file, err)
			continue
		}

		var seedData []tribunalSeedEstado
		if err := json.Unmarshal(bytes, &seedData); err != nil {
			log.Printf("❌ Error decodificando tribunales JSON en %s: %v\n", file, err)
			continue
		}

		for _, estadoData := range seedData {
			// Buscar el estado en la BD
			var estado models.Estado
			if err := db.Where("nombre = ?", estadoData.Estado).First(&estado).Error; err != nil {
				log.Printf("⚠️ Estado '%s' no encontrado, omitiendo...\n", estadoData.Estado)
				continue
			}

			// Insertar tribunales estadales (Categoría A y B)
			for _, t := range estadoData.TribunalesEstadales {
				estadoID := estado.ID
				tribunal := models.Tribunal{
					Nombre:    t.Nombre,
					Categoria: t.Categoria,
					EstadoID:  &estadoID,
				}
				if err := db.Create(&tribunal).Error; err != nil {
					log.Printf("⚠️ Error creando tribunal estadal '%s': %v\n", t.Nombre, err)
				}
			}

			// Insertar tribunales municipales (Categoría C)
			for _, group := range estadoData.TribunalesMunicipales {
				// Resolver nombres de municipios a registros de la BD
				var municipios []models.Municipio
				for _, munName := range group.Municipios {
					var mun models.Municipio
					if err := db.Where("nombre = ? AND estado_id = ?", munName, estado.ID).First(&mun).Error; err != nil {
						log.Printf("⚠️ Municipio '%s' no encontrado en estado '%s'\n", munName, estadoData.Estado)
						continue
					}
					municipios = append(municipios, mun)
				}

				for _, t := range group.Tribunales {
					tribunal := models.Tribunal{
						Nombre:     t.Nombre,
						Categoria:  "Municipio",
						Municipios: municipios,
					}
					if err := db.Create(&tribunal).Error; err != nil {
						log.Printf("⚠️ Error creando tribunal municipal '%s': %v\n", t.Nombre, err)
					}
				}
			}

			log.Printf("✅ Tribunales del estado '%s' cargados desde %s.\n", estadoData.Estado, filepath.Base(file))
		}
	}

	log.Println("✅ Seeder completado: Tribunales cargados exitosamente en la base de datos.")
}

// ==========================================
// SEEDER DE CIUDADES
// ==========================================

type ciudadSeedEstado struct {
	Estado   string   `json:"estado"`
	Ciudades []string `json:"ciudades"`
}

// SeedCiudades lee el JSON de ciudades e inserta los datos
func SeedCiudades(db *gorm.DB) {
	var count int64
	db.Model(&models.Ciudad{}).Count(&count)

	if count > 0 {
		log.Println("⚡ Base de datos ya poblada. Omitiendo Seeder de ciudades.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando ciudades de Venezuela...")

	bytes, err := os.ReadFile("internal/database/seeds/ciudades.json")
	if err != nil {
		log.Printf("❌ Error leyendo archivo ciudades.json: %v\n", err)
		return
	}

	var seedData []ciudadSeedEstado
	if err := json.Unmarshal(bytes, &seedData); err != nil {
		log.Printf("❌ Error decodificando ciudades JSON: %v\n", err)
		return
	}

	for _, estadoData := range seedData {
		var estado models.Estado
		if err := db.Where("nombre = ?", estadoData.Estado).First(&estado).Error; err != nil {
			log.Printf("⚠️ Estado '%s' no encontrado para insertar ciudades, omitiendo...\n", estadoData.Estado)
			continue
		}

		for _, cityName := range estadoData.Ciudades {
			ciudad := models.Ciudad{
				Nombre:   cityName,
				EstadoID: estado.ID,
			}
			if err := db.Create(&ciudad).Error; err != nil {
				log.Printf("⚠️ Error creando ciudad '%s': %v\n", cityName, err)
			}
		}
	}

	log.Println("✅ Seeder completado: Ciudades cargadas exitosamente en la base de datos.")
}

// ==========================================
// SEEDER DE HISTÓRICO BCV
// ==========================================

// SeedBCVHistorico carga tasas USD/EUR históricas desde CSV a indicadores_economicos
func SeedBCVHistorico(db *gorm.DB) {
	var count int64
	db.Model(&models.IndicadorEconomico{}).Where("tipo = ?", "USD_BCV").Count(&count)

	// Si ya hay un histórico cargado (no solo 1-2 días de scrape), omitir
	if count >= 1000 {
		log.Println("⚡ Histórico BCV ya poblado. Omitiendo Seeder de BCV.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando histórico BCV...")

	file, err := os.Open("internal/database/seeds/bcv_historico.csv")
	if err != nil {
		log.Printf("❌ Error leyendo bcv_historico.csv: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("❌ Error parseando bcv_historico.csv: %v\n", err)
		return
	}

	if len(records) < 2 {
		log.Println("⚠️ bcv_historico.csv no tiene datos")
		return
	}

	var batch []models.IndicadorEconomico
	for i, row := range records {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 {
			continue
		}

		fecha, err := time.Parse("2006-01-02", strings.TrimSpace(row[0]))
		if err != nil {
			log.Printf("⚠️ Fecha inválida en fila %d: %v\n", i+1, err)
			continue
		}

		if usdStr := strings.TrimSpace(row[1]); usdStr != "" {
			usd, err := strconv.ParseFloat(usdStr, 64)
			if err == nil {
				batch = append(batch, models.IndicadorEconomico{
					Tipo:  "USD_BCV",
					Valor: usd,
					Fecha: fecha,
				})
			}
		}

		if eurStr := strings.TrimSpace(row[2]); eurStr != "" {
			eur, err := strconv.ParseFloat(eurStr, 64)
			if err == nil {
				batch = append(batch, models.IndicadorEconomico{
					Tipo:  "EUR_BCV",
					Valor: eur,
					Fecha: fecha,
				})
			}
		}
	}

	if len(batch) == 0 {
		log.Println("⚠️ No se generaron registros BCV para insertar")
		return
	}

	// DoNothing evita fallar si ya existe la tasa de hoy por un scrape previo
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(batch, 500).Error; err != nil {
		log.Printf("❌ Error insertando histórico BCV: %v\n", err)
		return
	}

	log.Printf("✅ Seeder completado: %d registros BCV procesados hacia indicadores_economicos.\n", len(batch))
}

// ==========================================
// SEEDER DE CÓDIGOS DE ÁREA
// ==========================================

// SeedCodigosArea carga los códigos de área telefónicos desde JSON
func SeedCodigosArea(db *gorm.DB) {
	var count int64
	db.Model(&models.CodigoArea{}).Count(&count)

	if count > 0 {
		log.Println("⚡ Códigos de área ya poblados. Omitiendo Seeder.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando códigos de área...")

	bytes, err := os.ReadFile("internal/database/seeds/codigos_area.json")
	if err != nil {
		log.Printf("❌ Error leyendo codigos_area.json: %v\n", err)
		return
	}

	var codigos []string
	if err := json.Unmarshal(bytes, &codigos); err != nil {
		log.Printf("❌ Error decodificando codigos_area.json: %v\n", err)
		return
	}

	var batch []models.CodigoArea
	for _, codigo := range codigos {
		batch = append(batch, models.CodigoArea{Codigo: codigo})
	}

	if err := db.CreateInBatches(batch, 100).Error; err != nil {
		log.Printf("❌ Error insertando códigos de área: %v\n", err)
		return
	}

	log.Printf("✅ Seeder completado: %d códigos de área cargados.\n", len(batch))
}

