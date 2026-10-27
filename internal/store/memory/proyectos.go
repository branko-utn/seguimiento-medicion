package memory

import (
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// Proyectos guarda los proyectos (épica E1) en un archivo JSON.
type Proyectos struct {
	store *Store[[]project.Proyecto]
}

// NuevoProyectos abre el repositorio. Con ruta "" funciona solo en memoria.
func NuevoProyectos(ruta string) (*Proyectos, error) {
	s, err := Open[[]project.Proyecto](ruta)
	if err != nil {
		return nil, err
	}
	return &Proyectos{store: s}, nil
}

// Agregar guarda un proyecto nuevo con el siguiente ID (P1, P2, …).
func (r *Proyectos) Agregar(p project.Proyecto) (project.Proyecto, error) {
	err := r.store.Update(func(lista *[]project.Proyecto) error {
		p.ID = domain.ProyectoID("P" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, p)
		return nil
	})
	if err != nil {
		return project.Proyecto{}, err
	}
	return p, nil
}

// Buscar devuelve el proyecto con ese ID, o project.ErrNoExiste.
func (r *Proyectos) Buscar(id domain.ProyectoID) (project.Proyecto, error) {
	var (
		encontrado project.Proyecto
		existe     bool
	)
	r.store.View(func(lista []project.Proyecto) {
		for _, p := range lista {
			if p.ID == id {
				encontrado, existe = p, true
				return
			}
		}
	})
	if !existe {
		return project.Proyecto{}, project.ErrNoExiste
	}
	return encontrado, nil
}

// Actualizar reemplaza un proyecto ya guardado (HU-02). Todavía sin implementar (RED).
func (r *Proyectos) Actualizar(p project.Proyecto) error {
	return nil
}

// Listar devuelve una copia de todos los proyectos.
func (r *Proyectos) Listar() []project.Proyecto {
	var copia []project.Proyecto
	r.store.View(func(lista []project.Proyecto) {
		copia = append(copia, lista...)
	})
	return copia
}
