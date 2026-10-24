package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServicioDefectos arma el servicio con dos proyectos: "Demo" (P1), con la historia HU-1 y el sprint S1, y
// "Otro" (P2), con el sprint S2. Devuelve también los repositorios para preparar historias completadas y sprints
// cerrados, que todavía no tienen caso de uso.
func nuevoServicioDefectos(t *testing.T) (*Defectos, *memory.Historias, *memory.Sprints) {
	t.Helper()
	repoProy, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoHist, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	repoSpr, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	repoDef, err := memory.NuevoDefectos("")
	if err != nil {
		t.Fatal(err)
	}
	proy := NuevoProyectos(repoProy)
	hist := NuevoBacklog(repoHist, repoProy)
	sprints := NuevoSprints(repoSpr, repoProy)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proy.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := hist.Crear("P1", "Crear proyecto", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	for _, p := range []domain.ProyectoID{"P1", "P2"} {
		if _, err := sprints.Crear(p, "Sprint de prueba", octubre(21), octubre(23)); err != nil {
			t.Fatal(err)
		}
	}
	return NuevoDefectos(repoDef, repoHist, repoSpr), repoHist, repoSpr
}

func TestDefectos_Registrar(t *testing.T) {
	s, _, _ := nuevoServicioDefectos(t)

	d, err := s.Registrar("Acepta fechas invertidas", defect.Alta, "HU-1", "S1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.ID != "D1" || d.Estado != defect.Abierto || d.SprintResolucion != "" {
		t.Errorf("defecto = %+v, se esperaba D1 Abierto y sin sprint de resolución", d)
	}
	guardado, ok := s.Buscar("D1")
	if !ok || guardado != d {
		t.Errorf("Buscar(D1) = %+v (encontrado %v), se esperaba %+v", guardado, ok, d)
	}
	if _, ok := s.Buscar("D99"); ok {
		t.Errorf("Buscar(D99) no debería encontrar nada")
	}
}

func TestDefectos_Registrar_casosLimite(t *testing.T) {
	t.Run("historia completada", func(t *testing.T) {
		s, repoHist, _ := nuevoServicioDefectos(t)
		h, err := repoHist.Buscar("HU-1")
		if err != nil {
			t.Fatal(err)
		}
		h.Estado = backlog.Completada
		if err := repoHist.Guardar(h); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Registrar("Acepta fechas invertidas", defect.Critica, "HU-1", "S1"); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
	})
	t.Run("sprint cerrado", func(t *testing.T) {
		s, _, repoSpr := nuevoServicioDefectos(t)
		sp, err := repoSpr.Buscar("S1")
		if err != nil {
			t.Fatal(err)
		}
		sp.Estado = sprint.Cerrado
		if err := repoSpr.Actualizar(sp); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Registrar("Acepta fechas invertidas", defect.Alta, "HU-1", "S1"); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
	})
}

func TestDefectos_Registrar_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		descripcion string
		severidad   defect.Severidad
		historia    domain.HistoriaID
		sprint      domain.SprintID
		quiereError error
	}{
		{"descripción vacía", "  ", defect.Alta, "HU-1", "S1", defect.ErrDescripcionVacia},
		{"severidad inválida", "Acepta fechas invertidas", 0, "HU-1", "S1", defect.ErrSeveridadInvalida},
		{"historia inexistente", "Acepta fechas invertidas", defect.Alta, "HU-99", "S1", backlog.ErrNoExiste},
		{"sprint inexistente", "Acepta fechas invertidas", defect.Alta, "HU-1", "S99", sprint.ErrNoExiste},
		{"sprint de otro proyecto", "Acepta fechas invertidas", defect.Alta, "HU-1", "S2", defect.ErrSprintDeOtroProyecto},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _, _ := nuevoServicioDefectos(t)
			_, err := s.Registrar(c.descripcion, c.severidad, c.historia, c.sprint)
			if !errors.Is(err, c.quiereError) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiereError)
			}
			if _, ok := s.Buscar("D1"); ok {
				t.Errorf("el defecto rechazado no debe guardarse")
			}
		})
	}
}
