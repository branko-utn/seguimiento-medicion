package features

import "github.com/cucumber/godog"

// registrarPasosHU06 registra los pasos de "Editar y eliminar ítem del backlog" (HU-06).
func registrarPasosHU06(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^edito la historia "([^"]*)" con título "([^"]*)", descripción "([^"]*)" y prioridad "([^"]*)"$`, m.editoLaHistoria)
	sc.Step(`^la historia "([^"]*)" tiene el título "([^"]*)" y prioridad "([^"]*)"$`, m.laHistoriaTieneElTituloYPrioridad)
	sc.Step(`^elimino la historia "([^"]*)"$`, m.eliminoLaHistoria)
	sc.Step(`^el backlog de "([^"]*)" no tiene la historia "([^"]*)"$`, m.elBacklogNoTieneLaHistoria)
	sc.Step(`^el backlog de "([^"]*)" tiene (\d+) historias$`, m.elBacklogTieneNHistorias)
}

func (m *mundo) editoLaHistoria(historia, titulo, descripcion, prioridad string) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaTieneElTituloYPrioridad(historia, titulo, prioridad string) error {
	return godog.ErrPending
}

func (m *mundo) eliminoLaHistoria(historia string) error {
	return godog.ErrPending
}

func (m *mundo) elBacklogNoTieneLaHistoria(proyecto, historia string) error {
	return godog.ErrPending
}

func (m *mundo) elBacklogTieneNHistorias(proyecto string, n int) error {
	return godog.ErrPending
}
