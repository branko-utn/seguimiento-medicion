package project

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func fecha(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.Local)
}

func TestNuevo(t *testing.T) {
	inicio := fecha(2026, time.October, 21)
	fin := fecha(2026, time.November, 2)

	casos := []struct {
		titulo       string
		nombre       string
		descripcion  string
		inicio       time.Time
		fin          time.Time
		quiereError  error
		quiereNombre string
		quiereInicio time.Time
		quiereFin    time.Time
	}{
		{
			titulo:       "crea un proyecto válido",
			nombre:       "Software Metrics",
			descripcion:  "Trabajo Práctico Integrador",
			inicio:       inicio,
			fin:          fin,
			quiereNombre: "Software Metrics",
			quiereInicio: inicio,
			quiereFin:    fin,
		},
		{
			titulo:       "quita los espacios sobrantes del nombre",
			nombre:       "  Demo  ",
			inicio:       inicio,
			fin:          fin,
			quiereNombre: "Demo",
			quiereInicio: inicio,
			quiereFin:    fin,
		},
		{
			titulo:       "acepta un nombre de 3 caracteres",
			nombre:       strings.Repeat("a", 3),
			inicio:       inicio,
			fin:          fin,
			quiereNombre: strings.Repeat("a", 3),
			quiereInicio: inicio,
			quiereFin:    fin,
		},
		{
			titulo:       "acepta un nombre de 100 caracteres",
			nombre:       strings.Repeat("a", 100),
			inicio:       inicio,
			fin:          fin,
			quiereNombre: strings.Repeat("a", 100),
			quiereInicio: inicio,
			quiereFin:    fin,
		},
		{
			titulo:       "cuenta runas y no bytes",
			nombre:       "Año",
			inicio:       inicio,
			fin:          fin,
			quiereNombre: "Año",
			quiereInicio: inicio,
			quiereFin:    fin,
		},
		{
			titulo:      "rechaza un nombre de 2 caracteres",
			nombre:      "ab",
			inicio:      inicio,
			fin:         fin,
			quiereError: ErrNombreInvalido,
		},
		{
			titulo:      "rechaza un nombre de 101 caracteres",
			nombre:      strings.Repeat("a", 101),
			inicio:      inicio,
			fin:         fin,
			quiereError: ErrNombreInvalido,
		},
		{
			titulo:      "rechaza un nombre vacío",
			nombre:      "",
			inicio:      inicio,
			fin:         fin,
			quiereError: ErrNombreInvalido,
		},
		{
			titulo:      "rechaza un nombre con solo espacios",
			nombre:      "   ",
			inicio:      inicio,
			fin:         fin,
			quiereError: ErrNombreInvalido,
		},
		{
			titulo:      "rechaza un fin anterior al inicio",
			nombre:      "Demo",
			inicio:      fin,
			fin:         inicio,
			quiereError: ErrFechasInvalidas,
		},
		{
			titulo:       "acepta inicio y fin el mismo día",
			nombre:       "Demo",
			inicio:       inicio,
			fin:          inicio,
			quiereNombre: "Demo",
			quiereInicio: inicio,
			quiereFin:    inicio,
		},
		{
			titulo:      "rechaza un proyecto sin fecha de inicio",
			nombre:      "Demo",
			inicio:      time.Time{},
			fin:         fin,
			quiereError: ErrFechasInvalidas,
		},
		{
			titulo:      "rechaza un proyecto sin fecha de fin",
			nombre:      "Demo",
			inicio:      inicio,
			fin:         time.Time{},
			quiereError: ErrFechasInvalidas,
		},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			p, err := Nuevo(caso.nombre, caso.descripcion, caso.inicio, caso.fin)

			if caso.quiereError != nil {
				if !errors.Is(err, caso.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if p.Nombre != caso.quiereNombre {
				t.Errorf("nombre = %q, se esperaba %q", p.Nombre, caso.quiereNombre)
			}
			if p.Estado != Planificado {
				t.Errorf("estado = %q, se esperaba %q", p.Estado, Planificado)
			}
			if !p.Inicio.Equal(caso.quiereInicio) {
				t.Errorf("inicio = %v, se esperaba %v", p.Inicio, caso.quiereInicio)
			}
			if !p.Fin.Equal(caso.quiereFin) {
				t.Errorf("fin = %v, se esperaba %v", p.Fin, caso.quiereFin)
			}
		})
	}
}

// proyectoGuardado devuelve el proyecto "Demo" (P1) del 21/10 al 02/11 en el estado indicado.
func proyectoGuardado(estado Estado) Proyecto {
	return Proyecto{
		ID:          "P1",
		Nombre:      "Demo",
		Descripcion: "Original",
		Inicio:      fecha(2026, time.October, 21),
		Fin:         fecha(2026, time.November, 2),
		Estado:      estado,
	}
}

