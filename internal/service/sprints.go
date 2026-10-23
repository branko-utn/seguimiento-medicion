package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// RepoSprints es lo que el servicio de sprints necesita del almacenamiento.
type RepoSprints interface {
	// Agregar guarda un sprint nuevo y lo devuelve con el ID generado.
	Agregar(s sprint.Sprint) (sprint.Sprint, error)
	// Buscar y Actualizar devuelven sprint.ErrNoExiste si no hay un sprint con ese ID.
	Buscar(id domain.SprintID) (sprint.Sprint, error)
	Actualizar(s sprint.Sprint) error
	// DelProyecto devuelve los sprints de un proyecto en el orden en que se crearon.
	DelProyecto(id domain.ProyectoID) []sprint.Sprint
}

// Sprints son los casos de uso de la épica E3.
type Sprints struct {
	sprints   RepoSprints
	proyectos RepoProyectos
}

func NuevoSprints(sprints RepoSprints, proyectos RepoProyectos) *Sprints {
	return &Sprints{sprints: sprints, proyectos: proyectos}
}

// Crear da de alta un sprint en un proyecto (HU-10). El stub existe solo para que los tests de la fase RED compilen.
func (s *Sprints) Crear(proyectoID domain.ProyectoID, goal string, inicio, fin time.Time) (sprint.Sprint, error) {
	return sprint.Sprint{}, nil
}

// Iniciar pasa un sprint a Activo (HU-10). El stub existe solo para la fase RED.
func (s *Sprints) Iniciar(id domain.SprintID) error {
	return nil
}

// Buscar devuelve un sprint, o sprint.ErrNoExiste. El stub existe solo para la fase RED.
func (s *Sprints) Buscar(id domain.SprintID) (sprint.Sprint, error) {
	return sprint.Sprint{}, nil
}

// DelProyecto devuelve los sprints de un proyecto, o project.ErrNoExiste. El stub existe solo para la fase RED.
func (s *Sprints) DelProyecto(proyectoID domain.ProyectoID) ([]sprint.Sprint, error) {
	return nil, nil
}
