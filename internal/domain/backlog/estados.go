package backlog

import "errors"

// Errores del cambio de estado (spec HU-08).
var (
	// ErrEstadoInvalido: RN-08.1.
	ErrEstadoInvalido = errors.New("el estado debe ser Pendiente, En progreso o Completada")
	// ErrTransicionInvalida: RN-08.2 a RN-08.4. El mensaje concreto, con los estados, lo arma transicionInvalida.
	ErrTransicionInvalida = errors.New("transición de estado no permitida")
	// ErrSinSprint: RN-08.2.
	ErrSinSprint = errors.New("la historia debe estar asignada a un sprint")
)

// CambiarEstado deja la historia en el estado nuevo si la transición es válida (spec HU-08).
func (h *Historia) CambiarEstado(nuevo Estado) error {
	return nil // stub RED
}

// Completar pasa la historia a Completada con la misma regla de CambiarEstado; la usa HU-12 (spec HU-08).
func (h *Historia) Completar() error {
	return nil // stub RED
}
