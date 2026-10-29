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
	// Buscar devuelve una copia de la historia, o backlog.ErrNoExiste.
	Buscar(id domain.HistoriaID) (backlog.Historia, error)
	// Guardar reemplaza una historia ya guardada, o devuelve backlog.ErrNoExiste.
	Guardar(h backlog.Historia) error
	// Estimar asigna Story Points a una historia (HU-15, HU-18).
	Estimar(id domain.HistoriaID, sp int) error
}

// CompletadorHistorias es lo que el servicio de backlog necesita de HU-12: registrar el paso de una historia a
// Completada también en el sprint. El servicio de HU-08 delega acá cuando el estado nuevo es Completada
// (spec HU-08, restricciones).
type CompletadorHistorias interface {
	CompletarHistoria(sprintID domain.SprintID, historiaID domain.HistoriaID) error
}

// Backlog son los casos de uso del Product Backlog (épica E2).
type Backlog struct {
	historias RepoHistorias
	proyectos RepoProyectos
	completar CompletadorHistorias // HU-12: registra el paso a Completada en el sprint
}

// NuevoBacklog arma el servicio del Product Backlog. El completador es el servicio de sprints (HU-12).
func NuevoBacklog(historias RepoHistorias, proyectos RepoProyectos, completar CompletadorHistorias) *Backlog {
	return &Backlog{historias: historias, proyectos: proyectos, completar: completar}
}

// Crear agrega una historia al Product Backlog de un proyecto (HU-05). Primero verifica que el proyecto exista;
// después el dominio valida título y prioridad (RN-05.2, RN-05.3) y el repositorio genera el ID (RN-05.1).
// La historia nueva queda en estado Pendiente, sin Story Points, sin criterios y sin sprint (RN-05.4).
func (s *Backlog) Crear(proyectoID domain.ProyectoID, titulo, descripcion string, prioridad backlog.Prioridad) (backlog.Historia, error) {
	if _, err := s.proyectos.Buscar(proyectoID); err != nil {
		return backlog.Historia{}, err
	}
	h, err := backlog.Nueva(proyectoID, titulo, descripcion, prioridad)
	if err != nil {
		return backlog.Historia{}, err
	}
	return s.historias.Agregar(h)
}

// Listar devuelve todas las historias del backlog, en el orden en que se crearon.
func (s *Backlog) Listar() ([]backlog.Historia, error) {
	return s.historias.Listar(), nil
}

// Buscar devuelve la historia con ese ID, o backlog.ErrNoExiste ("la historia no existe").
func (s *Backlog) Buscar(id domain.HistoriaID) (backlog.Historia, error) {
	return s.historias.Buscar(id)
}

// CambiarEstado deja la historia en el estado nuevo (HU-08). Para Completada delega en HU-12 antes de tocar
// nada: así las reglas del sprint hablan primero (RN-12.1 a RN-12.3, convención de HU-11) y el sprint también
// registra la historia completada. Si la historia no tiene sprint, devuelve ErrSinSprint de HU-08: no hay
// sprint al cual delegar.
func (s *Backlog) CambiarEstado(id domain.HistoriaID, nuevo backlog.Estado) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if nuevo == backlog.Completada {
		if h.SprintID == "" {
			return backlog.ErrSinSprint
		}
		return s.completar.CompletarHistoria(h.SprintID, id)
	}
	if err := h.CambiarEstado(nuevo); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}

// AgregarCriterio agrega un criterio de aceptación a la historia (HU-07). Busca la historia, aplica la regla del
// dominio (RN-07.1) y guarda el resultado (ADR 0002, sección 4).
func (s *Backlog) AgregarCriterio(id domain.HistoriaID, texto string) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.AgregarCriterio(texto); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}

// QuitarCriterio quita de la historia el criterio en la posición numero (empezando en 1) y guarda el resultado
// (HU-07, RN-07.3). Si la historia o el criterio no existen, devuelve el error del dominio o del repositorio.
func (s *Backlog) QuitarCriterio(id domain.HistoriaID, numero int) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.QuitarCriterio(numero); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}

// Estimar asigna Story Points a la historia (HU-15). Busca la historia, aplica la regla del dominio
// (RN-15.1 a RN-15.3) y guarda el resultado.
func (s *Backlog) Estimar(id domain.HistoriaID, sp int) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.Estimar(sp); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}
