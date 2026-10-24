package features

import "github.com/cucumber/godog"

// registrarPasosHU20 registra los pasos de "Registrar esfuerzo trabajado" (HU-20).
func registrarPasosHU20(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^"([^"]*)" registra (-?\d+(?:\.\d+)?) horas de "([^"]*)" en "([^"]*)" el "([^"]*)"$`, m.registraHoras)
	sc.Step(`^que "([^"]*)" registró (\d+(?:\.\d+)?) horas de "([^"]*)" en "([^"]*)" el "([^"]*)"$`, m.queRegistroHoras)
	sc.Step(`^la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas reales$`, m.laHistoriaTieneHorasReales)
	sc.Step(`^"([^"]*)" tiene (\d+(?:\.\d+)?) horas registradas el "([^"]*)"$`, m.tieneHorasRegistradasEl)
}

func (m *mundo) registraHoras(integrante string, horas float64, actividad, historia, fecha string) error {
	return godog.ErrPending
}

func (m *mundo) queRegistroHoras(integrante string, horas float64, actividad, historia, fecha string) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaTieneHorasReales(historia string, horas float64) error {
	return godog.ErrPending
}

func (m *mundo) tieneHorasRegistradasEl(integrante string, horas float64, fecha string) error {
	return godog.ErrPending
}
