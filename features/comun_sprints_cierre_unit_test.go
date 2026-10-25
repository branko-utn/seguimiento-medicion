package features

import (
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// Todavía ningún .feature activo usa "que el sprint … está cerrado" (HU-11 y HU-12), así que se prueba acá.

func TestQueElSprintEstaCerrado(t *testing.T) {
	m := nuevoMundo()
	if err := m.queExisteUnProyecto("Demo"); err != nil {
		t.Fatal(err)
	}
	if err := m.existeUnSprintActivoEn("S1", "Demo"); err != nil {
		t.Fatal(err)
	}

	if err := m.queElSprintEstaCerrado("S1"); err != nil {
		t.Fatalf("queElSprintEstaCerrado: %v", err)
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", "S1")))
	if err != nil {
		t.Fatal(err)
	}
	if sp.Estado != sprint.Cerrado {
		t.Errorf("el sprint quedó en estado %q, se esperaba %q", sp.Estado, sprint.Cerrado)
	}
}

func TestQueElSprintEstaCerrado_inexistente(t *testing.T) {
	m := nuevoMundo()
	if err := m.queElSprintEstaCerrado("S9"); err == nil {
		t.Error("cerrar un sprint que no existe debería dar error")
	}
}

func TestQueLaHistoriaEstaCompletada(t *testing.T) {
	m := nuevoMundo()
	if err := m.queExisteUnProyecto("Demo"); err != nil {
		t.Fatal(err)
	}
	if err := m.queElBacklogTieneLaHistoria("Demo", "HU-1", "Crear proyecto"); err != nil {
		t.Fatal(err)
	}

	if err := m.queLaHistoriaEstaCompletada("HU-1"); err != nil {
		t.Fatalf("queLaHistoriaEstaCompletada: %v", err)
	}
	if err := m.laHistoriaEstaEnEstado("HU-1", "Completada"); err != nil {
		t.Error(err)
	}
	if err := m.queLaHistoriaEstaCompletada("HU-99"); err == nil {
		t.Error("una historia que no existe debería dar error")
	}
}

func TestQueLaHistoriaYaEstaCompletadaEnElSprint(t *testing.T) {
	m := nuevoMundo()
	if err := m.queExisteUnProyecto("Demo"); err != nil {
		t.Fatal(err)
	}
	if err := m.queElBacklogTieneLaHistoria("Demo", "HU-1", "Crear proyecto"); err != nil {
		t.Fatal(err)
	}
	if err := m.existeUnSprintActivoEn("S1", "Demo"); err != nil {
		t.Fatal(err)
	}

	// La historia no estaba asignada: el paso la asigna al sprint y la completa.
	if err := m.queLaHistoriaYaEstaCompletadaEnElSprint("HU-1", "S1"); err != nil {
		t.Fatalf("queLaHistoriaYaEstaCompletadaEnElSprint: %v", err)
	}
	// Repetirlo no la duplica en las listas del sprint.
	if err := m.queLaHistoriaYaEstaCompletadaEnElSprint("HU-1", "S1"); err != nil {
		t.Fatal(err)
	}

	h, err := m.buscarHistoria("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", "S1")))
	if err != nil {
		t.Fatal(err)
	}
	if h.Estado.String() != "Completada" || h.SprintID != sp.ID {
		t.Errorf("historia = %+v, se esperaba Completada y asignada a %s", h, sp.ID)
	}
	if len(sp.Historias) != 1 || len(sp.Completadas) != 1 || sp.Completadas[0] != h.ID {
		t.Errorf("sprint: historias = %v, completadas = %v; se esperaba [%s] en las dos", sp.Historias, sp.Completadas, h.ID)
	}
}
