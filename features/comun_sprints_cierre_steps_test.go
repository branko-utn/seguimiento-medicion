package features

import "github.com/cucumber/godog"

// registrarPasosCierreSprints registra los pasos compartidos de sprints que le tocan a Branko (features/README.md).
func registrarPasosCierreSprints(sc *godog.ScenarioContext, m *mundo) {}

func (m *mundo) queElSprintEstaCerrado(nombre string) error { return nil }
