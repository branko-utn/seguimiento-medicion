package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// La pantalla del sprint muestra sus métricas (TEC-17): SP planificados y completados (HU-24), porcentaje de
// historias completadas (HU-25) y la velocidad del equipo, que es la del proyecto (promedio de los sprints cerrados).
func TestSprints_metricas(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)
	// HU-3 (5 SP) necesita un criterio para poder asignarse (HU-11).
	h3, err := e.repoHistorias.Buscar("HU-3")
	if err != nil {
		t.Fatal(err)
	}
	h3.Criterios = []backlog.Criterio{{Texto: "Funciona"}}
	if err := e.repoHistorias.Guardar(h3); err != nil {
		t.Fatal(err)
	}
	ver := func(t *testing.T, id string) string {
		t.Helper()
		rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/"+id, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /sprints/%s: código = %d, se esperaba 200", id, rec.Code)
		}
		return rec.Body.String()
	}
	accion := func(t *testing.T, ruta string) {
		t.Helper()
		if rec := servir(e.router, pedidoFormulario(http.MethodPost, ruta, url.Values{})); rec.Code != http.StatusSeeOther {
			t.Fatalf("POST %s: código = %d, se esperaba 303", ruta, rec.Code)
		}
	}

	t.Run("sprint sin historias", func(t *testing.T) {
		contieneTodo(t, ver(t, "S1"), "Métricas del sprint", tarjeta("0", "Story Points planificados"), tarjeta("0", "Story Points completados"),
			tarjeta("0,0 %", "historias completadas (0 de 0)"), tarjeta("0,0", "velocidad del equipo en SP por sprint (sin sprints cerrados todavía)"))
	})

	for _, id := range []string{"HU-1", "HU-3"} {
		if err := e.svc.Sprints.Asignar("S1", domain.HistoriaID(id)); err != nil {
			t.Fatal(err)
		}
	}
	accion(t, "/sprints/S1/iniciar")
	accion(t, "/sprints/S1/historias/HU-1/completar")
	t.Run("activo con una de dos historias completadas", func(t *testing.T) {
		contieneTodo(t, ver(t, "S1"), tarjeta("8", "Story Points planificados"), tarjeta("3", "Story Points completados"),
			tarjeta("50,0 %", "historias completadas (1 de 2)"), "(sin sprints cerrados todavía)")
	})

	accion(t, "/sprints/S1/cerrar")
	t.Run("cerrado: cuenta para la velocidad", func(t *testing.T) {
		cuerpo := ver(t, "S1")
		contieneTodo(t, cuerpo, tarjeta("8", "Story Points planificados"), tarjeta("3", "Story Points completados"),
			tarjeta("50,0 %", "historias completadas (1 de 2)"), tarjeta("3,0", "velocidad del equipo en SP por sprint (1 sprint cerrado)"))
	})

	// Segundo sprint con HU-3, que volvió al backlog al cerrar S1: la velocidad pasa a (3 + 5) / 2.
	if _, err := e.svc.Sprints.Crear("P1", "Cierre", fechaDeFormulario("2026-10-26"), fechaDeFormulario("2026-10-28")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Asignar("S2", "HU-3"); err != nil {
		t.Fatal(err)
	}
	accion(t, "/sprints/S2/iniciar")
	accion(t, "/sprints/S2/historias/HU-3/completar")
	t.Run("activo: todavía no cuenta para la velocidad", func(t *testing.T) {
		contieneTodo(t, ver(t, "S2"), tarjeta("100,0 %", "historias completadas (1 de 1)"), tarjeta("3,0", "velocidad del equipo en SP por sprint (1 sprint cerrado)"))
	})
	accion(t, "/sprints/S2/cerrar")
	t.Run("dos sprints cerrados", func(t *testing.T) {
		contieneTodo(t, ver(t, "S2"), tarjeta("4,0", "velocidad del equipo en SP por sprint (2 sprints cerrados)"))
	})

	// La velocidad es por proyecto: un sprint de Otro no ve los sprints cerrados de Demo.
	if _, err := e.svc.Sprints.Crear("P2", "Métricas", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	t.Run("otro proyecto", func(t *testing.T) {
		cuerpo := ver(t, "S3")
		if !strings.Contains(cuerpo, "(sin sprints cerrados todavía)") {
			t.Error("la velocidad de Otro no debe contar los sprints cerrados de Demo")
		}
	})
}
