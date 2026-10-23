package memory

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// Integrantes guarda los integrantes de todos los proyectos (épica E1) en un archivo JSON.
type Integrantes struct {
	store *Store[[]project.Integrante]
}

// NuevoIntegrantes abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoIntegrantes(ruta string) (*Integrantes, error) {
	s, err := Open[[]project.Integrante](ruta)
	if err != nil {
		return nil, err
	}
	return &Integrantes{store: s}, nil
}

// Agregar guarda un integrante nuevo con el siguiente ID (I1, I2, …). El stub existe solo para la fase RED.
func (r *Integrantes) Agregar(i project.Integrante) (project.Integrante, error) {
	return project.Integrante{}, nil
}

// DelProyecto devuelve una copia de los integrantes de un proyecto, en el orden en que se registraron.
// El stub existe solo para la fase RED.
func (r *Integrantes) DelProyecto(id domain.ProyectoID) []project.Integrante {
	return nil
}
