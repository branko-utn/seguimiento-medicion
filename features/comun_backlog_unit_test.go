package features

import (
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// Tests de los pasos compartidos de backlog: todavía ningún .feature activo los usa (HU-08, HU-11, HU-12 y HU-15
// siguen con @pendiente), así que se prueban acá.

func mundoConProyecto(t *testing.T) *mundo {
	t.Helper()
	m := nuevoMundo()
	if err := m.queExisteUnProyecto("Demo"); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestCargarHistorias(t *testing.T) {
	m := mundoConProyecto(t)
	err := m.cargarHistorias([]filaHistoria{
		{historia: "HU-1", titulo: "Crear proyecto", storyPoints: 3, criterios: 2},
		{historia: "HU-2", titulo: "Crear sprint", storyPoints: 0, criterios: 1},
		{historia: "HU-3", titulo: "Cerrar sprint", storyPoints: 5, criterios: 0},
	})
	if err != nil {
		t.Fatal(err)
	}

	h, err := m.backlog.Buscar(domain.HistoriaID(m.id("historia", "HU-1")))
	if err != nil {
		t.Fatal(err)
	}
	if h.Titulo != "Crear proyecto" || h.StoryPoints != 3 || len(h.Criterios) != 2 {
		t.Errorf("HU-1 = %+v, se esperaba 3 SP y 2 criterios", h)
	}
	h3, _ := m.backlog.Buscar(domain.HistoriaID(m.id("historia", "HU-3")))
	if h3.StoryPoints != 5 || len(h3.Criterios) != 0 {
		t.Errorf("HU-3 = %+v, se esperaba 5 SP y sin criterios", h3)
	}
}

// TestCargarHistorias_prioridad: la columna prioridad es opcional y vale Media si no está (spec HU-09, nota para los
// pasos).
func TestCargarHistorias_prioridad(t *testing.T) {
	m := mundoConProyecto(t)
	err := m.cargarHistorias([]filaHistoria{
		{historia: "HU-1", titulo: "Registrar defecto", storyPoints: 2, criterios: 1, prioridad: backlog.PrioridadBaja},
		{historia: "HU-2", titulo: "Crear sprint", storyPoints: 5, criterios: 1, prioridad: backlog.PrioridadAlta},
		{historia: "HU-3", titulo: "Ver dashboard", storyPoints: 8, criterios: 1},
	})
	if err != nil {
		t.Fatal(err)
	}

	for historia, want := range map[string]backlog.Prioridad{
		"HU-1": backlog.PrioridadBaja, "HU-2": backlog.PrioridadAlta, "HU-3": backlog.PrioridadMedia,
	} {
		h, err := m.buscarHistoria(historia)
		if err != nil {
			t.Fatal(err)
		}
		if h.Prioridad != want {
			t.Errorf("%s tiene prioridad %v, se esperaba %v", historia, h.Prioridad, want)
		}
	}
}

func TestQueLaHistoriaEstaEnEstado(t *testing.T) {
	m := mundoConProyecto(t)
	if err := m.queElBacklogTieneLaHistoria("Demo", "HU-1", "Login"); err != nil {
		t.Fatal(err)
	}
	if err := m.queLaHistoriaEstaEnEstado("HU-1", "En progreso"); err != nil {
		t.Fatal(err)
	}
	if err := m.laHistoriaEstaEnEstado("HU-1", "En progreso"); err != nil {
		t.Error(err)
	}
	if err := m.queLaHistoriaEstaEnEstado("HU-1", "Bloqueada"); err == nil {
		t.Error("un estado que no existe debería dar error")
	}
}

func TestStoryPointsDeUnaHistoria(t *testing.T) {
	m := mundoConProyecto(t)
	if err := m.queElBacklogTieneLaHistoria("Demo", "HU-1", "Login"); err != nil {
		t.Fatal(err)
	}
	if err := m.queLaHistoriaTieneStoryPoints("HU-1", 3); err != nil {
		t.Fatal(err)
	}
	if err := m.laHistoriaTieneStoryPoints("HU-1", 3); err != nil {
		t.Error(err)
	}
	if err := m.laHistoriaTieneStoryPoints("HU-1", 5); err == nil {
		t.Error("con otro valor, el paso Entonces debería fallar")
	}
}

func TestBuscarHistoria_inexistente(t *testing.T) {
	m := mundoConProyecto(t)
	if _, err := m.buscarHistoria("HU-99"); err == nil {
		t.Error("buscar una historia que no existe debería dar error")
	}
}
