package backlog

import "errors"

// Errores de HU-06 (los mensajes son los de la spec).
var (
	// ErrCompletadaNoSeEdita: RN-06.2.
	ErrCompletadaNoSeEdita = errors.New("una historia completada no se puede editar")
	// ErrAsignadaNoSeElimina: RN-06.3.
	ErrAsignadaNoSeElimina = errors.New("no se puede eliminar una historia asignada a un sprint")
	// ErrCompletadaNoSeElimina: RN-06.3.
	ErrCompletadaNoSeElimina = errors.New("no se puede eliminar una historia completada")
)

// Editar cambia el título, la descripción y la prioridad de la historia (spec HU-06).
func (h *Historia) Editar(titulo, descripcion string, prioridad Prioridad) error {
	return nil
}

// PuedeEliminarse informa si la historia se puede sacar del backlog (spec HU-06).
func (h Historia) PuedeEliminarse() error {
	return nil
}