func TestModificar(t *testing.T) {
	inicio := fecha(2026, time.October, 22)
	fin := fecha(2026, time.November, 3)

	casos := []struct {
		titulo       string
		estado       Estado
		nombre       string
		descripcion  string
		quiereNombre string
	}{
		{titulo: "cambia nombre, descripción y fechas", estado: Planificado, nombre: "Software Metrics", descripcion: "TPI", quiereNombre: "Software Metrics"},
		{titulo: "quita los espacios sobrantes del nombre", estado: Planificado, nombre: "  Software Metrics  ", quiereNombre: "Software Metrics"},
		{titulo: "permite dejar la descripción vacía", estado: Planificado, nombre: "Demo", descripcion: "", quiereNombre: "Demo"},
		{titulo: "permite cambiar solo las mayúsculas del nombre", estado: Planificado, nombre: "DEMO", quiereNombre: "DEMO"},
		{titulo: "se puede modificar un proyecto en curso", estado: EnCurso, nombre: "Demo", quiereNombre: "Demo"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			p := proyectoGuardado(caso.estado)

			if err := p.Modificar(caso.nombre, caso.descripcion, inicio, fin); err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if p.Nombre != caso.quiereNombre || p.Descripcion != caso.descripcion {
				t.Errorf("nombre y descripción = %q y %q, se esperaba %q y %q", p.Nombre, p.Descripcion, caso.quiereNombre, caso.descripcion)
			}
			if !p.Inicio.Equal(inicio) || !p.Fin.Equal(fin) {
				t.Errorf("fechas = %v al %v, se esperaba %v al %v", p.Inicio, p.Fin, inicio, fin)
			}
			// RN-02.4: el ID y el estado no cambian.
			if p.ID != "P1" || p.Estado != caso.estado {
				t.Errorf("ID y estado = %q y %v, se esperaba P1 y %v", p.ID, p.Estado, caso.estado)
			}
		})
	}
}

func TestModificar_sinCambios(t *testing.T) {
	p := proyectoGuardado(Planificado)
	antes := p

	if err := p.Modificar(p.Nombre, p.Descripcion, p.Inicio, p.Fin); err != nil {
		t.Fatalf("guardar sin cambios es válido, pero se obtuvo %v", err)
	}
	if p != antes {
		t.Errorf("Modificar sin cambios dejó %+v, se esperaba %+v", p, antes)
	}
}

func TestModificar_errores(t *testing.T) {
	inicio := fecha(2026, time.October, 21)
	fin := fecha(2026, time.November, 2)

	casos := []struct {
		titulo      string
		estado      Estado
		nombre      string
		inicio, fin time.Time
		quiereError error
	}{
		{"nombre de 2 caracteres", Planificado, "AB", inicio, fin, ErrNombreInvalido},
		{"nombre de 101 caracteres", Planificado, strings.Repeat("a", 101), inicio, fin, ErrNombreInvalido},
		{"nombre con solo espacios", Planificado, "   ", inicio, fin, ErrNombreInvalido},
		{"fin anterior al inicio", Planificado, "Demo", fin, inicio, ErrFechasInvalidas},
		{"sin fecha de inicio", Planificado, "Demo", time.Time{}, fin, ErrFechasInvalidas},
		{"proyecto finalizado", Finalizado, "Demo 2", inicio, fin, ErrProyectoFinalizado},
		{"primero se controla que no esté finalizado", Finalizado, "AB", fin, inicio, ErrProyectoFinalizado},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			p := proyectoGuardado(caso.estado)
			antes := p

			err := p.Modificar(caso.nombre, "Cambiada", caso.inicio, caso.fin)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if p != antes {
				t.Errorf("el proyecto quedó %+v, se esperaba %+v: un error no debe cambiarlo", p, antes)
			}
		})
	}
}

// TestEstadoString verifica el texto de cada estado (ADR 0003: el cero no es válido).
func TestEstadoString(t *testing.T) {
	casos := []struct {
		titulo string
		estado Estado
		quiere string
	}{
		{titulo: "planificado", estado: Planificado, quiere: "Planificado"},
		{titulo: "en curso", estado: EnCurso, quiere: "En curso"},
		{titulo: "finalizado", estado: Finalizado, quiere: "Finalizado"},
		{titulo: "valor cero no es un estado válido", estado: Estado(0), quiere: ""},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			if got := caso.estado.String(); got != caso.quiere {
				t.Errorf("String() = %q, se esperaba %q", got, caso.quiere)
			}
		})
	}
}
