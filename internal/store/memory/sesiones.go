package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
)

// Sesiones guarda las votaciones de Planning Poker, una por historia (épica E4).
type Sesiones struct {
	store *Store[[]estimation.Sesion]
}

// NuevoSesiones abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoSesiones(ruta string) (*Sesiones, error) {
	s, err := Open[[]estimation.Sesion](ruta)
	if err != nil {
		return nil, err
	}
	return &Sesiones{store: s}, nil
}

// Guardar crea o reemplaza la votación de la historia de la sesión (RN-16.5: una historia tiene como
// máximo una votación abierta, así que guarda una sola por historia).
func (r *Sesiones) Guardar(s estimation.Sesion) error {
	return nil
}

// Buscar devuelve la votación de una historia y si existe. A diferencia de otros repositorios, no devuelve
// ErrNoExiste: que una historia todavía no tenga una votación abierta no es un error, sino un estado válido.
func (r *Sesiones) Buscar(id domain.HistoriaID) (estimation.Sesion, bool) {
	return estimation.Sesion{}, false
}
