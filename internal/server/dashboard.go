package server

import (
	"errors"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

func (a *app) rutasDashboard(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", a.dashboard)
}

func (a *app) dashboard(w http.ResponseWriter, r *http.Request) {
	// FASE RED: sin implementación completa. Solo para compilar.
	proyecto := r.URL.Query().Get("proyecto")
	if proyecto == "" {
		a.renderDashboardSelector(w, r)
		return
	}
	pid := r.URL.Query().Get("proyecto")
	_, err := a.svc.Dashboard.Del(domain.ProyectoID(pid), hoy())
	if err != nil {
		if errors.Is(err, project.ErrNoExiste) {
			http.Error(w, "el proyecto no existe", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Para compilar y que los tests RED fallen en assertions, devolvemos 500 o página mínima
	http.Error(w, "no implementado", http.StatusInternalServerError)
}

func (a *app) renderDashboardSelector(w http.ResponseWriter, r *http.Request) {
	lista, err := a.svc.Proyectos.Listar()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.render(w, r, "dashboard_seleccionar.html", lista)
}
