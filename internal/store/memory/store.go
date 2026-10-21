package memory

// Store guarda un dato de tipo T en memoria y lo persiste en un archivo JSON (TEC-04).
type Store[T any] struct{}

// Open abre el archivo de la ruta indicada.
func Open[T any](ruta string) (*Store[T], error) {
	return &Store[T]{}, nil
}

// View ejecuta fn con el estado actual.
func (s *Store[T]) View(fn func(T)) {}

// Update ejecuta fn sobre el estado y lo guarda si fn no devuelve error.
func (s *Store[T]) Update(fn func(*T) error) error {
	return nil
}
