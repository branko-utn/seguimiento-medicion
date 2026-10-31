package server

import "github.com/branko-utn/seguimiento-medicion/internal/domain/defect"

// Tonos de las insignias (TEC-21): la clase CSS la decide el handler con el enum (ADR 0003), así las plantillas no
// comparan textos de estado. Cada parte de HU-32 agrega a sus filas un campo Tono con estas funciones.
const (
	tonoNeutro  = "tono-neutro"
	tonoInfo    = "tono-info"
	tonoExito   = "tono-exito"
	tonoAlerta  = "tono-alerta"
	tonoPeligro = "tono-peligro"
)

// tonoEstado devuelve el tono de la insignia de un estado de proyecto, sprint, historia o defecto. Todavía sin
// implementar (RED).
func tonoEstado(estado any) string {
	return ""
}

// tonoSeveridad devuelve el tono de la insignia de la severidad de un defecto. Todavía sin implementar (RED).
func tonoSeveridad(s defect.Severidad) string {
	return ""
}
