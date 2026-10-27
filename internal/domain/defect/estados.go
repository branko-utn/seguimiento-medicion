package defect

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

var (
	ErrNoExiste                 = errors.New("el defecto no existe")
	ErrTransicionInvalida       = errors.New("transición de estado no permitida")
	ErrResolverSinSprint        = errors.New("para resolver un defecto hay que indicar el sprint de resolución")
	ErrResolucionDeOtroProyecto = errors.New("el sprint de resolución no pertenece al proyecto del defecto")
	ErrResolucionAnterior       = errors.New("el sprint de resolución no puede ser anterior al de detección")
)

func (d *Defecto) CambiarEstado(nuevo Estado) error {
	return nil
}

func (d *Defecto) Resolver(sprintID domain.SprintID) error {
	return nil
}
