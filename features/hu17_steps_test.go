package features

import (
	"github.com/cucumber/godog"
)

// registrarPasosHU17 registra los pasos de "Revelar votos y detectar diferencias" (HU-17). Reusa los
// "Dado" de HU-16 (votos) y los compartidos de proyectos, integrantes e historias (features/README.md).
func registrarPasosHU17(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^revelo los votos de la historia "([^"]*)"$`, func(historia string) error {
		return godog.ErrPending
	})

	sc.Step(`^que se revelaron los votos de la historia "([^"]*)"$`, func(historia string) error {
		return godog.ErrPending
	})

	sc.Step(`^los votos de la historia "([^"]*)" son:$`, func(historia string, tabla *godog.Table) error {
		return godog.ErrPending
	})

	sc.Step(`^hay consenso en (\d+)$`, func(valor int) error {
		return godog.ErrPending
	})

	sc.Step(`^no hay consenso$`, func() error {
		return godog.ErrPending
	})

	sc.Step(`^el voto más alto es (\d+), de "([^"]*)"$`, func(valor int, de string) error {
		return godog.ErrPending
	})

	sc.Step(`^el voto más bajo es (\d+), de "([^"]*)"$`, func(valor int, de string) error {
		return godog.ErrPending
	})
}
