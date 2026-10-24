package defect

import (
	"errors"

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

func (s Severidad) String() string {
	return ""
}

// Estado es la situación de un defecto (ADR 0003: el cero no es válido).
type Estado int

const (
	Abierto Estado = iota + 1
	EnProgreso
	Resuelto
)

func (e Estado) String() string {
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

func Nuevo(descripcion string, severidad Severidad, historiaID domain.HistoriaID, sprintDeteccion domain.SprintID) (Defecto, error) {
	return Defecto{}, nil
}
