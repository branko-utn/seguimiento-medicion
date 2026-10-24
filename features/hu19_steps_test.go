package features

import "github.com/cucumber/godog"

// registrarPasosHU19 registra los pasos de "Registrar horas estimadas de una historia" (HU-19).
func registrarPasosHU19(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^estimo la historia "([^"]*)" en (-?\d+(?:\.\d+)?) horas$`, m.estimoLaHistoria)
	sc.Step(`^que la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas estimadas$`, m.queLaHistoriaTieneHorasEstimadas)
	sc.Step(`^la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas estimadas$`, m.laHistoriaTieneHorasEstimadas)
}

func (m *mundo) estimoLaHistoria(historia string, horas float64) error {
	return godog.ErrPending
}

func (m *mundo) queLaHistoriaTieneHorasEstimadas(historia string, horas float64) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaTieneHorasEstimadas(historia string, horas float64) error {
	return godog.ErrPending
}
