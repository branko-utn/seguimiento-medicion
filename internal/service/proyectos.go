package service

import (
	"strings"
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

// Crear da de alta un proyecto (HU-01). El dominio valida nombre y fechas; el servicio controla que el nombre
// no esté repetido sin distinguir mayúsculas (RN-01.3).
func (s *Proyectos) Crear(nombre, descripcion string, inicio, fin time.Time) (project.Proyecto, error) {
	p, err := project.Nuevo(nombre, descripcion, inicio, fin)
	if err != nil {
		return project.Proyecto{}, err
	}
	for _, otro := range s.repo.Listar() {
		if strings.EqualFold(otro.Nombre, p.Nombre) {
			return project.Proyecto{}, project.ErrNombreDuplicado
		}
	}
	return s.repo.Agregar(p)
}

// Listar devuelve todos los proyectos.
func (s *Proyectos) Listar() ([]project.Proyecto, error) {
	return s.repo.Listar(), nil
}
