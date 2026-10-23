package memory

import (
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
