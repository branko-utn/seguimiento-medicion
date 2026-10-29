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
	// Eliminar elimina la votación de una historia (HU-18: al acordar la estimación, la votación se cierra).
	Eliminar(id domain.HistoriaID) error
}

// RepoIntegrantesPoker es lo que el servicio de Planning Poker necesita de los integrantes: saber si el
// que vota existe (y de qué proyecto es). No reutiliza RepoIntegrantes porque ese no expone Buscar.
type RepoIntegrantesPoker interface {
	// Buscar devuelve el integrante, o project.ErrIntegranteNoExiste.
	Buscar(id domain.IntegranteID) (project.Integrante, error)
}

// PlanningPoker son los casos de uso de las votaciones con cartas ocultas (épica E4, HU-16, HU-17, HU-18).
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

// Revelar cierra la ronda mostrando todos los votos y las diferencias (HU-17). El estado revelado queda
// guardado: la ronda no acepta votos nuevos (RN-17.4).
func (s *PlanningPoker) Revelar(historiaID domain.HistoriaID) (estimation.Resultado, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return estimation.Resultado{}, estimation.ErrSinVotacionAbierta
	}
	res, err := sesion.Revelar()
	if err != nil {
		return estimation.Resultado{}, err
	}
	// La copia revelada reemplaza la guardada (buscar → mutar → guardar, como Votar).
	if err := s.sesiones.Guardar(sesion); err != nil {
		return estimation.Resultado{}, err
	}
	return res, nil
}

// Resultado devuelve el resultado de la ronda revelada, sin repetir la revelada: lo consulta la pantalla
// al recargar. Mientras la votación está abierta, los valores no se pueden ver (RN-16.4).
func (s *PlanningPoker) Resultado(historiaID domain.HistoriaID) (estimation.Resultado, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return estimation.Resultado{}, estimation.ErrSinVotacionAbierta
	}
	return sesion.Resultado()
}

// NuevaRonda abre una nueva ronda de votación (HU-18): la historia debe tener una votación abierta
// con la ronda actual revelada. Archiva la ronda actual y abre la siguiente sin votos.
func (s *PlanningPoker) NuevaRonda(historiaID domain.HistoriaID) (estimation.Sesion, error) {
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return estimation.Sesion{}, err
	}
	if h.Estado == backlog.Completada {
		return estimation.Sesion{}, estimation.ErrHistoriaCompletada
	}
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return estimation.Sesion{}, estimation.ErrSinVotacionAbierta
	}
	if !sesion.Revelada() {
		return estimation.Sesion{}, estimation.ErrRondaSinRevelar
	}
	// Archivar la ronda actual y abrir la siguiente (el dominio lo hace en NuevaRonda)
	if err := sesion.NuevaRonda(); err != nil {
		return estimation.Sesion{}, err
	}
	if err := s.sesiones.Guardar(sesion); err != nil {
		return estimation.Sesion{}, err
	}
	return sesion, nil
}

// Acordar registra la estimación acordada y cierra la votación (HU-18). La historia debe tener una
// votación abierta con la ronda actual revelada, y el valor debe ser de la escala.
func (s *PlanningPoker) Acordar(historiaID domain.HistoriaID, sp int) error {
	h, err := s.historias.Buscar(historiaID)
	if err != nil {
		return err
	}
	if h.Estado == backlog.Completada {
		return backlog.ErrHistoriaCompletada
	}
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return estimation.ErrSinVotacionAbierta
	}
	if !sesion.Revelada() {
		return estimation.ErrRondaSinRevelar
	}
	if err := estimation.ValidarStoryPoints(sp); err != nil {
		return err
	}
	// Acordar en el dominio (valida revelada + escala)
	if err := sesion.Acordar(sp); err != nil {
		return err
	}
	// Estimar la historia en el backlog (HU-15): usa el dominio directamente
	if err := h.Estimar(sp); err != nil {
		return err
	}
	if err := s.historias.Guardar(h); err != nil {
		return err
	}
	// Eliminar la votación (cerrada): el historial se ve solo mientras está abierta (spec HU-18)
	if err := s.sesiones.Eliminar(historiaID); err != nil {
		return err
	}
	return nil
}

// Rondas devuelve el historial de rondas de la votación (HU-18). La historia debe tener una votación abierta.
func (s *PlanningPoker) Rondas(historiaID domain.HistoriaID) ([]estimation.Ronda, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return nil, estimation.ErrSinVotacionAbierta
	}
	return sesion.Rondas(), nil
}

// RondaActual devuelve la ronda actual de la votación (HU-18, pantalla).
func (s *PlanningPoker) RondaActual(historiaID domain.HistoriaID) (int, error) {
	sesion, ok := s.sesiones.Buscar(historiaID)
	if !ok {
		return 0, estimation.ErrSinVotacionAbierta
	}
	return sesion.Ronda, nil
}