package memory

import (
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

// Agregar guarda un sprint nuevo con el siguiente ID (S1, S2, …). El stub existe solo para la fase RED.
func (r *Sprints) Agregar(s sprint.Sprint) (sprint.Sprint, error) {
	return sprint.Sprint{}, nil
}

// Buscar devuelve una copia del sprint con ese ID, o sprint.ErrNoExiste. El stub existe solo para la fase RED.
func (r *Sprints) Buscar(id domain.SprintID) (sprint.Sprint, error) {
	return sprint.Sprint{}, nil
}

// Actualizar reemplaza un sprint ya guardado, o devuelve sprint.ErrNoExiste. El stub existe solo para la fase RED.
func (r *Sprints) Actualizar(s sprint.Sprint) error {
	return nil
}

// DelProyecto devuelve una copia de los sprints de un proyecto, en el orden en que se crearon.
// El stub existe solo para la fase RED.
func (r *Sprints) DelProyecto(id domain.ProyectoID) []sprint.Sprint {
	return nil
}
