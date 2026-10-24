package backlog

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
)

// ErrHistoriaCompletada: RN-15.2, la usa HU-15 (Estimar).
var ErrHistoriaCompletada = errors.New("no se puede re-estimar una historia completada")

// Estimar asigna Story Points a la historia y valida RN-15.1 a RN-15.3 (spec HU-15). Re-estimar reemplaza el
// valor anterior, mientras la historia no esté completada. El orden de validación sigue la tabla de errores de
// la spec: primero la escala (RN-15.1) y después el estado (RN-15.2).
func (h *Historia) Estimar(sp int) error {
	if err := estimation.ValidarStoryPoints(sp); err != nil {
		return err
	}
	if h.Estado == Completada {
		return ErrHistoriaCompletada
	}
	h.StoryPoints = sp
	return nil
}
