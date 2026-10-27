package estimation

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// Errores de Planning Poker (HU-16, HU-17).
var (
	ErrHistoriaCompletada = errors.New("no se puede votar una historia completada")
	ErrVotacionAbierta    = errors.New("la historia ya tiene una votación abierta")
	ErrNoEsDelProyecto    = errors.New("el integrante no pertenece al proyecto de la historia")
	ErrYaVoto             = errors.New("el integrante ya votó en esta ronda")
	ErrVotosOcultos       = errors.New("los votos se muestran recién al revelar las cartas")
	ErrYaRevelada         = errors.New("la ronda ya fue revelada")
	ErrSinVotos           = errors.New("no hay votos para revelar")
)

// Sesion representa una votación de Planning Poker para una historia.
type Sesion struct {
	HistoriaID domain.HistoriaID
	votos      []Voto
	votantes   map[string]bool
}

// Voto representa el voto de un integrante en la ronda.
type Voto struct {
	IntegranteID domain.IntegranteID
	Carta        int
}

// NuevaSesion crea una nueva sesión de votación para la historia.
func NuevaSesion(historiaID domain.HistoriaID) (*Sesion, error) {
	return &Sesion{
		HistoriaID: historiaID,
		votos:      []Voto{},
		votantes:   map[string]bool{},
	}, nil
}

// Votar registra el voto de un integrante en la ronda actual.
func (s *Sesion) Votar(integranteID domain.IntegranteID, carta int) error {
	if s.votantes[string(integranteID)] {
		return ErrYaVoto
	}
	if err := ValidarStoryPoints(carta); err != nil {
		return err
	}
	s.votantes[string(integranteID)] = true
	s.votos = append(s.votos, Voto{IntegranteID: integranteID, Carta: carta})
	return nil
}

// QuienesVotaron devuelve la lista de integrantes que ya votaron, sin exponer los valores.
func (s *Sesion) QuienesVotaron() []domain.IntegranteID {
	res := make([]domain.IntegranteID, 0, len(s.votos))
	for _, v := range s.votos {
		res = append(res, v.IntegranteID)
	}
	return res
}

// Votos devuelve los votos solo después de revelar; mientras está abierta devuelve ErrVotosOcultos.
func (s *Sesion) Votos() ([]Voto, error) {
	return nil, ErrVotosOcultos
}
