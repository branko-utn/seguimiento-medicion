package features

import "github.com/cucumber/godog"

// registrarPasosComunesSprints registra los pasos de sprints que usan varias historias (features/README.md).
// Responsable: Federico Sosa. Los escenarios los escriben con y sin "que": "Y existe un sprint …" y
// "Dado que existe un sprint …".
func registrarPasosComunesSprints(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^(?:que )?existe un sprint "([^"]*)" en "([^"]*)"$`, m.existeUnSprintEn)
	sc.Step(`^(?:que )?existe un sprint activo "([^"]*)" en "([^"]*)"$`, m.existeUnSprintActivoEn)
}

func (m *mundo) existeUnSprintEn(nombre, proyecto string) error {
	return godog.ErrPending
}

func (m *mundo) existeUnSprintActivoEn(nombre, proyecto string) error {
	return godog.ErrPending
}
