package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

// ResumenEsfuerzo es el esfuerzo registrado en un proyecto (HU-21).
type ResumenEsfuerzo struct {
	Registros     []RegistroEsfuerzo
	PorHistoria   []TotalHistoria
	PorIntegrante []TotalIntegrante
	Total         float64
}

// RegistroEsfuerzo es un registro de esfuerzo con el nombre de su integrante.
type RegistroEsfuerzo struct {
	effort.Registro
	Integrante string
}

// TotalHistoria es el total de horas registradas en una historia.
type TotalHistoria struct {
	Historia domain.HistoriaID
	Horas    float64
}

// TotalIntegrante es el total de horas registradas por un integrante.
type TotalIntegrante struct {
	Integrante domain.IntegranteID
	Nombre     string
	Horas      float64
}

// DelProyecto devuelve el resumen del esfuerzo registrado en un proyecto (HU-21).
func (s *Registros) DelProyecto(id domain.ProyectoID) (ResumenEsfuerzo, error) {
	return ResumenEsfuerzo{}, nil
}
