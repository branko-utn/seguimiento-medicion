package service

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

var (
	inicio = time.Date(2026, time.October, 21, 0, 0, 0, 0, time.Local)
	fin    = time.Date(2026, time.November, 2, 0, 0, 0, 0, time.Local)
)

func nuevoServicioProyectos(t *testing.T) *Proyectos {
	t.Helper()
	s, _ := armarServicioProyectos(t)
	return s
}

// armarServicioProyectos es nuevoServicioProyectos, pero devuelve también el repositorio para preparar estados que
// todavía no tienen caso de uso (Finalizado).
func armarServicioProyectos(t *testing.T) (*Proyectos, *memory.Proyectos) {
	t.Helper()
	repo, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	return NuevoProyectos(repo), repo
}

func TestProyectos_Crear(t *testing.T) {
	s := nuevoServicioProyectos(t)

	p, err := s.Crear("  Software Metrics  ", "TPI", inicio, fin)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.ID != "P1" {
		t.Errorf("ID = %q, se esperaba %q", p.ID, "P1")
	}
	if p.Nombre != "Software Metrics" {
		t.Errorf("nombre = %q, se esperaba %q", p.Nombre, "Software Metrics")
	}
	if p.Estado != project.Planificado {
		t.Errorf("estado = %q, se esperaba %q", p.Estado, project.Planificado)
	}

	lista, err := s.Listar()
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 1 || lista[0].ID != "P1" {
		t.Errorf("Listar() = %+v, se esperaba solo el proyecto P1", lista)
	}
}

func TestProyectos_Crear_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		nombre      string
		inicio, fin time.Time
		quiereError error
	}{
		{"nombre duplicado sin distinguir mayúsculas", "demo", inicio, fin, project.ErrNombreDuplicado},
		{"nombre duplicado con espacios sobrantes", "  DEMO ", inicio, fin, project.ErrNombreDuplicado},
		{"nombre inválido del dominio", "ab", inicio, fin, project.ErrNombreInvalido},
		{"fechas inválidas del dominio", "Otro", fin, inicio, project.ErrFechasInvalidas},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := nuevoServicioProyectos(t)
			if _, err := s.Crear("Demo", "", inicio, fin); err != nil {
				t.Fatalf("no se pudo crear el proyecto inicial: %v", err)
			}

			_, err := s.Crear(caso.nombre, "", caso.inicio, caso.fin)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}

			lista, _ := s.Listar()
			if len(lista) != 1 {
				t.Errorf("hay %d proyectos guardados, se esperaba 1: no se debe guardar uno inválido", len(lista))
			}
		})
	}
}

// prepararModificacion crea "Demo" (P1) y "Otro" (P2), del 21/10 al 02/11.
func prepararModificacion(t *testing.T) (*Proyectos, *memory.Proyectos) {
	t.Helper()
	s, repo := armarServicioProyectos(t)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := s.Crear(nombre, "Original", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	return s, repo
}

func TestProyectos_Modificar(t *testing.T) {
	s, repo := prepararModificacion(t)
	nuevoInicio, nuevoFin := inicio.AddDate(0, 0, 1), fin.AddDate(0, 0, 1)

	p, err := s.Modificar("P1", "  Software Metrics  ", "TPI", nuevoInicio, nuevoFin)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	guardado, _ := repo.Buscar("P1")
	for _, got := range []project.Proyecto{p, guardado} {
		if got.ID != "P1" || got.Nombre != "Software Metrics" || got.Descripcion != "TPI" || got.Estado != project.Planificado {
			t.Errorf("proyecto = %+v, se esperaba P1 · Software Metrics · TPI · Planificado", got)
		}
		if !got.Inicio.Equal(nuevoInicio) || !got.Fin.Equal(nuevoFin) {
			t.Errorf("fechas = %v al %v, se esperaba %v al %v", got.Inicio, got.Fin, nuevoInicio, nuevoFin)
		}
	}
	if otro, _ := repo.Buscar("P2"); otro.Nombre != "Otro" {
		t.Errorf("P2 = %+v: modificar P1 no debe cambiar otro proyecto", otro)
	}
}

func TestProyectos_Modificar_mismoNombre(t *testing.T) {
	for _, nombre := range []string{"Demo", "DEMO", " demo "} {
		t.Run(nombre, func(t *testing.T) {
			s, _ := prepararModificacion(t)

			if _, err := s.Modificar("P1", nombre, "", inicio, fin); err != nil {
				t.Errorf("mantener el propio nombre (o cambiarle las mayúsculas) es válido, pero se obtuvo %v", err)
			}
		})
	}
}

func TestProyectos_Modificar_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		proyecto    domain.ProyectoID
		nombre      string
		inicio, fin time.Time
		quiereError error
	}{
		{"proyecto inexistente", "P9", "Nuevo", inicio, fin, project.ErrNoExiste},
		{"primero se busca el proyecto", "P9", "AB", fin, inicio, project.ErrNoExiste},
		{"nombre de otro proyecto", "P1", "otro", inicio, fin, project.ErrNombreDuplicado},
		{"nombre de otro proyecto con espacios", "P1", "  OTRO ", inicio, fin, project.ErrNombreDuplicado},
		{"nombre inválido", "P1", "AB", inicio, fin, project.ErrNombreInvalido},
		{"fin anterior al inicio", "P1", "Demo", fin, inicio, project.ErrFechasInvalidas},
		{"proyecto finalizado", "P3", "Final 2", inicio, fin, project.ErrProyectoFinalizado},
		{"finalizado antes que el nombre repetido", "P3", "Otro", inicio, fin, project.ErrProyectoFinalizado},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, repo := prepararModificacion(t)
			// P3 "Final" en estado Finalizado: todavía no hay caso de uso que finalice un proyecto.
			if _, err := repo.Agregar(project.Proyecto{Nombre: "Final", Descripcion: "Original", Inicio: inicio, Fin: fin, Estado: project.Finalizado}); err != nil {
				t.Fatal(err)
			}
			antes, _ := s.Listar()

			_, err := s.Modificar(caso.proyecto, caso.nombre, "Cambiada", caso.inicio, caso.fin)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if despues, _ := s.Listar(); !slices.Equal(despues, antes) {
				t.Errorf("los proyectos quedaron %+v, antes %+v: un error no debe cambiar nada", despues, antes)
			}
		})
	}
}

func TestProyectos_Buscar(t *testing.T) {
	s, _ := prepararModificacion(t)

	p, err := s.Buscar("P2")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.ID != "P2" || p.Nombre != "Otro" {
		t.Errorf("Buscar(P2) = %+v, se esperaba el proyecto Otro", p)
	}

	if _, err := s.Buscar("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("Buscar(P9): se esperaba %v, se obtuvo %v", project.ErrNoExiste, err)
	}
}
