package backlog

import (
	"reflect"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// conPrioridad arma una historia con solo el ID y la prioridad, lo único que mira Ordenar.
func conPrioridad(id string, p Prioridad) Historia {
	return Historia{ID: domain.HistoriaID(id), Prioridad: p}
}

// idsDe devuelve los IDs de las historias, en orden.
func idsDe(hs []Historia) []domain.HistoriaID {
	ids := []domain.HistoriaID{}
	for _, h := range hs {
		ids = append(ids, h.ID)
	}
	return ids
}

func TestOrdenar(t *testing.T) {
	casos := []struct {
		nombre    string
		historias []Historia
		want      []domain.HistoriaID
	}{
		{
			// RN-09.1: Alta, Media, Baja y, a igual prioridad, por identificador.
			nombre: "por prioridad y despues por identificador",
			historias: []Historia{
				conPrioridad("HU-1", PrioridadBaja), conPrioridad("HU-2", PrioridadAlta),
				conPrioridad("HU-3", PrioridadMedia), conPrioridad("HU-4", PrioridadAlta),
			},
			want: []domain.HistoriaID{"HU-2", "HU-4", "HU-3", "HU-1"},
		},
		{
			// RN-09.1: el identificador se compara como número, HU-2 antes que HU-10.
			nombre: "identificador numerico",
			historias: []Historia{
				conPrioridad("HU-10", PrioridadAlta), conPrioridad("HU-1", PrioridadMedia),
				conPrioridad("HU-2", PrioridadAlta),
			},
			want: []domain.HistoriaID{"HU-2", "HU-10", "HU-1"},
		},
		{
			nombre: "todas con la misma prioridad",
			historias: []Historia{
				conPrioridad("HU-12", PrioridadMedia), conPrioridad("HU-3", PrioridadMedia),
				conPrioridad("HU-7", PrioridadMedia),
			},
			want: []domain.HistoriaID{"HU-3", "HU-7", "HU-12"},
		},
		{nombre: "backlog vacio", historias: nil, want: []domain.HistoriaID{}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := idsDe(Ordenar(c.historias)); !reflect.DeepEqual(got, c.want) {
				t.Errorf("orden = %v, se esperaba %v", got, c.want)
			}
		})
	}
}

func TestOrdenar_noModificaLaEntrada(t *testing.T) {
	entrada := []Historia{conPrioridad("HU-1", PrioridadBaja), conPrioridad("HU-2", PrioridadAlta)}

	Ordenar(entrada)

	// Restricción de la spec: función pura que devuelve una copia ordenada.
	if got := idsDe(entrada); !reflect.DeepEqual(got, []domain.HistoriaID{"HU-1", "HU-2"}) {
		t.Errorf("la entrada quedó en el orden %v, no se tenía que modificar", got)
	}
}
