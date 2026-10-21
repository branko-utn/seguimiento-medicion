package server

import "net/http"

func (a *app) inicio(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "inicio.html", nil)
}

// enConstruccion avisa en qué sprint llega una sección que todavía no está implementada.
func (a *app) enConstruccion(s seccion) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, r, "en_construccion.html", s)
	}
}
