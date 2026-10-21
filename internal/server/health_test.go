package server

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	casos := []struct {
		nombre string
		metodo string
		ruta   string
		codigo int
		cuerpo string
		tipo   string
	}{
		{nombre: "GET responde ok", metodo: http.MethodGet, ruta: "/health", codigo: http.StatusOK, cuerpo: `{"estado":"ok"}`, tipo: "application/json"},
		{nombre: "POST no está permitido", metodo: http.MethodPost, ruta: "/health", codigo: http.StatusMethodNotAllowed},
		{nombre: "ruta inexistente", metodo: http.MethodGet, ruta: "/salud", codigo: http.StatusNotFound},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, c.metodo, c.ruta)

			if rec.Code != c.codigo {
				t.Fatalf("código = %d, se esperaba %d", rec.Code, c.codigo)
			}
			if c.cuerpo != "" && rec.Body.String() != c.cuerpo {
				t.Errorf("cuerpo = %q, se esperaba %q", rec.Body.String(), c.cuerpo)
			}
			if c.tipo != "" && rec.Header().Get("Content-Type") != c.tipo {
				t.Errorf("Content-Type = %q, se esperaba %q", rec.Header().Get("Content-Type"), c.tipo)
			}
		})
	}
}
