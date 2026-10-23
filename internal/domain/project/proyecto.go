package project

import (
	"errors"
	"time"

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

// Nuevo crea un proyecto válido. El stub existe solo para que los tests de la fase RED compilen.
func Nuevo(nombre, descripcion string, inicio, fin time.Time) (Proyecto, error) {
	return Proyecto{}, nil
}
