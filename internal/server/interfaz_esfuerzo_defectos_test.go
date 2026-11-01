package server

import (
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// Tests de HU-32.3: la interfaz de esfuerzo y defectos con los componentes de TEC-21.

// RN-32.1: la severidad y el estado de cada defecto se ven con insignias; el tono lo decide el handler.
func TestDefectos_insignias(t *testing.T) {
	// D1 es Alta y está Abierto; D2 es Crítica y queda Resuelto.
	router, svc := escenarioDefectos(t, true)
	if _, err := svc.Defectos.Registrar("No guarda", defect.Critica, "HU-1", "S1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Defectos.Resolver("D2", "S2"); err != nil {
		t.Fatal(err)
	}

	cuerpo := ver(router, "/defectos").Body.String()

	contiene(t, cuerpo,
		`<span class="insignia tono-alerta">Alta</span>`, `<span class="insignia tono-peligro">Abierto</span>`,
		`<span class="insignia tono-peligro">Crítica</span>`, `<span class="insignia tono-exito">Resuelto</span>`)
}

// RN-32.4: la lista vacía usa .vacio con un enlace al formulario que la llena.
func TestDefectos_vacioConAccion(t *testing.T) {
	router, _ := escenarioDefectos(t, false)

	cuerpo := ver(router, "/defectos").Body.String()

	contiene(t, cuerpo, `class="vacio"`, "Todavía no hay defectos registrados.", `href="#registrar-defecto"`,
		`id="registrar-defecto"`)
}

// RN-32.2 y RN-32.3: los totales de esfuerzo de cada proyecto en tarjetas y las horas con coma decimal.
func TestEsfuerzo_tarjetasYComaDecimal(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)
	if _, err := svc.Backlog.Crear("P1", "Crear sprint", "", backlog.PrioridadMedia); err != nil {
		t.Fatal(err)
	}
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	for _, r := range []struct {
		historia domain.HistoriaID
		horas    float64
	}{{"HU-1", 3}, {"HU-1", 2.5}, {"HU-2", 4}} {
		if _, err := svc.Registros.Registrar("I1", r.historia, dia(22), "Desarrollo", r.horas, dia(24)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.Estimaciones.Estimar("HU-1", 4.5); err != nil {
		t.Fatal(err)
	}

	cuerpo := ver(router, "/esfuerzo").Body.String()

	demo := resumenDe(t, cuerpo, "P1")
	contiene(t, demo, `class="tarjetas"`,
		tarjeta("9,5", "horas registradas"), tarjeta("2", "historias con esfuerzo"), tarjeta("1", "integrantes con esfuerzo"),
		`<td class="num">5,5</td>`, `<td class="num">4,0</td>`, `<td class="num">9,5</td>`)
	// Las horas de cada registro y las estimadas también con un decimal y coma.
	contiene(t, cuerpo, `<td class="num">3,0</td>`, `<td class="num">2,5</td>`, `<td class="num">4,5</td>`)
	// Las tablas usan el estilo de TEC-21.
	if strings.Contains(cuerpo, "<table>") {
		t.Error("todas las tablas tienen que usar la clase tabla")
	}
}

// RN-32.4: cada lista vacía de la pantalla de esfuerzo usa .vacio con su acción.
func TestEsfuerzo_vaciosConAccion(t *testing.T) {
	router, _ := escenarioEsfuerzo(t)

	cuerpo := ver(router, "/esfuerzo").Body.String()

	contiene(t, cuerpo, `class="vacio"`, "Todavía no hay horas registradas.", `href="#registrar-horas"`,
		`id="registrar-horas"`)
	contiene(t, resumenDe(t, cuerpo, "P2"), `class="vacio"`, "Todavía no hay esfuerzo registrado.",
		`href="#registrar-horas"`)
}

// RN-32.4: sin historias, la tabla de horas estimadas deja lugar a .vacio con el enlace para crear una.
func TestEsfuerzo_sinHistoriasConAccion(t *testing.T) {
	proyectos, _ := memory.NuevoProyectos("")
	integrantes, _ := memory.NuevoIntegrantes("")
	historias, _ := memory.NuevoHistorias("")
	sprints, _ := memory.NuevoSprints("")
	estimaciones, _ := memory.NuevoEstimaciones("")
	registros, _ := memory.NuevoRegistros("")
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	router, err := NewRouter(Servicios{
		Proyectos:    service.NuevoProyectos(proyectos),
		Integrantes:  service.NuevoIntegrantes(integrantes, proyectos),
		Backlog:      service.NuevoBacklog(historias, proyectos, sprintsSvc, nil, nil),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(estimaciones, historias),
		Registros:    service.NuevoRegistros(registros, integrantes, historias, proyectos),
	})
	if err != nil {
		t.Fatal(err)
	}

	cuerpo := ver(router, "/esfuerzo").Body.String()

	contiene(t, cuerpo, "Todavía no hay historias en el backlog.", `href="/backlog/nueva"`)
	if strings.Count(cuerpo, `class="vacio"`) < 2 {
		t.Error("sin historias, la sección de horas estimadas tiene que usar .vacio")
	}
}
