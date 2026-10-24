package service

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// octubre devuelve un día de octubre de 2026; 32 = 1 de noviembre, 33 = 2 de noviembre (time.Date normaliza).
func octubre(d int) time.Time {
	return time.Date(2026, time.October, d, 0, 0, 0, 0, time.Local)
}

// nuevoServicioSprints arma el servicio con dos proyectos del 21/10 al 02/11: "Demo" (P1) y "Otro" (P2).
// Devuelve también el repositorio de sprints para preparar estados que todavía no tienen caso de uso (Cerrado).
func nuevoServicioSprints(t *testing.T) (*Sprints, *memory.Sprints) {
	t.Helper()
	s, repoSprints, _ := armarServicioSprints(t)
	return s, repoSprints
}

// armarServicioSprints es nuevoServicioSprints, pero devuelve también el repositorio de historias (HU-11).
func armarServicioSprints(t *testing.T) (*Sprints, *memory.Sprints, *memory.Historias) {
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
	proyectos := NuevoProyectos(repoProyectos)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	return NuevoSprints(repoSprints, repoProyectos, repoHistorias), repoSprints, repoHistorias
}

func TestSprints_Crear(t *testing.T) {
	s, _ := nuevoServicioSprints(t)

	sp, err := s.Crear("P1", "  Primer MVP  ", octubre(23), octubre(25))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if sp.ID != "S1" || sp.ProyectoID != "P1" || sp.Goal != "Primer MVP" || sp.Estado != sprint.Planificado {
		t.Errorf("Crear() = %+v, se esperaba S1 · P1 · Primer MVP · Planificado", sp)
	}

	lista, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 1 || lista[0].ID != "S1" {
		t.Errorf("DelProyecto(P1) = %+v, se esperaba solo el sprint S1", lista)
	}
}

func TestSprints_Crear_validos(t *testing.T) {
	casos := []struct {
		titulo      string
		proyecto    domain.ProyectoID
		inicio, fin time.Time
	}{
		{"consecutivo al anterior", "P1", octubre(26), octubre(27)},
		{"anterior al existente", "P1", octubre(21), octubre(22)},
		{"mismas fechas en otro proyecto", "P2", octubre(23), octubre(25)},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, _ := nuevoServicioSprints(t)
			if _, err := s.Crear("P1", "Primero", octubre(23), octubre(25)); err != nil {
				t.Fatalf("no se pudo crear el sprint inicial: %v", err)
			}

			if _, err := s.Crear(caso.proyecto, "Segundo", caso.inicio, caso.fin); err != nil {
				t.Errorf("error inesperado: %v", err)
			}
		})
	}
}

func TestSprints_Crear_ocupaTodoElProyecto(t *testing.T) {
	s, _ := nuevoServicioSprints(t)

	if _, err := s.Crear("P1", "Todo", inicio, fin); err != nil {
		t.Errorf("un sprint del inicio al fin del proyecto es válido, pero se obtuvo %v", err)
	}
}

func TestSprints_Crear_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		proyecto    domain.ProyectoID
		goal        string
		inicio, fin time.Time
		quiereError error
	}{
		{"proyecto inexistente", "P9", "Interfaz", octubre(26), octubre(27), project.ErrNoExiste},
		{"primero se verifica el proyecto", "P9", " ", octubre(27), octubre(26), project.ErrNoExiste},
		{"Sprint Goal vacío", "P1", "  ", octubre(26), octubre(27), sprint.ErrGoalVacio},
		{"fin anterior al inicio", "P1", "Interfaz", octubre(27), octubre(26), sprint.ErrFechasInvalidas},
		{"termina después del proyecto", "P1", "Cierre", octubre(30), octubre(34), sprint.ErrFueraDelProyecto},
		{"empieza antes del proyecto", "P1", "Previo", octubre(20), octubre(22), sprint.ErrFueraDelProyecto},
		{"se superpone con otro sprint", "P1", "Interfaz", octubre(25), octubre(27), sprint.ErrSuperpuesto},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, _ := nuevoServicioSprints(t)
			if _, err := s.Crear("P1", "Primero", octubre(23), octubre(25)); err != nil {
				t.Fatalf("no se pudo crear el sprint inicial: %v", err)
			}

			_, err := s.Crear(caso.proyecto, caso.goal, caso.inicio, caso.fin)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}

			lista, _ := s.DelProyecto("P1")
			if len(lista) != 1 {
				t.Errorf("hay %d sprints en P1, se esperaba 1: no se debe guardar uno inválido", len(lista))
			}
		})
	}
}

