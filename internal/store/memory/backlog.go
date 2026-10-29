package memory

import (
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// Historias guarda el Product Backlog (épica E2) en un archivo JSON.
type Historias struct {
	store *Store[[]historiaGuardada]
}

// historiaGuardada es como se guarda cada historia en el archivo. Una historia eliminada (HU-06) queda guardada con
// la marca Eliminada: así su ID no se reutiliza (RN-06.4), porque el siguiente ID se calcula con el largo de la
// lista. Como Historia va embebida, el JSON es el mismo de antes más el campo Eliminada (solo si es true), y los
// archivos guardados antes de HU-06 se siguen leyendo.
type historiaGuardada struct {
	backlog.Historia
	Eliminada bool `json:",omitempty"`
}

// NuevoHistorias abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoHistorias(ruta string) (*Historias, error) {
	s, err := Open[[]historiaGuardada](ruta)
	if err != nil {
		return nil, err
	}
	return &Historias{store: s}, nil
}

// Agregar guarda una historia nueva con el siguiente ID (RN-05.1: HU-1, HU-2, … únicos en todo el sistema). Las
// eliminadas siguen en la lista, así que su ID no se vuelve a usar (RN-06.4).
func (r *Historias) Agregar(h backlog.Historia) (backlog.Historia, error) {
	h.Criterios = clonarCriterios(h.Criterios)
	err := r.store.Update(func(lista *[]historiaGuardada) error {
		h.ID = domain.HistoriaID("HU-" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, historiaGuardada{Historia: h})
		return nil
	})
	if err != nil {
		return backlog.Historia{}, err
	}
	return h, nil
}

// Listar devuelve una copia de todas las historias no eliminadas, en el orden en que se crearon.
func (r *Historias) Listar() []backlog.Historia {
	var copia []backlog.Historia
	r.store.View(func(lista []historiaGuardada) {
		for _, g := range lista {
			if g.Eliminada {
				continue
			}
			h := g.Historia
			h.Criterios = clonarCriterios(h.Criterios)
			copia = append(copia, h)
		}
	})
	return copia
}

// Buscar devuelve una copia de la historia con ese ID, o backlog.ErrNoExiste (también si fue eliminada). Los cambios
// a la copia no afectan lo guardado hasta llamar a Guardar.
func (r *Historias) Buscar(id domain.HistoriaID) (backlog.Historia, error) {
	var (
		encontrada backlog.Historia
		ok         bool
	)
	r.store.View(func(lista []historiaGuardada) {
		if i := indiceHistoria(lista, id); i >= 0 {
			encontrada, ok = lista[i].Historia, true
			encontrada.Criterios = clonarCriterios(encontrada.Criterios)
		}
	})
	if !ok {
		return backlog.Historia{}, backlog.ErrNoExiste
	}
	return encontrada, nil
}

// Guardar reemplaza una historia ya guardada (por su ID), o devuelve backlog.ErrNoExiste (también si fue eliminada).
// Lo usan los servicios que modifican historias: edición (HU-06), criterios (HU-07), estado (HU-08), asignación a
// sprint (HU-11), estimación (HU-15).
func (r *Historias) Guardar(h backlog.Historia) error {
	h.Criterios = clonarCriterios(h.Criterios)
	return r.store.Update(func(lista *[]historiaGuardada) error {
		i := indiceHistoria(*lista, h.ID)
		if i < 0 {
			return backlog.ErrNoExiste
		}
		(*lista)[i].Historia = h
		return nil
	})
}

// Eliminar saca una historia del backlog (HU-06), o devuelve backlog.ErrNoExiste (también si ya fue eliminada). La
// historia queda marcada como eliminada en el archivo para no liberar su ID (RN-06.4). Las reglas para poder
// eliminarla (RN-06.3 y RN-06.5) las valida el servicio antes de llamar acá.
func (r *Historias) Eliminar(id domain.HistoriaID) error {
	return r.store.Update(func(lista *[]historiaGuardada) error {
		i := indiceHistoria(*lista, id)
		if i < 0 {
			return backlog.ErrNoExiste
		}
		(*lista)[i].Eliminada = true
		return nil
	})
}

// indiceHistoria devuelve la posición de la historia no eliminada con ese ID, o -1 si no está.
func indiceHistoria(lista []historiaGuardada, id domain.HistoriaID) int {
	for i, g := range lista {
		if g.ID == id && !g.Eliminada {
			return i
		}
	}
	return -1
}

// clonarCriterios copia la lista para que la historia guardada no comparta memoria con la que se devuelve.
func clonarCriterios(c []backlog.Criterio) []backlog.Criterio {
	if c == nil {
		return nil
	}
	return append([]backlog.Criterio(nil), c...)
}
