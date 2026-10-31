package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
)

// Los tests generan el PDF sin compresión para poder buscar los textos (spec HU-31, restricciones).
func init() { comprimir = false }

var hoy = time.Date(2026, 10, 31, 0, 0, 0, 0, time.Local)

func dia(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }

// enPDF escribe un texto como queda dentro del PDF: las fuentes estándar usan la codificación cp1252, así que las
// letras con acento y la eñe ocupan un byte (RN-31.3).
func enPDF(s string) string {
	return strings.NewReplacer("á", "\xe1", "é", "\xe9", "í", "\xed", "ó", "\xf3", "ú", "\xfa", "ñ", "\xf1",
		"Á", "\xc1", "É", "\xc9", "Í", "\xcd", "Ó", "\xd3", "Ú", "\xda", "Ñ", "\xd1").Replace(s)
}

// contiene verifica que el PDF tenga cada texto, ya codificado como en el PDF.
func contiene(t *testing.T, doc []byte, textos ...string) {
	t.Helper()
	for _, s := range textos {
		if !bytes.Contains(doc, []byte(enPDF(s))) {
			t.Errorf("el PDF no contiene %q", s)
		}
	}
}

// generar arma el PDF con la función dada y verifica que sea un PDF.
func generar(t *testing.T, f func(*bytes.Buffer) error) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := f(&buf); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF-")) {
		t.Fatalf("la salida no es un PDF: empieza con %q", buf.Bytes()[:min(8, buf.Len())])
	}
	return buf.Bytes()
}

// paginas cuenta las páginas del PDF.
func paginas(doc []byte) int {
	return bytes.Count(doc, []byte("/Type /Page\n"))
}

func sprintCerrado() report.ReporteSprint {
	return report.ReporteSprint{
		SprintID: "S1", Goal: "Primer MVP", Inicio: dia(21), Fin: dia(24), Estado: "Cerrado",
		Historias: []report.FilaHistoria{
			{ID: "HU-1", Titulo: "Crear proyecto", StoryPoints: 3, HorasEstimadas: 4.5, HorasReales: 5, Completada: true},
			{ID: "HU-2", Titulo: "Registrar defecto", StoryPoints: 2, HorasEstimadas: 2, HorasReales: 0},
		},
		SPPlanificados: 5, SPCompletados: 3, Porcentaje: 50, HorasEstimadas: 6.5, HorasReales: 5,
		Desviacion: -23.1, DesviacionCalculable: true,
		Registros: []report.FilaRegistro{
			{Fecha: dia(22), Integrante: "Ana Pérez", HistoriaID: "HU-1", Actividad: "Desarrollo", Horas: 5},
		},
		PorIntegrante: []report.TotalIntegrante{{Integrante: "Ana Pérez", Horas: 5}},
		Defectos: []report.FilaDefecto{
			{ID: "D1", Descripcion: "No guarda", Severidad: "Crítica", Estado: "Abierto", HistoriaID: "HU-1"},
		},
		Detectados: 1,
	}
}

func TestSprint_contenido(t *testing.T) {
	doc := generar(t, func(b *bytes.Buffer) error { return Sprint(b, sprintCerrado(), hoy) })

	// RN-31.1: las mismas secciones que la pantalla, con los datos del reporte.
	contiene(t, doc, "Reporte del sprint", "Primer MVP", "Cerrado", "del 21/10/2026 al 24/10/2026",
		"Generado el 31/10/2026",
		"Historias planificadas", "HU-1", "Crear proyecto", "Registrar defecto",
		"Métricas", "5 Story Points planificados", "3 Story Points completados", "50,0 % de historias completadas",
		"Esfuerzo registrado", "Ana Pérez", "Desarrollo",
		"Defectos", "No guarda", "Crítica", "Abierto")
	// RN-31.2: coma decimal.
	contiene(t, doc, "4,5", "6,5 horas estimadas", "5,0 horas reales", "desviación de -23,1 %")
	if bytes.Contains(doc, []byte(enPDF("Reporte parcial"))) {
		t.Error("un sprint cerrado no es un reporte parcial")
	}
}

func TestSprint_activoEsParcial(t *testing.T) {
	r := sprintCerrado()
	r.Estado, r.Parcial = "Activo", true

	doc := generar(t, func(b *bytes.Buffer) error { return Sprint(b, r, hoy) })

	contiene(t, doc, "Reporte parcial")
}

