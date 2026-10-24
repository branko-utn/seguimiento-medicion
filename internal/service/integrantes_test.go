package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServicioIntegrantes arma el servicio con dos proyectos ya creados: "Demo" (P1) y "Otro" (P2).
func nuevoServicioIntegrantes(t *testing.T) *Integrantes {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoIntegrantes, err := memory.NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}
	proyectos := NuevoProyectos(repoProyectos)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	return NuevoIntegrantes(repoIntegrantes, repoProyectos)
}

func TestIntegrantes_Registrar(t *testing.T) {
	s := nuevoServicioIntegrantes(t)

	i, err := s.Registrar("P1", "  Ana Pérez  ", project.ProductBuilder)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if i.ID != "I1" || i.ProyectoID != "P1" || i.Nombre != "Ana Pérez" || i.Rol != project.ProductBuilder {
		t.Errorf("Registrar() = %+v, se esperaba I1 · P1 · Ana Pérez · Product Builder", i)
	}
}

func TestIntegrantes_DelProyecto(t *testing.T) {
	s := nuevoServicioIntegrantes(t)

	vacio, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatalf("proyecto sin integrantes: error inesperado: %v", err)
	}
	if len(vacio) != 0 {
		t.Errorf("proyecto sin integrantes: se obtuvieron %d, se esperaba 0", len(vacio))
	}

	altas := []struct {
		proyecto domain.ProyectoID
		nombre   string
		rol      project.Rol
	}{
		{"P1", "Ana Pérez", project.AgileEnabler},
		{"P2", "Diego Paz", project.ProductBuilder},
		{"P1", "Beto Gómez", project.ProductBuilder},
		{"P1", "Carla Ruiz", project.ProductArchitect},
	}
	for _, a := range altas {
		if _, err := s.Registrar(a.proyecto, a.nombre, a.rol); err != nil {
			t.Fatalf("Registrar %s: %v", a.nombre, err)
		}
	}

	lista, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatal(err)
	}
	quiere := []string{"Ana Pérez", "Beto Gómez", "Carla Ruiz"}
	if len(lista) != len(quiere) {
		t.Fatalf("DelProyecto(P1) = %+v, se esperaban %v", lista, quiere)
	}
	for n, nombre := range quiere {
		if lista[n].Nombre != nombre {
			t.Errorf("integrante %d = %q, se esperaba %q (en orden de registro)", n, lista[n].Nombre, nombre)
		}
	}
}

func TestIntegrantes_DelProyecto_proyectoInexistente(t *testing.T) {
	s := nuevoServicioIntegrantes(t)

	if _, err := s.DelProyecto("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("se esperaba el error %v, se obtuvo %v", project.ErrNoExiste, err)
	}
}

func TestIntegrantes_MismaPersonaEnOtroProyecto(t *testing.T) {
	s := nuevoServicioIntegrantes(t)
	if _, err := s.Registrar("P1", "Ana Pérez", project.ProductBuilder); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Registrar("P2", "Ana Pérez", project.ProductBuilder); err != nil {
		t.Errorf("RN-03.4: la misma persona puede estar en otro proyecto, pero se obtuvo %v", err)
	}
}

func TestIntegrantes_Registrar_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		proyecto    domain.ProyectoID
		nombre      string
		rol         project.Rol
		quiereError error
	}{
		{"proyecto inexistente", "P9", "Beto Gómez", project.ProductBuilder, project.ErrNoExiste},
		{"primero se verifica el proyecto", "P9", "", project.Rol(99), project.ErrNoExiste},
		{"duplicado sin distinguir mayúsculas", "P1", "ana pérez", project.ProductBuilder, project.ErrIntegranteDuplicado},
		{"duplicado con espacios sobrantes", "P1", "  ANA PÉREZ ", project.AgileEnabler, project.ErrIntegranteDuplicado},
		{"nombre vacío del dominio", "P1", "   ", project.ProductBuilder, project.ErrIntegranteSinNombre},
		{"rol inválido del dominio", "P1", "Beto Gómez", project.Rol(99), project.ErrRolInvalido},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := nuevoServicioIntegrantes(t)
			if _, err := s.Registrar("P1", "Ana Pérez", project.ProductBuilder); err != nil {
				t.Fatalf("no se pudo registrar el integrante inicial: %v", err)
			}

			_, err := s.Registrar(caso.proyecto, caso.nombre, caso.rol)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}

			lista, _ := s.DelProyecto("P1")
			if len(lista) != 1 {
				t.Errorf("hay %d integrantes en P1, se esperaba 1: no se debe guardar uno inválido", len(lista))
			}
		})
	}
}
