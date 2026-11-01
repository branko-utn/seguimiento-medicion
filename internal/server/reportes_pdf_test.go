package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests de HU-31: la descarga del informe en PDF. El contenido del PDF lo verifican los tests de internal/report/pdf;
// acá se verifican los encabezados y que el cuerpo sea un PDF (spec HU-31, restricciones).

func TestReportes_pdf(t *testing.T) {
	router := escenarioReportes(t)
	casos := []struct {
		ruta, archivo string
	}{
		{"/reportes/proyecto.pdf?proyecto=P1", "reporte-proyecto-P1.pdf"},
		{"/reportes/sprints/S1/pdf", "reporte-sprint-S1.pdf"},
	}
	for _, c := range casos {
		t.Run(c.archivo, func(t *testing.T) {
			rec := servir(router, httptest.NewRequest(http.MethodGet, c.ruta, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("código = %d, se esperaba 200", rec.Code)
			}
			if tipo := rec.Header().Get("Content-Type"); tipo != "application/pdf" {
				t.Errorf("Content-Type = %q, se esperaba application/pdf", tipo)
			}
			if disp, want := rec.Header().Get("Content-Disposition"), `attachment; filename="`+c.archivo+`"`; disp != want {
				t.Errorf("Content-Disposition = %q, se esperaba %q", disp, want)
			}
			if !strings.HasPrefix(rec.Body.String(), "%PDF-") {
				t.Error("el cuerpo no es un PDF")
			}
		})
	}
}

// Sin reporte no hay PDF: los mismos errores que la pantalla (spec HU-31, condiciones de error).
func TestReportes_pdf_errores(t *testing.T) {
	router := escenarioReportes(t)
	casos := []struct {
		ruta    string
		codigo  int
		mensaje string
	}{
		{"/reportes/proyecto.pdf?proyecto=P9", http.StatusNotFound, "el proyecto no existe"},
		{"/reportes/sprints/S9/pdf", http.StatusNotFound, "el sprint no existe"},
		{"/reportes/sprints/S2/pdf", http.StatusBadRequest, "el sprint todavía no empezó"},
	}
	for _, c := range casos {
		rec := servir(router, httptest.NewRequest(http.MethodGet, c.ruta, nil))
		if rec.Code != c.codigo || !strings.Contains(rec.Body.String(), c.mensaje) {
			t.Errorf("GET %s: código = %d, se esperaba %d con %q", c.ruta, rec.Code, c.codigo, c.mensaje)
		}
		if strings.HasPrefix(rec.Body.String(), "%PDF-") {
			t.Errorf("GET %s: con error no se descarga un PDF", c.ruta)
		}
	}
}

// Los reportes ofrecen "Descargar PDF" y siguen ofreciendo imprimir.
func TestReportes_ofrecenDescargarPDF(t *testing.T) {
	router := escenarioReportes(t)
	casos := []struct{ pagina, enlace string }{
		{"/reportes/proyecto?proyecto=P1", `href="/reportes/proyecto.pdf?proyecto=P1"`},
		{"/reportes/sprints/S1", `href="/reportes/sprints/S1/pdf"`},
	}
	for _, c := range casos {
		cuerpo := servir(router, httptest.NewRequest(http.MethodGet, c.pagina, nil)).Body.String()
		contieneTodo(t, cuerpo, c.enlace, "Descargar PDF", "Imprimir o guardar como PDF")
	}
}
