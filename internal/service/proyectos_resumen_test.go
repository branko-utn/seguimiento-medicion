package service

import (
	"errors"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// entornoResumen tiene el servicio de proyectos configurado para el resumen y los repositorios para preparar casos.
type entornoResumen struct {
	proyectos *Proyectos
	sprints   *Sprints
	historias *memory.Historias
}

// nuevoEntornoResumen arma "Demo" (P1) y "Otro" (P2), los dos del 21/10 al 02/11, sin sprints ni historias.
func nuevoEntornoResumen(t *testing.T) entornoResumen {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoSprints, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	repoHistorias, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	e := entornoResumen{
		proyectos: NuevoProyectos(repoProyectos).ConResumen(repoSprints, repoHistorias),
		sprints:   NuevoSprints(repoSprints, repoProyectos, repoHistorias),
		historias: repoHistorias,
	}
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := e.proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	return e
}

// agregarHistorias carga historias en el proyecto: una por cada estado de la lista, directo en el repositorio.
func (e entornoResumen) agregarHistorias(t *testing.T, proyecto domain.ProyectoID, estados ...backlog.Estado) {
	t.Helper()
	for _, estado := range estados {
		if _, err := e.historias.Agregar(backlog.Historia{ProyectoID: proyecto, Titulo: "Historia", Prioridad: backlog.PrioridadMedia, Estado: estado}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestProyectos_Resumen(t *testing.T) {
	e := nuevoEntornoResumen(t)
	// Demo: 1 de 4 historias completadas, S1 activo y S2 planificado. Otro: datos que no deben contarse.
	e.agregarHistorias(t, "P1", backlog.Completada, backlog.EnProgreso, backlog.Pendiente, backlog.Pendiente)
	e.agregarHistorias(t, "P2", backlog.Completada, backlog.Completada)
	for _, sp := range []struct {
		proyecto    domain.ProyectoID
		goal        string
		inicio, fin int
	}{{"P1", "Primer MVP", 23, 25}, {"P1", "Interfaz", 26, 27}, {"P2", "De otro proyecto", 23, 25}} {
		if _, err := e.sprints.Crear(sp.proyecto, sp.goal, octubre(sp.inicio), octubre(sp.fin)); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []domain.SprintID{"S1", "S3"} {
		if err := e.sprints.Iniciar(id); err != nil {
			t.Fatal(err)
		}
	}

	r, err := e.proyectos.Resumen("P1", octubre(27))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.Proyecto.ID != "P1" || !r.Proyecto.Inicio.Equal(inicio) || !r.Proyecto.Fin.Equal(fin) {
		t.Errorf("proyecto = %+v, se esperaba P1 del 21/10 al 02/11", r.Proyecto)
	}
	if r.Estado != project.EnCurso || r.DiasRestantes != 6 {
		t.Errorf("estado y días = %v y %d, se esperaba En curso y 6", r.Estado, r.DiasRestantes)
	}
	if r.SprintActivo == nil || r.SprintActivo.ID != "S1" || r.SprintActivo.Goal != "Primer MVP" {
		t.Errorf("sprint activo = %+v, se esperaba S1 · Primer MVP", r.SprintActivo)
	}
	if r.PorcentajeCompletadas != 25 {
		t.Errorf("porcentaje = %v, se esperaba 25 (1 de 4, sin contar las de otro proyecto)", r.PorcentajeCompletadas)
	}
}

func TestProyectos_Resumen_estadoSegunLaFecha(t *testing.T) {
	casos := []struct {
		titulo     string
		hoy        time.Time
		quiere     project.Estado
		quiereDias int
	}{
		{"antes del inicio", octubre(20), project.Planificado, 13},
		{"el día de inicio", octubre(21), project.EnCurso, 12},
		{"el día de fin", octubre(33), project.EnCurso, 0},
		{"después del fin", octubre(34), project.Finalizado, 0},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			e := nuevoEntornoResumen(t)

			r, err := e.proyectos.Resumen("P1", caso.hoy)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if r.Estado != caso.quiere || r.DiasRestantes != caso.quiereDias {
				t.Errorf("estado y días = %v y %d, se esperaba %v y %d", r.Estado, r.DiasRestantes, caso.quiere, caso.quiereDias)
			}
		})
	}
}

func TestProyectos_Resumen_sinSprintActivo(t *testing.T) {
	e := nuevoEntornoResumen(t)
	// S1 cerrado y S2 planificado: ninguno está activo.
	for _, f := range [][2]int{{23, 25}, {26, 27}} {
		if _, err := e.sprints.Crear("P1", "Sprint", octubre(f[0]), octubre(f[1])); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := e.sprints.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	r, err := e.proyectos.Resumen("P1", octubre(27))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.SprintActivo != nil {
		t.Errorf("sprint activo = %+v, se esperaba ninguno", r.SprintActivo)
	}
	if sp, _ := e.sprints.Buscar("S1"); sp.Estado != sprint.Cerrado {
		t.Fatalf("la preparación falló: S1 quedó %v", sp.Estado)
	}
}

func TestProyectos_Resumen_sinHistorias(t *testing.T) {
	e := nuevoEntornoResumen(t)

	r, err := e.proyectos.Resumen("P1", octubre(27))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.PorcentajeCompletadas != 0 || r.SprintActivo != nil {
		t.Errorf("resumen = %+v, se esperaba 0 %% y sin sprint activo", r)
	}
}

func TestProyectos_Resumen_errores(t *testing.T) {
	e := nuevoEntornoResumen(t)
	if _, err := e.proyectos.Resumen("P9", octubre(27)); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("proyecto inexistente: se esperaba %v, se obtuvo %v", project.ErrNoExiste, err)
	}

	// Sin ConResumen el servicio no tiene sprints ni historias: devuelve un error en vez de romperse.
	repo, _ := memory.NuevoProyectos("")
	sinConfigurar := NuevoProyectos(repo)
	if _, err := sinConfigurar.Crear("Demo", "", inicio, fin); err != nil {
		t.Fatal(err)
	}
	if _, err := sinConfigurar.Resumen("P1", octubre(27)); err == nil {
		t.Error("sin ConResumen se esperaba un error")
	}
}
