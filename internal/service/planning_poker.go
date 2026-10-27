package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoSesiones es lo que el servicio de Planning Poker necesita del almacenamiento de votaciones.
type RepoSesiones interface {
	// Guardar crea o reemplaza la votación de la historia de la sesión (RN-16.5: una por historia).
	Guardar(s estimation.Sesion) error
	// Buscar devuelve la votación de una historia y si existe: que no haya votación no es un error,
	// sino un estado válido (como las estimaciones de horas).
	Buscar(id domain.HistoriaID) (estimation.Sesion, bool)
}

// RepoIntegrantesPoker es lo que el servicio de Planning Poker necesita de los integrantes: saber si el
// que vota existe (y de qué proyecto es). No reutiliza RepoIntegrantes porque ese no expone Buscar.
type RepoIntegrantesPoker interface {
	// Buscar devuelve el integrante, o project.ErrIntegranteNoExiste.
	Buscar(id domain.IntegranteID) (project.Integrante, error)
}

// PlanningPoker son los casos de uso de las votaciones con cartas ocultas (épica E4, HU-16).
type PlanningPoker struct {
	sesiones    RepoSesiones
	historias   RepoHistorias
	integrantes RepoIntegrantesPoker
}

// NuevoPlanningPoker arma el servicio. Usa los repositorios de historias (E2) e integrantes (E1) para
// validar que la historia exista y que quien vota sea del proyecto de la historia (RN-16.1).
func NuevoPlanningPoker(sesiones RepoSesiones, historias RepoHistorias, integrantes RepoIntegrantesPoker) *PlanningPoker {
	return &PlanningPoker{sesiones: sesiones, historias: historias, integrantes: integrantes}
}

// IniciarVotacion abre una votación para la historia (HU-16): la historia tiene que existir, no estar
// completada ni tener otra votación abierta (RN-16.5). La sesión queda abierta en la ronda 1, sin votos.
func (s *PlanningPoker) IniciarVotacion(historiaID domain.HistoriaID) (estimation.Sesion, error) {
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return estimation.Sesion{}, err
	}
	if h.Estado == backlog.Completada {
		return estimation.Sesion{}, estimation.ErrHistoriaCompletada
	}
	if _, ok := s.sesiones.Buscar(historiaID); ok {
		return estimation.Sesion{}, estimation.ErrVotacionAbierta
	}
	sesion, err := estimation.NuevaSesion(historiaID)
	if err != nil {
		return estimation.Sesion{}, err
	}
	if err := s.sesiones.Guardar(sesion); err != nil {
		return estimation.Sesion{}, err
	}
	return sesion, nil
}

// Votar registra el voto del integrante en la votación abierta de la historia (HU-16). Valida que la
// historia y el integrante existan, que el integrante sea del proyecto de la historia (RN-16.1) y las
// reglas de la sesión: carta de la escala (RN-16.3) y un voto por integrante y ronda (RN-16.2).
// El error de una votación que no existe es el acordado por el equipo (no está en la tabla de la spec).
func (s *PlanningPoker) Votar(historiaID domain.HistoriaID, integranteID domain.IntegranteID, carta int) error {
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return err
	}
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return estimation.ErrSinVotacionAbierta
	}
	i, err := s.integrantes.Buscar(integranteID)
	if err != nil {
		return err
	}
	if i.ProyectoID != h.ProyectoID {
		return estimation.ErrNoEsDelProyecto
	}
	if err := sesion.Votar(integranteID, carta); err != nil {
		return err
	}
	return s.sesiones.Guardar(sesion)
}

// QuienesVotaron devuelve los integrantes que ya votaron en la votación abierta de la historia, sin los
// valores de sus votos (RN-16.4: mientras está abierta, solo se consulta quién votó).
func (s *PlanningPoker) QuienesVotaron(historiaID domain.HistoriaID) ([]domain.IntegranteID, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return nil, estimation.ErrSinVotacionAbierta
	}
	return sesion.QuienesVotaron(), nil
}

// ConsultarVotos devuelve los votos de la ronda. Mientras la votación está abierta no se pueden ver
// (RN-16.4, estimation.ErrVotosOcultos); HU-17 los expone al revelar.
func (s *PlanningPoker) ConsultarVotos(historiaID domain.HistoriaID) ([]estimation.Voto, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return nil, estimation.ErrSinVotacionAbierta
	}
	return sesion.Votos()
}
