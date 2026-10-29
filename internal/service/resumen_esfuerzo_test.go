package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// servicioResumen arma el servicio con dos proyectos:
//   - "Demo" (P1): los integrantes Beto Gómez (I1) y Ana Pérez (I2), cargados en ese orden para que el orden por ID
//     no coincida con el orden por nombre, y las historias HU-1 a HU-10.
//   - "Otro" (P2): la integrante Carla (I3) y la historia HU-11.
func servicioResumen(t *testing.T) *Registros {
	t.Helper()
	repoProy, _ := memory.NuevoProyectos("")
	repoInt, _ := memory.NuevoIntegrantes("")
	repoHist, _ := memory.NuevoHistorias("")
	repoReg, _ := memory.NuevoRegistros("")
	repoSpr, _ := memory.NuevoSprints("")
	proyectos := NuevoProyectos(repoProy)
	integrantes := NuevoIntegrantes(repoInt, repoProy)
	hist := NuevoBacklog(repoHist, repoProy, NuevoSprints(repoSpr, repoProy, repoHist), nil, nil)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	for _, i := range []struct {
		proyecto domain.ProyectoID
		nombre   string
	}{{"P1", "Beto Gómez"}, {"P1", "Ana Pérez"}, {"P2", "Carla"}} {
		if _, err := integrantes.Registrar(i.proyecto, i.nombre, project.ProductBuilder); err != nil {
			t.Fatal(err)
		}
	}
	for n := 1; n <= 10; n++ {
		if _, err := hist.Crear("P1", "Historia", "", backlog.PrioridadMedia); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := hist.Crear("P2", "De otro proyecto", "", backlog.PrioridadMedia); err != nil {
		t.Fatal(err)
	}
	return NuevoRegistros(repoReg, repoInt, repoHist, repoProy)
}

// registrar carga un registro de esfuerzo válido; hoy es el 27 de octubre.
func registrar(t *testing.T, s *Registros, integrante domain.IntegranteID, historia domain.HistoriaID, dia int, actividad string, horas float64) {
	t.Helper()
	if _, err := s.Registrar(integrante, historia, octubre(dia), actividad, horas, octubre(27)); err != nil {
		t.Fatalf("registrar %s en %s: %v", integrante, historia, err)
	}
}

func TestRegistros_DelProyecto_totales(t *testing.T) {
	s := servicioResumen(t)
	registrar(t, s, "I2", "HU-10", 22, "Desarrollo", 3) // E1 · Ana
	registrar(t, s, "I1", "HU-10", 23, "Pruebas", 2.5)  // E2 · Beto
	registrar(t, s, "I2", "HU-2", 24, "Desarrollo", 4)  // E3 · Ana
	registrar(t, s, "I3", "HU-11", 24, "Desarrollo", 8) // E4 · Carla, de otro proyecto (RN-21.1)

	got, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	// RN-21.3: por historia, en orden numérico (HU-2 antes que HU-10); RN-21.4: sin las historias sin registros.
	wantHistorias := []TotalHistoria{{"HU-2", 4}, {"HU-10", 5.5}}
	if !reflect.DeepEqual(got.PorHistoria, wantHistorias) {
		t.Errorf("PorHistoria = %+v, se esperaba %+v", got.PorHistoria, wantHistorias)
	}
	// RN-21.3: por integrante, por nombre; un integrante con registros en varias historias aparece una vez.
	wantIntegrantes := []TotalIntegrante{{"I2", "Ana Pérez", 7}, {"I1", "Beto Gómez", 2.5}}
	if !reflect.DeepEqual(got.PorIntegrante, wantIntegrantes) {
		t.Errorf("PorIntegrante = %+v, se esperaba %+v", got.PorIntegrante, wantIntegrantes)
	}
	// RN-21.5: el total es la suma de cada lista de totales; las horas con decimales no se redondean.
	if got.Total != 9.5 {
		t.Errorf("Total = %v, se esperaba 9.5", got.Total)
	}
	if len(got.Registros) != 3 {
		t.Errorf("Registros tiene %d, se esperaban 3 (sin el de otro proyecto)", len(got.Registros))
	}
}

func TestRegistros_DelProyecto_ordenDeLosRegistros(t *testing.T) {
	s := servicioResumen(t)
	registrar(t, s, "I2", "HU-1", 22, "Desarrollo", 3) // E1
	registrar(t, s, "I1", "HU-2", 25, "Pruebas", 2)    // E2
	registrar(t, s, "I1", "HU-1", 24, "Revisión", 1)   // E3
	registrar(t, s, "I2", "HU-2", 24, "Desarrollo", 1) // E4: misma fecha que E3, cargado después

	got, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatal(err)
	}

	// RN-21.2: del más reciente al más antiguo; en la misma fecha, por orden de carga.
	var ids []domain.RegistroID
	for _, r := range got.Registros {
		ids = append(ids, r.ID)
	}
	if want := []domain.RegistroID{"E2", "E3", "E4", "E1"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("orden = %v, se esperaba %v", ids, want)
	}
	// Cada registro trae el nombre de su integrante, para mostrarlo.
	if got.Registros[0].Integrante != "Beto Gómez" || got.Registros[3].Integrante != "Ana Pérez" {
		t.Errorf("nombres = %q y %q, se esperaba Beto Gómez y Ana Pérez", got.Registros[0].Integrante, got.Registros[3].Integrante)
	}
}

func TestRegistros_DelProyecto_sinRegistros(t *testing.T) {
	s := servicioResumen(t)
	registrar(t, s, "I3", "HU-11", 24, "Desarrollo", 8) // solo hay esfuerzo en otro proyecto

	got, err := s.DelProyecto("P1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// Proyecto sin registros: listas vacías y total 0.
	if len(got.Registros) != 0 || len(got.PorHistoria) != 0 || len(got.PorIntegrante) != 0 || got.Total != 0 {
		t.Errorf("resumen = %+v, se esperaba vacío con total 0", got)
	}
}

func TestRegistros_DelProyecto_proyectoInexistente(t *testing.T) {
	s := servicioResumen(t)

	if _, err := s.DelProyecto("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("error = %v, se esperaba project.ErrNoExiste", err)
	}
}
