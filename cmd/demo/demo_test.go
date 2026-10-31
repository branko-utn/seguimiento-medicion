package main

import (
	"errors"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/app"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

var hoyDemo = time.Date(2026, 10, 31, 12, 0, 0, 0, time.Local)

// TestCargar verifica que la demo cargue nuestro propio proyecto (TEC-06) con los datos de los Sprints 1 a 3 y el
// Sprint 4 en curso, pasando por los mismos servicios (y validaciones) que la aplicación.
func TestCargar(t *testing.T) {
	svc, err := app.NuevosServicios("")
	if err != nil {
		t.Fatal(err)
	}
	p, err := cargar(svc, hoyDemo)
	if err != nil {
		t.Fatalf("cargar: %v", err)
	}

	if p.Nombre != "Software Metrics & Estimation" {
		t.Errorf("proyecto = %q", p.Nombre)
	}
	integrantes, _ := svc.Integrantes.DelProyecto(p.ID)
	if len(integrantes) != 5 {
		t.Errorf("integrantes = %d, se esperaban 5", len(integrantes))
	}

	// Los SP de cada sprint son los del tablero: 44 en el Sprint 1, 23 en el 2 (sin los ítems técnicos) y 35 en el 3.
	sprints, _ := svc.Sprints.DelProyecto(p.ID)
	if len(sprints) != 4 {
		t.Fatalf("sprints = %d, se esperaban 4", len(sprints))
	}
	quiere := []struct {
		estado      sprint.Estado
		completados int
	}{{sprint.Cerrado, 44}, {sprint.Cerrado, 23}, {sprint.Cerrado, 35}, {sprint.Activo, 0}}
	for i, q := range quiere {
		if sprints[i].Estado != q.estado || sprints[i].SPCompletados != q.completados {
			t.Errorf("sprint %d = %s con %d SP completados, se esperaba %s con %d", i+1, sprints[i].Estado, sprints[i].SPCompletados, q.estado, q.completados)
		}
	}

	r, err := svc.Reportes.Proyecto(p.ID, hoyDemo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Historias != 31 || r.Completadas != 30 || r.Velocidad != 34 {
		t.Errorf("historias %d, completadas %d, velocidad %v; se esperaba 31, 30 y 34", r.Historias, r.Completadas, r.Velocidad)
	}
	if r.HorasEstimadas == 0 || r.HorasReales == 0 || !r.DesviacionCalculable {
		t.Errorf("horas = %v/%v: la demo tiene que cargar horas estimadas y reales", r.HorasEstimadas, r.HorasReales)
	}
	if r.Detectados != 9 || r.Resueltos != 8 || r.Abiertos != 1 {
		t.Errorf("defectos = %d/%d/%d, se esperaba 9 detectados, 8 resueltos y 1 abierto", r.Detectados, r.Resueltos, r.Abiertos)
	}

	// El dashboard grafica los tres sprints cerrados.
	d, err := svc.Dashboard.Del(p.ID, hoyDemo)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Puntos) != 3 {
		t.Errorf("el dashboard tiene %d puntos, se esperaban los 3 sprints cerrados", len(d.Puntos))
	}
}

// Si el proyecto ya está cargado, la demo no lo duplica: devuelve el error del alta.
func TestCargar_dosVeces(t *testing.T) {
	svc, err := app.NuevosServicios("")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cargar(svc, hoyDemo); err != nil {
		t.Fatal(err)
	}
	if _, err := cargar(svc, hoyDemo); !errors.Is(err, project.ErrNombreDuplicado) {
		t.Errorf("error = %v, se esperaba project.ErrNombreDuplicado", err)
	}
}
