// Package app arma la aplicación: crea los repositorios y los servicios y los conecta entre sí (TEC-18).
package app

import "github.com/branko-utn/seguimiento-medicion/internal/server"

// NuevosServicios crea todos los servicios con sus repositorios.
func NuevosServicios(dir string) (server.Servicios, error) { return server.Servicios{}, nil }
