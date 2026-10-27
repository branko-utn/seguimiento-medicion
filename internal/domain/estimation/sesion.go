package estimation

import (
	"encoding/json"
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
	// ErrSinVotacionAbierta: votar o consultar los votos de una historia sin una votación abierta. No está
	// en la tabla de errores de la spec HU-16: lo acordó el equipo en la planning del Sprint 2 y se
	// documenta como aclaración en el PR.
	ErrSinVotacionAbierta = errors.New("la historia no tiene una votación abierta")
)

// Sesion representa una votación de Planning Poker para una historia.
type Sesion struct {
	HistoriaID domain.HistoriaID
	// Ronda es la ronda actual de la votación; una votación nueva abre en la ronda 1 (spec HU-16).
	Ronda int
	// votos son los votos de la ronda actual, en el orden en que se emitieron. No está exportado: mientras la
	// votación está abierta, la sesión no expone los valores de los votos (RN-16.4), solo quién votó.
	// Se persisten igual en el JSON con MarshalJSON/UnmarshalJSON: los votos se ocultan en la interfaz,
	// no en el archivo privado de data/.
	votos []Voto
}

// Voto representa el voto de un integrante en la ronda.
type Voto struct {
	IntegranteID domain.IntegranteID
	Carta        int
}

// NuevaSesion crea una votación abierta en la ronda 1 y sin votos (spec HU-16, salidas esperadas).
// Devuelve la sesión por valor, como el resto del dominio (ADR 0002): el servicio la busca, la muta
// con sus métodos y guarda la copia resultante.
func NuevaSesion(historiaID domain.HistoriaID) (Sesion, error) {
	return Sesion{
		HistoriaID: historiaID,
		Ronda:      1,
	}, nil
}

// Votar registra el voto de un integrante en la ronda actual. La carta tiene que ser de la escala
// (RN-16.3) y cada integrante vota una sola vez por ronda (RN-16.2).
func (s *Sesion) Votar(integranteID domain.IntegranteID, carta int) error {
	if s.yaVoto(integranteID) {
		return ErrYaVoto
	}
	if err := ValidarStoryPoints(carta); err != nil {
		return err
	}
	s.votos = append(s.votos, Voto{IntegranteID: integranteID, Carta: carta})
	return nil
}

// yaVoto informa si el integrante ya votó en la ronda actual (RN-16.2).
func (s *Sesion) yaVoto(integranteID domain.IntegranteID) bool {
	for _, v := range s.votos {
		if v.IntegranteID == integranteID {
			return true
		}
	}
	return false
}

// QuienesVotaron devuelve los integrantes que ya votaron, en el orden en que votaron, sin los valores
// de sus votos (RN-16.4: mientras la votación está abierta, solo se puede consultar quién votó).
func (s *Sesion) QuienesVotaron() []domain.IntegranteID {
	votantes := make([]domain.IntegranteID, 0, len(s.votos))
	for _, v := range s.votos {
		votantes = append(votantes, v.IntegranteID)
	}
	return votantes
}

// Votos devuelve los votos de la ronda. Mientras la votación está abierta, los valores no se pueden
// consultar (RN-16.4); HU-17 los expone al revelar.
func (s *Sesion) Votos() ([]Voto, error) {
	return nil, ErrVotosOcultos
}

// sesionJSON es la forma serializable de Sesion: la usa Store para guardar data/sesiones.json. No es
// parte de la API del dominio; la sesión se manipula con sus métodos.
type sesionJSON struct {
	HistoriaID domain.HistoriaID
	Ronda      int
	Votos      []Voto
}

// MarshalJSON persiste también los votos, que no están exportados: sin este método, el JSON guardado
// por Store perdería la votación completa.
func (s Sesion) MarshalJSON() ([]byte, error) {
	return json.Marshal(sesionJSON{HistoriaID: s.HistoriaID, Ronda: s.Ronda, Votos: s.votos})
}

// UnmarshalJSON reconstruye la sesión guardada, incluidos los votos.
func (s *Sesion) UnmarshalJSON(contenido []byte) error {
	var guardada sesionJSON
	if err := json.Unmarshal(contenido, &guardada); err != nil {
		return err
	}
	s.HistoriaID, s.Ronda, s.votos = guardada.HistoriaID, guardada.Ronda, guardada.Votos
	return nil
}
