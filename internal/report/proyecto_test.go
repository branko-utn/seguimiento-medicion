package report

import (
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// datosProyecto arma el escenario de la spec HU-30: "Demo" del 21/10 al 02/11 con HU-1 (3 SP, completada en S1),
// HU-2 (5 SP, volvió del S1 y se completó en S2) y HU-3 (2 SP). S1 cerrado (21–23/10), S2 activo (24–26/10) y S3
// planificado. Hay horas dentro y fuera de los sprints y defectos de las tres severidades.
func datosProyecto() DatosProyecto {
	return DatosProyecto{
		Proyecto: project.Proyecto{ID: "P1", Nombre: "Demo", Inicio: dia(21), Fin: time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local)},
		Hoy:      dia(27),
		Sprints: []sprint.Sprint{
			// Desordenados a propósito: el resumen va por fecha de inicio.
			{ID: "S2", Goal: "Segundo", Inicio: dia(24), Fin: dia(26), Estado: sprint.Activo,
				Historias: []domain.HistoriaID{"HU-2", "HU-3"}, Completadas: []domain.HistoriaID{"HU-2"}},
			{ID: "S1", Goal: "Primero", Inicio: dia(21), Fin: dia(23), Estado: sprint.Cerrado, SPPlanificados: 8, SPCompletados: 3,
				Historias: []domain.HistoriaID{"HU-1", "HU-2"}, Completadas: []domain.HistoriaID{"HU-1"}},
			{ID: "S3", Goal: "Tercero", Inicio: dia(27), Fin: dia(29), Estado: sprint.Planificado},
		},
		Historias: []backlog.Historia{
			{ID: "HU-1", Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Completada},
			{ID: "HU-2", Titulo: "Crear sprint", StoryPoints: 5, Estado: backlog.Completada},
			{ID: "HU-3", Titulo: "Registrar defecto", StoryPoints: 2, Estado: backlog.Pendiente},
		},
		Estimadas: map[domain.HistoriaID]float64{"HU-1": 4, "HU-2": 6, "HU-3": 2},
		Registros: []effort.Registro{
			{IntegranteID: "I1", HistoriaID: "HU-1", Fecha: dia(22), Horas: 3},
			{IntegranteID: "I1", HistoriaID: "HU-2", Fecha: dia(25), Horas: 4},
			{IntegranteID: "I1", HistoriaID: "HU-3", Fecha: dia(27), Horas: 1}, // fuera de los sprints con resultado
		},
		Integrantes: map[domain.IntegranteID]string{"I1": "Ana Pérez"},
		Defectos: []defect.Defecto{
			{ID: "D1", Descripcion: "Texto cortado", Severidad: defect.Baja, Estado: defect.Abierto, HistoriaID: "HU-1", SprintDeteccion: "S1"},
			{ID: "D2", Descripcion: "No guarda", Severidad: defect.Critica, Estado: defect.EnProgreso, HistoriaID: "HU-2", SprintDeteccion: "S2"},
			{ID: "D3", Descripcion: "Falla rara", Severidad: defect.Alta, Estado: defect.Resuelto, HistoriaID: "HU-1", SprintDeteccion: "S1", SprintResolucion: "S2"},
			{ID: "D4", Descripcion: "Botón mal", Severidad: defect.Baja, Estado: defect.Abierto, HistoriaID: "HU-3", SprintDeteccion: "S2"},
		},
	}
}

func TestNuevoReporteProyecto(t *testing.T) {
	r, err := NuevoReporteProyecto(datosProyecto())
	if err != nil {
		t.Fatal(err)
	}
	if r.Nombre != "Demo" || r.Estado != "En curso" || r.DiasRestantes != 6 {
		t.Errorf("proyecto = %q %q %d días, se esperaba Demo, En curso y 6 días", r.Nombre, r.Estado, r.DiasRestantes)
	}

	// RN-30.1 y RN-30.2: S1 y S2 en orden de inicio, con las reglas de HU-29; S3 (planificado) no.
	if len(r.Sprints) != 2 || r.Sprints[0].SprintID != "S1" || r.Sprints[1].SprintID != "S2" {
		t.Fatalf("resumen de sprints = %+v", r.Sprints)
	}
	s1, s2 := r.Sprints[0], r.Sprints[1]
	if s1.SPPlanificados != 8 || s1.SPCompletados != 3 || !cerca(s1.Porcentaje, 50) || !cerca(s1.HorasReales, 3) || s1.Detectados != 2 {
		t.Errorf("S1 = %+v", s1)
	}
	if !s2.Parcial || s2.SPPlanificados != 7 || s2.SPCompletados != 5 || !cerca(s2.HorasReales, 4) || s2.Resueltos != 1 {
		t.Errorf("S2 = %+v", s2)
	}

	// RN-30.3: la velocidad usa solo S1 (cerrado). SP completados en total: los de las filas.
	if !cerca(r.Velocidad, 3) || r.SPCompletados != 8 {
		t.Errorf("velocidad = %v y SP completados = %d, se esperaba 3 y 8", r.Velocidad, r.SPCompletados)
	}
	// HU-25 por proyecto: 2 de 3 historias completadas.
	if r.Completadas != 2 || r.Historias != 3 || !cerca(r.Porcentaje, 200.0/3) {
		t.Errorf("completadas = %d de %d (%v %%)", r.Completadas, r.Historias, r.Porcentaje)
	}
	// RN-30.4: horas de todo el proyecto, también las registradas fuera de un sprint.
	if !cerca(r.HorasEstimadas, 12) || !cerca(r.HorasReales, 8) || !r.DesviacionCalculable || !cerca(r.Desviacion, -100.0/3) {
		t.Errorf("horas = %v/%v, desviación %v (calculable %v)", r.HorasEstimadas, r.HorasReales, r.Desviacion, r.DesviacionCalculable)
	}

	// RN-30.5 y HU-27: abiertos = no resueltos, por severidad (Crítica, Alta, Media, Baja) y en el orden de registro.
	if r.Detectados != 4 || r.Resueltos != 1 || r.Abiertos != 3 {
		t.Errorf("defectos = %d/%d/%d, se esperaba 4 detectados, 1 resuelto y 3 abiertos", r.Detectados, r.Resueltos, r.Abiertos)
	}
	var abiertos []domain.DefectoID
	for _, d := range r.DefectosAbiertos {
		abiertos = append(abiertos, d.ID)
	}
	if len(abiertos) != 3 || abiertos[0] != "D2" || abiertos[1] != "D1" || abiertos[2] != "D4" {
		t.Errorf("defectos abiertos = %v, se esperaba [D2 D1 D4]", abiertos)
	}
}

// Casos límite de la spec: proyecto sin sprints, sin historias, sin horas y sin defectos.
func TestNuevoReporteProyecto_vacio(t *testing.T) {
	d := DatosProyecto{Proyecto: datosProyecto().Proyecto, Hoy: dia(20)}
	r, err := NuevoReporteProyecto(d)
	if err != nil {
		t.Fatal(err)
	}
	if r.Estado != "Planificado" || len(r.Sprints) != 0 || r.Velocidad != 0 || r.Porcentaje != 0 || r.DesviacionCalculable ||
		r.Detectados != 0 || len(r.DefectosAbiertos) != 0 {
		t.Errorf("reporte vacío = %+v", r)
	}
}
