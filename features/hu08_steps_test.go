package features

import "github.com/cucumber/godog"

// registrarPasosHU08 registra los pasos de "Cambiar el estado de una historia" (HU-08). El "Dado" de asignación
// es compartido de Sosa, los de estado y completada son compartidos de backlog y de Branko (features/README.md).
func registrarPasosHU08(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^cambio el estado de "([^"]*)" a "([^"]*)"$`, m.cambioElEstado)
}

func (m *mundo) cambioElEstado(historia, estado string) error {
	return godog.ErrPending
}
