package memory

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

func registro(integrante domain.IntegranteID, historia domain.HistoriaID, horas float64) effort.Registro {
	return effort.Registro{
		IntegranteID: integrante,
		HistoriaID:   historia,
		Fecha:        time.Date(2026, time.October, 24, 0, 0, 0, 0, time.Local),
		Actividad:    "Desarrollo",
		Horas:        horas,
	}
}

func TestRegistros_AgregarGeneraIDsYFiltra(t *testing.T) {
	r, err := NuevoRegistros("")
	if err != nil {
		t.Fatal(err)
	}
	altas := []struct {
		reg      effort.Registro
		quiereID domain.RegistroID
	}{
		{registro("I1", "HU-1", 3), "E1"},
		{registro("I2", "HU-1", 2), "E2"},
		{registro("I1", "HU-2", 1), "E3"},
	}
	for _, a := range altas {
		guardado, err := r.Agregar(a.reg)
		if err != nil {
			t.Fatalf("Agregar: %v", err)
		}
		if guardado.ID != a.quiereID {
			t.Errorf("ID = %q, se esperaba %q", guardado.ID, a.quiereID)
		}
	}

	casos := []struct {
		titulo string
		lista  []effort.Registro
		ids    []domain.RegistroID
	}{
		{"DelIntegrante(I1)", r.DelIntegrante("I1"), []domain.RegistroID{"E1", "E3"}},
		{"DelIntegrante(I9)", r.DelIntegrante("I9"), nil},
		{"DeHistoria(HU-1)", r.DeHistoria("HU-1"), []domain.RegistroID{"E1", "E2"}},
		{"DeHistoria(HU-99)", r.DeHistoria("HU-99"), nil},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			if len(c.lista) != len(c.ids) {
				t.Fatalf("devolvió %d registros, se esperaban %d", len(c.lista), len(c.ids))
			}
			for i, reg := range c.lista {
				if reg.ID != c.ids[i] {
					t.Errorf("registro %d: ID = %q, se esperaba %q (en el orden en que se registraron)", i, reg.ID, c.ids[i])
				}
			}
		})
	}
}

func TestRegistros_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "registros.json")
	r, err := NuevoRegistros(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(registro("I1", "HU-1", 3)); err != nil {
		t.Fatal(err)
	}

	otra, err := NuevoRegistros(ruta)
	if err != nil {
		t.Fatal(err)
	}
	lista := otra.DeHistoria("HU-1")
	if len(lista) != 1 || lista[0].ID != "E1" || lista[0].Horas != 3 || lista[0].Actividad != "Desarrollo" {
		t.Errorf("al reabrir: %+v, se esperaba el registro E1 de 3 horas", lista)
	}
}
