package project

import (
	"errors"
	"testing"
)

func TestNuevoIntegrante(t *testing.T) {
	casos := []struct {
		titulo       string
		nombre       string
		rol          Rol
		quiereError  error
		quiereNombre string
	}{
		{titulo: "registra un Agile Enabler", nombre: "Ana Pérez", rol: AgileEnabler, quiereNombre: "Ana Pérez"},
		{titulo: "registra un Product Builder", nombre: "Beto Gómez", rol: ProductBuilder, quiereNombre: "Beto Gómez"},
		{titulo: "registra un Product Architect", nombre: "Carla Ruiz", rol: ProductArchitect, quiereNombre: "Carla Ruiz"},
		{titulo: "quita los espacios sobrantes del nombre", nombre: "  Ana Pérez  ", rol: ProductBuilder, quiereNombre: "Ana Pérez"},
		{titulo: "rechaza un nombre vacío", nombre: "", rol: ProductBuilder, quiereError: ErrIntegranteSinNombre},
		{titulo: "rechaza un nombre con solo espacios", nombre: "   ", rol: ProductBuilder, quiereError: ErrIntegranteSinNombre},
		{titulo: "rechaza un rol fuera de la lista", nombre: "Ana Pérez", rol: Rol("Tester"), quiereError: ErrRolInvalido},
		{titulo: "rechaza un rol vacío", nombre: "Ana Pérez", rol: Rol(""), quiereError: ErrRolInvalido},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			i, err := NuevoIntegrante("P1", caso.nombre, caso.rol)

			if caso.quiereError != nil {
				if !errors.Is(err, caso.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if i.Nombre != caso.quiereNombre {
				t.Errorf("nombre = %q, se esperaba %q", i.Nombre, caso.quiereNombre)
			}
			if i.Rol != caso.rol {
				t.Errorf("rol = %q, se esperaba %q", i.Rol, caso.rol)
			}
			if i.ProyectoID != "P1" {
				t.Errorf("proyecto = %q, se esperaba %q", i.ProyectoID, "P1")
			}
		})
	}
}
