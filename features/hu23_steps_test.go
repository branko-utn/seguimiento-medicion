package features

import "github.com/cucumber/godog"

// registrarPasosHU23 registra los pasos de "Cambiar estado de un defecto y registrar su resolución" (HU-23).
func registrarPasosHU23(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que existe el defecto "([^"]*)" en "([^"]*)" detectado en "([^"]*)"$`, m.queExisteElDefectoDetectadoEn)
	sc.Step(`^que el defecto "([^"]*)" está en estado "([^"]*)"$`, m.queElDefectoEstaEnEstado)
	sc.Step(`^que el defecto "([^"]*)" se resolvió en el sprint "([^"]*)"$`, m.queElDefectoSeResolvioEn)
	sc.Step(`^cambio el estado del defecto "([^"]*)" a "([^"]*)"$`, m.cambioElEstadoDelDefecto)
	sc.Step(`^resuelvo el defecto "([^"]*)" en el sprint "([^"]*)"$`, m.resuelvoElDefecto)
	sc.Step(`^el defecto "([^"]*)" se resolvió en el sprint "([^"]*)"$`, m.elDefectoSeResolvioEn)
}

func (m *mundo) queExisteElDefectoDetectadoEn(defecto, historia, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) queElDefectoEstaEnEstado(defecto, estado string) error {
	return godog.ErrPending
}

func (m *mundo) queElDefectoSeResolvioEn(defecto, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) cambioElEstadoDelDefecto(defecto, estado string) error {
	return godog.ErrPending
}

func (m *mundo) resuelvoElDefecto(defecto, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) elDefectoSeResolvioEn(defecto, sprint string) error {
	return godog.ErrPending
}
