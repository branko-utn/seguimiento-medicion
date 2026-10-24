package backlog

import "errors"

// ErrHistoriaCompletada: RN-15.2, la usa HU-15 (Estimar).
var ErrHistoriaCompletada = errors.New("no se puede re-estimar una historia completada")

// Estimar asigna Story Points a la historia y valida RN-15.1 a RN-15.3 (spec HU-15). Re-estimar reemplaza el
// valor anterior, mientras la historia no esté completada.
func (h *Historia) Estimar(sp int) error {
	return nil // stub RED
}
