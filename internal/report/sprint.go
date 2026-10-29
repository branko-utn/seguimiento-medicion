package report

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// ErrSprintSinEmpezar indica que se pidió el reporte de un sprint planificado, que todavía no tiene resultado (RN-29.1).
var ErrSprintSinEmpezar = errors.New("el sprint todavía no empezó")

// DatosSprint es lo que el servicio junta para armar el reporte de un sprint (HU-29). El paquete no accede a
// repositorios (ADR 0002): recibe todo armado.
type DatosSprint struct {
	Sprint      sprint.Sprint
	Historias   map[domain.HistoriaID]backlog.Historia // las del sprint
	Estimadas   map[domain.HistoriaID]float64          // horas estimadas (HU-19); falta = 0
	Registros   []effort.Registro                      // de las historias del sprint, en cualquier fecha
	Integrantes map[domain.IntegranteID]string         // nombre de cada integrante
	Defectos    []defect.Defecto                       // los del proyecto
}

// ReporteSprint es el reporte de un sprint para la Sprint Review (HU-29).
type ReporteSprint struct {
	SprintID       domain.SprintID
	Goal           string
	Inicio, Fin    time.Time
	Estado         string
	Parcial        bool // el sprint está activo: el resultado todavía puede cambiar
	Historias      []FilaHistoria
	SPPlanificados int
	SPCompletados  int
	Porcentaje     float64 // historias completadas sobre planificadas (HU-25)
	HorasEstimadas float64
	HorasReales    float64
	// Desviacion solo vale si DesviacionCalculable; sin horas estimadas es "no calculable" (RN-29.4).
	Desviacion           float64
	DesviacionCalculable bool
	Registros            []FilaRegistro
	PorIntegrante        []TotalIntegrante
	Defectos             []FilaDefecto // los detectados en el sprint
	Detectados           int
	Resueltos            int
}

// FilaHistoria es una historia planificada en el sprint.
type FilaHistoria struct {
	ID             domain.HistoriaID
	Titulo         string
	StoryPoints    int
	HorasEstimadas float64
	HorasReales    float64 // solo las registradas entre el inicio y el fin del sprint (RN-29.3)
	Completada     bool
}

// FilaRegistro es un registro de esfuerzo del sprint, con el nombre del integrante.
type FilaRegistro struct {
	Fecha      time.Time
	Integrante string
	HistoriaID domain.HistoriaID
	Actividad  string
	Horas      float64
}

// TotalIntegrante son las horas de un integrante en el sprint.
type TotalIntegrante struct {
	Integrante string
	Horas      float64
}

// FilaDefecto es un defecto detectado en el sprint.
type FilaDefecto struct {
	ID          domain.DefectoID
	Descripcion string
	Severidad   string
	Estado      string
	HistoriaID  domain.HistoriaID
}

