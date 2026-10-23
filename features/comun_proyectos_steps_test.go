package features

import "github.com/cucumber/godog"

// registrarPasosComunesProyectos registra los pasos de proyectos que usan varias historias (features/README.md).
// Responsable: Federico Sosa.
func registrarPasosComunesProyectos(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que existe un proyecto "([^"]*)"$`, m.queExisteUnProyecto)
	sc.Step(`^que existe un proyecto "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queExisteUnProyectoDelAl)
}

func (m *mundo) queExisteUnProyecto(nombre string) error {
	return godog.ErrPending
}

func (m *mundo) queExisteUnProyectoDelAl(nombre, inicio, fin string) error {
	return godog.ErrPending
}
