package server

import (
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// resumenDe devuelve el resumen de un proyecto en la pantalla de esfuerzo: desde su id hasta el cierre del artículo.
func resumenDe(t *testing.T, cuerpo, proyecto string) string {
	t.Helper()
	inicio := strings.Index(cuerpo, `id="esfuerzo-`+proyecto+`"`)
	if inicio < 0 {
		t.Fatalf("la página no tiene el resumen de esfuerzo de %s", proyecto)
	}
	fin := strings.Index(cuerpo[inicio:], "</article>")
	if fin < 0 {
		t.Fatalf("el resumen de %s no se cierra", proyecto)
	}
	return cuerpo[inicio : inicio+fin]
}

func TestEsfuerzo_resumenSinRegistros(t *testing.T) {
	router, _ := escenarioEsfuerzo(t)

	cuerpo := ver(router, "/esfuerzo").Body.String()

	// Un resumen por proyecto, debajo de la lista de registros; sin registros muestra el mensaje de la spec.
	contiene(t, cuerpo, "<h2>Esfuerzo por proyecto</h2>")
	if strings.Index(cuerpo, "<h2>Esfuerzo por proyecto</h2>") < strings.Index(cuerpo, "<h2>Horas registradas</h2>") {
		t.Error("el resumen por proyecto tiene que ir debajo de la lista de registros")
	}
	for _, p := range []string{"P1", "P2"} {
		contiene(t, resumenDe(t, cuerpo, p), "Todavía no hay esfuerzo registrado.")
	}
}

func TestEsfuerzo_resumenConTotales(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	if _, err := svc.Backlog.Crear("P1", "Crear sprint", "", backlog.PrioridadMedia); err != nil {
		t.Fatal(err)
	}
	for _, r := range []struct {
		historia string
		dia      int
		horas    float64
	}{{"HU-1", 22, 3}, {"HU-1", 23, 2.5}, {"HU-2", 24, 4}} {
		if _, err := svc.Registros.Registrar("I1", domain.HistoriaID(r.historia), dia(r.dia), "Desarrollo", r.horas, dia(24)); err != nil {
			t.Fatal(err)
		}
	}

	cuerpo := ver(router, "/esfuerzo").Body.String()

	// Demo: totales por historia y por integrante, con un decimal, y el total del proyecto.
	demo := resumenDe(t, cuerpo, "P1")
	contiene(t, demo, "Demo", "9,5 horas",
		`<a href="/backlog/HU-1">HU-1</a>`, "<td>5,5</td>",
		`<a href="/backlog/HU-2">HU-2</a>`, "<td>4,0</td>",
		"<td>Ana Pérez</td>", "<td>9,5</td>")
	if strings.Contains(demo, "Todavía no hay esfuerzo registrado") {
		t.Error("Demo tiene esfuerzo registrado: no debería mostrar el mensaje de vacío")
	}
	// Otro no tiene registros.
	contiene(t, resumenDe(t, cuerpo, "P2"), "Otro", "Todavía no hay esfuerzo registrado.")
}
