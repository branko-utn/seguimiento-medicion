package features

import "github.com/cucumber/godog"

// registrarPasosHU02 registra los pasos de "Modificar proyecto" (HU-02).
func registrarPasosHU02(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^modifico el proyecto "([^"]*)" con el nombre "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.modificoElProyecto)
	sc.Step(`^que el proyecto "([^"]*)" está finalizado$`, m.queElProyectoEstaFinalizado)
}

func (m *mundo) modificoElProyecto(proyecto, nombre, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) queElProyectoEstaFinalizado(proyecto string) error {
	return godog.ErrPending
}
