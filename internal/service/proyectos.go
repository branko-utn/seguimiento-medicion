package service

import (
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoProyectos es lo que el servicio de proyectos necesita del almacenamiento.
type RepoProyectos interface {
	// Agregar guarda un proyecto nuevo y lo devuelve con el ID generado.
	Agregar(p project.Proyecto) (project.Proyecto, error)
	Listar() []project.Proyecto
	// Buscar y Actualizar devuelven project.ErrNoExiste si no hay un proyecto con ese ID.
	Buscar(id domain.ProyectoID) (project.Proyecto, error)
	Actualizar(p project.Proyecto) error
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
	// Un proyecto nuevo todavía no tiene ID: con "" no se excluye a ninguno.
	if s.nombreEnUso(p.Nombre, "") {
		return project.Proyecto{}, project.ErrNombreDuplicado
	}
	return s.repo.Agregar(p)
}

// Modificar cambia nombre, descripción y fechas de un proyecto (HU-02). Busca el proyecto, el dominio valida que no
// esté finalizado y los datos (RN-02.1, RN-02.3), y el servicio controla que el nombre no sea el de otro proyecto
// sin distinguir mayúsculas (RN-02.2): el propio nombre se puede mantener.
func (s *Proyectos) Modificar(id domain.ProyectoID, nombre, descripcion string, inicio, fin, hoy time.Time) (project.Proyecto, error) {
	p, err := s.repo.Buscar(id)
	if err != nil {
		return project.Proyecto{}, err
	}
	if err := p.Modificar(nombre, descripcion, inicio, fin, hoy); err != nil {
		return project.Proyecto{}, err
	}
	if s.nombreEnUso(p.Nombre, p.ID) {
		return project.Proyecto{}, project.ErrNombreDuplicado
	}
	if err := s.repo.Actualizar(p); err != nil {
		return project.Proyecto{}, err
	}
	return p, nil
}

// nombreEnUso informa si otro proyecto (distinto de excepto) ya tiene ese nombre, sin distinguir mayúsculas.
func (s *Proyectos) nombreEnUso(nombre string, excepto domain.ProyectoID) bool {
	for _, otro := range s.repo.Listar() {
		if otro.ID != excepto && strings.EqualFold(otro.Nombre, nombre) {
			return true
		}
	}
	return false
}

// Listar devuelve todos los proyectos.
func (s *Proyectos) Listar() ([]project.Proyecto, error) {
	return s.repo.Listar(), nil
}

// Buscar devuelve un proyecto, o project.ErrNoExiste (TEC-13: detalle y edición).
func (s *Proyectos) Buscar(id domain.ProyectoID) (project.Proyecto, error) {
	return s.repo.Buscar(id)
}
