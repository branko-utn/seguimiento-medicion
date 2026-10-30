package memory

import (
	"strconv"

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

// Agregar guarda un registro nuevo con el siguiente ID (E1, E2, …).
func (r *Registros) Agregar(reg effort.Registro) (effort.Registro, error) {
	err := r.store.Update(func(lista *[]effort.Registro) error {
		reg.ID = domain.RegistroID("E" + strconv.Itoa(len(*lista)+1))
		*lista = append(*lista, reg)
		return nil
	})
	if err != nil {
		return effort.Registro{}, err
	}
	return reg, nil
}

// DelIntegrante devuelve una copia de los registros de un integrante, en el orden en que se registraron.
func (r *Registros) DelIntegrante(id domain.IntegranteID) []effort.Registro {
	return r.filtrar(func(reg effort.Registro) bool { return reg.IntegranteID == id })
}

// DeHistoria devuelve una copia de los registros de una historia, en el orden en que se registraron.
func (r *Registros) DeHistoria(id domain.HistoriaID) []effort.Registro {
	return r.filtrar(func(reg effort.Registro) bool { return reg.HistoriaID == id })
}

func (r *Registros) filtrar(incluir func(effort.Registro) bool) []effort.Registro {
	var filtrados []effort.Registro
	r.store.View(func(lista []effort.Registro) {
		for _, reg := range lista {
			if incluir(reg) {
				filtrados = append(filtrados, reg)
			}
		}
	})
	return filtrados
}

// Listar devuelve una copia de todos los registros.
func (r *Registros) Listar() []effort.Registro {
	var lista []effort.Registro
	r.store.View(func(l []effort.Registro) {
		lista = append(lista, l...)
	})
	return lista
}
