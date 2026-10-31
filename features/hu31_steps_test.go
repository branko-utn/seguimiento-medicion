package features

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// registrarPasosHU31 registra los pasos de "Exportar reporte a PDF" (HU-31).
func registrarPasosHU31(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^la respuesta es un PDF llamado "([^"]*)"$`, respuestaEsPDF)
}

// respuestaEsPDF verifica que la respuesta sea un archivo PDF para descargar con ese nombre ("Entonces").
func respuestaEsPDF(ctx context.Context, archivo string) error {
	rec, err := respuesta(ctx)
	if err != nil {
		return err
	}
	if tipo := rec.Header().Get("Content-Type"); tipo != "application/pdf" {
		return fmt.Errorf("Content-Type = %q, se esperaba application/pdf", tipo)
	}
	if disp, want := rec.Header().Get("Content-Disposition"), `attachment; filename="`+archivo+`"`; disp != want {
		return fmt.Errorf("Content-Disposition = %q, se esperaba %q", disp, want)
	}
	if !strings.HasPrefix(rec.Body.String(), "%PDF-") {
		return fmt.Errorf("el cuerpo de la respuesta no es un PDF")
	}
	return nil
}