func TestSprint_sinDatos(t *testing.T) {
	r := report.ReporteSprint{SprintID: "S2", Goal: "Vacío", Inicio: dia(25), Fin: dia(28), Estado: "Activo", Parcial: true}

	doc := generar(t, func(b *bytes.Buffer) error { return Sprint(b, r, hoy) })

	contiene(t, doc, "El sprint no tuvo historias planificadas.", "No se registró esfuerzo en las fechas del sprint.",
		"No se registraron defectos en este sprint.", "desviación no calculable")
}

func TestProyecto_contenido(t *testing.T) {
	r := report.ReporteProyecto{
		ProyectoID: "P1", Nombre: "Demo", Inicio: dia(21), Fin: time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local),
		Estado: "En curso", DiasRestantes: 2,
		Sprints:   []report.ReporteSprint{sprintCerrado()},
		Velocidad: 12.5, SPCompletados: 3, Completadas: 1, Historias: 2, Porcentaje: 50,
		HorasEstimadas: 6.5, HorasReales: 7.5, Desviacion: 11.5, DesviacionCalculable: true,
		Detectados: 1, Abiertos: 1,
		DefectosAbiertos: []report.FilaDefecto{
			{ID: "D1", Descripcion: "No guarda", Severidad: "Crítica", Estado: "Abierto", HistoriaID: "HU-1"},
		},
	}

	doc := generar(t, func(b *bytes.Buffer) error { return Proyecto(b, r, hoy) })

	contiene(t, doc, "Reporte del proyecto", "Demo", "En curso", "2 días restantes", "Generado el 31/10/2026",
		"Resumen de sprints", "Primer MVP", "50,0 %",
		"Métricas acumuladas", "Velocidad: 12,5 Story Points por sprint", "3 Story Points completados",
		"50,0 % de historias completadas (1 de 2)", "7,5 horas reales", "desviación de 11,5 %",
		"Defectos", "1 detectado", "1 abierto", "No guarda")
}

func TestProyecto_sinSprints(t *testing.T) {
	r := report.ReporteProyecto{ProyectoID: "P1", Nombre: "Demo", Inicio: dia(21), Fin: dia(31), Estado: "En curso"}

	doc := generar(t, func(b *bytes.Buffer) error { return Proyecto(b, r, hoy) })

	contiene(t, doc, "Todavía no hay sprints con resultado.", "No hay defectos abiertos.")
}

// RN-31.4: una tabla que no entra en una página sigue en la siguiente, repitiendo el encabezado.
func TestSprint_tablaLargaRepiteElEncabezado(t *testing.T) {
	r := sprintCerrado()
	r.Historias = nil
	for i := 1; i <= 80; i++ {
		r.Historias = append(r.Historias, report.FilaHistoria{ID: domain.HistoriaID(fmt.Sprintf("HU-%d", i)), Titulo: "Historia", StoryPoints: 1})
	}

	doc := generar(t, func(b *bytes.Buffer) error { return Sprint(b, r, hoy) })

	n := paginas(doc)
	if n < 2 {
		t.Fatalf("el PDF tiene %d página, se esperaban al menos 2", n)
	}
	if enc := bytes.Count(doc, []byte("(Story Points)")); enc < 2 {
		t.Errorf("el encabezado de la tabla aparece %d vez, se esperaba en cada página", enc)
	}
	contiene(t, doc, "HU-80")
}

// Caso límite: un título largo se corta en varias líneas dentro de su celda, sin perder texto.
func TestSprint_tituloLargoEnVariasLineas(t *testing.T) {
	r := sprintCerrado()
	largo := strings.Repeat("Registrar el esfuerzo de cada integrante ", 4)
	r.Historias = []report.FilaHistoria{{ID: "HU-1", Titulo: largo, StoryPoints: 3}}

	doc := generar(t, func(b *bytes.Buffer) error { return Sprint(b, r, hoy) })

	if bytes.Contains(doc, []byte("("+strings.TrimSpace(largo)+")")) {
		t.Error("el título largo quedó en una sola línea")
	}
	if n := bytes.Count(doc, []byte("Registrar el esfuerzo")); n != 4 {
		t.Errorf("el título aparece partido de forma que se pierde texto: %d repeticiones, se esperaban 4", n)
	}
}
