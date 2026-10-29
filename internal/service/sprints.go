package service

import (
	"slices"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
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

// Asignar agrega una historia al Sprint Backlog (HU-11). Busca el sprint y la historia, y valida en el orden de la
// spec: primero Sprint.Asignar (RN-11.1, 11.2, 11.4, 11.5 y 11.6) y después Historia.AsignarASprint (RN-11.3), así
// una historia completada informa "ya está completada" y no "ya está asignada". Guarda los dos solo si ambas
// validaciones pasan.
func (s *Sprints) Asignar(sprintID domain.SprintID, historiaID domain.HistoriaID) error {
	sp, err := s.sprints.Buscar(sprintID)
	if err != nil {
		return err
	}
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return err
	}
	if err := sp.Asignar(h); err != nil {
		return err
	}
	if err := h.AsignarASprint(sp.ID); err != nil {
		return err
	}
	if err := s.sprints.Actualizar(sp); err != nil {
		return err
	}
	return s.historias.Guardar(h)
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

// DetalleSprint es un sprint cerrado con las historias que tuvo (HU-14).
type DetalleSprint struct {
	Sprint    sprint.Sprint
	Historias []HistoriaDetalle
}

// HistoriaDetalle es una historia del sprint en el detalle de HU-14.
type HistoriaDetalle struct {
	ID          domain.HistoriaID
	Titulo      string
	StoryPoints int
	Completada  bool
}

// Cerrados devuelve los sprints cerrados de un proyecto (HU-14).
func (s *Sprints) Cerrados(proyectoID domain.ProyectoID) ([]sprint.Sprint, error) { return nil, nil }

// Detalle devuelve un sprint cerrado con sus historias (HU-14).
func (s *Sprints) Detalle(sprintID domain.SprintID) (DetalleSprint, error) {
	return DetalleSprint{}, nil
}

// CompletarHistoria registra una historia del sprint como completada (HU-12). Busca el sprint y la historia, y valida
// en orden: primero Sprint.Completar (RN-12.1 a RN-12.3) y después Historia.Completar (RN-12.4, la regla de HU-08).
// Guarda los dos solo si ambas validaciones pasan; los SP completados (RN-12.5) se calculan con las métricas (HU-24).
func (s *Sprints) CompletarHistoria(sprintID domain.SprintID, historiaID domain.HistoriaID) error {
	sp, err := s.sprints.Buscar(sprintID)
	if err != nil {
		return err
	}
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return err
	}
	if err := sp.Completar(h.ID); err != nil {
		return err
	}
	if err := h.Completar(); err != nil {
		return err
	}
	if err := s.sprints.Actualizar(sp); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}

// CerrarSprint cierra un sprint (HU-13): fija los SP planificados y completados, calculados con las métricas de HU-24
// (RN-13.2), y devuelve al backlog las historias que no se completaron (RN-13.3). Valida el cierre antes de tocar
// las historias: si el sprint no se puede cerrar, no cambia nada.
func (s *Sprints) CerrarSprint(sprintID domain.SprintID) error {
	sp, err := s.sprints.Buscar(sprintID)
	if err != nil {
		return err
	}
	historias := make([]backlog.Historia, 0, len(sp.Historias))
	entradas := make([]metrics.HistoriaSprint, 0, len(sp.Historias))
	for _, id := range sp.Historias {
		h, err := s.historias.Buscar(id)
		if err != nil {
			return err
		}
		historias = append(historias, h)
		entradas = append(entradas, metrics.HistoriaSprint{StoryPoints: h.StoryPoints, Completada: slices.Contains(sp.Completadas, id)})
	}
	if err := sp.Cerrar(metrics.SPPlanificados(entradas), metrics.SPCompletados(entradas)); err != nil {
		return err
	}
	if err := s.sprints.Actualizar(sp); err != nil {
		return err
	}
	for _, h := range historias {
		if slices.Contains(sp.Completadas, h.ID) {
			continue
		}
		h.VolverAlBacklog()
		if err := s.historias.Guardar(h); err != nil {
			return err
		}
	}
	return nil
}
