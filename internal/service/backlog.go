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
}

// Backlog son los casos de uso del Product Backlog (épica E2).
type Backlog struct {
	historias RepoHistorias
	proyectos RepoProyectos
}

// NuevoBacklog arma el servicio del Product Backlog.
func NuevoBacklog(historias RepoHistorias, proyectos RepoProyectos) *Backlog {
	return &Backlog{historias: historias, proyectos: proyectos}
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
