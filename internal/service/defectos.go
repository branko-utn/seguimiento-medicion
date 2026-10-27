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
	// Actualizar reemplaza un defecto ya guardado, o devuelve defect.ErrNoExiste.
	Actualizar(d defect.Defecto) error
	// Listar devuelve todos los defectos, en el orden en que se registraron.
	Listar() []defect.Defecto
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

// Buscar devuelve el defecto, o defect.ErrNoExiste (HU-23).
func (s *Defectos) Buscar(id domain.DefectoID) (defect.Defecto, error) {
	d, ok := s.defectos.Buscar(id)
	if !ok {
		return defect.Defecto{}, defect.ErrNoExiste
	}
	return d, nil
}

// CambiarEstado pasa el defecto a Abierto o En progreso (HU-23). Las transiciones las valida el dominio (RN-23.1 y
// RN-23.2); para resolverlo se usa Resolver.
func (s *Defectos) CambiarEstado(id domain.DefectoID, nuevo defect.Estado) error {
	d, err := s.Buscar(id)
	if err != nil {
		return err
	}
	if err := d.CambiarEstado(nuevo); err != nil {
		return err
	}
	return s.defectos.Actualizar(d)
}

// Resolver deja el defecto Resuelto en el sprint indicado (HU-23). El dominio valida que no esté resuelto y que haya
// sprint (RN-23.1 a RN-23.3); el servicio, con los dos sprints, que el de resolución exista, sea del mismo proyecto
// que el de detección (RN-23.3) y no empiece antes que él (RN-23.4). Puede ser el mismo sprint.
func (s *Defectos) Resolver(id domain.DefectoID, sprintID domain.SprintID) error {
	d, err := s.Buscar(id)
	if err != nil {
		return err
	}
	if err := d.Resolver(sprintID); err != nil {
		return err
	}
	resolucion, err := s.sprints.Buscar(sprintID)
	if err != nil {
		return err
	}
	deteccion, err := s.sprints.Buscar(d.SprintDeteccion)
	if err != nil {
		return err
	}
	if resolucion.ProyectoID != deteccion.ProyectoID {
		return defect.ErrResolucionDeOtroProyecto
	}
	if resolucion.Inicio.Before(deteccion.Inicio) {
		return defect.ErrResolucionAnterior
	}
	return s.defectos.Actualizar(d)
}

// Listar devuelve todos los defectos, en el orden en que se registraron.
func (s *Defectos) Listar() []defect.Defecto {
	return s.defectos.Listar()
}
