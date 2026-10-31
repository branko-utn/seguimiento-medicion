package memory

import (
	"slices"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// Sprints guarda los sprints de todos los proyectos (épica E3) en un archivo JSON.
type Sprints struct {
	store *Store[[]sprint.Sprint]
}

// NuevoSprints abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoSprints(ruta string) (*Sprints, error) {
	s, err := Open[[]sprint.Sprint](ruta)
	if err != nil {
		return nil, err
	}
	return &Sprints{store: s}, nil
}

// Agregar guarda un sprint nuevo con el siguiente ID (S1, S2, …).
func (r *Sprints) Agregar(s sprint.Sprint) (sprint.Sprint, error) {
	s = copiarSprint(s)
	err := r.store.Update(func(lista *[]sprint.Sprint) error {
		s.ID = domain.SprintID("S" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, s)
		return nil
	})
	if err != nil {
		return sprint.Sprint{}, err
	}
	return copiarSprint(s), nil
}

// Buscar devuelve una copia del sprint con ese ID, o sprint.ErrNoExiste.
func (r *Sprints) Buscar(id domain.SprintID) (sprint.Sprint, error) {
	var (
		encontrado sprint.Sprint
		existe     bool
	)
	r.store.View(func(lista []sprint.Sprint) {
		for _, s := range lista {
			if s.ID == id {
				encontrado, existe = copiarSprint(s), true
				return
			}
		}
	})
	if !existe {
		return sprint.Sprint{}, sprint.ErrNoExiste
	}
	return encontrado, nil
}

// Actualizar reemplaza un sprint ya guardado, o devuelve sprint.ErrNoExiste.
func (r *Sprints) Actualizar(s sprint.Sprint) error {
	return r.store.Update(func(lista *[]sprint.Sprint) error {
		for n := range *lista {
			if (*lista)[n].ID == s.ID {
				(*lista)[n] = copiarSprint(s)
				return nil
			}
		}
		return sprint.ErrNoExiste
	})
}

// DelProyecto devuelve una copia de los sprints de un proyecto, en el orden en que se crearon.
func (r *Sprints) DelProyecto(id domain.ProyectoID) []sprint.Sprint {
	var delProyecto []sprint.Sprint
	r.store.View(func(lista []sprint.Sprint) {
		for _, s := range lista {
			if s.ProyectoID == id {
				delProyecto = append(delProyecto, copiarSprint(s))
			}
		}
	})
	return delProyecto
}

// Listar devuelve todos los sprints.
func (r *Sprints) Listar() []sprint.Sprint {
	var lista []sprint.Sprint
	r.store.View(func(l []sprint.Sprint) {
		for _, s := range l {
			lista = append(lista, copiarSprint(s))
		}
	})
	return lista
}

// copiarSprint copia también las listas de historias, así el sprint guardado y el que se devuelve no las comparten.
func copiarSprint(s sprint.Sprint) sprint.Sprint {
	s.Historias = slices.Clone(s.Historias)
	s.Completadas = slices.Clone(s.Completadas)
	return s
}
