package memory

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

func TestHistorias_AgregarGeneraIDsUnicos(t *testing.T) {
	// RN-05.1: el ID es único en todo el sistema y correlativo; no se reinicia por proyecto.
	r, err := NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}

	altas := []struct {
		proyecto domain.ProyectoID
		titulo   string
		quiereID domain.HistoriaID
	}{
		{"P1", "Crear proyecto", "HU-1"},
		{"P1", "Crear sprint", "HU-2"},
		{"P2", "Login", "HU-3"},
		{"P1", "Cerrar sprint", "HU-4"},
	}
	for _, a := range altas {
		h, err := r.Agregar(backlog.Historia{ProyectoID: a.proyecto, Titulo: a.titulo})
		if err != nil {
			t.Fatalf("Agregar %s: %v", a.titulo, err)
		}
		if h.ID != a.quiereID {
			t.Errorf("Agregar %s: ID = %q, se esperaba %q", a.titulo, h.ID, a.quiereID)
		}
	}

	lista := r.Listar()
	if len(lista) != len(altas) {
		t.Fatalf("Listar() tiene %d historias, se esperaban %d", len(lista), len(altas))
	}
	for n, a := range altas {
		if lista[n].Titulo != a.titulo {
			t.Errorf("historia %d = %q, se esperaba %q (en orden de creación)", n, lista[n].Titulo, a.titulo)
		}
	}
}

func TestHistorias_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "backlog.json")
	r, err := NuevoHistorias(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login", Prioridad: backlog.PrioridadAlta, Estado: backlog.Pendiente}); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevoHistorias(ruta)
	if err != nil {
		t.Fatal(err)
	}
	lista := reabierto.Listar()
	if len(lista) != 1 || lista[0].ID != "HU-1" || lista[0].Titulo != "Login" ||
		lista[0].Prioridad != backlog.PrioridadAlta || lista[0].Estado != backlog.Pendiente {
		t.Fatalf("al reabrir se obtuvo %+v", lista)
	}

	// Después de reabrir sigue numerando donde quedó.
	h, err := reabierto.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Crear sprint"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != "HU-2" {
		t.Errorf("después de reabrir, ID = %q, se esperaba %q", h.ID, "HU-2")
	}
}

func TestHistorias_ListarDevuelveCopia(t *testing.T) {
	r, _ := NuevoHistorias("")
	if _, err := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login"}); err != nil {
		t.Fatal(err)
	}

	lista := r.Listar()
	if len(lista) != 1 {
		t.Fatalf("Listar() tiene %d historias, se esperaba 1", len(lista))
	}
	lista[0].Titulo = "Cambiado"

	if got := r.Listar()[0].Titulo; got != "Login" {
		t.Errorf("modificar lo que devuelve Listar cambió el repositorio: título = %q", got)
	}
}

func TestHistorias_Buscar(t *testing.T) {
	r, _ := NuevoHistorias("")
	h, _ := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login", Criterios: []backlog.Criterio{{Texto: "A"}}})

	got, err := r.Buscar(h.ID)
	if err != nil {
		t.Fatalf("Buscar(%s): %v", h.ID, err)
	}
	if got.Titulo != "Login" {
		t.Errorf("Buscar(%s) = %+v", h.ID, got)
	}

	if _, err := r.Buscar("HU-99"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Buscar(HU-99) = %v, se esperaba backlog.ErrNoExiste", err)
	}
}

func TestHistorias_BuscarDevuelveCopiaIndependiente(t *testing.T) {
	// Modificar los criterios de la historia devuelta no debe cambiar lo guardado sin pasar por Guardar.
	r, _ := NuevoHistorias("")
	h, _ := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login", Criterios: []backlog.Criterio{{Texto: "A"}, {Texto: "B"}}})

	copia, _ := r.Buscar(h.ID)
	copia.Criterios[0].Texto = "cambiado"

	guardada, _ := r.Buscar(h.ID)
	if guardada.Criterios[0].Texto != "A" {
		t.Errorf("lo guardado cambió sin Guardar: %+v", guardada.Criterios)
	}
}

func TestHistorias_Guardar(t *testing.T) {
	r, _ := NuevoHistorias("")
	h, _ := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login"})

	h.StoryPoints = 5
	h.Criterios = []backlog.Criterio{{Texto: "A"}}
	if err := r.Guardar(h); err != nil {
		t.Fatalf("Guardar: %v", err)
	}
	got, _ := r.Buscar(h.ID)
	if got.StoryPoints != 5 || len(got.Criterios) != 1 {
		t.Errorf("después de Guardar, Buscar = %+v", got)
	}

	if err := r.Guardar(backlog.Historia{ID: "HU-99"}); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Guardar(HU-99) = %v, se esperaba backlog.ErrNoExiste", err)
	}
}
