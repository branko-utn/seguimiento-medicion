package memory

import (
	"strconv"

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

// Agregar guarda un integrante nuevo con el siguiente ID (I1, I2, …).
func (r *Integrantes) Agregar(i project.Integrante) (project.Integrante, error) {
	err := r.store.Update(func(lista *[]project.Integrante) error {
		i.ID = domain.IntegranteID("I" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, i)
		return nil
	})
	if err != nil {
		return project.Integrante{}, err
	}
	return i, nil
}

// DelProyecto devuelve una copia de los integrantes de un proyecto, en el orden en que se registraron.
func (r *Integrantes) DelProyecto(id domain.ProyectoID) []project.Integrante {
	var delProyecto []project.Integrante
	r.store.View(func(lista []project.Integrante) {
		for _, i := range lista {
			if i.ProyectoID == id {
				delProyecto = append(delProyecto, i)
			}
		}
	})
	return delProyecto
}
