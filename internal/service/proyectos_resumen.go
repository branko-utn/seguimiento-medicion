package service

import (
	"errors"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// errResumenSinConfigurar es un error de programación: el servicio se armó sin ConResumen. No es una regla de la spec.
var errResumenSinConfigurar = errors.New("el servicio de proyectos no tiene sprints ni historias: falta ConResumen")

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
// proyecto. Devuelve el mismo servicio, para encadenarlo con NuevoProyectos.
func (s *Proyectos) ConResumen(sprints RepoSprints, historias RepoHistorias) *Proyectos {
	s.sprints = sprints
	s.historias = historias
	return s
}

// Resumen arma el estado del proyecto al día de hoy (HU-04): el estado y los días restantes salen de las fechas
// (Proyecto.EstadoAl y DiasRestantes, TEC-20), el sprint activo de sus sprints (RN-04.3) y el porcentaje de las
// historias del backlog del proyecto con la métrica de HU-25 (RN-04.4). No modifica nada.
func (s *Proyectos) Resumen(id domain.ProyectoID, hoy time.Time) (ResumenProyecto, error) {
	if s.sprints == nil || s.historias == nil {
		return ResumenProyecto{}, errResumenSinConfigurar
	}
	p, err := s.repo.Buscar(id)
	if err != nil {
		return ResumenProyecto{}, err
	}
	r := ResumenProyecto{Proyecto: p, Estado: p.EstadoAl(hoy), DiasRestantes: p.DiasRestantes(hoy)}

	for _, sp := range s.sprints.DelProyecto(id) {
		if sp.Estado == sprint.Activo {
			r.SprintActivo = &sp
			break
		}
	}

	total, completadas := 0, 0
	for _, h := range s.historias.Listar() {
		if h.ProyectoID != id {
			continue
		}
		total++
		if h.Estado == backlog.Completada {
			completadas++
		}
	}
	if r.PorcentajeCompletadas, err = metrics.PorcentajeCompletadas(completadas, total); err != nil {
		return ResumenProyecto{}, err
	}
	return r, nil
}
