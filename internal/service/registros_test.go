package service

import (
	"errors"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServicioRegistros arma el servicio con dos proyectos: "Demo" (P1), con la integrante Ana (I1) y la historia
// HU-1, y "Otro" (P2), con el integrante Beto (I2).
func nuevoServicioRegistros(t *testing.T) *Registros {
	t.Helper()
	repoProy, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoInt, err := memory.NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}
	repoHist, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	repoReg, err := memory.NuevoRegistros("")
	if err != nil {
		t.Fatal(err)
	}
	repoSpr, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	proy := NuevoProyectos(repoProy)
	integrantes := NuevoIntegrantes(repoInt, repoProy)
	sprints := NuevoSprints(repoSpr, repoProy, repoHist)
	hist := NuevoBacklog(repoHist, repoProy, sprints, nil, nil)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proy.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := integrantes.Registrar("P1", "Ana", project.ProductBuilder); err != nil {
		t.Fatal(err)
	}
	if _, err := integrantes.Registrar("P2", "Beto", project.ProductBuilder); err != nil {
		t.Fatal(err)
	}
	if _, err := hist.Crear("P1", "Crear proyecto", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	return NuevoRegistros(repoReg, repoInt, repoHist, repoProy)
}

func TestRegistros_Registrar(t *testing.T) {
	s := nuevoServicioRegistros(t)
	hoy := octubre(24)

	r, err := s.Registrar("I1", "HU-1", hoy, "Desarrollo", 3, hoy)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.ID != "E1" || r.IntegranteID != "I1" || r.HistoriaID != "HU-1" || r.Horas != 3 || r.Actividad != "Desarrollo" {
		t.Errorf("registro = %+v, se esperaba E1 de Ana en HU-1 por 3 horas de Desarrollo", r)
	}
	if got := s.DeHistoria("HU-1"); len(got) != 1 || got[0].ID != "E1" {
		t.Errorf("DeHistoria(HU-1) = %+v, se esperaba [E1]", got)
	}
	if got := s.DelIntegrante("I1"); len(got) != 1 || got[0].ID != "E1" {
		t.Errorf("DelIntegrante(I1) = %+v, se esperaba [E1]", got)
	}
}

func TestRegistros_Registrar_topeDiario(t *testing.T) {
	hoy := octubre(24)
	casos := []struct {
		titulo      string
		previos     []float64 // horas ya registradas por Ana el 23/10
		fecha       time.Time
		horas       float64
		quiereError error
	}{
		{titulo: "10 + 14 horas el mismo día: 24 justas", previos: []float64{10}, fecha: octubre(23), horas: 14},
		{titulo: "24 horas y 0,5 más el mismo día", previos: []float64{24}, fecha: octubre(23), horas: 0.5, quiereError: effort.ErrSuperaHorasDiarias},
		{titulo: "10 + 14 + 0,5 el mismo día", previos: []float64{10, 14}, fecha: octubre(23), horas: 0.5, quiereError: effort.ErrSuperaHorasDiarias},
		{titulo: "24 horas un día y 1 al día siguiente", previos: []float64{24}, fecha: octubre(24), horas: 1},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s := nuevoServicioRegistros(t)
			for _, h := range c.previos {
				if _, err := s.Registrar("I1", "HU-1", octubre(23), "Desarrollo", h, hoy); err != nil {
					t.Fatal(err)
				}
			}
			_, err := s.Registrar("I1", "HU-1", c.fecha, "Tests", c.horas, hoy)
			if !errors.Is(err, c.quiereError) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiereError)
			}
			if c.quiereError != nil && len(s.DelIntegrante("I1")) != len(c.previos) {
				t.Errorf("el registro rechazado no debe guardarse")
			}
		})
	}
}

func TestRegistros_Registrar_errores(t *testing.T) {
	hoy := octubre(24)
	casos := []struct {
		titulo      string
		integrante  domain.IntegranteID
		historia    domain.HistoriaID
		fecha       time.Time
		actividad   string
		horas       float64
		quiereError error
	}{
		{"integrante inexistente", "I99", "HU-1", hoy, "Desarrollo", 2, project.ErrIntegranteNoExiste},
		{"historia inexistente", "I1", "HU-99", hoy, "Desarrollo", 2, backlog.ErrNoExiste},
		{"integrante de otro proyecto", "I2", "HU-1", hoy, "Desarrollo", 2, effort.ErrIntegranteDeOtroProyecto},
		{"actividad vacía", "I1", "HU-1", hoy, "  ", 3, effort.ErrActividadVacia},
		{"horas en cero", "I1", "HU-1", hoy, "Desarrollo", 0, effort.ErrHorasInvalidas},
		{"fecha futura", "I1", "HU-1", octubre(25), "Desarrollo", 2, effort.ErrFechaFutura},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s := nuevoServicioRegistros(t)
			_, err := s.Registrar(c.integrante, c.historia, c.fecha, c.actividad, c.horas, hoy)
			if !errors.Is(err, c.quiereError) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiereError)
			}
			if n := len(s.DeHistoria("HU-1")); n != 0 {
				t.Errorf("se guardaron %d registros, no debía guardarse ninguno", n)
			}
		})
	}
}
