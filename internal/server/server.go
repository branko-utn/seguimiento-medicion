// Package server expone la aplicación por HTTP: rutas, handlers y pantallas.
package server

import (
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/web"
)

// Servicios agrupa los casos de uso que usan los handlers. Los crea y conecta internal/app (TEC-18): cada pantalla
// solo registra sus rutas y usa el servicio que necesita.
type Servicios struct {
	Proyectos    *service.Proyectos     // E1: HU-01, HU-02
	Integrantes  *service.Integrantes   // E1: HU-03
	Backlog      *service.Backlog       // E2 y E4: HU-05, 07, 08, 15
	Sprints      *service.Sprints       // E3: HU-10, 11, 12, 13
	Estimaciones *service.Estimaciones  // E5: HU-19
	Registros    *service.Registros     // E5: HU-20
	Defectos     *service.Defectos      // E6: HU-22, 23
	Poker        *service.PlanningPoker // E4: HU-16, HU-17
}

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
	a.rutasProyectos(mux)
	a.rutasPlanningPoker(mux)
	a.rutasAccionesSprint(mux) // HU-12 y HU-13: completar historias y cerrar el sprint
	a.rutasDefectos(mux)       // HU-22 y HU-23: lista y alta de defectos, cambiar estado y resolver
	a.rutasEsfuerzo(mux)       // HU-19 y HU-20: horas estimadas y registro de esfuerzo

	return mux, nil
}
