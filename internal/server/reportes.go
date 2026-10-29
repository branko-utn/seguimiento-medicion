package server

import (
	"errors"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// rutasReportes registra la sección Reportes (épica E9): la lista de reportes disponibles y el reporte de sprint (HU-29).
func (a *app) rutasReportes(mux *http.ServeMux) {
	mux.HandleFunc("GET /reportes", a.listarReportes)
	mux.HandleFunc("GET /reportes/sprints/{sprint}", a.verReporteSprint)
}

// grupoReportes es un proyecto con los sprints que tienen reporte: los activos y los cerrados (RN-29.1).
type grupoReportes struct {
	Proyecto project.Proyecto
	Sprints  []sprint.Sprint
}

func (a *app) listarReportes(w http.ResponseWriter, r *http.Request) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los proyectos.")
		return
	}
	grupos := make([]grupoReportes, 0, len(proyectos))
	for _, p := range proyectos {
		sprints, err := a.svc.Sprints.DelProyecto(p.ID)
		if err != nil {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints.")
			return
		}
		g := grupoReportes{Proyecto: p}
		for _, s := range sprints {
			if s.Estado != sprint.Planificado {
				g.Sprints = append(g.Sprints, s)
			}
		}
		grupos = append(grupos, g)
	}
	a.render(w, r, "reportes.html", grupos)
}

// datosReporteSprint es lo que recibe reporte_sprint.html: el reporte y los números ya formateados con un decimal.
type datosReporteSprint struct {
	report.ReporteSprint
	Filas          []filaReporte
	Porcentaje     string
	HorasEstimadas string
	HorasReales    string
	Desviacion     string // vacío si no es calculable (RN-29.4)
	Registros      []registroReporte
	Totales        []totalReporte
}

type filaReporte struct {
	report.FilaHistoria
	HorasEstimadas, HorasReales string
}

type registroReporte struct {
	report.FilaRegistro
	Horas string
}

type totalReporte struct {
	Integrante, Horas string
}

// verReporteSprint muestra el reporte imprimible de un sprint (HU-29). Un sprint planificado todavía no tiene
// reporte (400) y uno inexistente da 404.
func (a *app) verReporteSprint(w http.ResponseWriter, r *http.Request) {
	rep, err := a.svc.Reportes.Sprint(domain.SprintID(r.PathValue("sprint")))
	switch {
	case errors.Is(err, sprint.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, report.ErrSprintSinEmpezar):
		a.renderError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo armar el reporte del sprint.")
		return
	}
	datos := datosReporteSprint{
		ReporteSprint:  rep,
		Porcentaje:     unDecimal(rep.Porcentaje),
		HorasEstimadas: unDecimal(rep.HorasEstimadas),
		HorasReales:    unDecimal(rep.HorasReales),
	}
	if rep.DesviacionCalculable {
		datos.Desviacion = unDecimal(rep.Desviacion)
	}
	for _, h := range rep.Historias {
		datos.Filas = append(datos.Filas, filaReporte{FilaHistoria: h, HorasEstimadas: unDecimal(h.HorasEstimadas), HorasReales: unDecimal(h.HorasReales)})
	}
	for _, reg := range rep.Registros {
		datos.Registros = append(datos.Registros, registroReporte{FilaRegistro: reg, Horas: unDecimal(reg.Horas)})
	}
	for _, t := range rep.PorIntegrante {
		datos.Totales = append(datos.Totales, totalReporte{Integrante: t.Integrante, Horas: unDecimal(t.Horas)})
	}
	a.render(w, r, "reporte_sprint.html", datos)
}
