package backlog

import (
	"errors"
)

// Errores de los criterios de aceptación (spec HU-07).
var (
	// ErrCriterioVacio: RN-07.1.
	ErrCriterioVacio = errors.New("el criterio de aceptación no puede estar vacío")
	// ErrCriterioNoExiste: RN-07.3.
	ErrCriterioNoExiste = errors.New("el criterio no existe")
)

// AgregarCriterio agrega un criterio de aceptación al final de la lista y valida RN-07.1 (spec HU-07).
func (h *Historia) AgregarCriterio(texto string) error {
	return nil // stub RED
}

// QuitarCriterio quita el criterio en la posición numero (empezando en 1) y renumera los siguientes (RN-07.3).
func (h *Historia) QuitarCriterio(numero int) error {
	return nil // stub RED
}

// TieneCriterios informa si la historia tiene al menos un criterio de aceptación. La usa HU-11: una historia
// sin criterios no puede asignarse a un sprint (RN-07.2).
func (h *Historia) TieneCriterios() bool {
	return false // stub RED
}
