package backlog

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
)

func TestEstimar(t *testing.T) {
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)

	// RN-15.1: solo valores de la escala. El 3 es el caso intermedio.
	if err := h.Estimar(3); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if h.StoryPoints != 3 {
		t.Errorf("StoryPoints = %d, se esperaba 3", h.StoryPoints)
	}
}

func TestEstimar_aceptaTodaLaEscala(t *testing.T) {
	// Los casos límite son 1 y 21, los extremos de la escala.
	for _, sp := range estimation.Escala {
		h, _ := Nueva("P1", "Login", "", PrioridadMedia)
		if err := h.Estimar(sp); err != nil {
			t.Errorf("Estimar(%d): %v", sp, err)
			continue
		}
		if h.StoryPoints != sp {
			t.Errorf("Estimar(%d) dejó StoryPoints = %d", sp, h.StoryPoints)
		}
	}
}

func TestEstimar_reemplazaElValorAnterior(t *testing.T) {
	// RN-15.2: se puede re-estimar mientras no esté completada; re-estimar una historia de un sprint activo
	// es válido (cambian los SP planificados del sprint).
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)
	if err := h.Estimar(3); err != nil {
		t.Fatal(err)
	}

	if err := h.Estimar(5); err != nil {
		t.Fatalf("re-estimar debería ser válido: %v", err)
	}
	if h.StoryPoints != 5 {
		t.Errorf("StoryPoints = %d, se esperaba 5 (el valor anterior se reemplaza)", h.StoryPoints)
	}
}

func TestEstimar_fueraDeEscala(t *testing.T) {
	casos := []struct {
		nombre string
		sp     int
	}{
		{"cero es sin estimar (RN-15.3)", 0},
		{"fuera de la escala", 4},
		{"negativo", -1},
		{"muy grande", 34},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			h, _ := Nueva("P1", "Login", "", PrioridadMedia)

			if err := h.Estimar(caso.sp); !errors.Is(err, estimation.ErrStoryPointsInvalidos) {
				t.Fatalf("se esperaba ErrStoryPointsInvalidos, se obtuvo %v", err)
			}
			if h.StoryPoints != 0 {
				t.Errorf("StoryPoints = %d, se esperaba 0: un valor inválido no se asigna", h.StoryPoints)
			}
		})
	}
}

func TestEstimar_historiaCompletada(t *testing.T) {
	// RN-15.2: no se puede re-estimar una historia completada.
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)
	if err := h.Estimar(3); err != nil {
		t.Fatal(err)
	}
	h.Estado = Completada

	if err := h.Estimar(5); !errors.Is(err, ErrHistoriaCompletada) {
		t.Fatalf("se esperaba ErrHistoriaCompletada, se obtuvo %v", err)
	}
	if h.StoryPoints != 3 {
		t.Errorf("StoryPoints = %d, se esperaba 3: la estimación anterior no cambia", h.StoryPoints)
	}
}
