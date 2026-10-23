package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoIntegrantes es lo que el servicio de integrantes necesita del almacenamiento.
type RepoIntegrantes interface {
	// Agregar guarda un integrante nuevo y lo devuelve con el ID generado.
	Agregar(i project.Integrante) (project.Integrante, error)
	// DelProyecto devuelve los integrantes de un proyecto en el orden en que se registraron.
	DelProyecto(id domain.ProyectoID) []project.Integrante
}

// Integrantes son los casos de uso de los integrantes de un proyecto (HU-03).
type Integrantes struct {
	integrantes RepoIntegrantes
	proyectos   RepoProyectos
}

func NuevoIntegrantes(integrantes RepoIntegrantes, proyectos RepoProyectos) *Integrantes {
	return &Integrantes{integrantes: integrantes, proyectos: proyectos}
}

// Registrar agrega un integrante a un proyecto. El stub existe solo para que los tests de la fase RED compilen.
func (s *Integrantes) Registrar(proyectoID domain.ProyectoID, nombre string, rol project.Rol) (project.Integrante, error) {
	return project.Integrante{}, nil
}

// DelProyecto devuelve los integrantes de un proyecto. El stub existe solo para la fase RED.
func (s *Integrantes) DelProyecto(proyectoID domain.ProyectoID) ([]project.Integrante, error) {
	return nil, nil
}
