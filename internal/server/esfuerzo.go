package server

import (
	"net/http"
	"time"
)

var hoy = func() time.Time { return time.Time{} }

func (a *app) rutasEsfuerzo(mux *http.ServeMux) {
}

func (a *app) verEsfuerzo(w http.ResponseWriter, r *http.Request) {
}

func (a *app) registrarEsfuerzo(w http.ResponseWriter, r *http.Request) {
}

func (a *app) estimarHistoria(w http.ResponseWriter, r *http.Request) {
}
