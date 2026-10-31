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

	datos := struct {
		dashboard.Dashboard
		SeriesJSON template.JS
	}{
		Dashboard:  d,
		SeriesJSON: template.JS(raw),
	}

	a.render(w, r, "dashboard.html", datos)
}

func (a *app) renderDashboardSelector(w http.ResponseWriter, r *http.Request) {
	lista, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	a.render(w, r, "dashboard_seleccionar.html", lista)
}
