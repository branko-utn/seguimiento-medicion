package server

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/dashboard"
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

func (a *app) rutasDashboard(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", a.dashboard)
}

func (a *app) dashboard(w http.ResponseWriter, r *http.Request) {
	proyecto := r.URL.Query().Get("proyecto")
	if proyecto == "" {
		a.renderDashboardSelector(w, r)
		return
	}
	pid := domain.ProyectoID(proyecto)
	d, err := a.svc.Dashboard.Del(pid, hoy())
	if err != nil {
		if errors.Is(err, project.ErrNoExiste) {
			a.renderError(w, r, http.StatusNotFound, err.Error())
			return
		}
		a.renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	type seriesData struct {
		Puntos    []dashboard.PuntoSprint `json:"puntos"`
		Mensajes  dashboard.Mensajes      `json:"mensajes"`
		Velocidad float64                 `json:"velocidad"` // línea de referencia del gráfico de velocidad
	}
	puntos := d.Puntos
	if puntos == nil {
		puntos = []dashboard.PuntoSprint{}
	}
	series := seriesData{
		Puntos:    puntos,
		Mensajes:  d.Mensajes,
		Velocidad: d.Velocidad,
	}
	raw, err := json.Marshal(series)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// El tono de la insignia sale del enum del estado (RN-32.1, ADR 0003), no del texto que trae el dashboard.
	p, err := a.svc.Proyectos.Buscar(pid)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	datos := datosDashboard{
		Dashboard:  d,
		SeriesJSON: template.JS(raw),
		Tono:       tonoEstado(p.EstadoAl(hoy())),
		Resumen: resumenDashboard{
			Velocidad:      unDecimal(d.Velocidad),
			Porcentaje:     unDecimal(d.PorcentajeCompletadas),
			HorasEstimadas: unDecimal(d.HorasEstimadas),
			HorasReales:    unDecimal(d.HorasReales),
		},
	}
	if d.DesviacionCalculable {
		datos.Resumen.Desviacion = unDecimal(d.Desviacion)
	}

	a.render(w, r, "dashboard.html", datos)
}

// datosDashboard es lo que recibe dashboard.html: el dashboard, sus series para los gráficos, el tono de la insignia
// del estado y los números del resumen ya formateados con coma decimal (HU-32.4).
type datosDashboard struct {
	dashboard.Dashboard
	SeriesJSON template.JS
	Tono       string
	Resumen    resumenDashboard
}

// resumenDashboard tiene las métricas con decimales del resumen, formateadas con unDecimal (RN-32.3). Desviacion
// queda vacía cuando no se puede calcular.
type resumenDashboard struct {
	Velocidad      string
	Porcentaje     string
	HorasEstimadas string
	HorasReales    string
	Desviacion     string
}

func (a *app) renderDashboardSelector(w http.ResponseWriter, r *http.Request) {
	lista, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	a.render(w, r, "dashboard_seleccionar.html", lista)
}
