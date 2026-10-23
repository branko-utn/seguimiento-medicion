package features

import "github.com/cucumber/godog"

// registrarPasosHU01 registra los pasos de "Crear proyecto" (HU-01).
func registrarPasosHU01(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^creo el proyecto "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.creoElProyecto)
	sc.Step(`^creo un proyecto con un nombre de (\d+) caracteres del "([^"]*)" al "([^"]*)"$`, m.creoUnProyectoConNombreDeLargo)
	sc.Step(`^el proyecto "([^"]*)" existe$`, m.elProyectoExiste)
	sc.Step(`^el proyecto "([^"]*)" está en estado "([^"]*)"$`, m.elProyectoEstaEnEstado)
	sc.Step(`^el proyecto se crea$`, m.elProyectoSeCrea)
}

func (m *mundo) creoElProyecto(nombre, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) creoUnProyectoConNombreDeLargo(largo int, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoExiste(nombre string) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoEstaEnEstado(nombre, estado string) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoSeCrea() error {
	return godog.ErrPending
}
