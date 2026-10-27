package memory

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

func TestProyectos_AgregarGeneraIDs(t *testing.T) {
	r, err := NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}

	for i, quiereID := range []string{"P1", "P2", "P3"} {
		p, err := r.Agregar(project.Proyecto{Nombre: "Proyecto", Estado: project.Planificado})
		if err != nil {
			t.Fatalf("Agregar #%d: %v", i+1, err)
		}
		if string(p.ID) != quiereID {
			t.Errorf("Agregar #%d: ID = %q, se esperaba %q", i+1, p.ID, quiereID)
		}
	}
	if n := len(r.Listar()); n != 3 {
		t.Errorf("Listar() tiene %d proyectos, se esperaban 3", n)
	}
}

func TestProyectos_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "proyectos.json")
	r, err := NuevoProyectos(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(project.Proyecto{Nombre: "Demo", Estado: project.Planificado}); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevoProyectos(ruta)
	if err != nil {
		t.Fatal(err)
	}
	lista := reabierto.Listar()
	if len(lista) != 1 || lista[0].ID != "P1" || lista[0].Nombre != "Demo" {
		t.Fatalf("al reabrir se obtuvo %+v", lista)
	}

	p, err := reabierto.Agregar(project.Proyecto{Nombre: "Otro"})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "P2" {
		t.Errorf("después de reabrir, ID = %q, se esperaba %q", p.ID, "P2")
	}
}

func TestProyectos_ListarDevuelveCopia(t *testing.T) {
	r, _ := NuevoProyectos("")
	if _, err := r.Agregar(project.Proyecto{Nombre: "Demo"}); err != nil {
		t.Fatal(err)
	}

	lista := r.Listar()
	if len(lista) != 1 {
		t.Fatalf("Listar() tiene %d proyectos, se esperaba 1", len(lista))
	}
	lista[0].Nombre = "Cambiado"

	if got := r.Listar()[0].Nombre; got != "Demo" {
		t.Errorf("modificar lo que devuelve Listar cambió el repositorio: nombre = %q", got)
	}
}

func TestProyectos_Buscar(t *testing.T) {
	r, _ := NuevoProyectos("")
	if _, err := r.Agregar(project.Proyecto{Nombre: "Demo"}); err != nil {
		t.Fatal(err)
	}

	p, err := r.Buscar("P1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.ID != "P1" || p.Nombre != "Demo" {
		t.Errorf("Buscar(P1) = %+v, se esperaba el proyecto Demo", p)
	}

	if _, err := r.Buscar("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("Buscar(P9): se esperaba %v, se obtuvo %v", project.ErrNoExiste, err)
	}
}

func TestProyectos_Actualizar(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "proyectos.json")
	r, _ := NuevoProyectos(ruta)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := r.Agregar(project.Proyecto{Nombre: nombre, Estado: project.Planificado}); err != nil {
			t.Fatal(err)
		}
	}

	p, _ := r.Buscar("P1")
	p.Nombre = "Software Metrics"
	if err := r.Actualizar(p); err != nil {
		t.Fatalf("Actualizar(P1): error inesperado: %v", err)
	}

	reabierto, err := NuevoProyectos(ruta)
	if err != nil {
		t.Fatal(err)
	}
	lista := reabierto.Listar()
	if len(lista) != 2 || lista[0].ID != "P1" || lista[0].Nombre != "Software Metrics" || lista[1].Nombre != "Otro" {
		t.Errorf("después de Actualizar y reabrir = %+v, se esperaba [P1 Software Metrics, P2 Otro]", lista)
	}

	if err := r.Actualizar(project.Proyecto{ID: "P9", Nombre: "Nadie"}); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("Actualizar(P9): se esperaba %v, se obtuvo %v", project.ErrNoExiste, err)
	}
	if n := len(r.Listar()); n != 2 {
		t.Errorf("Actualizar(P9) dejó %d proyectos, se esperaban 2: no debe agregar", n)
	}
}