func TestSprints_Iniciar(t *testing.T) {
	s, _ := nuevoServicioSprints(t)
	if _, err := s.Crear("P1", "Primero", octubre(23), octubre(25)); err != nil {
		t.Fatal(err)
	}

	if err := s.Iniciar("S1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	sp, err := s.Buscar("S1")
	if err != nil {
		t.Fatal(err)
	}
	if sp.Estado != sprint.Activo {
		t.Errorf("estado = %v, se esperaba %v", sp.Estado, sprint.Activo)
	}
}

func TestSprints_Iniciar_activoEnOtroProyecto(t *testing.T) {
	s, _ := nuevoServicioSprints(t)
	for _, proyecto := range []domain.ProyectoID{"P1", "P2"} {
		if _, err := s.Crear(proyecto, "Primero", octubre(23), octubre(25)); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}

	if err := s.Iniciar("S2"); err != nil {
		t.Errorf("un sprint activo en otro proyecto no impide iniciar, pero se obtuvo %v", err)
	}
}

func TestSprints_Iniciar_errores(t *testing.T) {
	// En P1: S1 (23–25) activo, S2 (26–27) planificado y S3 (28–29) cerrado.
	casos := []struct {
		titulo      string
		sprint      domain.SprintID
		quiereError error
		quiere      sprint.Estado
	}{
		{"sprint inexistente", "S9", sprint.ErrNoExiste, 0},
		{"ya hay otro sprint activo", "S2", sprint.ErrOtroActivo, sprint.Planificado},
		{"el sprint ya está activo", "S1", sprint.ErrNoPlanificado, sprint.Activo},
		{"primero se verifica el estado del sprint", "S3", sprint.ErrNoPlanificado, sprint.Cerrado},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, repo := nuevoServicioSprints(t)
			for _, f := range [][2]int{{23, 25}, {26, 27}, {28, 29}} {
				if _, err := s.Crear("P1", "Sprint", octubre(f[0]), octubre(f[1])); err != nil {
					t.Fatal(err)
				}
			}
			if err := s.Iniciar("S1"); err != nil {
				t.Fatal(err)
			}
			s3, _ := repo.Buscar("S3")
			s3.Estado = sprint.Cerrado
			if err := repo.Actualizar(s3); err != nil {
				t.Fatal(err)
			}

			err := s.Iniciar(caso.sprint)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if caso.quiere != 0 {
				if sp, _ := s.Buscar(caso.sprint); sp.Estado != caso.quiere {
					t.Errorf("estado = %v, se esperaba %v: un error no debe cambiar el sprint", sp.Estado, caso.quiere)
				}
			}
		})
	}
}

