package memory

import (
	"strconv"

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

// Agregar guarda un defecto nuevo con el siguiente ID (D1, D2, …).
func (r *Defectos) Agregar(d defect.Defecto) (defect.Defecto, error) {
	err := r.store.Update(func(lista *[]defect.Defecto) error {
		d.ID = domain.DefectoID("D" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, d)
		return nil
	})
	if err != nil {
		return defect.Defecto{}, err
	}
	return d, nil
}

// Buscar devuelve el defecto y si se encontró. No devuelve un error porque la spec de HU-22 no define uno para un
// defecto inexistente: lo usan los escenarios para verificar lo guardado.
func (r *Defectos) Buscar(id domain.DefectoID) (defect.Defecto, bool) {
	var (
		encontrado defect.Defecto
		existe     bool
	)
	r.store.View(func(lista []defect.Defecto) {
		for _, d := range lista {
			if d.ID == id {
				encontrado, existe = d, true
				return
			}
		}
	})
	return encontrado, existe
}

// Actualizar reemplaza un defecto ya guardado, o devuelve defect.ErrNoExiste.
func (r *Defectos) Actualizar(d defect.Defecto) error {
	return r.store.Update(func(lista *[]defect.Defecto) error {
		for i := range *lista {
			if (*lista)[i].ID == d.ID {
				(*lista)[i] = d
				return nil
			}
		}
		return defect.ErrNoExiste
	})
}

// Listar devuelve una copia de todos los defectos, en el orden en que se registraron.
func (r *Defectos) Listar() []defect.Defecto {
	var todos []defect.Defecto
	r.store.View(func(lista []defect.Defecto) {
		todos = append(todos, lista...)
	})
	return todos
}
