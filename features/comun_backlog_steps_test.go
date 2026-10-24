package features

import "github.com/cucumber/godog"

// registrarPasosComunesBacklog registra los pasos de historias que usan varias épicas (features/README.md).
// Responsable: Ignacio Bianchi.
func registrarPasosComunesBacklog(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que el backlog de "([^"]*)" tiene la historia "([^"]*)" "([^"]*)"$`, m.queElBacklogTieneLaHistoria)
	sc.Step(`^la historia "([^"]*)" está en estado "([^"]*)"$`, m.laHistoriaEstaEnEstado)
}

func (m *mundo) queElBacklogTieneLaHistoria(proyecto, historia, titulo string) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaEstaEnEstado(historia, estado string) error {
	return godog.ErrPending
}
