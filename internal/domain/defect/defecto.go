package defect

import (
	"errors"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Severidad es la gravedad de un defecto (ADR 0003: el cero no es válido).
type Severidad int

const (
	Critica Severidad = iota + 1
	Alta
	Media
	Baja
)

// String devuelve el texto que usan las specs y los escenarios BDD.
func (s Severidad) String() string {
	switch s {
	case Critica:
		return "Crítica"
	case Alta:
		return "Alta"
	case Media:
		return "Media"
	case Baja:
		return "Baja"
	}
	return ""
}

// Estado es la situación de un defecto (ADR 0003: el cero no es válido).
type Estado int

const (
	Abierto Estado = iota + 1
	EnProgreso
	Resuelto
)

// String devuelve el texto que usan las specs y los escenarios BDD.
func (e Estado) String() string {
	switch e {
	case Abierto:
		return "Abierto"
	case EnProgreso:
		return "En progreso"
	case Resuelto:
		return "Resuelto"
	}
	return ""
}

// Defecto es un defecto encontrado en una historia (ADR 0002).
type Defecto struct {
	ID               domain.DefectoID
	Descripcion      string
	Severidad        Severidad
	Estado           Estado
	HistoriaID       domain.HistoriaID
	SprintDeteccion  domain.SprintID
	SprintResolucion domain.SprintID // "" mientras no esté resuelto
}

var (
	ErrDescripcionVacia     = errors.New("la descripción es obligatoria")
	ErrSeveridadInvalida    = errors.New("la severidad debe ser Crítica, Alta, Media o Baja")
	ErrSprintDeOtroProyecto = errors.New("el sprint de detección no pertenece al proyecto de la historia")
)

// Nuevo crea un defecto válido: descripción no vacía (RN-22.1) y severidad de la lista (RN-22.2). Queda Abierto y sin
// sprint de resolución (RN-22.3). Que la historia y el sprint existan y sean del mismo proyecto (RN-22.4) lo controla
// el servicio.
func Nuevo(descripcion string, severidad Severidad, historiaID domain.HistoriaID, sprintDeteccion domain.SprintID) (Defecto, error) {
	descripcion = strings.TrimSpace(descripcion)
	if descripcion == "" {
		return Defecto{}, ErrDescripcionVacia
	}
	if severidad.String() == "" {
		return Defecto{}, ErrSeveridadInvalida
	}
	return Defecto{
		Descripcion:     descripcion,
		Severidad:       severidad,
		Estado:          Abierto,
		HistoriaID:      historiaID,
		SprintDeteccion: sprintDeteccion,
	}, nil
}
