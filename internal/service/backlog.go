package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// RepoHistorias es lo que el servicio de backlog necesita del almacenamiento.
type RepoHistorias interface {
	// Agregar guarda una historia nueva y la devuelve con el ID generado (RN-05.1).
	Agregar(h backlog.Historia) (backlog.Historia, error)
	// Listar devuelve todas las historias, en el orden en que se crearon.
	Listar() []backlog.Historia
}

// Backlog son los casos de uso del Product Backlog (épica E2).
type Backlog struct {
	historias RepoHistorias
	proyectos RepoProyectos
}

func NuevoBacklog(historias RepoHistorias, proyectos RepoProyectos) *Backlog {
	return &Backlog{historias: historias, proyectos: proyectos}
}

// Crear agrega una historia al Product Backlog de un proyecto (HU-05).
func (s *Backlog) Crear(proyectoID domain.ProyectoID, titulo, descripcion string, prioridad backlog.Prioridad) (backlog.Historia, error) {
	return backlog.Historia{}, nil // stub RED
}

// Listar devuelve todas las historias del backlog.
func (s *Backlog) Listar() ([]backlog.Historia, error) {
	return nil, nil // stub RED
}
