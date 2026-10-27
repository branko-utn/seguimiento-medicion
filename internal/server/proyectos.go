package server

import "net/http"

// rutasProyectos registra las rutas de la sección Proyectos (épica E1, TEC-13). Todavía sin implementar (RED).
func (a *app) rutasProyectos(mux *http.ServeMux) {
	mux.HandleFunc("GET /proyectos", a.listarProyectos)
	mux.HandleFunc("GET /proyectos/nuevo", a.nuevoProyecto)
	mux.HandleFunc("POST /proyectos", a.crearProyecto)
}

func (a *app) listarProyectos(w http.ResponseWriter, r *http.Request) {}

func (a *app) nuevoProyecto(w http.ResponseWriter, r *http.Request) {}

func (a *app) crearProyecto(w http.ResponseWriter, r *http.Request) {}
