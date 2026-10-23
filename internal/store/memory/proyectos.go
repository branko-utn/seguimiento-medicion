package memory

import "github.com/branko-utn/seguimiento-medicion/internal/domain/project"

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

// Agregar guarda un proyecto nuevo con el siguiente ID (P1, P2, …). El stub existe solo para la fase RED.
func (r *Proyectos) Agregar(p project.Proyecto) (project.Proyecto, error) {
	return project.Proyecto{}, nil
}

// Listar devuelve una copia de todos los proyectos.
func (r *Proyectos) Listar() []project.Proyecto {
	return nil
}
