package service

import (
	"strings"

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

// Registrar agrega un integrante a un proyecto. Primero verifica que el proyecto exista; después el dominio valida
// nombre y rol (RN-03.1, RN-03.2) y el servicio controla que no esté repetido en ese proyecto (RN-03.3).
func (s *Integrantes) Registrar(proyectoID domain.ProyectoID, nombre string, rol project.Rol) (project.Integrante, error) {
	if _, err := s.proyectos.Buscar(proyectoID); err != nil {
		return project.Integrante{}, err
	}
	i, err := project.NuevoIntegrante(proyectoID, nombre, rol)
	if err != nil {
		return project.Integrante{}, err
	}
	for _, otro := range s.integrantes.DelProyecto(proyectoID) {
		if strings.EqualFold(otro.Nombre, i.Nombre) {
			return project.Integrante{}, project.ErrIntegranteDuplicado
		}
	}
	return s.integrantes.Agregar(i)
}

// DelProyecto devuelve los integrantes de un proyecto en el orden en que se registraron, o project.ErrNoExiste.
func (s *Integrantes) DelProyecto(proyectoID domain.ProyectoID) ([]project.Integrante, error) {
	if _, err := s.proyectos.Buscar(proyectoID); err != nil {
		return nil, err
	}
	return s.integrantes.DelProyecto(proyectoID), nil
}
