// Package pdf genera el informe de un proyecto y el de un sprint en un archivo PDF (HU-31), con la biblioteca
// github.com/go-pdf/fpdf (ADR 0001). Arma el PDF con el mismo report.ReporteProyecto o report.ReporteSprint que la
// pantalla, sin recalcular nada (RN-31.1), y no accede a servicios (ADR 0002).
package pdf

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// comprimir indica si el contenido del PDF se comprime. Los tests lo apagan para buscar los textos.
var comprimir = true

const (
	margen     = 15.0 // mm, en los cuatro lados de la hoja A4
	altoLinea  = 5.0  // mm, de cada línea de texto en las tablas
	formatoDia = "02/01/2006"
)

// Proyecto escribe en w el informe del proyecto en PDF, con las secciones de la pantalla de HU-30: estado y días
// restantes, resumen de sprints, métricas acumuladas y defectos abiertos.
func Proyecto(w io.Writer, r report.ReporteProyecto, hoy time.Time) error {
	d := nuevoDocumento()
	d.titulo("Reporte del proyecto")
	d.subtitulo(fmt.Sprintf("%s · %s · del %s al %s · %d días restantes",
		r.Nombre, r.Estado, r.Inicio.Format(formatoDia), r.Fin.Format(formatoDia), r.DiasRestantes))
	d.generado(hoy)

	d.seccion("Resumen de sprints")
	cerrados := 0
	if len(r.Sprints) == 0 {
		d.parrafo("Todavía no hay sprints con resultado.")
	} else {
		var filas [][]string
		for _, s := range r.Sprints {
			filas = append(filas, []string{s.Goal, s.Estado, strconv.Itoa(s.SPPlanificados), strconv.Itoa(s.SPCompletados),
				unDecimal(s.Porcentaje) + " %", unDecimal(s.HorasEstimadas), unDecimal(s.HorasReales),
				strconv.Itoa(s.Detectados), strconv.Itoa(s.Resueltos)})
			if !s.Parcial {
				cerrados++
			}
		}
		d.tabla([]float64{34, 17, 20, 21, 20, 18, 16, 18, 16},
			[]string{"Sprint", "Estado", "SP planificados", "SP completados", "% completadas", "Horas estimadas",
				"Horas reales", "Defectos detectados", "Defectos resueltos"}, filas)
	}

	d.seccion("Métricas acumuladas")
	d.parrafo(velocidad(r.Velocidad, cerrados))
	d.parrafo(fmt.Sprintf("%d Story Points completados", r.SPCompletados))
	d.parrafo(fmt.Sprintf("%s %% de historias completadas (%d de %d)", unDecimal(r.Porcentaje), r.Completadas, r.Historias))
	d.parrafo(horas(r.HorasEstimadas, r.HorasReales, r.Desviacion, r.DesviacionCalculable))
	d.nota("Las horas reales incluyen todo el esfuerzo del proyecto, también el registrado fuera de un sprint.")

	d.seccion("Defectos")
	d.parrafo(fmt.Sprintf("%s · %s · %s", contar(r.Detectados, "detectado"), contar(r.Resueltos, "resuelto"),
		contar(r.Abiertos, "abierto")))
	if len(r.DefectosAbiertos) == 0 {
		d.parrafo("No hay defectos abiertos.")
	} else {
		d.defectos(r.DefectosAbiertos)
	}
	return d.escribir(w)
}

// Sprint escribe en w el informe del sprint en PDF, con las secciones de la pantalla de HU-29: historias
// planificadas, métricas, esfuerzo por integrante y defectos, con la marca "Reporte parcial" si el sprint está activo.
func Sprint(w io.Writer, r report.ReporteSprint, hoy time.Time) error {
	d := nuevoDocumento()
	d.titulo("Reporte del sprint")
	d.subtitulo(fmt.Sprintf("%s · %s · del %s al %s", r.Goal, r.Estado, r.Inicio.Format(formatoDia), r.Fin.Format(formatoDia)))
	d.generado(hoy)
	if r.Parcial {
		d.destacado("Reporte parcial: el sprint sigue activo y el resultado todavía puede cambiar.")
	}

	d.seccion("Historias planificadas")
	if len(r.Historias) == 0 {
		d.parrafo("El sprint no tuvo historias planificadas.")
	} else {
		var filas [][]string
		for _, h := range r.Historias {
			completada := "No"
			if h.Completada {
				completada = "Sí"
			}
			filas = append(filas, []string{string(h.ID), h.Titulo, strconv.Itoa(h.StoryPoints),
				unDecimal(h.HorasEstimadas), unDecimal(h.HorasReales), completada})
		}
		d.tabla([]float64{16, 70, 22, 24, 24, 24},
			[]string{"ID", "Historia", "Story Points", "Horas estimadas", "Horas reales", "Completada"}, filas)
	}

	d.seccion("Métricas")
	d.parrafo(fmt.Sprintf("%d Story Points planificados", r.SPPlanificados))
	d.parrafo(fmt.Sprintf("%d Story Points completados", r.SPCompletados))
	d.parrafo(fmt.Sprintf("%s %% de historias completadas", unDecimal(r.Porcentaje)))
	d.parrafo(horas(r.HorasEstimadas, r.HorasReales, r.Desviacion, r.DesviacionCalculable))

	d.seccion("Esfuerzo registrado")
	if len(r.Registros) == 0 {
		d.parrafo("No se registró esfuerzo en las fechas del sprint.")
	} else {
		var filas [][]string
		for _, reg := range r.Registros {
			filas = append(filas, []string{reg.Fecha.Format(formatoDia), reg.Integrante, string(reg.HistoriaID),
				reg.Actividad, unDecimal(reg.Horas)})
		}
		d.tabla([]float64{25, 45, 25, 55, 30}, []string{"Fecha", "Integrante", "Historia", "Actividad", "Horas"}, filas)
		d.subseccion("Total por integrante")
		filas = nil
		for _, t := range r.PorIntegrante {
			filas = append(filas, []string{t.Integrante, unDecimal(t.Horas)})
		}
		d.tabla([]float64{90, 40}, []string{"Integrante", "Horas"}, filas)
	}

	d.seccion("Defectos")
	d.parrafo(fmt.Sprintf("%s en el sprint · %s en el sprint", contar(r.Detectados, "detectado"), contar(r.Resueltos, "resuelto")))
	if len(r.Defectos) == 0 {
		d.parrafo("No se registraron defectos en este sprint.")
	} else {
		d.defectos(r.Defectos)
	}
	return d.escribir(w)
}
