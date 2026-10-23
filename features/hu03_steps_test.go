package features

import "github.com/cucumber/godog"

// registrarPasosHU03 registra los pasos de "Registrar integrantes" (HU-03).
func registrarPasosHU03(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^registro a "([^"]*)" como "([^"]*)" en el proyecto "([^"]*)"$`, m.registroA)
	sc.Step(`^que el proyecto "([^"]*)" tiene los integrantes:$`, m.queElProyectoTieneLosIntegrantes)
	sc.Step(`^el proyecto "([^"]*)" tiene los integrantes:$`, m.elProyectoTieneLosIntegrantes)
	sc.Step(`^el proyecto "([^"]*)" tiene (\d+) integrantes$`, m.elProyectoTieneNIntegrantes)
}

func (m *mundo) registroA(nombre, rol, proyecto string) error {
	return godog.ErrPending
}

func (m *mundo) queElProyectoTieneLosIntegrantes(proyecto string, tabla *godog.Table) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoTieneLosIntegrantes(proyecto string, tabla *godog.Table) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoTieneNIntegrantes(proyecto string, n int) error {
	return godog.ErrPending
}
