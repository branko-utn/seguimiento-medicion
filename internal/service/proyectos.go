package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoProyectos es lo que el servicio de proyectos necesita del almacenamiento.
type RepoProyectos interface {
	// Agregar guarda un proyecto nuevo y lo devuelve con el ID generado.
	Agregar(p project.Proyecto) (project.Proyecto, error)
	Listar() []project.Proyecto
}

// Proyectos son los casos de uso de la épica E1.
type Proyectos struct {
	repo RepoProyectos
}

func NuevoProyectos(repo RepoProyectos) *Proyectos {
	return &Proyectos{repo: repo}
}

// Crear da de alta un proyecto (HU-01). El stub existe solo para que los tests de la fase RED compilen.
func (s *Proyectos) Crear(nombre, descripcion string, inicio, fin time.Time) (project.Proyecto, error) {
	return project.Proyecto{}, nil
}

// Listar devuelve todos los proyectos.
func (s *Proyectos) Listar() ([]project.Proyecto, error) {
	return nil, nil
}
