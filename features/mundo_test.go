package features

import "time"

// mundo es el estado de un escenario (TEC-10).
type mundo struct {
	hoy   time.Time
	err   error
	alias map[string]string
}

func nuevoMundo() *mundo { return &mundo{} }

func (m *mundo) recordar(tipo, nombre, id string) {}

func (m *mundo) id(tipo, nombre string) string { return "" }

func (m *mundo) resultado(err error) error { return nil }

func (m *mundo) sinError() error { return nil }

func (m *mundo) veoElError(mensaje string) error { return nil }

func (m *mundo) hoyEs(fecha string) error { return nil }

func parseFecha(s string) (time.Time, error) { return time.Time{}, nil }

func lista(s string) []string { return nil }
