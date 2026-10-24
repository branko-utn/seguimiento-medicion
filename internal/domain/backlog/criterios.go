package backlog

import (
	"errors"
	"strings"
)

// Errores de los criterios de aceptación (spec HU-07).
var (
	// ErrCriterioVacio: RN-07.1.
	ErrCriterioVacio = errors.New("el criterio de aceptación no puede estar vacío")
	// ErrCriterioNoExiste: RN-07.3.
	ErrCriterioNoExiste = errors.New("el criterio no existe")
)

// AgregarCriterio agrega un criterio de aceptación al final de la lista, sin espacios sobrantes, y valida
// RN-07.1 (spec HU-07).
func (h *Historia) AgregarCriterio(texto string) error {
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return ErrCriterioVacio
	}
	h.Criterios = append(h.Criterios, Criterio{Texto: texto})
	return nil
}

// QuitarCriterio quita el criterio en la posición numero (empezando en 1, RN-07.3) y los siguientes se renumeran.
// Si no hay criterio en esa posición, devuelve ErrCriterioNoExiste.
func (h *Historia) QuitarCriterio(numero int) error {
	if numero < 1 || numero > len(h.Criterios) {
		return ErrCriterioNoExiste
	}
	h.Criterios = append(h.Criterios[:numero-1], h.Criterios[numero:]...)
	return nil
}

// TieneCriterios informa si la historia tiene al menos un criterio de aceptación. La usa HU-11: una historia
// sin criterios no puede asignarse a un sprint (RN-07.2).
func (h *Historia) TieneCriterios() bool {
	return len(h.Criterios) > 0
}
