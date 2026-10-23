package project

import (
	"errors"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Rol es el rol de un integrante en el proyecto. Su valor es el texto que ve el usuario.
type Rol string

const (
	AgileEnabler     Rol = "Agile Enabler"
	ProductBuilder   Rol = "Product Builder"
	ProductArchitect Rol = "Product Architect"
)

// Integrante es una persona que forma parte de un proyecto (ADR 0002).
type Integrante struct {
	ID         domain.IntegranteID
	ProyectoID domain.ProyectoID
	Nombre     string
	Rol        Rol
}

var (
	ErrIntegranteSinNombre = errors.New("el nombre del integrante es obligatorio")
	ErrRolInvalido         = errors.New("el rol debe ser Agile Enabler, Product Builder o Product Architect")
	ErrIntegranteDuplicado = errors.New("el integrante ya está registrado en el proyecto")
)

// NuevoIntegrante crea un integrante válido (RN-03.1 y RN-03.2). Que el proyecto exista y que el integrante no
// esté repetido (RN-03.3) lo controla el servicio.
func NuevoIntegrante(proyectoID domain.ProyectoID, nombre string, rol Rol) (Integrante, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return Integrante{}, ErrIntegranteSinNombre
	}
	switch rol {
	case AgileEnabler, ProductBuilder, ProductArchitect:
	default:
		return Integrante{}, ErrRolInvalido
	}
	return Integrante{ProyectoID: proyectoID, Nombre: nombre, Rol: rol}, nil
}
