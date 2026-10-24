package features

import "github.com/cucumber/godog"

// registrarPasosHU07 registra los pasos de "Criterios de aceptación de una historia" (HU-07).
func registrarPasosHU07(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^agrego a "([^"]*)" el criterio "([^"]*)"$`, m.agregoElCriterio)
	sc.Step(`^quito de "([^"]*)" el criterio (\d+)$`, m.quitoElCriterio)
	sc.Step(`^que la historia "([^"]*)" tiene los criterios "([^"]*)"$`, m.queLaHistoriaTieneLosCriterios)
	sc.Step(`^la historia "([^"]*)" tiene los criterios:$`, m.laHistoriaTieneLosCriterios)
	sc.Step(`^la historia "([^"]*)" no tiene criterios$`, m.laHistoriaNoTieneCriterios)
}

func (m *mundo) agregoElCriterio(historia, criterio string) error {
	return godog.ErrPending
}

func (m *mundo) quitoElCriterio(historia string, numero int) error {
	return godog.ErrPending
}

func (m *mundo) queLaHistoriaTieneLosCriterios(historia, criterios string) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaTieneLosCriterios(historia string, tabla *godog.Table) error {
	return godog.ErrPending
}

func (m *mundo) laHistoriaNoTieneCriterios(historia string) error {
	return godog.ErrPending
}
