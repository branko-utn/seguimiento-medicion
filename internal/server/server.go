// Package server expone la aplicación por HTTP: rutas, handlers y pantallas.
package server

import "net/http"

// NewRouter arma el enrutador con todas las rutas de la aplicación.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	return mux
}
