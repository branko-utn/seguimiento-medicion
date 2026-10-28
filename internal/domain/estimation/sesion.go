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
	// revelada indica que la ronda actual ya se reveló: no acepta votos nuevos (RN-17.4).
	revelada bool
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

// Resultado es el resultado de revelar una ronda (spec HU-17): todos los votos, si hay consenso y,
// si no lo hay, el valor más alto y el más bajo con todos los integrantes que votaron cada uno.
type Resultado struct {
	Votos         []Voto
	Consenso      bool
	Maximo        int
	Minimo        int
	VotaronMaximo []domain.IntegranteID
	VotaronMinimo []domain.IntegranteID
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

// Votar registra el voto de un integrante en la ronda actual. Una ronda revelada no acepta votos nuevos
// (RN-17.4); la carta tiene que ser de la escala (RN-16.3) y cada integrante vota una sola vez por
// ronda (RN-16.2).
func (s *Sesion) Votar(integranteID domain.IntegranteID, carta int) error {
	if s.revelada {
		return ErrYaRevelada
	}
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
// consultar (RN-16.4); al revelar quedan a la vista (HU-17).
func (s *Sesion) Votos() ([]Voto, error) {
	if !s.revelada {
		return nil, ErrVotosOcultos
	}
	return append([]Voto{}, s.votos...), nil
}

// Revelar cierra la ronda mostrando todos sus votos y las diferencias (HU-17). Requiere al menos un voto
// (RN-17.4) y no se puede repetir.
func (s *Sesion) Revelar() (Resultado, error) {
	if s.revelada {
		return Resultado{}, ErrYaRevelada
	}
	if len(s.votos) == 0 {
		return Resultado{}, ErrSinVotos
	}
	res := s.calcularResultado()
	s.revelada = true
	return res, nil
}

// Resultado devuelve el resultado de una ronda revelada sin cambiar la sesión: lo consulta la pantalla al
// recargar, cuando Revelar ya no puede repetirse. Mientras la votación está abierta, los valores no se
// pueden consultar (RN-16.4).
func (s *Sesion) Resultado() (Resultado, error) {
	if !s.revelada {
		return Resultado{}, ErrVotosOcultos
	}
	return s.calcularResultado(), nil
}

// calcularResultado arma el resultado de la ronda (RN-17.1 a RN-17.3): todos los votos en orden, consenso
// cuando todos son iguales (con un solo voto también, RN-17.2) y, si no lo hay, el valor más alto y el
// más bajo con todos los integrantes que votaron cada uno, en orden de voto.
func (s *Sesion) calcularResultado() Resultado {
	res := Resultado{Votos: append([]Voto{}, s.votos...)}
	res.Maximo, res.Minimo = res.Votos[0].Carta, res.Votos[0].Carta
	for _, v := range res.Votos[1:] {
		if v.Carta > res.Maximo {
			res.Maximo = v.Carta
		}
		if v.Carta < res.Minimo {
			res.Minimo = v.Carta
		}
	}
	res.Consenso = res.Maximo == res.Minimo
	for _, v := range s.votos {
		if v.Carta == res.Maximo {
			res.VotaronMaximo = append(res.VotaronMaximo, v.IntegranteID)
		}
		if v.Carta == res.Minimo {
			res.VotaronMinimo = append(res.VotaronMinimo, v.IntegranteID)
		}
	}
	return res
}

// sesionJSON es la forma serializable de Sesion: la usa Store para guardar data/sesiones.json. No es
// parte de la API del dominio; la sesión se manipula con sus métodos.
type sesionJSON struct {
	HistoriaID domain.HistoriaID
	Ronda      int
	Revelada   bool
	Votos      []Voto
}

// MarshalJSON persiste también los votos, que no están exportados: sin este método, el JSON guardado
// por Store perdería la votación completa.
func (s Sesion) MarshalJSON() ([]byte, error) {
	return json.Marshal(sesionJSON{HistoriaID: s.HistoriaID, Ronda: s.Ronda, Revelada: s.revelada, Votos: s.votos})
}

// UnmarshalJSON reconstruye la sesión guardada, incluido el estado de la ronda y sus votos.
func (s *Sesion) UnmarshalJSON(contenido []byte) error {
	var guardada sesionJSON
	if err := json.Unmarshal(contenido, &guardada); err != nil {
		return err
	}
	s.HistoriaID, s.Ronda, s.revelada, s.votos = guardada.HistoriaID, guardada.Ronda, guardada.Revelada, guardada.Votos
	return nil
}
