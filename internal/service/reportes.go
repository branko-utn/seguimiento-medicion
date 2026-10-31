package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// Reportes son los casos de uso de los reportes (épica E9). Solo leen: juntan los datos de los repositorios y
// internal/report arma el reporte (ADR 0002).
type Reportes struct {
	proyectos    RepoProyectos
	sprints      RepoSprints
	historias    RepoHistorias
	estimaciones RepoEstimaciones
	registros    RepoRegistros
	integrantes  BuscadorIntegrantes
	defectos     RepoDefectos
}

func NuevoReportes(proyectos RepoProyectos, sprints RepoSprints, historias RepoHistorias, estimaciones RepoEstimaciones,
	registros RepoRegistros, integrantes BuscadorIntegrantes, defectos RepoDefectos) *Reportes {
	return &Reportes{proyectos: proyectos, sprints: sprints, historias: historias, estimaciones: estimaciones, registros: registros,
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
		datos.Registros = append(datos.Registros, s.registrosConNombres(hid, datos.Integrantes)...)
	}
	datos.Defectos = s.defectosDelProyecto(sp.ProyectoID)
	return report.NuevoReporteSprint(datos)
}

// Proyecto arma el reporte del proyecto completo (HU-30), con su estado al día de hoy. Devuelve project.ErrNoExiste si
// el proyecto no existe.
func (s *Reportes) Proyecto(id domain.ProyectoID, hoy time.Time) (report.ReporteProyecto, error) {
	p, err := s.proyectos.Buscar(id)
	if err != nil {
		return report.ReporteProyecto{}, err
	}
	datos := report.DatosProyecto{
		Proyecto:    p,
		Hoy:         hoy,
		Sprints:     s.sprints.DelProyecto(id),
		Estimadas:   make(map[domain.HistoriaID]float64),
		Integrantes: make(map[domain.IntegranteID]string),
		Defectos:    s.defectosDelProyecto(id),
	}
	for _, h := range s.historias.Listar() {
		if h.ProyectoID != id {
			continue
		}
		datos.Historias = append(datos.Historias, h)
		if e, ok := s.estimaciones.Buscar(h.ID); ok {
			datos.Estimadas[h.ID] = e.Horas
		}
		datos.Registros = append(datos.Registros, s.registrosConNombres(h.ID, datos.Integrantes)...)
	}
	return report.NuevoReporteProyecto(datos)
}

// registrosConNombres devuelve los registros de la historia y anota en nombres el de cada integrante que aparece.
func (s *Reportes) registrosConNombres(id domain.HistoriaID, nombres map[domain.IntegranteID]string) []effort.Registro {
	registros := s.registros.DeHistoria(id)
	for _, reg := range registros {
		if _, ok := nombres[reg.IntegranteID]; !ok {
			if i, err := s.integrantes.Buscar(reg.IntegranteID); err == nil {
				nombres[reg.IntegranteID] = i.Nombre
			}
		}
	}
	return registros
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
