package features

import (
	"github.com/cucumber/godog"
)

// registrarPasosHU18 registra los pasos de "Nueva ronda y estimación acordada" (HU-18). Reusa los
// "Dado" de HU-16/HU-17 (votos, revelar) y los compartidos de proyectos, integrantes e historias.
func registrarPasosHU18(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^abro una nueva ronda en la historia "([^"]*)"$`, func(historia string) error {
		return godog.ErrPending
	})

	sc.Step(`^la votación de la historia "([^"]*)" está en la ronda (\d+) sin votos$`, func(historia string, ronda int) error {
		return godog.ErrPending
	})

	sc.Step(`^el historial de la historia "([^"]*)" tiene la ronda (\d+) con los votos:$`, func(historia string, ronda int, tabla *godog.Table) error {
		return godog.ErrPending
	})

	sc.Step(`^registro (\d+) como estimación acordada de la historia "([^"]*)"$`, func(sp int, historia string) error {
		return godog.ErrPending
	})

	sc.Step(`^la historia "([^"]*)" no tiene una votación abierta$`, func(historia string) error {
		return godog.ErrPending
	})

	sc.Step(`^la historia "([^"]*)" no tiene Story Points$`, func(historia string) error {
		return godog.ErrPending
	})
}