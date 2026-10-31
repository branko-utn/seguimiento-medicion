package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// HU-32.1: las pantallas de proyectos y sprints usan los componentes de TEC-21, sin estilos en línea ni emojis.

// emojis son los símbolos que se usaban como íconos y que HU-32 reemplaza por insignias (RN-32.5).
var emojis = []string{"✔", "✓", "✅", "❌", "⚠", "☕"}

func TestHU32_1_pantallasConLosComponentes(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)
	if err := e.svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}

	// Las pantallas con tablas usan table.tabla; todas tienen el encabezado de página.
	pantallas := []struct {
		ruta      string
		conTablas bool
	}{
		{"/proyectos", true},
		{"/proyectos/nuevo", false},
		{"/proyectos/P1/editar", false},
		{"/sprints", true},
		{"/sprints/nuevo?proyecto=P1", false},
		{"/sprints/S1", true},
	}
	for _, p := range pantallas {
		t.Run(p.ruta, func(t *testing.T) {
			rec := servir(e.router, httptest.NewRequest(http.MethodGet, p.ruta, nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("código = %d, se esperaba 200", rec.Code)
			}
			cuerpo := rec.Body.String()
			contieneTodo(t, cuerpo, `class="encabezado"`)
			if p.conTablas && (!strings.Contains(cuerpo, `<table class="tabla">`) || strings.Contains(cuerpo, "<table>")) {
				t.Error("todas las tablas deben usar la clase tabla (TEC-21)")
			}
			if strings.Contains(cuerpo, `style="`) {
				t.Error("la página tiene estilos en línea (RN-32.5)")
			}
			for _, emoji := range emojis {
				if strings.Contains(cuerpo, emoji) {
					t.Errorf("la página usa %q como ícono (RN-32.5): va una insignia", emoji)
				}
			}
		})
	}
}

func TestHU32_1_listasVaciasConAccion(t *testing.T) {
	e := nuevoServidorSprints(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints", nil))

	// Demo y Otro no tienen sprints: cada uno muestra el estado vacío con el enlace para crear uno (RN-32.4).
	cuerpo := rec.Body.String()
	if n := strings.Count(cuerpo, `class="vacio"`); n < 2 {
		t.Errorf("hay %d estados vacíos, se esperaba al menos uno por proyecto sin sprints", n)
	}
	contieneTodo(t, cuerpo, "Todavía no hay sprints", `href="/sprints/nuevo?proyecto=P1"`, `href="/sprints/nuevo?proyecto=P2"`)
}
