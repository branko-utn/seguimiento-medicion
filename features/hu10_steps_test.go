package features

import "github.com/cucumber/godog"

// registrarPasosHU10 registra los pasos de "Crear sprint con Sprint Goal" (HU-10).
func registrarPasosHU10(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^creo en "([^"]*)" un sprint del "([^"]*)" al "([^"]*)" con el objetivo "([^"]*)"$`, m.creoUnSprint)
	sc.Step(`^inicio el sprint "([^"]*)"$`, m.inicioElSprint)
	sc.Step(`^que en "([^"]*)" existe un sprint del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintDelAl)
	sc.Step(`^que en "([^"]*)" existe un sprint "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintNombradoDelAl)
	sc.Step(`^que en "([^"]*)" existe un sprint activo "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintActivoDelAl)
	sc.Step(`^el proyecto "([^"]*)" tiene un sprint "([^"]*)" en estado "([^"]*)" con el objetivo "([^"]*)"$`, m.elProyectoTieneUnSprint)
	sc.Step(`^el proyecto "([^"]*)" tiene (\d+) sprints$`, m.elProyectoTieneNSprints)
	sc.Step(`^el sprint "([^"]*)" está en estado "([^"]*)"$`, m.elSprintEstaEnEstado)
}

func (m *mundo) creoUnSprint(proyecto, inicio, fin, goal string) error {
	return godog.ErrPending
}

func (m *mundo) inicioElSprint(nombre string) error {
	return godog.ErrPending
}

func (m *mundo) queEnExisteUnSprintDelAl(proyecto, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) queEnExisteUnSprintNombradoDelAl(proyecto, nombre, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) queEnExisteUnSprintActivoDelAl(proyecto, nombre, inicio, fin string) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoTieneUnSprint(proyecto, nombre, estado, goal string) error {
	return godog.ErrPending
}

func (m *mundo) elProyectoTieneNSprints(proyecto string, n int) error {
	return godog.ErrPending
}

func (m *mundo) elSprintEstaEnEstado(nombre, estado string) error {
	return godog.ErrPending
}
