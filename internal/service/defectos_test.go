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

// nuevoServicioDefectos arma el servicio con dos proyectos: "Demo" (P1), con la historia HU-1 y los sprints S1 (21 al
// 23/10) y S3 (24 al 26/10), y "Otro" (P2), con el sprint S2. Devuelve también los repositorios para preparar historias completadas y sprints
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
	sprints := NuevoSprints(repoSpr, repoProy, repoHist)
	hist := NuevoBacklog(repoHist, repoProy, sprints, nil, nil)
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
	if _, err := sprints.Crear("P1", "Sprint de prueba", octubre(24), octubre(26)); err != nil {
		t.Fatal(err)
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
	guardado, err := s.Buscar("D1")
	if err != nil || guardado != d {
		t.Errorf("Buscar(D1) = %+v (error %v), se esperaba %+v", guardado, err, d)
	}
	if _, err := s.Buscar("D99"); !errors.Is(err, defect.ErrNoExiste) {
		t.Errorf("Buscar(D99): error = %v, se esperaba defect.ErrNoExiste", err)
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
			if _, err := s.Buscar("D1"); err == nil {
				t.Errorf("el defecto rechazado no debe guardarse")
			}
		})
	}
}

// registrarD1 registra el defecto D1 en HU-1, detectado en el sprint indicado.
func registrarD1(t *testing.T, s *Defectos, deteccion domain.SprintID) {
	t.Helper()
	if _, err := s.Registrar("Acepta fechas invertidas", defect.Alta, "HU-1", deteccion); err != nil {
		t.Fatal(err)
	}
}

func TestDefectos_CambiarEstado(t *testing.T) {
	s, _, _ := nuevoServicioDefectos(t)
	registrarD1(t, s, "S1")

	if err := s.CambiarEstado("D1", defect.EnProgreso); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d, _ := s.Buscar("D1"); d.Estado != defect.EnProgreso {
		t.Errorf("Estado = %v, se esperaba En progreso", d.Estado)
	}
}

func TestDefectos_CambiarEstado_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		id          domain.DefectoID
		nuevo       defect.Estado
		quiereError error
	}{
		{"defecto inexistente", "D9", defect.EnProgreso, defect.ErrNoExiste},
		{"al mismo estado", "D1", defect.Abierto, defect.ErrTransicionInvalida},
		{"a Resuelto sin sprint", "D1", defect.Resuelto, defect.ErrResolverSinSprint},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _, _ := nuevoServicioDefectos(t)
			registrarD1(t, s, "S1")
			if err := s.CambiarEstado(c.id, c.nuevo); !errors.Is(err, c.quiereError) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiereError)
			}
			if d, _ := s.Buscar("D1"); d.Estado != defect.Abierto {
				t.Errorf("el defecto cambió a %v aunque hubo error", d.Estado)
			}
		})
	}
}

func TestDefectos_Resolver(t *testing.T) {
	casos := []struct {
		titulo     string
		deteccion  domain.SprintID
		resolucion domain.SprintID
	}{
		{"en un sprint posterior", "S1", "S3"},
		{"en el mismo sprint de detección", "S3", "S3"},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _, _ := nuevoServicioDefectos(t)
			registrarD1(t, s, c.deteccion)
			if err := s.Resolver("D1", c.resolucion); err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			d, _ := s.Buscar("D1")
			if d.Estado != defect.Resuelto || d.SprintResolucion != c.resolucion {
				t.Errorf("defecto = %+v, se esperaba Resuelto en %s", d, c.resolucion)
			}
		})
	}
}

func TestDefectos_Resolver_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		id          domain.DefectoID
		resolucion  domain.SprintID
		quiereError error
	}{
		{"defecto inexistente", "D9", "S3", defect.ErrNoExiste},
		{"sin sprint", "D1", "", defect.ErrResolverSinSprint},
		{"sprint inexistente", "D1", "S99", sprint.ErrNoExiste},
		{"sprint de otro proyecto", "D1", "S2", defect.ErrResolucionDeOtroProyecto},
		{"sprint anterior al de detección", "D1", "S1", defect.ErrResolucionAnterior},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _, _ := nuevoServicioDefectos(t)
			registrarD1(t, s, "S3")
			if err := s.Resolver(c.id, c.resolucion); !errors.Is(err, c.quiereError) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiereError)
			}
			if d, _ := s.Buscar("D1"); d.Estado != defect.Abierto || d.SprintResolucion != "" {
				t.Errorf("el defecto cambió aunque hubo error: %+v", d)
			}
		})
	}
}

func TestDefectos_Resolver_yaResuelto(t *testing.T) {
	s, _, _ := nuevoServicioDefectos(t)
	registrarD1(t, s, "S1")
	if err := s.Resolver("D1", "S1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Resolver("D1", "S3"); !errors.Is(err, defect.ErrTransicionInvalida) {
		t.Fatalf("error = %v, se esperaba defect.ErrTransicionInvalida", err)
	}
	if d, _ := s.Buscar("D1"); d.SprintResolucion != "S1" {
		t.Errorf("SprintResolucion = %q, debía seguir en S1", d.SprintResolucion)
	}
}

func TestDefectos_Listar(t *testing.T) {
	s, _, _ := nuevoServicioDefectos(t)
	if n := len(s.Listar()); n != 0 {
		t.Fatalf("Listar() sin defectos devolvió %d", n)
	}
	registrarD1(t, s, "S1")
	if _, err := s.Registrar("Texto cortado", defect.Baja, "HU-1", "S3"); err != nil {
		t.Fatal(err)
	}
	lista := s.Listar()
	if len(lista) != 2 || lista[0].ID != "D1" || lista[1].ID != "D2" {
		t.Errorf("Listar() = %+v, se esperaba [D1 D2] en ese orden", lista)
	}
}
