package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"api-global/internal/models"

	"gorm.io/gorm"
)

const territorySeedPath = "internal/database/seeds/venezuela.json"

type estadoSeed struct {
	Nombre     string          `json:"nombre"`
	Municipios []municipioSeed `json:"municipios"`
}

type municipioSeed struct {
	Nombre     string   `json:"nombre"`
	Parroquias []string `json:"parroquias"`
	Ciudades   []string `json:"ciudades"`
}

// SeedTerritories carga o completa los datos territoriales a partir del seed JSON.
func SeedTerritories(db *gorm.DB) {
	log.Println("Iniciando seeder territorial de Venezuela...")

	estados, warnings, err := loadTerritorySeeds(territorySeedPath)
	if err != nil {
		log.Printf("Error cargando seed territorial: %v\n", err)
		return
	}

	for _, warning := range warnings {
		log.Printf("Advertencia seed territorial: %s\n", warning)
	}

	var estadosInsertados int
	var municipiosInsertados int
	var parroquiasInsertadas int
	var ciudadesInsertadas int

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, estadoSeed := range estados {
			estado := models.Estado{Nombre: estadoSeed.Nombre}
			created, err := firstOrCreateEstado(tx, &estado)
			if err != nil {
				return err
			}
			if created {
				estadosInsertados++
			}

			for _, municipioSeed := range estadoSeed.Municipios {
				municipio := models.Municipio{
					EstadoID: estado.ID,
					Nombre:   municipioSeed.Nombre,
				}
				created, err = firstOrCreateMunicipio(tx, &municipio)
				if err != nil {
					return err
				}
				if created {
					municipiosInsertados++
				}

				for _, nombreParroquia := range municipioSeed.Parroquias {
					parroquia := models.Parroquia{
						MunicipioID: municipio.ID,
						Nombre:      nombreParroquia,
					}
					created, err = firstOrCreateParroquia(tx, &parroquia)
					if err != nil {
						return err
					}
					if created {
						parroquiasInsertadas++
					}
				}

				for _, nombreCiudad := range municipioSeed.Ciudades {
					ciudad := models.Ciudad{
						MunicipioID: municipio.ID,
						Nombre:      nombreCiudad,
					}
					created, err = firstOrCreateCiudad(tx, &ciudad)
					if err != nil {
						return err
					}
					if created {
						ciudadesInsertadas++
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error insertando seed territorial: %v\n", err)
		return
	}

	log.Printf(
		"Seeder territorial completado. Estados nuevos: %d, municipios nuevos: %d, parroquias nuevas: %d, ciudades nuevas: %d\n",
		estadosInsertados,
		municipiosInsertados,
		parroquiasInsertadas,
		ciudadesInsertadas,
	)
}

func loadTerritorySeeds(path string) ([]estadoSeed, []string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("leyendo %s: %w", path, err)
	}

	var estados []estadoSeed
	if err := json.Unmarshal(bytes, &estados); err != nil {
		return nil, nil, fmt.Errorf("decodificando %s: %w", path, err)
	}

	warnings, err := validateTerritorySeeds(estados)
	if err != nil {
		return nil, nil, err
	}

	return estados, warnings, nil
}

func validateTerritorySeeds(estados []estadoSeed) ([]string, error) {
	var warnings []string
	seenStates := make(map[string]struct{}, len(estados))

	for _, estado := range estados {
		estadoNombre := cleanSeedName(estado.Nombre)
		if estadoNombre == "" {
			return nil, fmt.Errorf("hay un estado sin nombre")
		}
		if _, exists := seenStates[normalizedSeedKey(estadoNombre)]; exists {
			return nil, fmt.Errorf("estado duplicado en seed: %s", estadoNombre)
		}
		seenStates[normalizedSeedKey(estadoNombre)] = struct{}{}

		seenMunicipios := make(map[string]struct{}, len(estado.Municipios))
		for _, municipio := range estado.Municipios {
			municipioNombre := cleanSeedName(municipio.Nombre)
			if municipioNombre == "" {
				return nil, fmt.Errorf("el estado %s tiene un municipio sin nombre", estadoNombre)
			}

			municipioKey := normalizedSeedKey(municipioNombre)
			if _, exists := seenMunicipios[municipioKey]; exists {
				return nil, fmt.Errorf("municipio duplicado en %s: %s", estadoNombre, municipioNombre)
			}
			seenMunicipios[municipioKey] = struct{}{}

			if len(municipio.Ciudades) == 0 {
				return nil, fmt.Errorf("%s > %s no tiene ciudades", estadoNombre, municipioNombre)
			}

			seenParroquias := map[string]struct{}{}
			for _, parroquia := range municipio.Parroquias {
				parroquiaNombre := cleanSeedName(parroquia)
				if parroquiaNombre == "" {
					return nil, fmt.Errorf("%s > %s tiene una parroquia vacia", estadoNombre, municipioNombre)
				}

				parroquiaKey := normalizedSeedKey(parroquiaNombre)
				if _, exists := seenParroquias[parroquiaKey]; exists {
					return nil, fmt.Errorf("%s > %s tiene la parroquia duplicada %s", estadoNombre, municipioNombre, parroquiaNombre)
				}
				seenParroquias[parroquiaKey] = struct{}{}
			}

			seenCiudades := map[string]struct{}{}
			for index, ciudad := range municipio.Ciudades {
				ciudadNombre := cleanSeedName(ciudad)
				if ciudadNombre == "" {
					return nil, fmt.Errorf("%s > %s tiene una ciudad vacia", estadoNombre, municipioNombre)
				}

				ciudadKey := normalizedSeedKey(ciudadNombre)
				if _, exists := seenCiudades[ciudadKey]; exists {
					return nil, fmt.Errorf("%s > %s tiene la ciudad duplicada %s", estadoNombre, municipioNombre, ciudadNombre)
				}
				seenCiudades[ciudadKey] = struct{}{}

				if index == 0 && ciudadKey == municipioKey {
					if _, parroquiaMatch := seenParroquias[ciudadKey]; !parroquiaMatch {
						warnings = append(
							warnings,
							fmt.Sprintf("%s > %s usa una ciudad principal homonima al municipio; revisar si esta documentada", estadoNombre, municipioNombre),
						)
					}
				}
			}
		}
	}

	return warnings, nil
}

func cleanSeedName(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func normalizedSeedKey(value string) string {
	return strings.ToLower(cleanSeedName(value))
}

func firstOrCreateEstado(tx *gorm.DB, estado *models.Estado) (bool, error) {
	return firstOrCreate(tx, estado, models.Estado{Nombre: estado.Nombre})
}

func firstOrCreateMunicipio(tx *gorm.DB, municipio *models.Municipio) (bool, error) {
	return firstOrCreate(tx, municipio, models.Municipio{EstadoID: municipio.EstadoID, Nombre: municipio.Nombre})
}

func firstOrCreateParroquia(tx *gorm.DB, parroquia *models.Parroquia) (bool, error) {
	return firstOrCreate(tx, parroquia, models.Parroquia{MunicipioID: parroquia.MunicipioID, Nombre: parroquia.Nombre})
}

func firstOrCreateCiudad(tx *gorm.DB, ciudad *models.Ciudad) (bool, error) {
	return firstOrCreate(tx, ciudad, models.Ciudad{MunicipioID: ciudad.MunicipioID, Nombre: ciudad.Nombre})
}

func firstOrCreate[T any](tx *gorm.DB, target *T, where T) (bool, error) {
	result := tx.Where(where).FirstOrCreate(target)
	return result.RowsAffected == 1, result.Error
}
