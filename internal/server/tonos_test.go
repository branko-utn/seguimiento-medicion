package server

import (
	"fmt"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// Los tonos siguen el avance: gris si todavía no empezó, azul si está en curso y verde si terminó. Un defecto abierto
// es rojo y uno en progreso, ámbar.
func TestTonoEstado(t *testing.T) {
	casos := []struct {
		estado any
		quiere string
	}{
		{project.Planificado, tonoNeutro},
		{project.EnCurso, tonoInfo},
		{project.Finalizado, tonoExito},
		{sprint.Planificado, tonoNeutro},
		{sprint.Activo, tonoInfo},
		{sprint.Cerrado, tonoExito},
		{backlog.Pendiente, tonoNeutro},
		{backlog.EnProgreso, tonoInfo},
		{backlog.Completada, tonoExito},
		{defect.Abierto, tonoPeligro},
		{defect.EnProgreso, tonoAlerta},
		{defect.Resuelto, tonoExito},
		// Valores que no son estados válidos: el cero (ADR 0003) y un tipo cualquiera quedan en neutro.
		{project.Estado(0), tonoNeutro},
		{sprint.Estado(0), tonoNeutro},
		{backlog.Estado(0), tonoNeutro},
		{defect.Estado(0), tonoNeutro},
		{"Activo", tonoNeutro},
	}

	for _, caso := range casos {
		t.Run(fmt.Sprintf("%T %v", caso.estado, caso.estado), func(t *testing.T) {
			if got := tonoEstado(caso.estado); got != caso.quiere {
				t.Errorf("tonoEstado(%T %v) = %q, se esperaba %q", caso.estado, caso.estado, got, caso.quiere)
			}
		})
	}
}

// La severidad va de más grave (rojo) a menos grave (gris).
func TestTonoSeveridad(t *testing.T) {
	casos := []struct {
		severidad defect.Severidad
		quiere    string
	}{
		{defect.Critica, tonoPeligro},
		{defect.Alta, tonoAlerta},
		{defect.Media, tonoInfo},
		{defect.Baja, tonoNeutro},
		{defect.Severidad(0), tonoNeutro},
	}

	for _, caso := range casos {
		t.Run(caso.severidad.String(), func(t *testing.T) {
			if got := tonoSeveridad(caso.severidad); got != caso.quiere {
				t.Errorf("tonoSeveridad(%v) = %q, se esperaba %q", caso.severidad, got, caso.quiere)
			}
		})
	}
}
