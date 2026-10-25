package memory

import (
	"path/filepath"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
)

func defecto(descripcion string) defect.Defecto {
	return defect.Defecto{Descripcion: descripcion, Severidad: defect.Alta, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"}
}

func TestDefectos_AgregarGeneraIDsYBusca(t *testing.T) {
	r, err := NuevoDefectos("")
	if err != nil {
		t.Fatal(err)
	}
	for i, desc := range []string{"Uno", "Dos"} {
		d, err := r.Agregar(defecto(desc))
		if err != nil {
			t.Fatalf("Agregar: %v", err)
		}
		quiere := domain.DefectoID([]string{"D1", "D2"}[i])
		if d.ID != quiere {
			t.Errorf("ID = %q, se esperaba %q", d.ID, quiere)
		}
	}

	casos := []struct {
		titulo           string
		id               domain.DefectoID
		quiereEncontrado bool
		quiereDesc       string
	}{
		{"busca el primero", "D1", true, "Uno"},
		{"busca el segundo", "D2", true, "Dos"},
		{"inexistente devuelve encontrado falso", "D99", false, ""},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			d, ok := r.Buscar(c.id)
			if ok != c.quiereEncontrado {
				t.Fatalf("encontrado = %v, se esperaba %v", ok, c.quiereEncontrado)
			}
			if d.Descripcion != c.quiereDesc {
				t.Errorf("Descripcion = %q, se esperaba %q", d.Descripcion, c.quiereDesc)
			}
		})
	}
}

func TestDefectos_persisteYRecupera(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "defectos.json")
	r, err := NuevoDefectos(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Agregar(defecto("Acepta fechas invertidas")); err != nil {
		t.Fatal(err)
	}

	otra, err := NuevoDefectos(ruta)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := otra.Buscar("D1")
	if !ok || d.Descripcion != "Acepta fechas invertidas" || d.Severidad != defect.Alta || d.Estado != defect.Abierto {
		t.Errorf("al reabrir: %+v (encontrado %v), se esperaba el defecto D1 Alta y Abierto", d, ok)
	}
}
