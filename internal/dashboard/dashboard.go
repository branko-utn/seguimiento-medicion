// Package dashboard arma la vista del dashboard del proyecto (épica E8, HU-28).
// No accede a repositorios (ADR 0002): el servicio junta los datos y este paquete los presenta.
package dashboard

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Mensajes que reemplazan a un gráfico cuando no hay datos para dibujar (RN-28.6, CA-28.3).
const (
	MensajeSinSprintsCerrados = "Todavía no hay sprints cerrados para graficar"
	MensajeSinHoras           = "Todavía no hay horas registradas"
	MensajeSinDefectos        = "No hay defectos registrados"
)

// PuntoSprint es el valor de un sprint cerrado en las series de los gráficos (RN-28.1 y RN-28.2).
type PuntoSprint struct {
	Sprint             domain.SprintID
	SPPlanificados     int
	SPCompletados      int
	HorasEstimadas     float64
	HorasReales        float64
	DefectosDetectados int
	DefectosResueltos  int
}

// Mensajes indica, por gráfico, el texto que se muestra en lugar del gráfico cuando no hay datos (RN-28.6).
// Un mensaje vacío significa que el gráfico se dibuja.
type Mensajes struct {
	Velocidad string
	SP        string
	Horas     string
	Defectos  string
}

// Dashboard es el resumen del proyecto que muestra la pantalla /dashboard (HU-28).
type Dashboard struct {
	ProyectoID    domain.ProyectoID
	Nombre        string
	Inicio, Fin   time.Time
	Estado        string
	DiasRestantes int

	SprintActivo          domain.SprintID
	Velocidad             float64
	PorcentajeCompletadas float64
	HorasEstimadas        float64
	HorasReales           float64
	Desviacion            float64
	DesviacionCalculable  bool
	DefectosDetectados    int
	DefectosResueltos     int
	DefectosAbiertos      int

	Puntos   []PuntoSprint
	Mensajes Mensajes
}
