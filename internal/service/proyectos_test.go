package service

import (
	"errors"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

var (
	inicio = time.Date(2026, time.October, 21, 0, 0, 0, 0, time.Local)
	fin    = time.Date(2026, time.November, 2, 0, 0, 0, 0, time.Local)
)

func nuevoServicioProyectos(t *testing.T) *Proyectos {
	t.Helper()
	repo, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	return NuevoProyectos(repo)
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
