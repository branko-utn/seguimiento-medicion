package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// Reportes son los casos de uso de los reportes (épica E9). Solo leen: juntan los datos de los repositorios y
// internal/report arma el reporte (ADR 0002).
type Reportes struct {
	sprints      RepoSprints
	historias    RepoHistorias
	estimaciones RepoEstimaciones
	registros    RepoRegistros
	integrantes  BuscadorIntegrantes
	defectos     RepoDefectos
}

func NuevoReportes(sprints RepoSprints, historias RepoHistorias, estimaciones RepoEstimaciones, registros RepoRegistros,
	integrantes BuscadorIntegrantes, defectos RepoDefectos) *Reportes {
	return &Reportes{sprints: sprints, historias: historias, estimaciones: estimaciones, registros: registros,
		integrantes: integrantes, defectos: defectos}
}

// Sprint arma el reporte de un sprint activo o cerrado (HU-29). Devuelve sprint.ErrNoExiste si no existe y
// report.ErrSprintSinEmpezar si está planificado (RN-29.1).
func (s *Reportes) Sprint(id domain.SprintID) (report.ReporteSprint, error) {
	sp, err := s.sprints.Buscar(id)
	if err != nil {
		return report.ReporteSprint{}, err
	}
	datos := report.DatosSprint{
		Sprint:      sp,
		Historias:   make(map[domain.HistoriaID]backlog.Historia, len(sp.Historias)),
		Estimadas:   make(map[domain.HistoriaID]float64, len(sp.Historias)),
		Integrantes: make(map[domain.IntegranteID]string),
	}
	for _, hid := range sp.Historias {
		h, err := s.historias.Buscar(hid)
		if err != nil {
			return report.ReporteSprint{}, err
		}
		datos.Historias[hid] = h
		if e, ok := s.estimaciones.Buscar(hid); ok {
			datos.Estimadas[hid] = e.Horas
		}
		for _, reg := range s.registros.DeHistoria(hid) {
			datos.Registros = append(datos.Registros, reg)
			if _, ok := datos.Integrantes[reg.IntegranteID]; !ok {
				if i, err := s.integrantes.Buscar(reg.IntegranteID); err == nil {
					datos.Integrantes[reg.IntegranteID] = i.Nombre
				}
			}
		}
	}
	datos.Defectos = s.defectosDelProyecto(sp.ProyectoID)
	return report.NuevoReporteSprint(datos)
}

// defectosDelProyecto devuelve los defectos de las historias del proyecto, en el orden en que se registraron.
func (s *Reportes) defectosDelProyecto(proyecto domain.ProyectoID) []defect.Defecto {
	delProyecto := make(map[domain.HistoriaID]bool)
	for _, h := range s.historias.Listar() {
		if h.ProyectoID == proyecto {
			delProyecto[h.ID] = true
		}
	}
	var ds []defect.Defecto
	for _, d := range s.defectos.Listar() {
		if delProyecto[d.HistoriaID] {
			ds = append(ds, d)
		}
	}
	return ds
}
