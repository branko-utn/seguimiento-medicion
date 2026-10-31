package server

import (
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

func TestDashboard_muestraSelectorYSinDatos(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoSesiones, _ := memory.NuevoSesiones("")

	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	svc := Servicios{
		Proyectos:    service.NuevoProyectos(repoProyectos),
		Integrantes:  service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		Backlog:      service.NuevoBacklog(repoHistorias, repoProyectos, sprintsSvc, nil, nil),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(repoEstimaciones, repoHistorias),
		Registros:    service.NuevoRegistros(repoRegistros, repoIntegrantes, repoHistorias, repoProyectos),
		Defectos:     service.NuevoDefectos(repoDefectos, repoHistorias, repoSprints),
		Poker:        service.NuevoPlanningPoker(repoSesiones, repoHistorias, repoIntegrantes),
		Dashboard:    service.NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos),
	}
	r, _ := NewRouter(svc)

	p, _ := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), Fin: time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local)})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/dashboard", nil))
	if rec.Code != 200 {
		t.Fatalf("codigo=%d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Demo") {
		t.Fatalf("debe mostrar selector con proyecto")
	}

	// con proyecto sin sprints cerrados
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if _, err := sprintsSvc.Crear(p.ID, "S1", time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), time.Date(2026, 10, 23, 0, 0, 0, 0, time.Local)); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S1", h1.ID); err != nil {
		t.Fatal(err)
	}

	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, httptest.NewRequest("GET", "/dashboard?proyecto="+string(p.ID), nil))
	if rec2.Code != 200 {
		t.Fatalf("codigo=%d", rec2.Code)
	}
	body := rec2.Body.String()
	if !strings.Contains(body, "Todavía no hay sprints cerrados para graficar") {
		t.Fatalf("falta mensaje")
	}
	if !strings.Contains(body, "series-dashboard") {
		t.Fatalf("falta script de series")
	}
	if !strings.Contains(body, "grafico-velocidad") || !strings.Contains(body, "grafico-sp") || !strings.Contains(body, "grafico-horas") || !strings.Contains(body, "grafico-defectos") {
		t.Fatalf("faltan canvas")
	}
}

func TestDashboard_conSprintsCerradosMuestraSeries(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoSesiones, _ := memory.NuevoSesiones("")

	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	svc := Servicios{
		Proyectos:    service.NuevoProyectos(repoProyectos),
		Integrantes:  service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		Backlog:      service.NuevoBacklog(repoHistorias, repoProyectos, sprintsSvc, nil, nil),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(repoEstimaciones, repoHistorias),
		Registros:    service.NuevoRegistros(repoRegistros, repoIntegrantes, repoHistorias, repoProyectos),
		Defectos:     service.NuevoDefectos(repoDefectos, repoHistorias, repoSprints),
		Poker:        service.NuevoPlanningPoker(repoSesiones, repoHistorias, repoIntegrantes),
		Dashboard:    service.NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos),
	}
	r, _ := NewRouter(svc)

	p, _ := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), Fin: time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local)})
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	h2, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear sprint", StoryPoints: 5, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if err := repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: h1.ID, Horas: 4}); err != nil {
		t.Fatal(err)
	}
	if err := repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: h2.ID, Horas: 6}); err != nil {
		t.Fatal(err)
	}
	if _, err := sprintsSvc.Crear(p.ID, "S1", time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), time.Date(2026, 10, 23, 0, 0, 0, 0, time.Local)); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S1", h1.ID); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S1", h2.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repoRegistros.Agregar(effort.Registro{HistoriaID: h1.ID, Fecha: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), Horas: 5}); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.CompletarHistoria("S1", h1.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repoDefectos.Agregar(defect.Defecto{Descripcion: "Falla", Severidad: defect.Alta, Estado: defect.Abierto, HistoriaID: h1.ID, SprintDeteccion: "S1"}); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/dashboard?proyecto="+string(p.ID), nil))
	if rec.Code != 200 {
		t.Fatalf("codigo=%d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "\"puntos\"") || !strings.Contains(body, "\"mensajes\"") {
		t.Fatalf("falta JSON de series")
	}
	// La línea de referencia del gráfico de velocidad (spec HU-28) viaja con las series.
	if !strings.Contains(body, "\"velocidad\":") {
		t.Error("las series no traen la velocidad promedio para la línea de referencia")
	}
	if !strings.Contains(body, "S1") {
		t.Fatalf("falta sprint en dashboard")
	}
}

func TestDashboard_proyectoInexistente404(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoSesiones, _ := memory.NuevoSesiones("")

	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	svc := Servicios{
		Proyectos:    service.NuevoProyectos(repoProyectos),
		Integrantes:  service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		Backlog:      service.NuevoBacklog(repoHistorias, repoProyectos, sprintsSvc, nil, nil),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(repoEstimaciones, repoHistorias),
		Registros:    service.NuevoRegistros(repoRegistros, repoIntegrantes, repoHistorias, repoProyectos),
		Defectos:     service.NuevoDefectos(repoDefectos, repoHistorias, repoSprints),
		Poker:        service.NuevoPlanningPoker(repoSesiones, repoHistorias, repoIntegrantes),
		Dashboard:    service.NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos),
	}
	r, _ := NewRouter(svc)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/dashboard?proyecto=Otro", nil))
	if rec.Code != 404 {
		t.Fatalf("codigo=%d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "el proyecto no existe") {
		t.Fatalf("falta mensaje de error")
	}
}
