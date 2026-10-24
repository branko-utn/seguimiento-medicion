package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// Historias guarda el Product Backlog (épica E2) en un archivo JSON.
type Historias struct {
	store *Store[[]backlog.Historia]
}

// NuevoHistorias abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoHistorias(ruta string) (*Historias, error) {
	s, err := Open[[]backlog.Historia](ruta)
	if err != nil {
		return nil, err
	}
	return &Historias{store: s}, nil
}

// Agregar guarda una historia nueva con el siguiente ID de su proyecto (RN-05.1).
func (r *Historias) Agregar(h backlog.Historia) (backlog.Historia, error) {
	return h, nil // stub RED
}

// Listar devuelve todas las historias, en el orden en que se crearon.
func (r *Historias) Listar() []backlog.Historia {
	return nil // stub RED
}
