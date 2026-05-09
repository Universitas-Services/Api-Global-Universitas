package database

import "testing"

func TestLoadTerritorySeeds_ProjectSeedUsesStringParishesAndCities(t *testing.T) {
	estados, warnings, err := loadTerritorySeeds("seeds/venezuela.json")
	if err != nil {
		t.Fatalf("loadTerritorySeeds returned error: %v", err)
	}
	if len(estados) != 24 {
		t.Fatalf("expected 24 estados, got %d", len(estados))
	}
	if len(warnings) == 0 {
		t.Log("seed loaded without homonym warnings")
	}
}

func TestLoadTerritorySeeds_KnownMunicipiosHaveExpectedPrimaryCities(t *testing.T) {
	estados, _, err := loadTerritorySeeds("seeds/venezuela.json")
	if err != nil {
		t.Fatalf("loadTerritorySeeds returned error: %v", err)
	}

	assertPrimaryCity(t, estados, "Vargas", "Vargas", "La Guaira")
	assertPrimaryCity(t, estados, "Bolívar", "Sifontes", "Tumeremo")
	assertPrimaryCity(t, estados, "Sucre", "Sucre", "Cumaná")
}

func TestValidateTerritorySeeds_RejectsMunicipioWithoutCities(t *testing.T) {
	_, err := validateTerritorySeeds([]estadoSeed{
		{
			Nombre: "Prueba",
			Municipios: []municipioSeed{
				{
					Nombre:     "Demo",
					Parroquias: []string{"Centro"},
					Ciudades:   nil,
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected error for municipio without ciudades")
	}
}

func assertPrimaryCity(t *testing.T, estados []estadoSeed, estadoNombre, municipioNombre, expected string) {
	t.Helper()

	for _, estado := range estados {
		if estado.Nombre != estadoNombre {
			continue
		}
		for _, municipio := range estado.Municipios {
			if municipio.Nombre != municipioNombre {
				continue
			}
			if len(municipio.Ciudades) == 0 {
				t.Fatalf("%s > %s has no ciudades", estadoNombre, municipioNombre)
			}
			if municipio.Ciudades[0] != expected {
				t.Fatalf("%s > %s expected primary city %q, got %q", estadoNombre, municipioNombre, expected, municipio.Ciudades[0])
			}
			return
		}
	}

	t.Fatalf("municipio not found: %s > %s", estadoNombre, municipioNombre)
}
