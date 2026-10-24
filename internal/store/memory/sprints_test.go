package memory

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

func TestSprints_AgregarGeneraIDsYFiltraPorProyecto(t *testing.T) {
	r, err := NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}

	altas := []struct {
		proyecto domain.ProyectoID
		goal     string
		quiereID domain.SprintID
	}{
		{"P1", "Primero", "S1"},
		{"P2", "Otro proyecto", "S2"},
		{"P1", "Segundo", "S3"},
	}
	for _, a := range altas {
		s, err := r.Agregar(sprint.Sprint{ProyectoID: a.proyecto, Goal: a.goal, Estado: sprint.Planificado})
		if err != nil {
			t.Fatalf("Agregar %s: %v", a.goal, err)
		}
		if s.ID != a.quiereID {
			t.Errorf("Agregar %s: ID = %q, se esperaba %q", a.goal, s.ID, a.quiereID)
		}
	}

	delP1 := r.DelProyecto("P1")
	if len(delP1) != 2 || delP1[0].Goal != "Primero" || delP1[1].Goal != "Segundo" {
		t.Errorf("DelProyecto(P1) = %+v, se esperaba [Primero Segundo] en ese orden", delP1)
	}
	if n := len(r.DelProyecto("P9")); n != 0 {
		t.Errorf("DelProyecto(P9) tiene %d sprints, se esperaba 0", n)
	}
}

func TestSprints_BuscarYActualizar(t *testing.T) {
	r, _ := NuevoSprints("")
	if _, err := r.Agregar(sprint.Sprint{ProyectoID: "P1", Goal: "Primero", Estado: sprint.Planificado}); err != nil {
		t.Fatal(err)
	}

	s, err := r.Buscar("S1")
	if err != nil {
		t.Fatalf("Buscar(S1): error inesperado: %v", err)
	}
	if s.ID != "S1" || s.Goal != "Primero" {
		t.Fatalf("Buscar(S1) = %+v, se esperaba el sprint Primero", s)
	}

	s.Estado = sprint.Activo
	if err := r.Actualizar(s); err != nil {
		t.Fatalf("Actualizar(S1): error inesperado: %v", err)
	}
	if got, _ := r.Buscar("S1"); got.Estado != sprint.Activo {
		t.Errorf("después de Actualizar, estado = %v, se esperaba %v", got.Estado, sprint.Activo)
	}

	if _, err := r.Buscar("S9"); !errors.Is(err, sprint.ErrNoExiste) {
		t.Errorf("Buscar(S9): se esperaba %v, se obtuvo %v", sprint.ErrNoExiste, err)
	}
	if err := r.Actualizar(sprint.Sprint{ID: "S9"}); !errors.Is(err, sprint.ErrNoExiste) {
		t.Errorf("Actualizar(S9): se esperaba %v, se obtuvo %v", sprint.ErrNoExiste, err)
	}
}

func TestSprints_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "sprints.json")
	r, err := NuevoSprints(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(sprint.Sprint{ProyectoID: "P1", Goal: "Primero", Estado: sprint.Planificado}); err != nil {
		t.Fatal(err)
	}
	s, _ := r.Buscar("S1")
	s.Estado = sprint.Activo
	if err := r.Actualizar(s); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevoSprints(ruta)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reabierto.Buscar("S1")
	if err != nil || got.Goal != "Primero" || got.Estado != sprint.Activo {
		t.Fatalf("al reabrir se obtuvo %+v (error %v)", got, err)
	}

	nuevo, err := reabierto.Agregar(sprint.Sprint{ProyectoID: "P1", Goal: "Segundo"})
	if err != nil {
		t.Fatal(err)
	}
	if nuevo.ID != "S2" {
		t.Errorf("después de reabrir, ID = %q, se esperaba %q", nuevo.ID, "S2")
	}
}

// Las listas del sprint (Historias, Completadas) no se comparten con lo guardado: modificar lo que devuelve el
// repositorio no lo cambia sin pasar por Actualizar.
func TestSprints_DevuelveCopias(t *testing.T) {
	r, _ := NuevoSprints("")
	if _, err := r.Agregar(sprint.Sprint{ProyectoID: "P1", Goal: "Primero", Historias: []domain.HistoriaID{"HU-1"}}); err != nil {
		t.Fatal(err)
	}

	s, err := r.Buscar("S1")
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Historias) != 1 {
		t.Fatalf("Buscar(S1) = %+v, se esperaba el sprint con la historia HU-1", s)
	}
	s.Historias[0] = "HU-9"
	lista := r.DelProyecto("P1")
	if len(lista) != 1 || len(lista[0].Historias) != 1 {
		t.Fatalf("DelProyecto(P1) = %+v, se esperaba 1 sprint con 1 historia", lista)
	}
	lista[0].Goal = "Cambiado"
	lista[0].Historias[0] = "HU-8"

	got, _ := r.Buscar("S1")
	if got.Goal != "Primero" || got.Historias[0] != "HU-1" {
		t.Errorf("modificar las copias cambió el repositorio: %+v", got)
	}
}
