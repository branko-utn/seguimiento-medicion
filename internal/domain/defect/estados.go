package defect

import (
	"errors"
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

var (
	// ErrNoExiste: el defecto no existe (lo devuelve el servicio).
	ErrNoExiste = errors.New("el defecto no existe")
	// ErrTransicionInvalida: RN-23.1 y RN-23.2. El mensaje concreto, con los estados, lo arma transicionInvalida.
	ErrTransicionInvalida = errors.New("transición de estado no permitida")
	// ErrResolverSinSprint: RN-23.3. CambiarEstado no acepta Resuelto: para resolver se usa Resolver.
	ErrResolverSinSprint = errors.New("para resolver un defecto hay que indicar el sprint de resolución")
	// ErrResolucionDeOtroProyecto: RN-23.3 (lo controla el servicio con los dos sprints).
	ErrResolucionDeOtroProyecto = errors.New("el sprint de resolución no pertenece al proyecto del defecto")
	// ErrResolucionAnterior: RN-23.4 (lo controla el servicio con los dos sprints).
	ErrResolucionAnterior = errors.New("el sprint de resolución no puede ser anterior al de detección")
)

// CambiarEstado pasa el defecto de Abierto a En progreso o de En progreso a Abierto (RN-23.1). Un defecto resuelto no
// cambia y no hay cambio al mismo estado (RN-23.2): ErrTransicionInvalida. Para pasar a Resuelto hay que usar
// Resolver, que exige el sprint: ErrResolverSinSprint.
//
// Como en las historias (HU-08), la transición se valida antes que el sprint: Resuelto → Resuelto es un error de
// transición.
func (d *Defecto) CambiarEstado(nuevo Estado) error {
	if d.Estado == Resuelto || d.Estado == nuevo {
		return transicionInvalida(d.Estado, nuevo)
	}
	if nuevo == Resuelto {
		return ErrResolverSinSprint
	}
	d.Estado = nuevo
	return nil
}

// Resolver deja el defecto Resuelto en el sprint indicado, desde Abierto o En progreso (RN-23.1). El sprint es
// obligatorio (RN-23.3). Que exista, sea del mismo proyecto y no sea anterior al de detección (RN-23.3 y RN-23.4) lo
// controla el servicio: defect no importa sprint (ADR 0002).
func (d *Defecto) Resolver(sprintID domain.SprintID) error {
	if d.Estado == Resuelto {
		return transicionInvalida(d.Estado, Resuelto)
	}
	if sprintID == "" {
		return ErrResolverSinSprint
	}
	d.Estado = Resuelto
	d.SprintResolucion = sprintID
	return nil
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

func transicionInvalida(desde, hacia Estado) error {
	return &errorTransicion{desde: desde, hacia: hacia}
}
