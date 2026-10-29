package features

import "github.com/cucumber/godog"

// registrarPasosHU21 registra los pasos de "Consultar el esfuerzo registrado" (HU-21).
func registrarPasosHU21(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^consulto el esfuerzo registrado en "([^"]*)"$`, func(proyecto string) error {
		return godog.ErrPending
	})
	sc.Step(`^los totales por historia son:$`, func(tabla *godog.Table) error {
		return godog.ErrPending
	})
	sc.Step(`^los totales por integrante son:$`, func(tabla *godog.Table) error {
		return godog.ErrPending
	})
	sc.Step(`^el esfuerzo total del proyecto es (\d+(?:\.\d+)?) horas$`, func(horas float64) error {
		return godog.ErrPending
	})
	sc.Step(`^los registros son:$`, func(tabla *godog.Table) error {
		return godog.ErrPending
	})
	sc.Step(`^no hay esfuerzo registrado$`, func() error {
		return godog.ErrPending
	})
}
