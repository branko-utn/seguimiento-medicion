package backlog

import (
	"errors"
	"strings"

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

// String devuelve el texto que ve el usuario (ADR 0003). Para un valor inválido devuelve cadena vacía.
func (p Prioridad) String() string {
	switch p {
	case PrioridadAlta:
		return "Alta"
	case PrioridadMedia:
		return "Media"
	case PrioridadBaja:
		return "Baja"
	}
	return ""
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

// String devuelve el texto que ve el usuario (ADR 0003). Para un valor inválido devuelve cadena vacía.
func (e Estado) String() string {
	switch e {
	case Pendiente:
		return "Pendiente"
	case EnProgreso:
		return "En progreso"
	case Completada:
		return "Completada"
	}
	return ""
}

// Criterio es un criterio de aceptación de una historia (spec HU-07).
type Criterio struct {
	Texto string
}

// Historia es un ítem del Product Backlog (ADR 0002).
type Historia struct {
	ID          domain.HistoriaID
	ProyectoID  domain.ProyectoID
	Titulo      string
	Descripcion string
	Prioridad   Prioridad
	Estado      Estado
	StoryPoints int
	Criterios   []Criterio
	SprintID    domain.SprintID
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
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return Historia{}, ErrTituloVacio
	}
	if !prioridadValida(prioridad) {
		return Historia{}, ErrPrioridadInvalida
	}
	return Historia{
		ProyectoID: proyectoID,
		Titulo:     titulo,
		// RN-05.2: la descripción es opcional.
		Descripcion: descripcion,
		Prioridad:   prioridad,
		// RN-05.4: nueva = Pendiente, sin Story Points (0), sin criterios y sin sprint.
		// El valor cero de esos campos ya lo garantiza.
		Estado: Pendiente,
	}, nil
}

// prioridadValida informa si la prioridad pertenece a la escala de la spec HU-05 (RN-05.3).
func prioridadValida(p Prioridad) bool {
	switch p {
	case PrioridadAlta, PrioridadMedia, PrioridadBaja:
		return true
	}
	return false
}

// AsignarASprint deja la historia asignada al sprint indicado y valida RN-11.3 (spec HU-11).
func (h *Historia) AsignarASprint(id domain.SprintID) error {
	if h.SprintID != "" {
		return ErrYaAsignada
	}
	h.SprintID = id
	return nil
}
