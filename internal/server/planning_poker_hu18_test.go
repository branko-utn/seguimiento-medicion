package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// votacionRevelada abre la votación de HU-1, Ana vota 5, Beto vota 8 y el Agile Enabler revela (HU-16 y HU-17).
func votacionRevelada(t *testing.T, router http.Handler) {
	t.Helper()
	pasos := []struct {
		ruta    string
		valores url.Values
	}{
		{"/planning-poker/abrir", url.Values{"historia": {"HU-1"}}},
		{"/planning-poker/votar", url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}}},
		{"/planning-poker/votar", url.Values{"historia": {"HU-1"}, "integrante": {"I2"}, "carta": {"8"}}},
		{"/planning-poker/revelar", url.Values{"historia": {"HU-1"}, "revelo": {"I3"}}},
	}
	for _, p := range pasos {
		if rec := servir(router, pedidoFormulario(http.MethodPost, p.ruta, p.valores)); rec.Code != http.StatusSeeOther {
			t.Fatalf("POST %s: código = %d, se esperaba 303", p.ruta, rec.Code)
		}
	}
}

// HU-18: después de revelar se abre una ronda nueva sin votos; al revelarla, la ronda anterior queda en el historial.
func TestPlanningPoker_nuevaRonda(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	votacionRevelada(t, router)

	rec := servir(router, pedidoFormulario(http.MethodPost, "/planning-poker/nueva-ronda", url.Values{"historia": {"HU-1"}}))
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/planning-poker" {
		t.Fatalf("nueva ronda: código = %d, Location = %q; se esperaba 303 a /planning-poker", rec.Code, rec.Header().Get("Location"))
	}
	cuerpo := servir(router, httptest.NewRequest(http.MethodGet, "/planning-poker", nil)).Body.String()
	contieneTodo(t, cuerpo, "<td>2</td>", "Votás como", "Votaron: todavía nadie")

	for _, v := range []url.Values{
		{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}},
		{"historia": {"HU-1"}, "integrante": {"I2"}, "carta": {"5"}},
	} {
		servir(router, pedidoFormulario(http.MethodPost, "/planning-poker/votar", v))
	}
	servir(router, pedidoFormulario(http.MethodPost, "/planning-poker/revelar", url.Values{"historia": {"HU-1"}, "revelo": {"I3"}}))
	cuerpo = servir(router, httptest.NewRequest(http.MethodGet, "/planning-poker", nil)).Body.String()
	contieneTodo(t, cuerpo, "Historial de rondas", "Ronda 1")
}

// HU-18: la estimación acordada queda como Story Points de la historia y la votación se cierra; la historia puede
// abrir otra votación (RN-18.4 y RN-18.5).
func TestPlanningPoker_acordar(t *testing.T) {
	router, historias := nuevoServidorPoker(t)
	votacionRevelada(t, router)

	rec := servir(router, pedidoFormulario(http.MethodPost, "/planning-poker/acordar", url.Values{"historia": {"HU-1"}, "estimacion": {"5"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("acordar: código = %d, se esperaba 303", rec.Code)
	}
	h, err := historias.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	if h.StoryPoints != 5 {
		t.Errorf("la historia tiene %d Story Points, se esperaban 5", h.StoryPoints)
	}
	cuerpo := servir(router, httptest.NewRequest(http.MethodGet, "/planning-poker", nil)).Body.String()
	contieneTodo(t, cuerpo, "Estimación acordada: 5 SP", "Abrir nueva votación")
}

func TestPlanningPoker_nuevaRondaYAcordar_errores(t *testing.T) {
	casos := []struct {
		nombre   string
		revelada bool
		abierta  bool
		ruta     string
		valores  url.Values
		mensaje  string
	}{
		{"nueva ronda sin revelar", false, true, "/planning-poker/nueva-ronda", url.Values{"historia": {"HU-1"}}, "primero hay que revelar los votos de la ronda"},
		{"acordar sin revelar", false, true, "/planning-poker/acordar", url.Values{"historia": {"HU-1"}, "estimacion": {"5"}}, "primero hay que revelar los votos de la ronda"},
		{"acordar fuera de la escala", true, true, "/planning-poker/acordar", url.Values{"historia": {"HU-1"}, "estimacion": {"4"}}, "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"},
		{"acordar sin votación abierta", false, false, "/planning-poker/acordar", url.Values{"historia": {"HU-1"}, "estimacion": {"5"}}, "la historia no tiene una votación abierta"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			router, _ := nuevoServidorPoker(t)
			switch {
			case c.revelada:
				votacionRevelada(t, router)
			case c.abierta:
				servir(router, pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))
			}
			rec := servir(router, pedidoFormulario(http.MethodPost, c.ruta, c.valores))
			if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), c.mensaje) {
				t.Errorf("código = %d, se esperaba 400 con %q", rec.Code, c.mensaje)
			}
		})
	}
}
