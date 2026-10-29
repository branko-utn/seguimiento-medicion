package backlog

import (
	"errors"
	"strings"
)

// Errores de HU-06 (los mensajes son los de la spec).
var (
	// ErrCompletadaNoSeEdita: RN-06.2.
	ErrCompletadaNoSeEdita = errors.New("una historia completada no se puede editar")
	// ErrAsignadaNoSeElimina: RN-06.3.
	ErrAsignadaNoSeElimina = errors.New("no se puede eliminar una historia asignada a un sprint")
	// ErrCompletadaNoSeElimina: RN-06.3.
	ErrCompletadaNoSeElimina = errors.New("no se puede eliminar una historia completada")
)

// Editar cambia el título, la descripción y la prioridad de la historia (spec HU-06). Aplica las validaciones de
// HU-05 (RN-06.1) y no edita una historia completada (RN-06.2). Si falla, la historia no cambia. El ID, el estado,
// los Story Points, los criterios y el sprint no se tocan.
func (h *Historia) Editar(titulo, descripcion string, prioridad Prioridad) error {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return ErrTituloVacio
	}
	if !prioridadValida(prioridad) {
		return ErrPrioridadInvalida
	}
	if h.Estado == Completada {
		return ErrCompletadaNoSeEdita
	}
	h.Titulo = titulo
	h.Descripcion = descripcion
	h.Prioridad = prioridad
	return nil
}

// PuedeEliminarse informa si la historia se puede sacar del backlog (spec HU-06, RN-06.3). Primero mira si está
// completada: una historia completada siempre tiene sprint, y así el error dice el motivo real.
func (h Historia) PuedeEliminarse() error {
	if h.Estado == Completada {
		return ErrCompletadaNoSeElimina
	}
	if h.SprintID != "" {
		return ErrAsignadaNoSeElimina
	}
	return nil
}
