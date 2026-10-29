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

// Agregar guarda una historia nueva con el siguiente ID (RN-05.1: HU-1, HU-2, … únicos en todo el sistema).
func (r *Historias) Agregar(h backlog.Historia) (backlog.Historia, error) {
	h.Criterios = clonarCriterios(h.Criterios)
	err := r.store.Update(func(lista *[]backlog.Historia) error {
		h.ID = domain.HistoriaID("HU-" + strconv.Itoa(len(*lista)+1))
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
		for _, h := range lista {
			h.Criterios = clonarCriterios(h.Criterios)
			copia = append(copia, h)
		}
	})
	return copia
}

// Buscar devuelve una copia de la historia con ese ID, o backlog.ErrNoExiste. Los cambios a la copia no afectan lo
// guardado hasta llamar a Guardar.
func (r *Historias) Buscar(id domain.HistoriaID) (backlog.Historia, error) {
	var (
		encontrada backlog.Historia
		ok         bool
	)
	r.store.View(func(lista []backlog.Historia) {
		for _, h := range lista {
			if h.ID == id {
				encontrada, ok = h, true
				encontrada.Criterios = clonarCriterios(h.Criterios)
				return
			}
		}
	})
	if !ok {
		return backlog.Historia{}, backlog.ErrNoExiste
	}
	return encontrada, nil
}

// Guardar reemplaza una historia ya guardada (por su ID), o devuelve backlog.ErrNoExiste. Lo usan los servicios que
// modifican historias: criterios (HU-07), estado (HU-08), asignación a sprint (HU-11), estimación (HU-15).
func (r *Historias) Guardar(h backlog.Historia) error {
	h.Criterios = clonarCriterios(h.Criterios)
	return r.store.Update(func(lista *[]backlog.Historia) error {
		for i := range *lista {
			if (*lista)[i].ID == h.ID {
				(*lista)[i] = h
				return nil
			}
		}
		return backlog.ErrNoExiste
	})
}

// Estimar asigna Story Points a una historia (HU-15, HU-18). Busca la historia, aplica Estimar del dominio
// y guarda el resultado.
func (r *Historias) Estimar(id domain.HistoriaID, sp int) error {
	h, err := r.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.Estimar(sp); err != nil {
		return err
	}
	return r.Guardar(h)
}

// clonarCriterios copia la lista para que la historia guardada no comparta memoria con la que se devuelve.
func clonarCriterios(c []backlog.Criterio) []backlog.Criterio {
	if c == nil {
		return nil
	}
	return append([]backlog.Criterio(nil), c...)
}
