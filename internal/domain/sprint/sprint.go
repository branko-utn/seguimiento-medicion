package sprint

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// Estado es la situación de un sprint (ADR 0003: el cero no es válido).
type Estado int

const (
	Planificado Estado = iota + 1
	Activo
	Cerrado
)

// String devuelve el texto que usan las specs y los escenarios BDD.
func (e Estado) String() string {
	switch e {
	case Planificado:
		return "Planificado"
	case Activo:
		return "Activo"
	case Cerrado:
		return "Cerrado"
	}
	return ""
}

// Sprint es una iteración de un proyecto (ADR 0002).
type Sprint struct {
	ID          domain.SprintID
	ProyectoID  domain.ProyectoID
	Goal        string
	Inicio      time.Time
	Fin         time.Time
	Estado      Estado
	Historias   []domain.HistoriaID // asignadas (HU-11)
	Completadas []domain.HistoriaID // completadas en este sprint (HU-12)
	// Al cerrar (HU-13) quedan fijos:
	SPPlanificados, SPCompletados int
}

var (
	ErrGoalVacio        = errors.New("el Sprint Goal es obligatorio")
	ErrFechasInvalidas  = errors.New("la fecha de fin no puede ser anterior a la de inicio")
	ErrFueraDelProyecto = errors.New("las fechas del sprint deben estar dentro de las del proyecto")
	ErrSuperpuesto      = errors.New("las fechas se superponen con otro sprint del proyecto")
	ErrOtroActivo       = errors.New("ya hay un sprint activo en el proyecto")
	ErrNoPlanificado    = errors.New("solo se puede iniciar un sprint planificado")
	ErrNoExiste         = errors.New("el sprint no existe")

	// Errores de HU-11 (asignar historias).
	ErrHistoriaSinEstimar   = errors.New("la historia no está estimada")
	ErrHistoriaSinCriterios = errors.New("la historia no tiene criterios de aceptación")
	ErrSprintCerrado        = errors.New("el sprint está cerrado y no se puede modificar")
	ErrHistoriaCompletada   = errors.New("la historia ya está completada")
	ErrOtroProyecto         = errors.New("la historia pertenece a otro proyecto")

	// Errores de HU-12 (completar historias).
	ErrHistoriaNoAsignada = errors.New("la historia no está asignada a este sprint")
)

// Nuevo crea un sprint válido en estado Planificado y sin historias (RN-10.1, RN-10.2 y RN-10.5).
// Que las fechas estén dentro del proyecto y no se superpongan con otro sprint lo coordina el servicio.
func Nuevo(proyectoID domain.ProyectoID, goal string, inicio, fin time.Time) (Sprint, error) {
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return Sprint{}, ErrGoalVacio
	}
	if fin.Before(inicio) {
		return Sprint{}, ErrFechasInvalidas
	}
	return Sprint{
		ProyectoID: proyectoID,
		Goal:       goal,
		Inicio:     inicio,
		Fin:        fin,
		Estado:     Planificado,
	}, nil
}

// ValidarFechasEnProyecto controla que el sprint esté dentro de las fechas del proyecto, inicio y fin incluidos (RN-10.3).
func (s Sprint) ValidarFechasEnProyecto(p project.Proyecto) error {
	if s.Inicio.Before(p.Inicio) || s.Fin.After(p.Fin) {
		return ErrFueraDelProyecto
	}
	return nil
}

// ValidarSinSuperposicion controla que el sprint no comparta ningún día con los otros sprints del proyecto (RN-10.4).
func (s Sprint) ValidarSinSuperposicion(otros []Sprint) error {
	for _, otro := range otros {
		if !s.Inicio.After(otro.Fin) && !otro.Inicio.After(s.Fin) {
			return ErrSuperpuesto
		}
	}
	return nil
}

// Iniciar pasa el sprint a Activo. Solo se inicia un sprint Planificado (RN-10.6); que no haya otro activo en el
// proyecto (RN-10.5) lo controla el servicio.
func (s *Sprint) Iniciar() error {
	if s.Estado != Planificado {
		return ErrNoPlanificado
	}
	s.Estado = Activo
	return nil
}

// Asignar agrega la historia al Sprint Backlog (HU-11). Valida primero el sprint y después la historia: que el
// sprint no esté cerrado (RN-11.4), que la historia sea del mismo proyecto (RN-11.6), que no esté completada
// (RN-11.5), que esté estimada (RN-11.1) y que tenga criterios (RN-11.2). Que la historia no esté ya en otro sprint
// (RN-11.3) lo valida Historia.AsignarASprint; el servicio coordina las dos.
func (s *Sprint) Asignar(h backlog.Historia) error {
	if s.Estado == Cerrado {
		return ErrSprintCerrado
	}
	if h.ProyectoID != s.ProyectoID {
		return ErrOtroProyecto
	}
	if h.Estado == backlog.Completada {
		return ErrHistoriaCompletada
	}
	if h.StoryPoints <= 0 {
		return ErrHistoriaSinEstimar
	}
	if !h.TieneCriterios() {
		return ErrHistoriaSinCriterios
	}
	s.Historias = append(s.Historias, h.ID)
	return nil
}

// Completar registra la historia como completada en el sprint (HU-12). Valida primero el sprint y después la
// historia: que el sprint no esté cerrado (RN-12.2), que la historia esté asignada (RN-12.1) y que no esté ya
// completada (RN-12.3). Que la historia pase a estado Completada lo hace Historia.Completar; el servicio coordina las dos.
func (s *Sprint) Completar(id domain.HistoriaID) error {
	if s.Estado == Cerrado {
		return ErrSprintCerrado
	}
	if !slices.Contains(s.Historias, id) {
		return ErrHistoriaNoAsignada
	}
	if slices.Contains(s.Completadas, id) {
		return ErrHistoriaCompletada
	}
	s.Completadas = append(s.Completadas, id)
	return nil
}
