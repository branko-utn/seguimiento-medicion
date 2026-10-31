// Package app arma la aplicación: crea los repositorios y los servicios y los conecta entre sí (TEC-18).
package app

import (
	"fmt"
	"path/filepath"

	"github.com/branko-utn/seguimiento-medicion/internal/server"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// NuevosServicios crea todos los repositorios y servicios. Cada colección se guarda en su archivo JSON dentro de dir
// (por ejemplo data/proyectos.json); con dir vacío todo queda solo en memoria, útil en los tests.
func NuevosServicios(dir string) (server.Servicios, error) {
	ruta := func(nombre string) string {
		if dir == "" {
			return ""
		}
		return filepath.Join(dir, nombre+".json")
	}

	proyectos, err := memory.NuevoProyectos(ruta("proyectos"))
	if err != nil {
		return server.Servicios{}, abrir("proyectos", err)
	}
	integrantes, err := memory.NuevoIntegrantes(ruta("integrantes"))
	if err != nil {
		return server.Servicios{}, abrir("integrantes", err)
	}
	historias, err := memory.NuevoHistorias(ruta("historias"))
	if err != nil {
		return server.Servicios{}, abrir("historias", err)
	}
	sprints, err := memory.NuevoSprints(ruta("sprints"))
	if err != nil {
		return server.Servicios{}, abrir("sprints", err)
	}
	estimaciones, err := memory.NuevoEstimaciones(ruta("estimaciones"))
	if err != nil {
		return server.Servicios{}, abrir("estimaciones", err)
	}
	registros, err := memory.NuevoRegistros(ruta("registros"))
	if err != nil {
		return server.Servicios{}, abrir("registros", err)
	}
	defectos, err := memory.NuevoDefectos(ruta("defectos"))
	if err != nil {
		return server.Servicios{}, abrir("defectos", err)
	}
	sesiones, err := memory.NuevoSesiones(ruta("sesiones"))
	if err != nil {
		return server.Servicios{}, abrir("sesiones", err)
	}

	// Backlog delega en Sprints para completar historias (HU-08 → HU-12).
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	// Backlog consulta registros y defectos para no eliminar historias con datos (HU-06, RN-06.5).
	registrosSvc := service.NuevoRegistros(registros, integrantes, historias, proyectos)
	defectosSvc := service.NuevoDefectos(defectos, historias, sprints)
	return server.Servicios{
		Proyectos:    service.NuevoProyectos(proyectos).ConResumen(sprints, historias),
		Integrantes:  service.NuevoIntegrantes(integrantes, proyectos),
		Backlog:      service.NuevoBacklog(historias, proyectos, sprintsSvc, registrosSvc, defectosSvc),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(estimaciones, historias),
		Registros:    registrosSvc,
		Defectos:     defectosSvc,
		Poker:        service.NuevoPlanningPoker(sesiones, historias, integrantes),
		Reportes:     service.NuevoReportes(proyectos, sprints, historias, estimaciones, registros, integrantes, defectos),
		Dashboard:    service.NuevoDashboard(proyectos, sprints, historias, estimaciones, registros, defectos),
	}, nil
}

func abrir(coleccion string, err error) error {
	return fmt.Errorf("abrir los datos de %s: %w", coleccion, err)
}
