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
