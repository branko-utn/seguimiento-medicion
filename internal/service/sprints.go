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
	historias RepoHistorias // para asignar historias al sprint (HU-11)
}

func NuevoSprints(sprints RepoSprints, proyectos RepoProyectos, historias RepoHistorias) *Sprints {
	return &Sprints{sprints: sprints, proyectos: proyectos, historias: historias}
}

// Crear da de alta un sprint en un proyecto (HU-10). Valida en este orden: que el proyecto exista, el Sprint Goal y
// las fechas (sprint.Nuevo), que las fechas estén dentro del proyecto y que no se superponga con otro sprint.
func (s *Sprints) Crear(proyectoID domain.ProyectoID, goal string, inicio, fin time.Time) (sprint.Sprint, error) {
	p, err := s.proyectos.Buscar(proyectoID)
	if err != nil {
		return sprint.Sprint{}, err
	}
	sp, err := sprint.Nuevo(proyectoID, goal, inicio, fin)
	if err != nil {
		return sprint.Sprint{}, err
	}
	if err := sp.ValidarFechasEnProyecto(p); err != nil {
		return sprint.Sprint{}, err
	}
	if err := sp.ValidarSinSuperposicion(s.sprints.DelProyecto(proyectoID)); err != nil {
		return sprint.Sprint{}, err
	}
	return s.sprints.Agregar(sp)
}

// Iniciar pasa un sprint a Activo (HU-10). Primero el sprint valida su propio estado (RN-10.6) y después se controla
// que no haya otro sprint activo en el proyecto (RN-10.5).
func (s *Sprints) Iniciar(id domain.SprintID) error {
	sp, err := s.sprints.Buscar(id)
	if err != nil {
		return err
	}
	if err := sp.Iniciar(); err != nil {
		return err
	}
	for _, otro := range s.sprints.DelProyecto(sp.ProyectoID) {
		if otro.ID != sp.ID && otro.Estado == sprint.Activo {
			return sprint.ErrOtroActivo
		}
	}
	return s.sprints.Actualizar(sp)
}

// Asignar agrega una historia al Sprint Backlog (HU-11). Todavía sin implementar (RED).
func (s *Sprints) Asignar(sprintID domain.SprintID, historiaID domain.HistoriaID) error {
	return nil
}

// Buscar devuelve un sprint, o sprint.ErrNoExiste.
func (s *Sprints) Buscar(id domain.SprintID) (sprint.Sprint, error) {
	return s.sprints.Buscar(id)
}

// DelProyecto devuelve los sprints de un proyecto en el orden en que se crearon, o project.ErrNoExiste.
func (s *Sprints) DelProyecto(proyectoID domain.ProyectoID) ([]sprint.Sprint, error) {
	if _, err := s.proyectos.Buscar(proyectoID); err != nil {
		return nil, err
	}
	return s.sprints.DelProyecto(proyectoID), nil
}
