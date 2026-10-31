package backlog

import (
	"errors"
	"fmt"
)

// Errores del cambio de estado (spec HU-08).
var (
	// ErrEstadoInvalido: RN-08.1.
	ErrEstadoInvalido = errors.New("el estado debe ser Pendiente, En progreso o Completada")
	// ErrTransicionInvalida: RN-08.2 a RN-08.4. El mensaje concreto, con los estados, lo arma transicionInvalida.
	ErrTransicionInvalida = errors.New("transición de estado no permitida")
	// ErrSinSprint: RN-08.2.
	ErrSinSprint = errors.New("la historia debe estar asignada a un sprint")
)

// CambiarEstado deja la historia en el estado nuevo si la transición es válida (spec HU-08):
//   - el estado nuevo debe ser Pendiente, En progreso o Completada (RN-08.1);
//   - una historia completada no cambia (RN-08.3) y no hay cambio al mismo estado (RN-08.4), antes de exigir
//     sprint: ErrTransicionInvalida;
//   - En progreso y Completada exigen estar asignada a un sprint (RN-08.2): ErrSinSprint.
//
// El orden de validación es el que fijan los escenarios BDD ("cambiar al mismo estado" no tiene sprint y espera
// el error de transición).
func (h *Historia) CambiarEstado(nuevo Estado) error {
	if !estadoValido(nuevo) {
		return ErrEstadoInvalido
	}
	if h.Estado == Completada || h.Estado == nuevo {
		return transicionInvalida(h.Estado, nuevo)
	}
	if (nuevo == EnProgreso || nuevo == Completada) && h.SprintID == "" {
		return ErrSinSprint
	}
	h.Estado = nuevo
	return nil
}

// Completar pasa la historia a Completada con la misma regla de CambiarEstado (spec HU-08); la usa HU-12.
func (h *Historia) Completar() error {
	return h.CambiarEstado(Completada)
}

// ValidarEstado devuelve ErrEstadoInvalido si el estado no es Pendiente, En progreso ni Completada (RN-08.1). Lo usa
// el filtro por estado del backlog (HU-09).
func ValidarEstado(e Estado) error {
	if !estadoValido(e) {
		return ErrEstadoInvalido
	}
	return nil
}

// estadoValido informa si el estado pertenece a los tres de la spec (RN-08.1).
func estadoValido(e Estado) bool {
	switch e {
	case Pendiente, EnProgreso, Completada:
		return true
	}
	return false
}

// errorTransicion es una transición no permitida con el mensaje exacto de la spec ("no se puede pasar de X a Y")
// y se desenvuelve en ErrTransicionInvalida, para que errors.Is la encuentre.
type errorTransicion struct {
	desde, hacia Estado
}

func (e *errorTransicion) Error() string {
	return fmt.Sprintf("no se puede pasar de %s a %s", e.desde, e.hacia)
}

func (e *errorTransicion) Unwrap() error {
	return ErrTransicionInvalida
}

// transicionInvalida construye el error de una transición no permitida (RN-08.2 a RN-08.4).
func transicionInvalida(desde, hacia Estado) error {
	return &errorTransicion{desde: desde, hacia: hacia}
}
