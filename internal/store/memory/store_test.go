package memory

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

type datos struct {
	Nombres []string       `json:"nombres"`
	Horas   map[string]int `json:"horas"`
}

func leer(t *testing.T, s *Store[datos]) datos {
	t.Helper()
	var copia datos
	s.View(func(d datos) { copia = d })
	return copia
}

func agregar(nombre string) func(*datos) error {
	return func(d *datos) error {
		d.Nombres = append(d.Nombres, nombre)
		return nil
	}
}

func TestStore_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "sub", "datos.json")

	s, err := Open[datos](ruta)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Update(func(d *datos) error {
		d.Nombres = []string{"Ana"}
		d.Horas = map[string]int{"Ana": 3}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	reabierto, err := Open[datos](ruta)
	if err != nil {
		t.Fatal(err)
	}
	got := leer(t, reabierto)
	if len(got.Nombres) != 1 || got.Nombres[0] != "Ana" || got.Horas["Ana"] != 3 {
		t.Errorf("al reabrir se obtuvo %+v", got)
	}
}

func TestOpen_archivoInexistente(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "nuevo.json")

	s, err := Open[datos](ruta)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got := leer(t, s); got.Nombres != nil || got.Horas != nil {
		t.Errorf("se esperaba el valor cero, se obtuvo %+v", got)
	}
	if _, err := os.Stat(ruta); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("el archivo no debería crearse hasta el primer Update")
	}
}

func TestUpdate_conErrorNoGuarda(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "datos.json")
	s, _ := Open[datos](ruta)
	if err := s.Update(agregar("Ana")); err != nil {
		t.Fatal(err)
	}

	errValidacion := errors.New("dato inválido")
	err := s.Update(func(d *datos) error {
		d.Nombres = append(d.Nombres, "Beto")
		return errValidacion
	})

	if !errors.Is(err, errValidacion) {
		t.Fatalf("se esperaba el error de la función, se obtuvo %v", err)
	}
	if got := leer(t, s); len(got.Nombres) != 1 {
		t.Errorf("el estado en memoria debería volver al anterior, se obtuvo %+v", got)
	}
	reabierto, _ := Open[datos](ruta)
	if got := leer(t, reabierto); len(got.Nombres) != 1 {
		t.Errorf("el archivo no debería cambiar, se obtuvo %+v", got)
	}
}

func TestOpen_archivoCorrupto(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "roto.json")
	if err := os.WriteFile(ruta, []byte("{esto no es json"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Open[datos](ruta); err == nil {
		t.Fatal("se esperaba un error al abrir un archivo corrupto")
	}
}

func TestStore_soloMemoria(t *testing.T) {
	s, err := Open[datos]("")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Update(agregar("Ana")); err != nil {
		t.Fatal(err)
	}
	if got := leer(t, s); len(got.Nombres) != 1 {
		t.Errorf("se obtuvo %+v", got)
	}
}

func TestUpdate_concurrente(t *testing.T) {
	s, _ := Open[datos](filepath.Join(t.TempDir(), "datos.json"))

	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Update(agregar("x"))
			leer(t, s)
		}()
	}
	wg.Wait()

	if got := leer(t, s); len(got.Nombres) != n {
		t.Errorf("se esperaban %d nombres, hay %d", n, len(got.Nombres))
	}
}
