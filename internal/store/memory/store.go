package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store guarda un dato de tipo T en memoria y lo persiste en un archivo JSON (TEC-04).
// Es seguro para uso concurrente. Cada repositorio usa su propio Store y su propio archivo.
type Store[T any] struct {
	mu       sync.RWMutex
	ruta     string
	dato     T
	guardado []byte // último JSON guardado, para restaurar si un Update falla
}

// Open abre el archivo de la ruta indicada. Si no existe, arranca con el valor cero de T.
// Con ruta vacía el Store funciona solo en memoria (útil en los tests).
func Open[T any](ruta string) (*Store[T], error) {
	s := &Store[T]{ruta: ruta}
	if ruta == "" {
		return s, nil
	}

	contenido, err := os.ReadFile(ruta)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("leer %s: %w", ruta, err)
	}
	if err := json.Unmarshal(contenido, &s.dato); err != nil {
		return nil, fmt.Errorf("el archivo %s no es un JSON válido: %w", ruta, err)
	}
	s.guardado = contenido
	return s, nil
}

// View ejecuta fn con el estado actual. fn no debe modificar el dato: para eso está Update.
func (s *Store[T]) View(fn func(T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.dato)
}

// Update ejecuta fn sobre el estado y, si no devuelve error, lo guarda en el archivo.
// Si fn o la escritura fallan, no se guarda nada y el estado vuelve a ser el anterior.
func (s *Store[T]) Update(fn func(*T) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := fn(&s.dato); err != nil {
		s.restaurar()
		return err
	}
	contenido, err := json.MarshalIndent(s.dato, "", "  ")
	if err != nil {
		s.restaurar()
		return fmt.Errorf("serializar datos: %w", err)
	}
	if err := s.escribir(contenido); err != nil {
		s.restaurar()
		return err
	}
	s.guardado = contenido
	return nil
}

// restaurar vuelve el estado en memoria al último guardado.
func (s *Store[T]) restaurar() {
	var anterior T
	if s.guardado != nil {
		// No puede fallar: s.guardado se generó con json.Marshal del mismo tipo.
		_ = json.Unmarshal(s.guardado, &anterior)
	}
	s.dato = anterior
}

// escribir guarda el contenido de forma atómica: escribe un temporal y lo renombra.
func (s *Store[T]) escribir(contenido []byte) error {
	if s.ruta == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.ruta), 0o755); err != nil {
		return fmt.Errorf("crear carpeta de %s: %w", s.ruta, err)
	}
	tmp := s.ruta + ".tmp"
	if err := os.WriteFile(tmp, contenido, 0o644); err != nil {
		return fmt.Errorf("escribir %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, s.ruta); err != nil {
		return fmt.Errorf("reemplazar %s: %w", s.ruta, err)
	}
	return nil
}
