package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
	"github.com/branko-utn/seguimiento-medicion/internal/report/pdf"
)

// rutasReportes registra la sección Reportes (épica E9): la lista de reportes disponibles, el reporte de sprint (HU-29)
// y el de proyecto (HU-30), y la descarga de cada uno en PDF (HU-31).
func (a *app) rutasReportes(mux *http.ServeMux) {
	mux.HandleFunc("GET /reportes", a.listarReportes)
	mux.HandleFunc("GET /reportes/sprints/{sprint}", a.verReporteSprint)
	mux.HandleFunc("GET /reportes/sprints/{sprint}/pdf", a.descargarReporteSprint)
	mux.HandleFunc("GET /reportes/proyecto", a.verReporteProyecto)
	mux.HandleFunc("GET /reportes/proyecto.pdf", a.descargarReporteProyecto)
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

// datosReporteProyecto es lo que recibe reporte_proyecto.html: el reporte y los números ya formateados con un decimal.
type datosReporteProyecto struct {
	report.ReporteProyecto
	Filas           []filaResumen
	Velocidad       string
	Porcentaje      string
	HorasEstimadas  string
	HorasReales     string
	Desviacion      string // vacío si no es calculable
	SprintsCerrados int
}

// filaResumen es una fila del resumen de sprints, con los números formateados.
type filaResumen struct {
	report.ReporteSprint
	Porcentaje, HorasEstimadas, HorasReales string
}

// verReporteProyecto muestra el reporte imprimible del proyecto elegido (?proyecto=P1, HU-30). Si el proyecto no
// existe, da 404.
func (a *app) verReporteProyecto(w http.ResponseWriter, r *http.Request) {
	rep, err := a.svc.Reportes.Proyecto(domain.ProyectoID(r.URL.Query().Get("proyecto")), hoy())
	switch {
	case errors.Is(err, project.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
		return
	case err != nil:
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo armar el reporte del proyecto.")
		return
	}
	datos := datosReporteProyecto{
		ReporteProyecto: rep,
		Velocidad:       unDecimal(rep.Velocidad),
		Porcentaje:      unDecimal(rep.Porcentaje),
		HorasEstimadas:  unDecimal(rep.HorasEstimadas),
		HorasReales:     unDecimal(rep.HorasReales),
	}
	if rep.DesviacionCalculable {
		datos.Desviacion = unDecimal(rep.Desviacion)
	}
	for _, s := range rep.Sprints {
		datos.Filas = append(datos.Filas, filaResumen{ReporteSprint: s, Porcentaje: unDecimal(s.Porcentaje),
			HorasEstimadas: unDecimal(s.HorasEstimadas), HorasReales: unDecimal(s.HorasReales)})
		if !s.Parcial {
			datos.SprintsCerrados++
		}
	}
	a.render(w, r, "reporte_proyecto.html", datos)
}

// descargarReporteSprint descarga el reporte del sprint en PDF (HU-31), con los mismos errores que la pantalla.
func (a *app) descargarReporteSprint(w http.ResponseWriter, r *http.Request) {
	rep, err := a.svc.Reportes.Sprint(domain.SprintID(r.PathValue("sprint")))
	if !a.reporteDisponible(w, r, err, "No se pudo armar el reporte del sprint.") {
		return
	}
	a.enviarPDF(w, r, "reporte-sprint-"+string(rep.SprintID)+".pdf", func(buf io.Writer) error {
		return pdf.Sprint(buf, rep, hoy())
	})
}

// descargarReporteProyecto descarga el reporte del proyecto elegido (?proyecto=P1) en PDF (HU-31).
func (a *app) descargarReporteProyecto(w http.ResponseWriter, r *http.Request) {
	rep, err := a.svc.Reportes.Proyecto(domain.ProyectoID(r.URL.Query().Get("proyecto")), hoy())
	if !a.reporteDisponible(w, r, err, "No se pudo armar el reporte del proyecto.") {
		return
	}
	a.enviarPDF(w, r, "reporte-proyecto-"+string(rep.ProyectoID)+".pdf", func(buf io.Writer) error {
		return pdf.Proyecto(buf, rep, hoy())
	})
}

// reporteDisponible muestra la página de error si no hay reporte (spec HU-31, condiciones de error): proyecto o
// sprint inexistente (404) y sprint planificado (400). Devuelve true si se puede seguir.
func (a *app) reporteDisponible(w http.ResponseWriter, r *http.Request, err error, mensaje string) bool {
	switch {
	case err == nil:
		return true
	case errors.Is(err, project.ErrNoExiste), errors.Is(err, sprint.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case errors.Is(err, report.ErrSprintSinEmpezar):
		a.renderError(w, r, http.StatusBadRequest, err.Error())
	default:
		a.renderError(w, r, http.StatusInternalServerError, mensaje)
	}
	return false
}

// enviarPDF genera el PDF en memoria y recién después lo envía, para no mandar un archivo a medias si falla.
func (a *app) enviarPDF(w http.ResponseWriter, r *http.Request, archivo string, generar func(io.Writer) error) {
	var buf bytes.Buffer
	if err := generar(&buf); err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo generar el PDF.")
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+archivo+`"`)
	_, _ = w.Write(buf.Bytes())
}
