package backlog

import (
	"errors"
	"testing"
)

func TestAgregarCriterio_losAgregaAlFinal(t *testing.T) {
	// RN-07.3: los criterios se identifican por su posición, empezando en 1.
	h, err := Nueva("P1", "Login", "", PrioridadMedia)
	if err != nil {
		t.Fatal(err)
	}
	if err := h.AgregarCriterio("El nombre es obligatorio"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if err := h.AgregarCriterio("  La fecha de fin no es anterior a la de inicio  "); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	quiere := []string{"El nombre es obligatorio", "La fecha de fin no es anterior a la de inicio"}
	if len(h.Criterios) != len(quiere) {
		t.Fatalf("hay %d criterios, se esperaban %d: %+v", len(h.Criterios), len(quiere), h.Criterios)
	}
	for n, texto := range quiere {
		if h.Criterios[n].Texto != texto {
			t.Errorf("criterio %d = %q, se esperaba %q (el nuevo va al final y sin espacios sobrantes)", n+1, h.Criterios[n].Texto, texto)
		}
	}
	if !h.TieneCriterios() {
		t.Error("con criterios, TieneCriterios() debería ser true")
	}
}

func TestAgregarCriterio_errores(t *testing.T) {
	casos := []struct {
		nombre string
		texto  string
	}{
		{"criterio vacío", ""},
		{"criterio con solo espacios", "   "},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			h, _ := Nueva("P1", "Login", "", PrioridadMedia)

			if err := h.AgregarCriterio(caso.texto); !errors.Is(err, ErrCriterioVacio) {
				t.Fatalf("RN-07.1: se esperaba ErrCriterioVacio, se obtuvo %v", err)
			}
			if len(h.Criterios) != 0 {
				t.Errorf("no se debe agregar un criterio inválido; hay %d criterios", len(h.Criterios))
			}
		})
	}
}

func TestQuitarCriterio_renumeraLosSiguientes(t *testing.T) {
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)
	for _, texto := range []string{"A", "B", "C"} {
		if err := h.AgregarCriterio(texto); err != nil {
			t.Fatal(err)
		}
	}

	if err := h.QuitarCriterio(2); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// RN-07.3: el 3 pasa a ser el 2.
	quiere := []string{"A", "C"}
	if len(h.Criterios) != len(quiere) {
		t.Fatalf("quedan %d criterios, se esperaban %d", len(h.Criterios), len(quiere))
	}
	for n, texto := range quiere {
		if h.Criterios[n].Texto != texto {
			t.Errorf("criterio %d = %q, se esperaba %q", n+1, h.Criterios[n].Texto, texto)
		}
	}
}

func TestQuitarCriterio_elUnico(t *testing.T) {
	// Caso límite: quitar el único criterio deja la historia sin criterios, y es válido.
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)
	if err := h.AgregarCriterio("A"); err != nil {
		t.Fatal(err)
	}
	if len(h.Criterios) != 1 {
		t.Fatalf("preparación: hay %d criterios, se esperaba 1", len(h.Criterios))
	}

	if err := h.QuitarCriterio(1); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(h.Criterios) != 0 {
		t.Errorf("quedaron %d criterios, se esperaba 0", len(h.Criterios))
	}
	if h.TieneCriterios() {
		t.Error("sin criterios, TieneCriterios() debería ser false (RN-07.2)")
	}
}

func TestQuitarCriterio_errores(t *testing.T) {
	casos := []struct {
		nombre       string
		conCriterios []string
		numero       int
	}{
		{"criterio 0", []string{"A"}, 0},
		{"mayor que la cantidad", []string{"A"}, 2},
		{"historia sin criterios", nil, 1},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			h, _ := Nueva("P1", "Login", "", PrioridadMedia)
			for _, texto := range caso.conCriterios {
				if err := h.AgregarCriterio(texto); err != nil {
					t.Fatal(err)
				}
			}

			if err := h.QuitarCriterio(caso.numero); !errors.Is(err, ErrCriterioNoExiste) {
				t.Fatalf("se esperaba ErrCriterioNoExiste, se obtuvo %v", err)
			}
			if len(h.Criterios) != len(caso.conCriterios) {
				t.Errorf("quedan %d criterios, se esperaban %d: no se debe quitar uno que no existe", len(h.Criterios), len(caso.conCriterios))
			}
		})
	}
}

func TestTieneCriterios(t *testing.T) {
	// Lo usa HU-11 (RN-07.2): una historia sin criterios no puede asignarse a un sprint.
	h, _ := Nueva("P1", "Login", "", PrioridadMedia)
	if h.TieneCriterios() {
		t.Error("una historia nueva no tiene criterios (RN-05.4)")
	}
	if err := h.AgregarCriterio("A"); err != nil {
		t.Fatal(err)
	}
	if !h.TieneCriterios() {
		t.Error("con un criterio, TieneCriterios() debería ser true")
	}
}
