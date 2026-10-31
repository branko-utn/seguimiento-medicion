// Package pdf genera el informe de un proyecto y el de un sprint en un archivo PDF (HU-31), con la biblioteca
// github.com/go-pdf/fpdf (ADR 0001). Arma el PDF con el mismo report.ReporteProyecto o report.ReporteSprint que la
// pantalla, sin recalcular nada (RN-31.1), y no accede a servicios (ADR 0002).
package pdf

import (
	"io"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// comprimir indica si el contenido del PDF se comprime. Los tests lo apagan para buscar los textos.
var comprimir = true

// Proyecto escribe en w el informe del proyecto en PDF (HU-30 y HU-31).
func Proyecto(w io.Writer, r report.ReporteProyecto, hoy time.Time) error {
	return nil
}

// Sprint escribe en w el informe del sprint en PDF (HU-29 y HU-31).
func Sprint(w io.Writer, r report.ReporteSprint, hoy time.Time) error {
	return nil
}
