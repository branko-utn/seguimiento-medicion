package sprint

import (
	"errors"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// Estado es la situación de un sprint (ADR 0003: el cero no es válido).
type Estado int

const (
	Planificado Estado = iota + 1
	Activo
	Cerrado
)

// String devuelve el texto que usan las specs y los escenarios BDD. El stub existe solo para la fase RED.
func (e Estado) String() string {
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
)

// Nuevo crea un sprint válido. El stub existe solo para que los tests de la fase RED compilen.
func Nuevo(proyectoID domain.ProyectoID, goal string, inicio, fin time.Time) (Sprint, error) {
	return Sprint{}, nil
}

// ValidarFechasEnProyecto controla RN-10.3. El stub existe solo para la fase RED.
func (s Sprint) ValidarFechasEnProyecto(p project.Proyecto) error {
	return nil
}

// ValidarSinSuperposicion controla RN-10.4 contra los otros sprints del proyecto. El stub existe solo para la fase RED.
func (s Sprint) ValidarSinSuperposicion(otros []Sprint) error {
	return nil
}

// Iniciar pasa el sprint a Activo (RN-10.6). El stub existe solo para la fase RED.
func (s *Sprint) Iniciar() error {
	return nil
}
