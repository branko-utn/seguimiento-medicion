package features

import (
	"github.com/cucumber/godog"
)

// registrarPasosHU16 registra los pasos de "Votar en Planning Poker con votos ocultos" (HU-16). Los
// "Dado" de proyectos, integrantes e historias son compartidos de Sosa y de backlog (features/README.md).
func registrarPasosHU16(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^se abrió la votación de la historia "([^"]*)"$`, m.seAbrioLaVotacion)
	sc.Step(`^abro la votación de la historia "([^"]*)"$`, m.abroLaVotacion)
	sc.Step(`^"([^"]*)" vota (\d+) en la historia "([^"]*)"$`, m.votaEnLaHistoria)
	sc.Step(`^que "([^"]*)" votó (\d+) en la historia "([^"]*)"$`, m.queVotoEnLaHistoria)
	sc.Step(`^en la votación de "([^"]*)" votaron "([^"]*)"$`, m.enLaVotacionVotaron)
	sc.Step(`^en la votación de "([^"]*)" no votó nadie$`, m.enLaVotacionNoVotoNadie)
	sc.Step(`^consulto los votos de la historia "([^"]*)"$`, m.consultoLosVotos)
}

func (m *mundo) seAbrioLaVotacion(historia string) error {
	return godog.ErrPending
}

func (m *mundo) abroLaVotacion(historia string) error {
	return godog.ErrPending
}

func (m *mundo) votaEnLaHistoria(integrante string, carta int, historia string) error {
	return godog.ErrPending
}

func (m *mundo) queVotoEnLaHistoria(integrante string, carta int, historia string) error {
	return godog.ErrPending
}

func (m *mundo) enLaVotacionVotaron(historia, votantes string) error {
	return godog.ErrPending
}

func (m *mundo) enLaVotacionNoVotoNadie(historia string) error {
	return godog.ErrPending
}

func (m *mundo) consultoLosVotos(historia string) error {
	return godog.ErrPending
}
