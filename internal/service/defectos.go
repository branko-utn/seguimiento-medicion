package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// RepoDefectos es lo que el servicio de defectos necesita del almacenamiento.
type RepoDefectos interface {
	// Agregar guarda un defecto nuevo y lo devuelve con el ID generado (D1, D2, …).
	Agregar(d defect.Defecto) (defect.Defecto, error)
	// Buscar devuelve (defecto, encontrado).
	Buscar(id domain.DefectoID) (defect.Defecto, bool)
}

// BuscadorSprints es lo que el servicio de defectos necesita del repositorio de sprints (épica E3).
type BuscadorSprints interface {
	// Buscar devuelve el sprint, o sprint.ErrNoExiste.
	Buscar(id domain.SprintID) (sprint.Sprint, error)
}

// Defectos son los casos de uso de los defectos (épica E6).
type Defectos struct {
	defectos  RepoDefectos
	historias RepoHistorias
	sprints   BuscadorSprints
}

// NuevoDefectos crea el servicio. Usa las historias (épica E2) y los sprints (épica E3) para validar el defecto.
func NuevoDefectos(defectos RepoDefectos, historias RepoHistorias, sprints BuscadorSprints) *Defectos {
	return &Defectos{defectos: defectos, historias: historias, sprints: sprints}
}

// Registrar guarda un defecto nuevo (HU-22). El dominio valida descripción y severidad y lo deja Abierto (RN-22.1 a
// RN-22.3); el servicio verifica que la historia y el sprint existan y sean del mismo proyecto (RN-22.4). Se acepta
// una historia completada y un sprint cerrado: es cuando suelen detectarse los defectos.
func (s *Defectos) Registrar(descripcion string, severidad defect.Severidad, historiaID domain.HistoriaID, sprintDeteccion domain.SprintID) (defect.Defecto, error) {
	d, err := defect.Nuevo(descripcion, severidad, historiaID, sprintDeteccion)
	if err != nil {
		return defect.Defecto{}, err
	}
	historia, err := s.historias.Buscar(historiaID)
	if err != nil {
		return defect.Defecto{}, err
	}
	sp, err := s.sprints.Buscar(sprintDeteccion)
	if err != nil {
		return defect.Defecto{}, err
	}
	if sp.ProyectoID != historia.ProyectoID {
		return defect.Defecto{}, defect.ErrSprintDeOtroProyecto
	}
	return s.defectos.Agregar(d)
}

// Buscar devuelve el defecto y si se encontró.
func (s *Defectos) Buscar(id domain.DefectoID) (defect.Defecto, bool) {
	return s.defectos.Buscar(id)
}
