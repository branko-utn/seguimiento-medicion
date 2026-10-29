package features

import "github.com/cucumber/godog"

// registrarPasosHU09 registra los pasos de "Ver el backlog ordenado por prioridad" (HU-09).
func registrarPasosHU09(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^veo el backlog de "([^"]*)"$`, func(proyecto string) error {
		return godog.ErrPending
	})
	sc.Step(`^veo el backlog de "([^"]*)" con estado "([^"]*)"$`, func(proyecto, estado string) error {
		return godog.ErrPending
	})
	sc.Step(`^el backlog se muestra en el orden "([^"]*)"$`, func(orden string) error {
		return godog.ErrPending
	})
	sc.Step(`^el backlog se muestra vacío$`, func() error {
		return godog.ErrPending
	})
	sc.Step(`^ordeno las historias:$`, func(tabla *godog.Table) error {
		return godog.ErrPending
	})
	sc.Step(`^quedan en el orden "([^"]*)"$`, func(orden string) error {
		return godog.ErrPending
	})
}
