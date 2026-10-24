package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
)

// Defectos guarda los defectos (épica E6) en un archivo JSON.
type Defectos struct {
	store *Store[[]defect.Defecto]
}

// NuevoDefectos abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoDefectos(ruta string) (*Defectos, error) {
	s, err := Open[[]defect.Defecto](ruta)
	if err != nil {
		return nil, err
	}
	return &Defectos{store: s}, nil
}

func (r *Defectos) Agregar(d defect.Defecto) (defect.Defecto, error) {
	return defect.Defecto{}, nil
}

func (r *Defectos) Buscar(id domain.DefectoID) (defect.Defecto, bool) {
	return defect.Defecto{}, false
}
