package server

import "net/http"

// rutasSprints registra las rutas de la sección Sprints (épica E3, TEC-14). Todavía sin implementar (RED).
func (a *app) rutasSprints(mux *http.ServeMux) {
	mux.HandleFunc("GET /sprints", a.listarSprints)
	mux.HandleFunc("GET /sprints/nuevo", a.nuevoSprint)
	mux.HandleFunc("POST /sprints", a.crearSprint)
}

func (a *app) listarSprints(w http.ResponseWriter, r *http.Request) {}

func (a *app) nuevoSprint(w http.ResponseWriter, r *http.Request) {}

func (a *app) crearSprint(w http.ResponseWriter, r *http.Request) {}
