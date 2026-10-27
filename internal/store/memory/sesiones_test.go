package memory

import (
	"strconv"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
)

// votacionDeArma una sesión de la historia con un voto por carta indicada, para preparar el caso.
// Los votos no se exportan (RN-16.4), así que la sesión se arma con sus métodos, no con literales.
func votacionDe(t *testing.T, historia domain.HistoriaID, cartas ...int) estimation.Sesion {
	t.Helper()
	s, err := estimation.NuevaSesion(historia)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	for i, carta := range cartas {
		if err := s.Votar(domain.IntegranteID("I"+strconv.Itoa(i+1)), carta); err != nil {
			t.Fatalf("voto %d: no se esperaba error: %v", i+1, err)
		}
	}
	return s
}

func TestSesionesGuardarBuscar(t *testing.T) {
	casos := []struct {
		titulo         string
		guardar        []domain.HistoriaID // historias de las sesiones que se guardan, en orden
		cartas         [][]int             // votos de cada sesión guardada
		buscar         domain.HistoriaID
		quiereOk       bool
		quiereVotantes int
	}{
		{
			titulo:         "guarda una votación",
			guardar:        []domain.HistoriaID{"HU-1"},
			cartas:         [][]int{{5, 8}},
			buscar:         domain.HistoriaID("HU-1"),
			quiereOk:       true,
			quiereVotantes: 2,
		},
		{
			titulo:         "reemplaza la votación de la misma historia",
			guardar:        []domain.HistoriaID{"HU-1", "HU-1"},
			cartas:         [][]int{{5}, {3, 13}},
			buscar:         domain.HistoriaID("HU-1"),
			quiereOk:       true,
			quiereVotantes: 2,
		},
		{
			titulo:         "historias distintas conviven",
			guardar:        []domain.HistoriaID{"HU-1", "HU-2"},
			cartas:         [][]int{{5}, {3, 13}},
			buscar:         domain.HistoriaID("HU-2"),
			quiereOk:       true,
			quiereVotantes: 2,
		},
		{
			titulo:   "buscar inexistente devuelve ok falso",
			guardar:  []domain.HistoriaID{"HU-1"},
			cartas:   [][]int{{5}},
			buscar:   domain.HistoriaID("HU-99"),
			quiereOk: false,
		},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			r, err := NuevoSesiones("")
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			for n, historia := range c.guardar {
				if err := r.Guardar(votacionDe(t, historia, c.cartas[n]...)); err != nil {
					t.Fatalf("no se esperaba error al guardar: %v", err)
				}
			}

			got, ok := r.Buscar(c.buscar)
			if ok != c.quiereOk {
				t.Fatalf("ok esperado %v, obtenido %v", c.quiereOk, ok)
			}
			if !c.quiereOk {
				return
			}
			if got.HistoriaID != c.buscar {
				t.Fatalf("HistoriaID esperada %s, obtenida %s", c.buscar, got.HistoriaID)
			}
			if got.Ronda != 1 {
				t.Fatalf("Ronda esperada 1, obtenida %d", got.Ronda)
			}
			if votantes := got.QuienesVotaron(); len(votantes) != c.quiereVotantes {
				t.Fatalf("se esperaban %d votantes, obtenidos %d", c.quiereVotantes, len(votantes))
			}
		})
	}
}
