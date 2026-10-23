package backlog

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Prioridad es la urgencia de una historia: Alta, Media o Baja (ADR 0002).
// Se representa como int con el cero reservado para "sin valor" (ADR 0003).
type Prioridad int

// Las prioridades que acepta la spec HU-05 (RN-05.3).
const (
	PrioridadAlta Prioridad = iota + 1
	PrioridadMedia
	PrioridadBaja
)

// String devuelve el texto que ve el usuario (ADR 0003).
func (p Prioridad) String() string {
	return "" // stub RED
}

// Estado es el estado del ciclo de vida de una historia: Pendiente, En progreso o Completada (ADR 0002).
// Se representa como int con el cero reservado para "sin valor" (ADR 0003).
type Estado int

// Los estados de la spec HU-08 (RN-08.1).
const (
	Pendiente Estado = iota + 1
	EnProgreso
	Completada
)

// String devuelve el texto que ve el usuario (ADR 0003).
func (e Estado) String() string {
	return "" // stub RED
}

// Criterio es un criterio de aceptación de una historia (spec HU-07).
type Criterio struct {
	Texto string `json:"texto"`
}

// Historia es un ítem del Product Backlog (ADR 0002).
type Historia struct {
	ID          domain.HistoriaID `json:"id"`
	ProyectoID  domain.ProyectoID `json:"proyecto_id"`
	Titulo      string            `json:"titulo"`
	Descripcion string            `json:"descripcion"`
	Prioridad   Prioridad         `json:"prioridad"`
	Estado      Estado            `json:"estado"`
	StoryPoints int               `json:"story_points"`
	Criterios   []Criterio        `json:"criterios,omitempty"`
	SprintID    domain.SprintID   `json:"sprint_id,omitempty"`
}

// Errores que devuelve este paquete (los mensajes son los de las specs).
var (
	// ErrTituloVacio: RN-05.2.
	ErrTituloVacio = errors.New("el título es obligatorio")
	// ErrPrioridadInvalida: RN-05.3.
	ErrPrioridadInvalida = errors.New("la prioridad debe ser Alta, Media o Baja")
	// ErrYaAsignada: RN-11.3, la usa HU-11 (AsignarASprint).
	ErrYaAsignada = errors.New("la historia ya está asignada a un sprint")
)

// Nueva crea una historia del backlog de un proyecto y valida RN-05.2, RN-05.3 y RN-05.4 (spec HU-05).
// El ID no se asigna acá: lo genera el repositorio al guardar (RN-05.1).
func Nueva(proyectoID domain.ProyectoID, titulo, descripcion string, prioridad Prioridad) (Historia, error) {
	return Historia{}, nil // stub RED
}

// AsignarASprint deja la historia asignada al sprint indicado y valida RN-11.3 (spec HU-11).
func (h *Historia) AsignarASprint(id domain.SprintID) error {
	return nil // stub RED
}
