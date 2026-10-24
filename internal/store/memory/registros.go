package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

// Registros guarda los registros de esfuerzo (épica E5) en un archivo JSON.
type Registros struct {
	store *Store[[]effort.Registro]
}

// NuevoRegistros abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoRegistros(ruta string) (*Registros, error) {
	s, err := Open[[]effort.Registro](ruta)
	if err != nil {
		return nil, err
	}
	return &Registros{store: s}, nil
}

func (r *Registros) Agregar(reg effort.Registro) (effort.Registro, error) {
	return effort.Registro{}, nil
}

func (r *Registros) DelIntegrante(id domain.IntegranteID) []effort.Registro {
	return nil
}

func (r *Registros) DeHistoria(id domain.HistoriaID) []effort.Registro {
	return nil
}
