package features

import "github.com/cucumber/godog"

// registrarPasosHU27 registra los pasos de "Defectos detectados y resueltos" (HU-27).
func registrarPasosHU27(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^los defectos:$`, func(tabla *godog.Table) error {
		return godog.ErrPending
	})
	sc.Step(`^en el sprint "([^"]*)" se detectaron (\d+) defectos y se resolvieron (\d+)$`, func(sprint string, detectados, resueltos int) error {
		return godog.ErrPending
	})
	sc.Step(`^en total se detectaron (\d+) defectos y se resolvieron (\d+)$`, func(detectados, resueltos int) error {
		return godog.ErrPending
	})
}
