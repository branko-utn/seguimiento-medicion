package report

import (
	"errors"
	"slices"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// DatosProyecto es lo que el servicio junta para armar el reporte de un proyecto (HU-30): todo lo del proyecto, en el
// mismo formato que DatosSprint, para que las filas del resumen se armen con NuevoReporteSprint (RN-30.2).
type DatosProyecto struct {
	Proyecto    project.Proyecto
	Hoy         time.Time
	Sprints     []sprint.Sprint    // todos los del proyecto
	Historias   []backlog.Historia // todas las del proyecto
	Estimadas   map[domain.HistoriaID]float64
	Registros   []effort.Registro // de las historias del proyecto, en cualquier fecha
	Integrantes map[domain.IntegranteID]string
	Defectos    []defect.Defecto // los del proyecto, en el orden en que se registraron
}

// ReporteProyecto es el reporte del proyecto completo (HU-30).
type ReporteProyecto struct {
	ProyectoID    domain.ProyectoID
	Nombre        string
	Inicio, Fin   time.Time
	Estado        string
	DiasRestantes int
	// Sprints es el resumen: una fila por sprint activo o cerrado, en orden de inicio (RN-30.1).
	Sprints       []ReporteSprint
	Velocidad     float64 // de los sprints cerrados (RN-30.3)
	SPCompletados int     // suma de las filas del resumen
	Completadas   int     // historias del proyecto en estado Completada
	Historias     int     // historias del proyecto
	Porcentaje    float64 // HU-25 por proyecto
	// Horas de todo el proyecto (RN-30.4); Desviacion solo vale si DesviacionCalculable.
	HorasEstimadas       float64
	HorasReales          float64
	Desviacion           float64
	DesviacionCalculable bool
	Detectados           int
	Resueltos            int
	Abiertos             int
	DefectosAbiertos     []FilaDefecto // por severidad (RN-30.5)
}

// NuevoReporteProyecto arma el reporte de un proyecto (HU-30). Cada fila del resumen es un ReporteSprint, así el
// reporte del proyecto muestra los mismos valores que el de cada sprint (RN-30.2).
func NuevoReporteProyecto(d DatosProyecto) (ReporteProyecto, error) {
	p := d.Proyecto
	r := ReporteProyecto{
		ProyectoID: p.ID, Nombre: p.Nombre, Inicio: p.Inicio, Fin: p.Fin,
		Estado: p.EstadoAl(d.Hoy).String(), DiasRestantes: p.DiasRestantes(d.Hoy),
		Historias: len(d.Historias),
	}

	if err := r.resumirSprints(d); err != nil {
		return ReporteProyecto{}, err
	}

	var estimadas, reales []float64
	for _, h := range d.Historias {
		if h.Estado == backlog.Completada {
			r.Completadas++
		}
		estimadas = append(estimadas, d.Estimadas[h.ID])
	}
	porcentaje, err := metrics.PorcentajeCompletadas(r.Completadas, r.Historias)
	if err != nil {
		return ReporteProyecto{}, err
	}
	r.Porcentaje = porcentaje

	for _, reg := range d.Registros {
		reales = append(reales, reg.Horas)
	}
	r.HorasEstimadas, r.HorasReales = metrics.HorasReales(estimadas), metrics.HorasReales(reales)
	if desviacion, err := metrics.Desviacion(r.HorasEstimadas, r.HorasReales); err == nil {
		r.Desviacion, r.DesviacionCalculable = desviacion, true
	} else if !errors.Is(err, metrics.ErrDesviacionNoCalculable) {
		return ReporteProyecto{}, err
	}

	r.resumirDefectos(d.Defectos)
	return r, nil
}

// resumirSprints arma una fila por sprint activo o cerrado, en orden de inicio, y la velocidad con los cerrados.
func (r *ReporteProyecto) resumirSprints(d DatosProyecto) error {
	sprints := slices.Clone(d.Sprints)
	slices.SortStableFunc(sprints, func(a, b sprint.Sprint) int { return a.Inicio.Compare(b.Inicio) })
	historias := make(map[domain.HistoriaID]backlog.Historia, len(d.Historias))
	for _, h := range d.Historias {
		historias[h.ID] = h
	}
	var cerrados []int
	for _, sp := range sprints {
		if sp.Estado == sprint.Planificado {
			continue
		}
		fila, err := NuevoReporteSprint(DatosSprint{
			Sprint: sp, Historias: historias, Estimadas: d.Estimadas,
			Registros: d.Registros, Integrantes: d.Integrantes, Defectos: d.Defectos,
		})
		if err != nil {
			return err
		}
		r.Sprints = append(r.Sprints, fila)
		r.SPCompletados += fila.SPCompletados
		if sp.Estado == sprint.Cerrado {
			cerrados = append(cerrados, sp.SPCompletados)
		}
	}
	r.Velocidad = metrics.Velocidad(cerrados)
	return nil
}

// resumirDefectos cuenta los defectos del proyecto (HU-27) y lista los abiertos, los que no están resueltos, por
// severidad y en el orden en que se registraron (RN-30.5).
func (r *ReporteProyecto) resumirDefectos(ds []defect.Defecto) {
	entradas := make([]metrics.DefectoSprint, 0, len(ds))
	for _, def := range ds {
		entradas = append(entradas, metrics.DefectoSprint{SprintDeteccion: string(def.SprintDeteccion), SprintResolucion: string(def.SprintResolucion)})
		if def.Estado != defect.Resuelto {
			r.DefectosAbiertos = append(r.DefectosAbiertos, FilaDefecto{
				ID: def.ID, Descripcion: def.Descripcion, Severidad: def.Severidad.String(), Estado: def.Estado.String(), HistoriaID: def.HistoriaID,
			})
		}
	}
	severidad := make(map[domain.DefectoID]defect.Severidad, len(ds))
	for _, def := range ds {
		severidad[def.ID] = def.Severidad
	}
	// Severidad sigue el orden de la spec con iota (Crítica = 1 … Baja = 4, ADR 0003).
	slices.SortStableFunc(r.DefectosAbiertos, func(a, b FilaDefecto) int { return int(severidad[a.ID]) - int(severidad[b.ID]) })
	r.Detectados, r.Resueltos = metrics.TotalDetectados(entradas), metrics.TotalResueltos(entradas)
	r.Abiertos = len(r.DefectosAbiertos)
}
