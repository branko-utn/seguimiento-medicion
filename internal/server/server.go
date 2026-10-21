// Package server expone la aplicación por HTTP: rutas, handlers y pantallas.
package server

import (
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/web"
)

// Servicios agrupa los casos de uso que usan los handlers. Cada épica agrega acá el campo de su servicio.
type Servicios struct{}

// NewRouter arma el enrutador con todas las rutas de la aplicación.
func NewRouter(svc Servicios) (http.Handler, error) {
	a, err := nuevaApp(svc)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.Handle("GET /static/", http.FileServerFS(web.FS))
	mux.HandleFunc("GET /{$}", a.inicio)
	for _, s := range secciones {
		if s.EnConstruccion {
			mux.HandleFunc("GET "+s.Ruta, a.enConstruccion(s))
		}
	}
	// Cada épica registra acá sus rutas, por ejemplo: a.rutasProyectos(mux)

	return mux, nil
}
