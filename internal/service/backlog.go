package service

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
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
	// Eliminar saca la historia del backlog sin liberar su ID (RN-06.4), o devuelve backlog.ErrNoExiste.
	Eliminar(id domain.HistoriaID) error
}

// CompletadorHistorias es lo que el servicio de backlog necesita de HU-12: registrar el paso de una historia a
// Completada también en el sprint. El servicio de HU-08 delega acá cuando el estado nuevo es Completada
// (spec HU-08, restricciones).
type CompletadorHistorias interface {
	CompletarHistoria(sprintID domain.SprintID, historiaID domain.HistoriaID) error
}

// ErrHistoriaConRegistros: RN-06.5, una historia con esfuerzo o defectos registrados no se elimina (HU-06).
var ErrHistoriaConRegistros = errors.New("no se puede eliminar una historia con esfuerzo o defectos registrados")

// EsfuerzoDeHistorias es lo que el servicio de backlog necesita de HU-20 para validar RN-06.5. Lo cumple
// service.Registros; así backlog no depende del paquete de esfuerzo (ADR 0002).
type EsfuerzoDeHistorias interface {
	DeHistoria(id domain.HistoriaID) []effort.Registro
}

// ListadorDefectos es lo que el servicio de backlog necesita de HU-22 para validar RN-06.5. Lo cumple
// service.Defectos.
type ListadorDefectos interface {
	Listar() []defect.Defecto
}

// Backlog son los casos de uso del Product Backlog (épica E2).
type Backlog struct {
	historias RepoHistorias
	proyectos RepoProyectos
	completar CompletadorHistorias // HU-12: registra el paso a Completada en el sprint
	esfuerzo  EsfuerzoDeHistorias  // HU-20: para no eliminar historias con horas registradas (RN-06.5)
	defectos  ListadorDefectos     // HU-22: para no eliminar historias con defectos (RN-06.5)
}

// NuevoBacklog arma el servicio del Product Backlog. El completador es el servicio de sprints (HU-12); esfuerzo y
// defectos son los servicios de registros (HU-20) y de defectos (HU-22), que usa Eliminar (HU-06).
func NuevoBacklog(historias RepoHistorias, proyectos RepoProyectos, completar CompletadorHistorias, esfuerzo EsfuerzoDeHistorias, defectos ListadorDefectos) *Backlog {
	return &Backlog{historias: historias, proyectos: proyectos, completar: completar, esfuerzo: esfuerzo, defectos: defectos}
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

// DelProyecto devuelve el backlog de un proyecto ordenado por prioridad e identificador (HU-09).
func (s *Backlog) DelProyecto(proyectoID domain.ProyectoID, estado *backlog.Estado) ([]backlog.Historia, error) {
	return nil, nil
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

// Editar cambia el título, la descripción y la prioridad de una historia (HU-06). Busca la historia, aplica las
// reglas del dominio (RN-06.1 y RN-06.2) y guarda el resultado; si algo falla, no se guarda nada.
func (s *Backlog) Editar(id domain.HistoriaID, titulo, descripcion string, prioridad backlog.Prioridad) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.Editar(titulo, descripcion, prioridad); err != nil {
		return err
	}
	return s.historias.Guardar(h)
}

// Eliminar saca una historia del backlog (HU-06). Valida en orden: que exista, las reglas del dominio (RN-06.3) y
// que no tenga esfuerzo ni defectos registrados (RN-06.5). El repositorio no libera su ID (RN-06.4).
func (s *Backlog) Eliminar(id domain.HistoriaID) error {
	h, err := s.historias.Buscar(id)
	if err != nil {
		return err
	}
	if err := h.PuedeEliminarse(); err != nil {
		return err
	}
	if s.tieneRegistros(id) {
		return ErrHistoriaConRegistros
	}
	return s.historias.Eliminar(id)
}

// tieneRegistros informa si la historia tiene horas registradas (HU-20) o algún defecto (HU-22), para RN-06.5.
func (s *Backlog) tieneRegistros(id domain.HistoriaID) bool {
	if len(s.esfuerzo.DeHistoria(id)) > 0 {
		return true
	}
	for _, d := range s.defectos.Listar() {
		if d.HistoriaID == id {
			return true
		}
	}
	return false
}
