package server

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// Tonos de las insignias (TEC-21): la clase CSS la decide el handler con el enum (ADR 0003), así las plantillas no
// comparan textos de estado. Cada parte de HU-32 agrega a sus filas un campo Tono con estas funciones.
const (
	tonoNeutro  = "tono-neutro"
	tonoInfo    = "tono-info"
	tonoExito   = "tono-exito"
	tonoAlerta  = "tono-alerta"
	tonoPeligro = "tono-peligro"
)

// tonoEstado devuelve el tono de la insignia de un estado de proyecto, sprint, historia o defecto, según el avance:
// neutro si todavía no empezó, info si está en curso y éxito si terminó. Un defecto abierto es peligro y uno en
// progreso, alerta. Cualquier otro valor (el cero de cada tipo, o algo que no es un estado) queda en neutro.
func tonoEstado(estado any) string {
	switch e := estado.(type) {
	case project.Estado:
		switch e {
		case project.EnCurso:
			return tonoInfo
		case project.Finalizado:
			return tonoExito
		}
	case sprint.Estado:
		switch e {
		case sprint.Activo:
			return tonoInfo
		case sprint.Cerrado:
			return tonoExito
		}
	case backlog.Estado:
		switch e {
		case backlog.EnProgreso:
			return tonoInfo
		case backlog.Completada:
			return tonoExito
		}
	case defect.Estado:
		switch e {
		case defect.Abierto:
			return tonoPeligro
		case defect.EnProgreso:
			return tonoAlerta
		case defect.Resuelto:
			return tonoExito
		}
	}
	return tonoNeutro
}

// tonoSeveridad devuelve el tono de la insignia de la severidad de un defecto, de más grave (peligro) a menos grave
// (neutro).
func tonoSeveridad(s defect.Severidad) string {
	switch s {
	case defect.Critica:
		return tonoPeligro
	case defect.Alta:
		return tonoAlerta
	case defect.Media:
		return tonoInfo
	}
	return tonoNeutro
}
