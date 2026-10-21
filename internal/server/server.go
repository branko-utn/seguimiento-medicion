// Package server expone la aplicación por HTTP: rutas, handlers y pantallas.
package server

import "net/http"

// Servicios agrupa los casos de uso que usan los handlers. Cada épica agrega acá el campo de su servicio.
type Servicios struct{}

// NewRouter arma el enrutador con todas las rutas de la aplicación.
func NewRouter(svc Servicios) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	return mux, nil
}
