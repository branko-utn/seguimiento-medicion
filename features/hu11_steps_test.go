package features

import "github.com/cucumber/godog"

// registrarPasosHU11 registra los pasos de "Asignar historias a un sprint" (HU-11).
func registrarPasosHU11(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^asigno la historia "([^"]*)" al sprint "([^"]*)"$`, m.asignoLaHistoriaAlSprint)
	sc.Step(`^el sprint "([^"]*)" tiene las historias "([^"]*)"$`, m.elSprintTieneLasHistorias)
	sc.Step(`^el sprint "([^"]*)" no tiene historias$`, m.elSprintNoTieneHistorias)
	sc.Step(`^el sprint "([^"]*)" tiene (\d+) Story Points planificados$`, m.elSprintTieneSPPlanificados)
}

func (m *mundo) asignoLaHistoriaAlSprint(historia, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) elSprintTieneLasHistorias(sprint, historias string) error {
	return godog.ErrPending
}

func (m *mundo) elSprintNoTieneHistorias(sprint string) error {
	return godog.ErrPending
}

func (m *mundo) elSprintTieneSPPlanificados(sprint string, sp int) error {
	return godog.ErrPending
}
