package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

func TestDashboard_Del_resumenYMetricasPrincipales(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")

	sprintsSvc := NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	dash := NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos)

	p, err := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: octubre(21), Fin: octubre(31)})
	if err != nil {
		t.Fatal(err)
	}
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	h2, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear sprint", StoryPoints: 5, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	
	if _, err := sprintsSvc.Crear(p.ID, "S1", octubre(21), octubre(23)); err != nil {
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
	if err := sprintsSvc.CompletarHistoria("S1", h1.ID); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	if _, err := sprintsSvc.Crear(p.ID, "S2", octubre(24), octubre(27)); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Iniciar("S2"); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S2", h2.ID); err != nil {
		t.Fatal(err)
	}

	dashRes, err := dash.Del(p.ID, octubre(27))
	if err != nil {
		t.Fatal(err)
	}
	if dashRes.Estado != "En curso" || dashRes.DiasRestantes != 6 {
		t.Errorf("estado=%s dias=%d", dashRes.Estado, dashRes.DiasRestantes)
	}
	if string(dashRes.SprintActivo) != "S2" {
		t.Errorf("sprint activo=%s", dashRes.SprintActivo)
	}
	if dashRes.Velocidad != 3.0 {
		t.Errorf("velocidad=%v", dashRes.Velocidad)
	}
	if dashRes.PorcentajeCompletadas != 33.3 {
		t.Errorf("porcentaje=%v", dashRes.PorcentajeCompletadas)
	}
}

func TestDashboard_Del_seriesPorSprintCerrado(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")

	sprintsSvc := NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	dash := NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos)

	p, _ := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: octubre(21), Fin: octubre(31)})
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	h2, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear sprint", StoryPoints: 5, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if err := repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: h1.ID, Horas: 4}); err != nil {
		t.Fatal(err)
	}
	if err := repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: h2.ID, Horas: 6}); err != nil {
		t.Fatal(err)
	}

	if _, err := sprintsSvc.Crear(p.ID, "S1", octubre(21), octubre(23)); err != nil {
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
	if _, err := repoRegistros.Agregar(effort.Registro{HistoriaID: h1.ID, Fecha: octubre(21), Horas: 5}); err != nil {
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

	res, err := dash.Del(p.ID, octubre(27))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Puntos) != 1 {
		t.Fatalf("puntos=%d", len(res.Puntos))
	}
	pt := res.Puntos[0]
	if string(pt.Sprint) != "S1" || pt.SPPlanificados != 8 || pt.SPCompletados != 3 || pt.HorasEstimadas != 10.0 || pt.HorasReales != 5.0 || pt.DefectosDetectados != 1 || pt.DefectosResueltos != 0 {
		t.Errorf("punto=%+v", pt)
	}
}

func TestDashboard_Del_sinSprintsCerrados(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")

	sprintsSvc := NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	dash := NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos)

	p, _ := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: octubre(21), Fin: octubre(31)})
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if _, err := sprintsSvc.Crear(p.ID, "S1", octubre(21), octubre(23)); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S1", h1.ID); err != nil {
		t.Fatal(err)
	}

	res, err := dash.Del(p.ID, octubre(27))
	if err != nil {
		t.Fatal(err)
	}
	if res.Mensajes.Velocidad == "" || res.Mensajes.Defectos == "" {
		t.Errorf("mensajes=%+v", res.Mensajes)
	}
	if res.DesviacionCalculable {
		t.Errorf("se esperaba desviación no calculable")
	}
}

func TestDashboard_Del_sprintsCerradosSinDefectosPeroConHoras(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")

	sprintsSvc := NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	dash := NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos)

	p, _ := repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: octubre(21), Fin: octubre(31)})
	h1, _ := repoHistorias.Agregar(backlog.Historia{ProyectoID: p.ID, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if err := repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: h1.ID, Horas: 4}); err != nil {
		t.Fatal(err)
	}
	if _, err := sprintsSvc.Crear(p.ID, "S1", octubre(21), octubre(23)); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.Asignar("S1", h1.ID); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.CompletarHistoria("S1", h1.ID); err != nil {
		t.Fatal(err)
	}
	if err := sprintsSvc.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	res, err := dash.Del(p.ID, octubre(27))
	if err != nil {
		t.Fatal(err)
	}
	if res.Mensajes.Defectos == "" {
		t.Errorf("se esperaba mensaje de defectos")
	}
	if res.Mensajes.Horas != "" {
		t.Errorf("se esperaba gráfico de horas dibujado")
	}
}

func TestDashboard_Del_proyectoInexistente(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")

	dash := NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos)

	_, err := dash.Del("Otro", octubre(27))
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if !errors.Is(err, project.ErrNoExiste) {
		t.Fatalf("error=%v", err)
	}
}
