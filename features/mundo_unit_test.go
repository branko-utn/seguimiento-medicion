package features

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMundo_alias(t *testing.T) {
	m := nuevoMundo()
	m.recordar("sprint", "S9", "S2")

	if got := m.id("sprint", "S9"); got != "S2" {
		t.Errorf("id(sprint, S9) = %q, se esperaba S2", got)
	}
	// Un nombre sin alias se usa tal cual: sirve para referirse a algo que no existe ("HU-99").
	if got := m.id("historia", "HU-99"); got != "HU-99" {
		t.Errorf("id(historia, HU-99) = %q, se esperaba HU-99", got)
	}
	// El mismo nombre en tipos distintos no se mezcla.
	if got := m.id("proyecto", "S9"); got != "S9" {
		t.Errorf("id(proyecto, S9) = %q, se esperaba S9", got)
	}
}

func TestMundo_veoElError(t *testing.T) {
	m := nuevoMundo()

	if err := m.veoElError("el nombre es obligatorio"); err == nil {
		t.Error("sin error previo, el paso debería fallar")
	}

	_ = m.resultado(errors.New("el nombre es obligatorio"))
	if err := m.veoElError("el nombre es obligatorio"); err != nil {
		t.Errorf("con el mismo mensaje, el paso debería pasar: %v", err)
	}
	if err := m.sinError(); err != nil {
		t.Errorf("después de verificar el error esperado, los pasos siguientes no deberían verlo: %v", err)
	}

	_ = m.resultado(errors.New("el nombre es obligatorio"))
	if err := m.veoElError("otro mensaje"); err == nil {
		t.Error("con otro mensaje, el paso debería fallar")
	}
}

func TestMundo_sinError(t *testing.T) {
	m := nuevoMundo()
	if err := m.sinError(); err != nil {
		t.Errorf("sin error previo: %v", err)
	}
	_ = m.resultado(errors.New("falló"))
	if err := m.sinError(); err == nil {
		t.Error("con un error previo, sinError debería fallar")
	}
}

func TestMundo_hoyEs(t *testing.T) {
	m := nuevoMundo()
	if err := m.hoyEs("2026-10-24"); err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 10, 24, 0, 0, 0, 0, time.Local); !m.hoy.Equal(want) {
		t.Errorf("hoy = %v, se esperaba %v", m.hoy, want)
	}
	if err := m.hoyEs("24/10/2026"); err == nil {
		t.Error("una fecha con otro formato debería dar error")
	}
}

func TestLista(t *testing.T) {
	casos := map[string][]string{
		"HU-1, HU-2": {"HU-1", "HU-2"},
		"HU-1":       {"HU-1"},
		"A,B , C":    {"A", "B", "C"},
		"":           nil,
	}
	for entrada, want := range casos {
		if got := lista(entrada); !reflect.DeepEqual(got, want) {
			t.Errorf("lista(%q) = %#v, se esperaba %#v", entrada, got, want)
		}
	}
}
