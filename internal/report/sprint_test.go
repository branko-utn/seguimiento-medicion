package report

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

func dia(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }

func cerca(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// datosEscenario arma el escenario de la spec HU-29: S1 del 21/10 al 23/10 con HU-1 (3 SP, 4 h estimadas,
// completada) y HU-2 (5 SP, 6 h estimadas). Ana registró 3 h el 21/10 y 2 h el 22/10 en HU-1, y Beto 4 h en HU-1
// el 26/10, fuera del sprint. Hay un defecto detectado en S1 y otro detectado en S0 y resuelto en S1.
func datosEscenario(estado sprint.Estado) DatosSprint {
	sp := sprint.Sprint{
		ID: "S1", Goal: "MVP", Inicio: dia(21), Fin: dia(23), Estado: estado,
		Historias: []domain.HistoriaID{"HU-1", "HU-2"}, Completadas: []domain.HistoriaID{"HU-1"},
	}
	if estado == sprint.Cerrado {
		sp.SPPlanificados, sp.SPCompletados = 8, 3
	}
	return DatosSprint{
		Sprint: sp,
		Historias: map[domain.HistoriaID]backlog.Historia{
			"HU-1": {ID: "HU-1", Titulo: "Crear proyecto", StoryPoints: 3},
			"HU-2": {ID: "HU-2", Titulo: "Crear sprint", StoryPoints: 5},
		},
		Estimadas: map[domain.HistoriaID]float64{"HU-1": 4, "HU-2": 6},
		Registros: []effort.Registro{
			{IntegranteID: "I2", HistoriaID: "HU-1", Fecha: dia(26), Actividad: "Desarrollo", Horas: 4},
			{IntegranteID: "I1", HistoriaID: "HU-1", Fecha: dia(22), Actividad: "Pruebas", Horas: 2},
			{IntegranteID: "I1", HistoriaID: "HU-1", Fecha: dia(21), Actividad: "Desarrollo", Horas: 3},
		},
		Integrantes: map[domain.IntegranteID]string{"I1": "Ana Pérez", "I2": "Beto Gómez"},
		Defectos: []defect.Defecto{
			{ID: "D-1", Descripcion: "Falla al guardar", Severidad: defect.Alta, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"},
			{ID: "D-2", Descripcion: "Texto cortado", Severidad: defect.Baja, Estado: defect.Resuelto, HistoriaID: "HU-2", SprintDeteccion: "S0", SprintResolucion: "S1"},
			{ID: "D-3", Descripcion: "De otro sprint", Severidad: defect.Media, Estado: defect.Abierto, HistoriaID: "HU-2", SprintDeteccion: "S2"},
		},
	}
}

func TestNuevoReporteSprint_cerrado(t *testing.T) {
	r, err := NuevoReporteSprint(datosEscenario(sprint.Cerrado))
	if err != nil {
		t.Fatal(err)
	}
	if r.Goal != "MVP" || r.Estado != "Cerrado" || r.Parcial {
		t.Errorf("datos del sprint = %q %q parcial=%v", r.Goal, r.Estado, r.Parcial)
	}
	// RN-29.5: en el orden de asignación; RN-29.3: las horas reales son solo las del período del sprint.
	quiere := []FilaHistoria{
		{ID: "HU-1", Titulo: "Crear proyecto", StoryPoints: 3, HorasEstimadas: 4, HorasReales: 5, Completada: true},
		{ID: "HU-2", Titulo: "Crear sprint", StoryPoints: 5, HorasEstimadas: 6, HorasReales: 0},
	}
	if len(r.Historias) != len(quiere) {
		t.Fatalf("historias = %+v", r.Historias)
	}
	for i, f := range quiere {
		if r.Historias[i] != f {
			t.Errorf("historia %d = %+v, se esperaba %+v", i, r.Historias[i], f)
		}
	}
	if r.SPPlanificados != 8 || r.SPCompletados != 3 {
		t.Errorf("SP = %d/%d, se esperaba 8/3", r.SPPlanificados, r.SPCompletados)
	}
	if !cerca(r.Porcentaje, 50) {
		t.Errorf("porcentaje = %v, se esperaba 50", r.Porcentaje)
	}
	if !cerca(r.HorasEstimadas, 10) || !cerca(r.HorasReales, 5) || !r.DesviacionCalculable || !cerca(r.Desviacion, -50) {
		t.Errorf("horas = %v/%v desviación %v (calculable %v), se esperaba 10/5 y -50", r.HorasEstimadas, r.HorasReales, r.Desviacion, r.DesviacionCalculable)
	}
	// Los registros del sprint, por fecha; el del 26/10 queda afuera.
	if len(r.Registros) != 2 || !r.Registros[0].Fecha.Equal(dia(21)) || r.Registros[1].Actividad != "Pruebas" || r.Registros[0].Integrante != "Ana Pérez" {
		t.Errorf("registros = %+v", r.Registros)
	}
	if len(r.PorIntegrante) != 1 || r.PorIntegrante[0] != (TotalIntegrante{Integrante: "Ana Pérez", Horas: 5}) {
		t.Errorf("por integrante = %+v", r.PorIntegrante)
	}
	// Defectos detectados en el sprint y cantidades de HU-27: D-2 se resolvió en S1 aunque se detectó en S0.
	if len(r.Defectos) != 1 || r.Defectos[0].ID != "D-1" || r.Defectos[0].Severidad != "Alta" || r.Defectos[0].Estado != "Abierto" {
		t.Errorf("defectos = %+v", r.Defectos)
	}
	if r.Detectados != 1 || r.Resueltos != 1 {
		t.Errorf("defectos detectados/resueltos = %d/%d, se esperaba 1/1", r.Detectados, r.Resueltos)
	}
}

// RN-29.2: en un sprint cerrado se usan los SP fijados al cerrar, aunque una historia haya cambiado después.
func TestNuevoReporteSprint_cerradoUsaLosSPFijados(t *testing.T) {
	d := datosEscenario(sprint.Cerrado)
	h := d.Historias["HU-2"]
	h.StoryPoints = 13
	d.Historias["HU-2"] = h
	r, err := NuevoReporteSprint(d)
	if err != nil {
		t.Fatal(err)
	}
	if r.SPPlanificados != 8 || r.SPCompletados != 3 {
		t.Errorf("SP = %d/%d, se esperaban los fijados al cerrar (8/3)", r.SPPlanificados, r.SPCompletados)
	}
}

// RN-29.1 y RN-29.2: un sprint activo da un reporte parcial con los SP calculados en el momento.
func TestNuevoReporteSprint_activo(t *testing.T) {
	r, err := NuevoReporteSprint(datosEscenario(sprint.Activo))
	if err != nil {
		t.Fatal(err)
	}
	if !r.Parcial || r.Estado != "Activo" {
		t.Errorf("parcial = %v, estado = %q", r.Parcial, r.Estado)
	}
	if r.SPPlanificados != 8 || r.SPCompletados != 3 {
		t.Errorf("SP = %d/%d, se esperaba 8/3", r.SPPlanificados, r.SPCompletados)
	}
}

// Casos límite de la spec: sin historias, sin horas y sin defectos.
func TestNuevoReporteSprint_vacio(t *testing.T) {
	r, err := NuevoReporteSprint(DatosSprint{Sprint: sprint.Sprint{ID: "S1", Inicio: dia(21), Fin: dia(23), Estado: sprint.Activo}})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Historias) != 0 || r.SPPlanificados != 0 || r.Porcentaje != 0 || r.DesviacionCalculable || len(r.Defectos) != 0 {
		t.Errorf("reporte vacío = %+v", r)
	}
}

func TestNuevoReporteSprint_planificado(t *testing.T) {
	if _, err := NuevoReporteSprint(datosEscenario(sprint.Planificado)); !errors.Is(err, ErrSprintSinEmpezar) {
		t.Errorf("error = %v, se esperaba ErrSprintSinEmpezar", err)
	}
}