// prepararAsignacion arma el escenario de HU-11. En P1: S1 (23–25) planificado, S2 (26–27) planificado con HU-6 y
// S3 (28–29) cerrado. Historias: HU-1 lista, HU-2 sin estimar, HU-3 sin criterios, HU-4 completada en S2, HU-5 de
// P2 y HU-6 lista pero ya asignada a S2. Los estados se cargan directo en los repositorios.
func prepararAsignacion(t *testing.T) (*Sprints, *memory.Sprints, *memory.Historias) {
	t.Helper()
	s, repoSprints, repoHistorias := armarServicioSprints(t)
	for _, f := range [][2]int{{23, 25}, {26, 27}, {28, 29}} {
		if _, err := s.Crear("P1", "Sprint", octubre(f[0]), octubre(f[1])); err != nil {
			t.Fatal(err)
		}
	}
	s2, _ := repoSprints.Buscar("S2")
	s2.Historias = []domain.HistoriaID{"HU-6"}
	s3, _ := repoSprints.Buscar("S3")
	s3.Estado = sprint.Cerrado
	for _, sp := range []sprint.Sprint{s2, s3} {
		if err := repoSprints.Actualizar(sp); err != nil {
			t.Fatal(err)
		}
	}

	criterio := []backlog.Criterio{{Texto: "funciona"}}
	historias := []backlog.Historia{
		{ProyectoID: "P1", Titulo: "Lista", Estado: backlog.Pendiente, StoryPoints: 3, Criterios: criterio},
		{ProyectoID: "P1", Titulo: "Sin estimar", Estado: backlog.Pendiente, Criterios: criterio},
		{ProyectoID: "P1", Titulo: "Sin criterios", Estado: backlog.Pendiente, StoryPoints: 5},
		{ProyectoID: "P1", Titulo: "Completada", Estado: backlog.Completada, StoryPoints: 2, Criterios: criterio, SprintID: "S2"},
		{ProyectoID: "P2", Titulo: "De otro proyecto", Estado: backlog.Pendiente, StoryPoints: 2, Criterios: criterio},
		{ProyectoID: "P1", Titulo: "Ya asignada", Estado: backlog.Pendiente, StoryPoints: 1, Criterios: criterio, SprintID: "S2"},
	}
	for _, h := range historias {
		if _, err := repoHistorias.Agregar(h); err != nil {
			t.Fatal(err)
		}
	}
	return s, repoSprints, repoHistorias
}

func TestSprints_Asignar(t *testing.T) {
	s, repoSprints, repoHistorias := prepararAsignacion(t)

	if err := s.Asignar("S1", "HU-1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if sp, _ := repoSprints.Buscar("S1"); !slices.Equal(sp.Historias, []domain.HistoriaID{"HU-1"}) {
		t.Errorf("historias de S1 = %v, se esperaba [HU-1]", sp.Historias)
	}
	if h, _ := repoHistorias.Buscar("HU-1"); h.SprintID != "S1" {
		t.Errorf("sprint de HU-1 = %q, se esperaba S1", h.SprintID)
	}
}

func TestSprints_Asignar_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		sprint      domain.SprintID
		historia    domain.HistoriaID
		quiereError error
	}{
		{"sprint inexistente", "S9", "HU-1", sprint.ErrNoExiste},
		{"historia inexistente", "S1", "HU-9", backlog.ErrNoExiste},
		{"primero se busca el sprint", "S9", "HU-9", sprint.ErrNoExiste},
		{"historia sin estimar", "S1", "HU-2", sprint.ErrHistoriaSinEstimar},
		{"historia sin criterios", "S1", "HU-3", sprint.ErrHistoriaSinCriterios},
		{"sprint cerrado", "S3", "HU-1", sprint.ErrSprintCerrado},
		{"historia de otro proyecto", "S1", "HU-5", sprint.ErrOtroProyecto},
		{"historia ya asignada a otro sprint", "S1", "HU-6", backlog.ErrYaAsignada},
		{"misma historia dos veces al mismo sprint", "S2", "HU-6", backlog.ErrYaAsignada},
		// Una completada siempre tiene sprint: valida primero el sprint, así no informa "ya está asignada".
		{"historia completada", "S1", "HU-4", sprint.ErrHistoriaCompletada},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, repoSprints, repoHistorias := prepararAsignacion(t)
			spAntes, _ := repoSprints.Buscar(caso.sprint)
			hAntes, _ := repoHistorias.Buscar(caso.historia)

			err := s.Asignar(caso.sprint, caso.historia)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if sp, _ := repoSprints.Buscar(caso.sprint); !slices.Equal(sp.Historias, spAntes.Historias) {
				t.Errorf("historias de %s = %v, antes %v: un error no debe cambiar el sprint", caso.sprint, sp.Historias, spAntes.Historias)
			}
			if h, _ := repoHistorias.Buscar(caso.historia); h.SprintID != hAntes.SprintID {
				t.Errorf("sprint de %s = %q, antes %q: un error no debe cambiar la historia", caso.historia, h.SprintID, hAntes.SprintID)
			}
		})
	}
}

func TestSprints_DelProyecto_proyectoInexistente(t *testing.T) {
	s, _ := nuevoServicioSprints(t)

	if _, err := s.DelProyecto("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("se esperaba el error %v, se obtuvo %v", project.ErrNoExiste, err)
	}
}