// NuevoReporteSprint arma el reporte de un sprint activo o cerrado (RN-29.1). Las métricas salen de internal/metrics:
// el paquete solo elige los datos y los ordena.
func NuevoReporteSprint(d DatosSprint) (ReporteSprint, error) {
	sp := d.Sprint
	if sp.Estado != sprint.Activo && sp.Estado != sprint.Cerrado {
		return ReporteSprint{}, ErrSprintSinEmpezar
	}
	r := ReporteSprint{
		SprintID: sp.ID, Goal: sp.Goal, Inicio: sp.Inicio, Fin: sp.Fin,
		Estado: sp.Estado.String(), Parcial: sp.Estado == sprint.Activo,
	}

	// Historias, en el orden en que se asignaron (RN-29.5), con sus horas dentro del sprint (RN-29.3).
	porHistoria := make(map[domain.HistoriaID][]metrics.RegistroHoras)
	for _, reg := range d.Registros {
		porHistoria[reg.HistoriaID] = append(porHistoria[reg.HistoriaID], metrics.RegistroHoras{Fecha: reg.Fecha, Horas: reg.Horas})
	}
	entradas := make([]metrics.HistoriaSprint, 0, len(sp.Historias))
	var estimadas, reales []float64
	for _, id := range sp.Historias {
		h := d.Historias[id]
		fila := FilaHistoria{
			ID: id, Titulo: h.Titulo, StoryPoints: h.StoryPoints,
			HorasEstimadas: d.Estimadas[id],
			HorasReales:    metrics.HorasEnPeriodo(porHistoria[id], sp.Inicio, sp.Fin),
			Completada:     slices.Contains(sp.Completadas, id),
		}
		r.Historias = append(r.Historias, fila)
		entradas = append(entradas, metrics.HistoriaSprint{StoryPoints: h.StoryPoints, Completada: fila.Completada})
		estimadas = append(estimadas, fila.HorasEstimadas)
		reales = append(reales, fila.HorasReales)
	}

	// SP: los fijados al cerrar o, en un sprint activo, los del momento (RN-29.2).
	if sp.Estado == sprint.Cerrado {
		r.SPPlanificados, r.SPCompletados = sp.SPPlanificados, sp.SPCompletados
	} else {
		r.SPPlanificados, r.SPCompletados = metrics.SPPlanificados(entradas), metrics.SPCompletados(entradas)
	}
	porcentaje, err := metrics.PorcentajeCompletadas(len(sp.Completadas), len(sp.Historias))
	if err != nil {
		return ReporteSprint{}, err
	}
	r.Porcentaje = porcentaje

	r.HorasEstimadas, r.HorasReales = metrics.HorasReales(estimadas), metrics.HorasReales(reales)
	if desviacion, err := metrics.Desviacion(r.HorasEstimadas, r.HorasReales); err == nil {
		r.Desviacion, r.DesviacionCalculable = desviacion, true
	} else if !errors.Is(err, metrics.ErrDesviacionNoCalculable) {
		return ReporteSprint{}, err
	}

	r.Registros, r.PorIntegrante = esfuerzoDelSprint(d, sp)
	r.Defectos, r.Detectados, r.Resueltos = defectosDelSprint(d.Defectos, sp.ID)
	return r, nil
}

// esfuerzoDelSprint devuelve los registros con fecha dentro del sprint, por fecha y en el orden de carga, y el total de
// cada integrante, por nombre.
func esfuerzoDelSprint(d DatosSprint, sp sprint.Sprint) ([]FilaRegistro, []TotalIntegrante) {
	var filas []FilaRegistro
	horas := make(map[string]float64)
	for _, reg := range d.Registros {
		if !enElPeriodo(reg.Fecha, sp.Inicio, sp.Fin) {
			continue
		}
		nombre := d.Integrantes[reg.IntegranteID]
		filas = append(filas, FilaRegistro{Fecha: reg.Fecha, Integrante: nombre, HistoriaID: reg.HistoriaID, Actividad: reg.Actividad, Horas: reg.Horas})
		horas[nombre] += reg.Horas
	}
	slices.SortStableFunc(filas, func(a, b FilaRegistro) int { return a.Fecha.Compare(b.Fecha) })
	totales := make([]TotalIntegrante, 0, len(horas))
	for nombre, h := range horas {
		totales = append(totales, TotalIntegrante{Integrante: nombre, Horas: h})
	}
	slices.SortFunc(totales, func(a, b TotalIntegrante) int { return strings.Compare(a.Integrante, b.Integrante) })
	return filas, totales
}

// defectosDelSprint devuelve los defectos detectados en el sprint y las cantidades de detectados y resueltos en él
// (HU-27).
func defectosDelSprint(ds []defect.Defecto, id domain.SprintID) ([]FilaDefecto, int, int) {
	var filas []FilaDefecto
	entradas := make([]metrics.DefectoSprint, 0, len(ds))
	for _, def := range ds {
		entradas = append(entradas, metrics.DefectoSprint{SprintDeteccion: string(def.SprintDeteccion), SprintResolucion: string(def.SprintResolucion)})
		if def.SprintDeteccion == id {
			filas = append(filas, FilaDefecto{ID: def.ID, Descripcion: def.Descripcion, Severidad: def.Severidad.String(), Estado: def.Estado.String(), HistoriaID: def.HistoriaID})
		}
	}
	return filas, metrics.DefectosDetectados(entradas, string(id)), metrics.DefectosResueltos(entradas, string(id))
}

// enElPeriodo indica si la fecha cae entre inicio y fin, ambos incluidos, comparando por día (la misma regla que
// metrics.HorasEnPeriodo).
func enElPeriodo(fecha, inicio, fin time.Time) bool {
	dia := func(t time.Time) time.Time { return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) }
	f := dia(fecha)
	return !f.Before(dia(inicio)) && !f.After(dia(fin))
}
