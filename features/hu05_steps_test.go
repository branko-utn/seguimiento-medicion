package features

import "github.com/cucumber/godog"

// registrarPasosHU05 registra los pasos de "Crear ítem del Product Backlog" (HU-05).
func registrarPasosHU05(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^agrego al proyecto "([^"]*)" la historia "([^"]*)" con prioridad "([^"]*)"$`, m.agregoLaHistoria)
	sc.Step(`^que el backlog de "([^"]*)" tiene (\d+) historias$`, m.queElBacklogTieneNHistorias)
	sc.Step(`^el backlog de "([^"]*)" tiene la historia "([^"]*)" "([^"]*)"$`, m.elBacklogTieneLaHistoria)
	sc.Step(`^la historia "([^"]*)" no tiene Story Points$`, m.laHistoriaNoTieneStoryPoints)
}

func (m *mundo) agregoLaHistoria(proyecto, titulo, prioridad string) error {
	return godog.ErrPending
}

func (m *mundo) queElBacklogTieneNHistorias(proyecto string, n int) error {
	return godog.ErrPending
}

func (m *mundo) elBacklogTieneLaHistoria(proyecto, historia, titulo string) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaNoTieneStoryPoints(historia string) error {
	return godog.ErrPending
}
