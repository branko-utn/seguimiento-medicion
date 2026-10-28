package server

import (
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// rutasBacklog registra las rutas de la sección Backlog (épica E2, TEC-16).
func (a *app) rutasBacklog(mux *http.ServeMux) {
	mux.HandleFunc("GET /backlog", a.listarBacklog)
}

// datosBacklog es lo que recibe la plantilla del backlog.
type datosBacklog struct {
	Historias []backlog.Historia
	Proyectos map[string]string
}

func (a *app) listarBacklog(w http.ResponseWriter, r *http.Request) {
	var historias []backlog.Historia
	if a.svc.Backlog != nil {
		if h, err := a.svc.Backlog.Listar(); err == nil {
			historias = h
		} else {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
			return
		}
	}
	nombres := make(map[string]string)
	if a.svc.Proyectos != nil {
		if ps, err := a.svc.Proyectos.Listar(); err == nil {
			for _, p := range ps {
				nombres[string(p.ID)] = p.Nombre
			}
		}
	}
	a.render(w, r, "backlog.html", datosBacklog{Historias: historias, Proyectos: nombres})
}
