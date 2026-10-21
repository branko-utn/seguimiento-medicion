// Package memory implementa los repositorios de la aplicación en memoria, con persistencia a archivos JSON.
//
// Store[T] resuelve la lectura y escritura del archivo; cada repositorio define qué guarda y cómo lo busca.
// Cada épica usa su propio archivo, así no se comparten estructuras. Ejemplo:
//
//	type Proyectos struct {
//		store *Store[[]project.Project]
//	}
//
//	func NuevoProyectos(ruta string) (*Proyectos, error) { // ruta: "data/proyectos.json" o "" en los tests
//		s, err := Open[[]project.Project](ruta)
//		if err != nil {
//			return nil, err
//		}
//		return &Proyectos{store: s}, nil
//	}
//
//	func (r *Proyectos) Guardar(p project.Project) error {
//		return r.store.Update(func(lista *[]project.Project) error {
//			*lista = append(*lista, p)
//			return nil
//		})
//	}
//
// Si la función que se pasa a Update devuelve un error, no se guarda nada.
package memory
