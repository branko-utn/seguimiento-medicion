package memory

import (
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

func TestIntegrantes_AgregarGeneraIDsYFiltraPorProyecto(t *testing.T) {
	r, err := NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}

	altas := []struct {
		proyecto domain.ProyectoID
		nombre   string
		quiereID domain.IntegranteID
	}{
		{"P1", "Ana", "I1"},
		{"P2", "Beto", "I2"},
		{"P1", "Carla", "I3"},
	}
	for _, a := range altas {
		i, err := r.Agregar(project.Integrante{ProyectoID: a.proyecto, Nombre: a.nombre})
		if err != nil {
			t.Fatalf("Agregar %s: %v", a.nombre, err)
		}
		if i.ID != a.quiereID {
			t.Errorf("Agregar %s: ID = %q, se esperaba %q", a.nombre, i.ID, a.quiereID)
		}
	}

	delP1 := r.DelProyecto("P1")
	if len(delP1) != 2 || delP1[0].Nombre != "Ana" || delP1[1].Nombre != "Carla" {
		t.Errorf("DelProyecto(P1) = %+v, se esperaba [Ana Carla] en ese orden", delP1)
	}
	if n := len(r.DelProyecto("P9")); n != 0 {
		t.Errorf("DelProyecto(P9) tiene %d integrantes, se esperaba 0", n)
	}
}

func TestIntegrantes_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "integrantes.json")
	r, err := NuevoIntegrantes(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(project.Integrante{ProyectoID: "P1", Nombre: "Ana", Rol: project.AgileEnabler}); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevoIntegrantes(ruta)
	if err != nil {
		t.Fatal(err)
	}
	lista := reabierto.DelProyecto("P1")
	if len(lista) != 1 || lista[0].ID != "I1" || lista[0].Rol != project.AgileEnabler {
		t.Fatalf("al reabrir se obtuvo %+v", lista)
	}

	i, err := reabierto.Agregar(project.Integrante{ProyectoID: "P1", Nombre: "Beto"})
	if err != nil {
		t.Fatal(err)
	}
	if i.ID != "I2" {
		t.Errorf("después de reabrir, ID = %q, se esperaba %q", i.ID, "I2")
	}
}

func TestIntegrantes_DelProyectoDevuelveCopia(t *testing.T) {
	r, _ := NuevoIntegrantes("")
	if _, err := r.Agregar(project.Integrante{ProyectoID: "P1", Nombre: "Ana"}); err != nil {
		t.Fatal(err)
	}

	lista := r.DelProyecto("P1")
	if len(lista) != 1 {
		t.Fatalf("DelProyecto(P1) tiene %d integrantes, se esperaba 1", len(lista))
	}
	lista[0].Nombre = "Cambiado"

	if got := r.DelProyecto("P1")[0].Nombre; got != "Ana" {
		t.Errorf("modificar lo que devuelve DelProyecto cambió el repositorio: nombre = %q", got)
	}
}
