package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// escenarioReportes arma "Demo" con HU-1 (3 SP, 4 h estimadas) y HU-2 (5 SP), el sprint activo S1 (23/10 al 25/10)
// con las dos historias y HU-1 completada, 3 h de Ana en HU-1 el 24/10 y un defecto detectado en S1, y el sprint
// planificado S2.
func escenarioReportes(t *testing.T) http.Handler {
	t.Helper()
	proyectos, _ := memory.NuevoProyectos("")
	integrantes, _ := memory.NuevoIntegrantes("")
	historias, _ := memory.NuevoHistorias("")
	sprints, _ := memory.NuevoSprints("")
	estimaciones, _ := memory.NuevoEstimaciones("")
	registros, _ := memory.NuevoRegistros("")
	defectos, _ := memory.NuevoDefectos("")
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(proyectos),
		Backlog:   service.NuevoBacklog(historias, proyectos, sprintsSvc),
		Sprints:   sprintsSvc,
		Reportes:  service.NuevoReportes(sprints, historias, estimaciones, registros, integrantes, defectos),
	}
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	p, err := svc.Proyectos.Crear("Demo", "", dia(21), dia(31))
	if err != nil {
		t.Fatal(err)
	}
	for i, titulo := range []string{"Crear proyecto", "Crear sprint"} {
		h, err := svc.Backlog.Crear(p.ID, titulo, "", backlog.PrioridadAlta)
		if err != nil {
			t.Fatal(err)
		}
		h.StoryPoints = []int{3, 5}[i]
		h.Criterios = []backlog.Criterio{{Texto: "funciona"}}
		if err := historias.Guardar(h); err != nil {
			t.Fatal(err)
		}
	}
	ana, _ := integrantes.Agregar(project.Integrante{ProyectoID: p.ID, Nombre: "Ana Pérez", Rol: project.ProductBuilder})
	pasos := []func() error{
		func() error { _, err := svc.Sprints.Crear(p.ID, "Primer MVP", dia(23), dia(25)); return err },
		func() error { return svc.Sprints.Iniciar("S1") },
		func() error { return svc.Sprints.Asignar("S1", "HU-1") },
		func() error { return svc.Sprints.Asignar("S1", "HU-2") },
		func() error { return svc.Sprints.CompletarHistoria("S1", "HU-1") },
		func() error { _, err := svc.Sprints.Crear(p.ID, "Segundo", dia(26), dia(28)); return err },
	}
	for _, paso := range pasos {
		if err := paso(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := service.NuevoEstimaciones(estimaciones, historias).Estimar("HU-1", 4); err != nil {
		t.Fatal(err)
	}
	if _, err := registros.Agregar(effort.Registro{IntegranteID: ana.ID, HistoriaID: "HU-1", Fecha: dia(24), Actividad: "Desarrollo", Horas: 3}); err != nil {
		t.Fatal(err)
	}
	if _, err := defectos.Agregar(defect.Defecto{Descripcion: "Falla al guardar", Severidad: defect.Alta, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"}); err != nil {
		t.Fatal(err)
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}
	return router
}

// La sección Reportes (HU-29) lista por proyecto los sprints con reporte: los activos y cerrados (RN-29.1).
func TestReportes_seccion(t *testing.T) {
	router := escenarioReportes(t)
	rec := servir(router, httptest.NewRequest(http.MethodGet, "/reportes", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	contieneTodo(t, cuerpo, "<h2>Demo</h2>", `href="/reportes/sprints/S1"`, "Primer MVP")
	if strings.Contains(cuerpo, `href="/reportes/sprints/S2"`) || strings.Contains(cuerpo, "en construcción") {
		t.Error("un sprint planificado no tiene reporte y la sección ya no está en construcción")
	}
}

func TestReportes_sprint(t *testing.T) {
	router := escenarioReportes(t)
	rec := servir(router, httptest.NewRequest(http.MethodGet, "/reportes/sprints/S1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		"Reporte del sprint", "Primer MVP", "del 23/10/2026 al 25/10/2026", "Reporte parcial",
		"<td>Crear proyecto</td>", "<td>4,0</td>", "<td>3,0</td>", "✔",
		"8 Story Points planificados", "3 Story Points completados", "50,0 % de historias completadas",
		"4,0 horas estimadas", "3,0 horas reales", "desviación de -25,0 %",
		"<td>Ana Pérez</td>", "<td>Desarrollo</td>",
		"<td>Falla al guardar</td>", "<td>Alta</td>", "1 detectado en el sprint · 0 resueltos en el sprint",
		"Imprimir o guardar como PDF", `class="no-imprimir"`)
}

func TestReportes_sprint_errores(t *testing.T) {
	router := escenarioReportes(t)
	casos := []struct {
		ruta    string
		codigo  int
		mensaje string
	}{
		{"/reportes/sprints/S2", http.StatusBadRequest, "el sprint todavía no empezó"},
		{"/reportes/sprints/S9", http.StatusNotFound, "el sprint no existe"},
	}
	for _, c := range casos {
		rec := servir(router, httptest.NewRequest(http.MethodGet, c.ruta, nil))
		if rec.Code != c.codigo || !strings.Contains(rec.Body.String(), c.mensaje) {
			t.Errorf("GET %s: código = %d, se esperaba %d con %q", c.ruta, rec.Code, c.codigo, c.mensaje)
		}
	}
}
