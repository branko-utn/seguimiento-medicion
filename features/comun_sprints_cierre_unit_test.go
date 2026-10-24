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
