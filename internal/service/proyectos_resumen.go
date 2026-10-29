package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// ResumenProyecto es el estado del proyecto que se muestra de un vistazo (HU-04).
type ResumenProyecto struct {
	Proyecto      project.Proyecto // con sus fechas de inicio y fin
	Estado        project.Estado   // calculado con las fechas (RN-04.1)
	DiasRestantes int              // RN-04.2
	// SprintActivo es el sprint del proyecto en estado Activo, o nil si no hay ("sin sprint activo", RN-04.3).
	SprintActivo          *sprint.Sprint
	PorcentajeCompletadas float64 // RN-04.4
}

// ConResumen le da al servicio lo que necesita para armar el resumen (HU-04): los sprints y las historias del
// proyecto. Todavía sin implementar (RED).
func (s *Proyectos) ConResumen(sprints RepoSprints, historias RepoHistorias) *Proyectos {
	return s
}

// Resumen arma el estado del proyecto al día de hoy (HU-04). Todavía sin implementar (RED).
func (s *Proyectos) Resumen(id domain.ProyectoID, hoy time.Time) (ResumenProyecto, error) {
	return ResumenProyecto{}, nil
}
