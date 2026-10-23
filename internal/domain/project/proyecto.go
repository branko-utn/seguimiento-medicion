package project

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Estado es la situación de un proyecto en su ciclo de vida.
type Estado string

const (
	Planificado Estado = "Planificado"
	EnCurso     Estado = "EnCurso"
	Finalizado  Estado = "Finalizado"
)

// Proyecto es un proyecto del equipo (ADR 0002).
type Proyecto struct {
	ID          domain.ProyectoID
	Nombre      string
	Descripcion string
	Inicio      time.Time
	Fin         time.Time
	Estado      Estado
}

var (
	ErrNombreInvalido  = errors.New("el nombre debe tener entre 3 y 100 caracteres")
	ErrFechasInvalidas = errors.New("la fecha de fin no puede ser anterior a la de inicio")
	ErrNombreDuplicado = errors.New("ya existe un proyecto con ese nombre")
)

// Nuevo crea un proyecto válido en estado Planificado (RN-01.1, RN-01.2 y RN-01.4).
// La unicidad del nombre (RN-01.3) la controla el servicio.
func Nuevo(nombre, descripcion string, inicio, fin time.Time) (Proyecto, error) {
	nombre = strings.TrimSpace(nombre)
	if largo := utf8.RuneCountInString(nombre); largo < 3 || largo > 100 {
		return Proyecto{}, ErrNombreInvalido
	}
	if inicio.IsZero() || fin.IsZero() || fin.Before(inicio) {
		return Proyecto{}, ErrFechasInvalidas
	}
	return Proyecto{
		Nombre:      nombre,
		Descripcion: descripcion,
		Inicio:      inicio,
		Fin:         fin,
		Estado:      Planificado,
	}, nil
}
