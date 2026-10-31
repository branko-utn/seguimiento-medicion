package pdf

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/report"
	"github.com/go-pdf/fpdf"
)

// documento es un PDF A4 vertical con los bloques que usan los dos informes: títulos, párrafos y tablas.
type documento struct {
	f *fpdf.Fpdf
	// tr pasa el texto de UTF-8 a la codificación de las fuentes estándar del PDF (cp1252), para que los acentos y
	// las eñes se vean bien (RN-31.3).
	tr func(string) string
}

func nuevoDocumento() *documento {
	f := fpdf.New("P", "mm", "A4", "")
	f.SetCompression(comprimir)
	f.SetMargins(margen, margen, margen)
	f.SetAutoPageBreak(true, margen)
	f.AddPage()
	return &documento{f: f, tr: f.UnicodeTranslatorFromDescriptor("")}
}

// ancho es el ancho útil de la hoja, entre los márgenes.
func (d *documento) ancho() float64 {
	w, _ := d.f.GetPageSize()
	return w - 2*margen
}

func (d *documento) titulo(s string) {
	d.f.SetFont("Helvetica", "B", 18)
	d.f.MultiCell(0, 9, d.tr(s), "", "L", false)
}

func (d *documento) subtitulo(s string) {
	d.f.SetFont("Helvetica", "", 11)
	d.f.MultiCell(0, 6, d.tr(s), "", "L", false)
}

// generado escribe la fecha en que se armó el informe.
func (d *documento) generado(hoy time.Time) {
	d.f.SetFont("Helvetica", "I", 9)
	d.f.SetTextColor(90, 90, 90)
	d.f.MultiCell(0, 5, d.tr("Generado el "+hoy.Format(formatoDia)), "", "L", false)
	d.f.SetTextColor(0, 0, 0)
	d.f.Ln(2)
}

// destacado es un aviso con fondo, como la marca de reporte parcial.
func (d *documento) destacado(s string) {
	d.f.SetFont("Helvetica", "B", 10)
	d.f.SetFillColor(254, 243, 199)
	d.f.MultiCell(0, 7, d.tr(s), "", "L", true)
	d.f.Ln(2)
}

func (d *documento) seccion(s string) {
	d.f.Ln(3)
	d.f.SetFont("Helvetica", "B", 13)
	d.f.MultiCell(0, 8, d.tr(s), "", "L", false)
}

func (d *documento) subseccion(s string) {
	d.f.Ln(2)
	d.f.SetFont("Helvetica", "B", 11)
	d.f.MultiCell(0, 7, d.tr(s), "", "L", false)
}

func (d *documento) parrafo(s string) {
	d.f.SetFont("Helvetica", "", 10)
	d.f.MultiCell(0, 6, d.tr(s), "", "L", false)
}

func (d *documento) nota(s string) {
	d.f.SetFont("Helvetica", "I", 8)
	d.f.MultiCell(0, 5, d.tr(s), "", "L", false)
}

// tabla dibuja una tabla con bordes. Cada celda corta su texto en varias líneas si no entra, y si una fila no
// entra en lo que queda de la página, la tabla sigue en la siguiente con el encabezado repetido (RN-31.4).
func (d *documento) tabla(anchos []float64, encabezado []string, filas [][]string) {
	d.fila(anchos, encabezado, true)
	for _, f := range filas {
		if d.f.GetY()+d.altoFila(anchos, f, false) > d.limite() {
			d.f.AddPage()
			d.fila(anchos, encabezado, true)
		}
		d.fila(anchos, f, false)
	}
	d.f.Ln(2)
}

// limite es la altura donde termina el área útil de la página.
func (d *documento) limite() float64 {
	_, h := d.f.GetPageSize()
	return h - margen
}

func (d *documento) fuenteFila(encabezado bool) {
	if encabezado {
		d.f.SetFont("Helvetica", "B", 7)
		return
	}
	d.f.SetFont("Helvetica", "", 9)
}

// lineas corta el texto de una celda en las líneas que entran en su ancho, ya traducidas para el PDF.
func (d *documento) lineas(texto string, ancho float64) []string {
	partes := d.f.SplitText(texto, ancho-2) // 1 mm de aire a cada lado
	if len(partes) == 0 {
		partes = []string{""}
	}
	for i, p := range partes {
		partes[i] = d.tr(p)
	}
	return partes
}

func (d *documento) altoFila(anchos []float64, celdas []string, encabezado bool) float64 {
	d.fuenteFila(encabezado)
	maximo := 1
	for i, c := range celdas {
		if n := len(d.lineas(c, anchos[i])); n > maximo {
			maximo = n
		}
	}
	return float64(maximo) * altoLinea
}

// fila dibuja una fila: un rectángulo por celda, del alto de la celda con más líneas, y sus líneas de texto.
func (d *documento) fila(anchos []float64, celdas []string, encabezado bool) {
	alto := d.altoFila(anchos, celdas, encabezado)
	x0, y := d.f.GetX(), d.f.GetY()
	estilo := "D"
	if encabezado {
		d.f.SetFillColor(229, 231, 235)
		estilo = "FD"
	}
	x := x0
	for i, c := range celdas {
		d.f.Rect(x, y, anchos[i], alto, estilo)
		for j, l := range d.lineas(c, anchos[i]) {
			d.f.SetXY(x+1, y+float64(j)*altoLinea)
			d.f.CellFormat(anchos[i]-2, altoLinea, l, "", 0, "L", false, 0, "")
		}
		x += anchos[i]
	}
	d.f.SetXY(x0, y+alto)
}

// defectos es la tabla de defectos que comparten los dos informes.
func (d *documento) defectos(ds []report.FilaDefecto) {
	var filas [][]string
	for _, def := range ds {
		filas = append(filas, []string{string(def.ID), def.Descripcion, def.Severidad, def.Estado, string(def.HistoriaID)})
	}
	d.tabla([]float64{14, 80, 25, 30, 31}, []string{"ID", "Descripción", "Severidad", "Estado", "Historia"}, filas)
}

func (d *documento) escribir(w io.Writer) error {
	return d.f.Output(w)
}

// unDecimal escribe un número con un decimal y coma, como la pantalla (RN-31.2).
func unDecimal(v float64) string {
	return strings.Replace(strconv.FormatFloat(v, 'f', 1, 64), ".", ",", 1)
}

// horas es la línea de horas estimadas, reales y desviación, igual que en la pantalla.
func horas(estimadas, reales, desviacion float64, calculable bool) string {
	s := fmt.Sprintf("%s horas estimadas · %s horas reales · ", unDecimal(estimadas), unDecimal(reales))
	if calculable {
		return s + "desviación de " + unDecimal(desviacion) + " %"
	}
	return s + "desviación no calculable (sin horas estimadas)"
}

// velocidad es la línea de velocidad del informe del proyecto, con la cantidad de sprints cerrados.
func velocidad(v float64, cerrados int) string {
	s := "Velocidad: " + unDecimal(v) + " Story Points por sprint"
	switch cerrados {
	case 0:
		return s + " (sin sprints cerrados todavía)"
	case 1:
		return s + " (1 sprint cerrado)"
	}
	return fmt.Sprintf("%s (%d sprints cerrados)", s, cerrados)
}

// contar escribe "1 detectado" o "3 detectados", como la pantalla.
func contar(n int, palabra string) string {
	if n == 1 {
		return "1 " + palabra
	}
	return fmt.Sprintf("%d %ss", n, palabra)
}
