package project

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Estado es la situación de un proyecto en su ciclo de vida (ADR 0003: el cero no es válido).
type Estado int

const (
	Planificado Estado = iota + 1
	EnCurso
	Finalizado
)

// String devuelve el texto que usan las specs y los escenarios BDD.
func (e Estado) String() string {
	switch e {
	case Planificado:
		return "Planificado"
	case EnCurso:
		return "En curso" // con espacio, como en HU-04
	case Finalizado:
		return "Finalizado"
	}
	return ""
}

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
	ErrNoExiste        = errors.New("el proyecto no existe")

	// ErrProyectoFinalizado: RN-02.3 (HU-02).
	ErrProyectoFinalizado = errors.New("un proyecto finalizado no se puede modificar")
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

// Modificar cambia nombre, descripción y fechas del proyecto (HU-02). Primero controla que no esté finalizado
// (RN-02.3) y después aplica las validaciones del alta reutilizando Nuevo (RN-02.1). El ID y el estado no cambian
// (RN-02.4). Si hay un error, el proyecto queda como estaba. Que el nombre no sea de otro proyecto (RN-02.2) lo
// controla el servicio.
//
// "Finalizado" se calcula con las fechas actuales del proyecto respecto de hoy (EstadoAl, decisión de la planning del
// Sprint 3 en HU-04), no con el campo Estado guardado: un proyecto que ya terminó no se reabre cambiándole el fin.
func (p *Proyecto) Modificar(nombre, descripcion string, inicio, fin, hoy time.Time) error {
	if p.EstadoAl(hoy) == Finalizado {
		return ErrProyectoFinalizado
	}
	datos, err := Nuevo(nombre, descripcion, inicio, fin)
	if err != nil {
		return err
	}
	p.Nombre = datos.Nombre
	p.Descripcion = datos.Descripcion
	p.Inicio = datos.Inicio
	p.Fin = datos.Fin
	return nil
}
