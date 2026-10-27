package features

import (
	"fmt"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
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

// seAbrioLaVotacion abre la votación de la historia ("Dado"): si falla, el escenario falla.
func (m *mundo) seAbrioLaVotacion(historia string) error {
	_, err := m.poker.IniciarVotacion(domain.HistoriaID(m.id("historia", historia)))
	return err
}

// abroLaVotacion ejecuta el caso de uso y guarda su resultado ("Cuando"): el error puede ser el esperado.
func (m *mundo) abroLaVotacion(historia string) error {
	_, err := m.poker.IniciarVotacion(domain.HistoriaID(m.id("historia", historia)))
	return m.resultado(err)
}

// votaEnLaHistoria registra el voto y guarda su resultado ("Cuando"): el error puede ser el esperado.
func (m *mundo) votaEnLaHistoria(integrante string, carta int, historia string) error {
	return m.resultado(m.votarEnLaHistoria(integrante, carta, historia))
}

// queVotoEnLaHistoria registra el voto ("Dado"): si falla, el escenario falla.
func (m *mundo) queVotoEnLaHistoria(integrante string, carta int, historia string) error {
	return m.votarEnLaHistoria(integrante, carta, historia)
}

// votarEnLaHistoria ejecuta el caso de uso de votar, traduciendo los nombres del escenario a IDs
// (los "Dado" de proyectos los registraron con m.recordar).
func (m *mundo) votarEnLaHistoria(integrante string, carta int, historia string) error {
	return m.poker.Votar(
		domain.HistoriaID(m.id("historia", historia)),
		domain.IntegranteID(m.id("integrante", integrante)),
		carta,
	)
}

// enLaVotacionVotaron verifica quiénes votaron, en el orden en que votaron ("Entonces"). La consulta
// no trae los valores de los votos: mientras la votación está abierta solo se ve quién votó (RN-16.4).
func (m *mundo) enLaVotacionVotaron(historia, votantes string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	got, err := m.poker.QuienesVotaron(domain.HistoriaID(m.id("historia", historia)))
	if err != nil {
		return err
	}
	var quiere []domain.IntegranteID
	for _, nombre := range lista(votantes) {
		quiere = append(quiere, domain.IntegranteID(m.id("integrante", nombre)))
	}
	if len(got) != len(quiere) {
		return fmt.Errorf("en la votación de %s votaron %s, se esperaba %q", historia, idsDeIntegrantes(got), votantes)
	}
	for n, id := range quiere {
		if got[n] != id {
			return fmt.Errorf("en la votación de %s votaron %s, se esperaba %q", historia, idsDeIntegrantes(got), votantes)
		}
	}
	return nil
}

// enLaVotacionNoVotoNadie verifica que la votación no tenga votos ("Entonces"): caso límite de la spec.
func (m *mundo) enLaVotacionNoVotoNadie(historia string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	got, err := m.poker.QuienesVotaron(domain.HistoriaID(m.id("historia", historia)))
	if err != nil {
		return err
	}
	if len(got) != 0 {
		return fmt.Errorf("en la votación de %s votaron %d integrantes, se esperaba nadie", historia, len(got))
	}
	return nil
}

// consultoLosVotos consulta los valores de los votos ("Cuando"): mientras la votación está abierta
// falla con ErrVotosOcultos (RN-16.4), que el "Entonces" verifica con "veo el error".
func (m *mundo) consultoLosVotos(historia string) error {
	_, err := m.poker.ConsultarVotos(domain.HistoriaID(m.id("historia", historia)))
	return m.resultado(err)
}

// idsDeIntegrantes junta los IDs para los mensajes de error de los pasos.
func idsDeIntegrantes(ids []domain.IntegranteID) string {
	textos := make([]string, 0, len(ids))
	for _, id := range ids {
		textos = append(textos, string(id))
	}
	return strings.Join(textos, ", ")
}
