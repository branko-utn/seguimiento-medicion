package features

import "github.com/cucumber/godog"

// registrarPasosHU22 registra los pasos de "Registrar defecto" (HU-22).
func registrarPasosHU22(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" detectado en "([^"]*)"$`, m.registroElDefecto)
	sc.Step(`^existe el defecto "([^"]*)" en estado "([^"]*)"$`, m.existeElDefectoEnEstado)
	sc.Step(`^el defecto "([^"]*)" no tiene sprint de resolución$`, m.elDefectoNoTieneSprintDeResolucion)
}

func (m *mundo) registroElDefecto(descripcion, severidad, historia, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) existeElDefectoEnEstado(defecto, estado string) error {
	return godog.ErrPending
}

func (m *mundo) elDefectoNoTieneSprintDeResolucion(defecto string) error {
	return godog.ErrPending
}
