package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// RepoDefectos es lo que el servicio de defectos necesita del almacenamiento.
type RepoDefectos interface {
	Agregar(d defect.Defecto) (defect.Defecto, error)
	Buscar(id domain.DefectoID) (defect.Defecto, bool)
}

// BuscadorSprints es lo que el servicio de defectos necesita del repositorio de sprints.
type BuscadorSprints interface {
	Buscar(id domain.SprintID) (sprint.Sprint, error)
}

// Defectos son los casos de uso de los defectos (épica E6).
type Defectos struct {
	defectos  RepoDefectos
	historias RepoHistorias
	sprints   BuscadorSprints
}

func NuevoDefectos(defectos RepoDefectos, historias RepoHistorias, sprints BuscadorSprints) *Defectos {
	return &Defectos{defectos: defectos, historias: historias, sprints: sprints}
}

func (s *Defectos) Registrar(descripcion string, severidad defect.Severidad, historiaID domain.HistoriaID, sprintDeteccion domain.SprintID) (defect.Defecto, error) {
	return defect.Defecto{}, nil
}

func (s *Defectos) Buscar(id domain.DefectoID) (defect.Defecto, bool) {
	return defect.Defecto{}, false
}
