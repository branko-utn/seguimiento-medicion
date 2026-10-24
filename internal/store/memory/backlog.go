package memory

import (
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// Historias guarda el Product Backlog (épica E2) en un archivo JSON.
type Historias struct {
	store *Store[[]backlog.Historia]
}

// NuevoHistorias abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoHistorias(ruta string) (*Historias, error) {
	s, err := Open[[]backlog.Historia](ruta)
	if err != nil {
		return nil, err
	}
	return &Historias{store: s}, nil
}

// Agregar guarda una historia nueva con el siguiente ID de su proyecto (RN-05.1: HU-1, HU-2, …; cada proyecto
// numera su backlog desde HU-1).
func (r *Historias) Agregar(h backlog.Historia) (backlog.Historia, error) {
	err := r.store.Update(func(lista *[]backlog.Historia) error {
		h.ID = domain.HistoriaID("HU-" + strconv.Itoa(cantidadDe(*lista, h.ProyectoID)+1))
		*lista = append(*lista, h)
		return nil
	})
	if err != nil {
		return backlog.Historia{}, err
	}
	return h, nil
}

// Listar devuelve una copia de todas las historias, en el orden en que se crearon.
func (r *Historias) Listar() []backlog.Historia {
	var copia []backlog.Historia
	r.store.View(func(lista []backlog.Historia) {
		copia = append(copia, lista...)
	})
	return copia
}

// cantidadDe cuenta las historias de un proyecto: cada proyecto numera su backlog desde HU-1 (RN-05.1).
func cantidadDe(historias []backlog.Historia, proyectoID domain.ProyectoID) int {
	cantidad := 0
	for _, h := range historias {
		if h.ProyectoID == proyectoID {
			cantidad++
		}
	}
	return cantidad
}
