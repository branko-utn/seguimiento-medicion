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

// entornoReportes es el servicio de reportes sobre repositorios en memoria, con "Demo" (P1) y "Otro" (P2) del 21/10 al
// 02/11, y los repositorios para preparar los datos.
type entornoReportes struct {
	sprints      *Sprints
	historias    *memory.Historias
	integrantes  *memory.Integrantes
	estimaciones *memory.Estimaciones
	registros    *memory.Registros
	defectos     *memory.Defectos
	reportes     *Reportes
}

func nuevoEntornoReportes(t *testing.T) entornoReportes {
	t.Helper()
	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	e := entornoReportes{}
	e.historias, _ = memory.NuevoHistorias("")
	e.integrantes, _ = memory.NuevoIntegrantes("")
	e.estimaciones, _ = memory.NuevoEstimaciones("")
	e.registros, _ = memory.NuevoRegistros("")
	e.defectos, _ = memory.NuevoDefectos("")
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := NuevoProyectos(repoProyectos).Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	e.sprints = NuevoSprints(repoSprints, repoProyectos, e.historias)
	e.reportes = NuevoReportes(repoProyectos, repoSprints, e.historias, e.estimaciones, e.registros, e.integrantes, e.defectos)
	return e
}

// TestReportes_Sprint verifica que el servicio junte los datos de los repositorios para el reporte (HU-29): las
// historias del sprint, sus horas estimadas, los registros con el nombre del integrante y solo los defectos del
// proyecto del sprint.
func TestReportes_Sprint(t *testing.T) {
	e := nuevoEntornoReportes(t)
	sprints, repoHistorias, repoIntegrantes, repoEstimaciones, repoRegistros, repoDefectos, reportes :=
		e.sprints, e.historias, e.integrantes, e.estimaciones, e.registros, e.defectos, e.reportes

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
	e := nuevoEntornoReportes(t)
	sprints, reportes := e.sprints, e.reportes
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

// TestReportes_Proyecto verifica que el reporte de proyecto (HU-30) junte solo los datos del proyecto pedido: sus
// historias, sus sprints, sus horas (también fuera de un sprint) y sus defectos.
func TestReportes_Proyecto(t *testing.T) {
	e := nuevoEntornoReportes(t)
	for _, h := range []backlog.Historia{
		{ProyectoID: "P1", Titulo: "Crear proyecto", StoryPoints: 3, Criterios: []backlog.Criterio{{Texto: "ok"}}},
		{ProyectoID: "P1", Titulo: "Crear sprint", StoryPoints: 5, Criterios: []backlog.Criterio{{Texto: "ok"}}},
		{ProyectoID: "P2", Titulo: "Ver métricas", StoryPoints: 2, Criterios: []backlog.Criterio{{Texto: "ok"}}},
	} {
		if _, err := e.historias.Agregar(h); err != nil {
			t.Fatal(err)
		}
	}
	for _, paso := range []func() error{
		func() error { _, err := e.sprints.Crear("P1", "MVP", octubre(21), octubre(23)); return err },
		func() error { return e.sprints.Iniciar("S1") },
		func() error { return e.sprints.Asignar("S1", "HU-1") },
		func() error { return e.sprints.CompletarHistoria("S1", "HU-1") },
		func() error { return e.sprints.CerrarSprint("S1") },
		func() error { _, err := e.sprints.Crear("P2", "Otro MVP", octubre(21), octubre(23)); return err },
		func() error { return e.sprints.Iniciar("S2") },
		func() error { return e.estimaciones.Guardar(effort.EstimacionHoras{HistoriaID: "HU-1", Horas: 4}) },
		func() error { return e.estimaciones.Guardar(effort.EstimacionHoras{HistoriaID: "HU-3", Horas: 9}) },
	} {
		if err := paso(); err != nil {
			t.Fatal(err)
		}
	}
	for _, reg := range []effort.Registro{
		{IntegranteID: "I1", HistoriaID: "HU-1", Fecha: octubre(22), Actividad: "Desarrollo", Horas: 3},
		{IntegranteID: "I1", HistoriaID: "HU-2", Fecha: octubre(27), Actividad: "Desarrollo", Horas: 2}, // fuera de un sprint
		{IntegranteID: "I1", HistoriaID: "HU-3", Fecha: octubre(22), Actividad: "Desarrollo", Horas: 7}, // de Otro
	} {
		if _, err := e.registros.Agregar(reg); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range []defect.Defecto{
		{Descripcion: "No guarda", Severidad: defect.Critica, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"},
		{Descripcion: "De otro proyecto", Severidad: defect.Baja, Estado: defect.Abierto, HistoriaID: "HU-3", SprintDeteccion: "S2"},
	} {
		if _, err := e.defectos.Agregar(d); err != nil {
			t.Fatal(err)
		}
	}

	r, err := e.reportes.Proyecto("P1", octubre(27))
	if err != nil {
		t.Fatal(err)
	}
	if r.Nombre != "Demo" || r.Estado != "En curso" || len(r.Sprints) != 1 || r.Sprints[0].SprintID != "S1" {
		t.Errorf("proyecto %q %q, sprints %+v: solo tiene que aparecer S1", r.Nombre, r.Estado, r.Sprints)
	}
	if r.Historias != 2 || r.Completadas != 1 || r.Velocidad != 3 {
		t.Errorf("historias %d, completadas %d, velocidad %v", r.Historias, r.Completadas, r.Velocidad)
	}
	if r.HorasEstimadas != 4 || r.HorasReales != 5 {
		t.Errorf("horas = %v/%v, se esperaba 4/5 (las de Otro no cuentan; las de fuera de un sprint sí)", r.HorasEstimadas, r.HorasReales)
	}
	if r.Detectados != 1 || len(r.DefectosAbiertos) != 1 || r.DefectosAbiertos[0].Descripcion != "No guarda" {
		t.Errorf("defectos = %d, abiertos %+v", r.Detectados, r.DefectosAbiertos)
	}
}

func TestReportes_Proyecto_inexistente(t *testing.T) {
	e := nuevoEntornoReportes(t)
	if _, err := e.reportes.Proyecto("P9", octubre(27)); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("error = %v, se esperaba project.ErrNoExiste", err)
	}
}
