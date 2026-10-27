package app

import (
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

func dia(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }

func TestNuevosServicios_todosConectados(t *testing.T) {
	svc, err := NuevosServicios("")
	if err != nil {
		t.Fatal(err)
	}
	if svc.Proyectos == nil || svc.Integrantes == nil || svc.Backlog == nil || svc.Sprints == nil ||
		svc.Estimaciones == nil || svc.Registros == nil || svc.Defectos == nil {
		t.Fatalf("hay servicios sin crear: %+v", svc)
	}

	// Recorre el flujo principal para comprobar que los servicios comparten los mismos repositorios.
	p, err := svc.Proyectos.Crear("Demo", "", dia(21), dia(31))
	if err != nil {
		t.Fatal(err)
	}
	h, err := svc.Backlog.Crear(p.ID, "Crear proyecto", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Backlog.AgregarCriterio(h.ID, "El nombre es obligatorio"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Backlog.Estimar(h.ID, 3); err != nil {
		t.Fatal(err)
	}
	sp, err := svc.Sprints.Crear(p.ID, "MVP", dia(23), dia(25))
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Sprints.Asignar(sp.ID, h.ID); err != nil {
		t.Fatal(err)
	}
	// Backlog delega en Sprints para completar (HU-08 → HU-12): tienen que estar conectados entre sí.
	if err := svc.Backlog.CambiarEstado(h.ID, backlog.Completada); err != nil {
		t.Fatalf("completar desde el backlog: %v", err)
	}
	if got, _ := svc.Backlog.Buscar(h.ID); got.Estado != backlog.Completada {
		t.Errorf("estado = %s, se esperaba Completada", got.Estado)
	}
}

func TestNuevosServicios_persisteEnLaCarpeta(t *testing.T) {
	dir := t.TempDir()
	svc, err := NuevosServicios(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Proyectos.Crear("Demo", "", dia(21), dia(31)); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevosServicios(dir)
	if err != nil {
		t.Fatal(err)
	}
	lista, err := reabierto.Proyectos.Listar()
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 1 || lista[0].Nombre != "Demo" {
		t.Errorf("al reabrir la carpeta se obtuvo %+v", lista)
	}
}
