package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/dashboard"
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Dashboard es el caso de uso del dashboard del proyecto (HU-28). Solo lee: junta los datos de los repositorios y
// los servicios que ya existen, y internal/dashboard arma la vista.
type Dashboard struct {
	proyectos    RepoProyectos
	sprints      RepoSprints
	historias    RepoHistorias
	estimaciones RepoEstimaciones
	registros    RepoRegistros
	defectos     RepoDefectos
}

func NuevoDashboard(proyectos RepoProyectos, sprints RepoSprints, historias RepoHistorias,
	estimaciones RepoEstimaciones, registros RepoRegistros, defectos RepoDefectos) *Dashboard {
	return &Dashboard{proyectos: proyectos, sprints: sprints, historias: historias,
		estimaciones: estimaciones, registros: registros, defectos: defectos}
}

// Del arma el dashboard de un proyecto (HU-28). Devuelve project.ErrNoExiste si el proyecto no existe.
// TODO(HU-28): implementar en la fase GREEN.
func (s *Dashboard) Del(proyectoID domain.ProyectoID, hoy time.Time) (dashboard.Dashboard, error) {
	return dashboard.Dashboard{}, nil
}
