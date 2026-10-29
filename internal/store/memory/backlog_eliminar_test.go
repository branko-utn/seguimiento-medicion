package memory

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// historiasDePrueba abre un repositorio con las historias HU-1, HU-2 y HU-3.
func historiasDePrueba(t *testing.T, ruta string) *Historias {
	t.Helper()
	r, err := NuevoHistorias(ruta)
	if err != nil {
		t.Fatal(err)
	}
	for _, titulo := range []string{"Crear proyecto", "Crear sprint", "Login"} {
		if _, err := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: titulo}); err != nil {
			t.Fatal(err)
		}
	}
	return r
}

func ids(lista []backlog.Historia) []domain.HistoriaID {
	var out []domain.HistoriaID
	for _, h := range lista {
		out = append(out, h.ID)
	}
	return out
}

func TestHistorias_Eliminar(t *testing.T) {
	r := historiasDePrueba(t, "")

	if err := r.Eliminar("HU-2"); err != nil {
		t.Fatalf("Eliminar: %v", err)
	}

	// La historia deja de estar en el backlog.
	if got := ids(r.Listar()); len(got) != 2 || got[0] != "HU-1" || got[1] != "HU-3" {
		t.Errorf("Listar() = %v, se esperaba [HU-1 HU-3]", got)
	}
	if _, err := r.Buscar("HU-2"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Buscar de una eliminada: error = %v, se esperaba ErrNoExiste", err)
	}
	if err := r.Guardar(backlog.Historia{ID: "HU-2", Titulo: "Otra"}); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Guardar de una eliminada: error = %v, se esperaba ErrNoExiste", err)
	}
}

func TestHistorias_EliminarInexistente(t *testing.T) {
	r := historiasDePrueba(t, "")

	if err := r.Eliminar("HU-99"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Eliminar HU-99: error = %v, se esperaba ErrNoExiste", err)
	}
	if err := r.Eliminar("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := r.Eliminar("HU-1"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Eliminar dos veces: error = %v, se esperaba ErrNoExiste", err)
	}
}

func TestHistorias_EliminarLaUnica(t *testing.T) {
	// Caso límite: eliminar la única historia deja el backlog vacío.
	r, _ := NuevoHistorias("")
	if _, err := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Login"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Eliminar("HU-1"); err != nil {
		t.Fatal(err)
	}
	if got := r.Listar(); len(got) != 0 {
		t.Errorf("Listar() = %v, se esperaba vacío", ids(got))
	}
}

func TestHistorias_EliminarNoLiberaElID(t *testing.T) {
	// RN-06.4: aunque se elimine la última, la historia nueva no reutiliza su ID.
	r := historiasDePrueba(t, "")
	if err := r.Eliminar("HU-3"); err != nil {
		t.Fatal(err)
	}
	h, err := r.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Cerrar sprint"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != "HU-4" {
		t.Errorf("ID = %q, se esperaba HU-4 (HU-3 no se reutiliza)", h.ID)
	}
}

func TestHistorias_EliminarPersiste(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "backlog.json")
	r := historiasDePrueba(t, ruta)
	if err := r.Eliminar("HU-3"); err != nil {
		t.Fatal(err)
	}

	reabierto, err := NuevoHistorias(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(reabierto.Listar()); len(got) != 2 {
		t.Errorf("después de reabrir, Listar() = %v, se esperaba [HU-1 HU-2]", got)
	}
	h, err := reabierto.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Cerrar sprint"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != "HU-4" {
		t.Errorf("después de reabrir, ID = %q, se esperaba HU-4 (RN-06.4)", h.ID)
	}
}

func TestHistorias_leeElFormatoAnterior(t *testing.T) {
	// Los archivos guardados antes de HU-06 (una lista de historias, sin la marca de eliminada) se siguen leyendo.
	ruta := filepath.Join(t.TempDir(), "backlog.json")
	viejo := `[{"ID":"HU-1","ProyectoID":"P1","Titulo":"Login","Descripcion":"","Prioridad":1,"Estado":1,"StoryPoints":3,"Criterios":null,"SprintID":""}]`
	if err := os.WriteFile(ruta, []byte(viejo), 0o644); err != nil {
		t.Fatal(err)
	}

	r, err := NuevoHistorias(ruta)
	if err != nil {
		t.Fatal(err)
	}
	h, err := r.Buscar("HU-1")
	if err != nil {
		t.Fatalf("Buscar HU-1: %v", err)
	}
	if h.Titulo != "Login" || h.StoryPoints != 3 || h.Prioridad != backlog.PrioridadAlta {
		t.Errorf("historia leída = %+v", h)
	}
}
