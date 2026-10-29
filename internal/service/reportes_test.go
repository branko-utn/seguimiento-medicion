package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// TestReportes_Sprint verifica que el servicio junte los datos de los repositorios para el reporte (HU-29): las
// historias del sprint, sus horas estimadas, los registros con el nombre del integrante y solo los defectos del
// proyecto del sprint.
func TestReportes_Sprint(t *testing.T) {
	sprints, repoSprints, repoHistorias := armarServicioSprints(t)
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	reportes := NuevoReportes(repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoIntegrantes, repoDefectos)

	// HU-1 y HU-2 en Demo (P1), HU-3 en Otro (P2).
	for _, h := range []backlog.Historia{
		{ProyectoID: "P1", Titulo: "Crear proyecto", StoryPoints: 3, Criterios: []backlog.Criterio{{Texto: "ok"}}},
		{ProyectoID: "P1", Titulo: "Crear sprint", StoryPoints: 5, Criterios: []backlog.Criterio{{Texto: "ok"}}},
		{ProyectoID: "P2", Titulo: "Ver métricas", StoryPoints: 2, Criterios: []backlog.Criterio{{Texto: "ok"}}},
	} {
		if _, err := repoHistorias.Agregar(h); err != nil {
			t.Fatal(err)
		}
	}
	ana, _ := repoIntegrantes.Agregar(project.Integrante{ProyectoID: "P1", Nombre: "Ana Pérez", Rol: project.ProductBuilder})
	if _, err := sprints.Crear("P1", "MVP", octubre(21), octubre(23)); err != nil {
		t.Fatal(err)
	}
	for _, paso := range []func() error{
		func() error { return sprints.Iniciar("S1") },
		func() error { return sprints.Asignar("S1", "HU-1") },
		func() error { return sprints.Asignar("S1", "HU-2") },
		func() error { return sprints.CompletarHistoria("S1", "HU-1") },
		func() error { return repoEstimaciones.Guardar(effort.EstimacionHoras{HistoriaID: "HU-1", Horas: 4}) },
	} {
		if err := paso(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repoRegistros.Agregar(effort.Registro{IntegranteID: ana.ID, HistoriaID: "HU-1", Fecha: octubre(22), Actividad: "Desarrollo", Horas: 3}); err != nil {
		t.Fatal(err)
	}
	for _, d := range []defect.Defecto{
		{Descripcion: "Falla al guardar", Severidad: defect.Alta, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"},
		{Descripcion: "De otro proyecto", Severidad: defect.Baja, Estado: defect.Abierto, HistoriaID: "HU-3", SprintDeteccion: "S1"},
	} {
		if _, err := repoDefectos.Agregar(d); err != nil {
			t.Fatal(err)
		}
	}

	r, err := reportes.Sprint("S1")
	if err != nil {
		t.Fatal(err)
	}
	if !r.Parcial || len(r.Historias) != 2 || r.Historias[0].Titulo != "Crear proyecto" || r.Historias[0].HorasEstimadas != 4 || !r.Historias[0].Completada {
		t.Errorf("historias = %+v (parcial %v)", r.Historias, r.Parcial)
	}
	if r.SPPlanificados != 8 || r.SPCompletados != 3 || r.HorasReales != 3 {
		t.Errorf("SP %d/%d, horas reales %v", r.SPPlanificados, r.SPCompletados, r.HorasReales)
	}
	if len(r.PorIntegrante) != 1 || r.PorIntegrante[0].Integrante != "Ana Pérez" {
		t.Errorf("por integrante = %+v", r.PorIntegrante)
	}
	if len(r.Defectos) != 1 || r.Detectados != 1 {
		t.Errorf("defectos = %+v (detectados %d): el de otro proyecto no cuenta", r.Defectos, r.Detectados)
	}
}

func TestReportes_Sprint_errores(t *testing.T) {
	sprints, repoSprints, repoHistorias := armarServicioSprints(t)
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	reportes := NuevoReportes(repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoIntegrantes, repoDefectos)
	if _, err := sprints.Crear("P1", "MVP", octubre(21), octubre(23)); err != nil {
		t.Fatal(err)
	}

	casos := []struct {
		sprint domain.SprintID
		quiere error
	}{
		{"S9", sprint.ErrNoExiste},
		{"S1", report.ErrSprintSinEmpezar},
	}
	for _, c := range casos {
		if _, err := reportes.Sprint(c.sprint); !errors.Is(err, c.quiere) {
			t.Errorf("Sprint(%s): error = %v, se esperaba %v", c.sprint, err, c.quiere)
		}
	}
}
