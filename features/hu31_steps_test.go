package features

import "github.com/cucumber/godog"

// registrarPasosHU31 registra los pasos de "Exportar reporte a PDF" (HU-31).
func registrarPasosHU31(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^la respuesta es un PDF llamado "([^"]*)"$`, func(archivo string) error {
		return godog.ErrPending
	})
}
