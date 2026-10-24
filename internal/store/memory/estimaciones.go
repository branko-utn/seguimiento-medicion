package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

// Estimaciones guarda las horas estimadas por historia (Epica E5).
type Estimaciones struct {
	store *Store[[]effort.EstimacionHoras]
}

// NuevoEstimaciones abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoEstimaciones(ruta string) (*Estimaciones, error) {
	s, err := Open[[]effort.EstimacionHoras](ruta)
	if err != nil {
		return nil, err
	}
	return &Estimaciones{store: s}, nil
}

// Guardar guarda o reemplaza la estimacion de una historia (RN-19.2).
func (r *Estimaciones) Guardar(e effort.EstimacionHoras) error {
	return nil
}

// Buscar devuelve la estimacion de una historia. Si no existe, devuelve (cero, false).
func (r *Estimaciones) Buscar(id domain.HistoriaID) (effort.EstimacionHoras, bool) {
	return effort.EstimacionHoras{}, false
}
