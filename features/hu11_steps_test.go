package features

import (
	"fmt"
	"slices"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU11 registra los pasos de "Asignar historias a un sprint" (HU-11).
func registrarPasosHU11(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^asigno la historia "([^"]*)" al sprint "([^"]*)"$`, m.asignoLaHistoriaAlSprint)
	sc.Step(`^el sprint "([^"]*)" tiene las historias "([^"]*)"$`, m.elSprintTieneLasHistorias)
	sc.Step(`^el sprint "([^"]*)" no tiene historias$`, m.elSprintNoTieneHistorias)
	sc.Step(`^el sprint "([^"]*)" tiene (\d+) Story Points planificados$`, m.elSprintTieneSPPlanificados)
}

func (m *mundo) asignoLaHistoriaAlSprint(historia, nombreSprint string) error {
	err := m.sprints.Asignar(domain.SprintID(m.id("sprint", nombreSprint)), domain.HistoriaID(m.id("historia", historia)))
	return m.resultado(err)
}

// elSprintTieneLasHistorias verifica el Sprint Backlog ("HU-1, HU-2", en ese orden). No llama a m.sinError porque los
// escenarios lo usan también después de "veo el error", para comprobar que el sprint no cambió.
func (m *mundo) elSprintTieneLasHistorias(nombreSprint, historias string) error {
	sp, err := m.buscarSprint(nombreSprint)
	if err != nil {
		return err
	}
	var esperadas []domain.HistoriaID
	for _, h := range lista(historias) {
		esperadas = append(esperadas, domain.HistoriaID(m.id("historia", h)))
	}
	if !slices.Equal(sp.Historias, esperadas) {
		return fmt.Errorf("el sprint %q tiene las historias %v, se esperaba %v", nombreSprint, sp.Historias, esperadas)
	}
	return nil
}

// elSprintNoTieneHistorias verifica que el Sprint Backlog esté vacío. Como elSprintTieneLasHistorias, se usa
// después de "veo el error" y por eso no llama a m.sinError.
func (m *mundo) elSprintNoTieneHistorias(nombreSprint string) error {
	sp, err := m.buscarSprint(nombreSprint)
	if err != nil {
		return err
	}
	if len(sp.Historias) != 0 {
		return fmt.Errorf("el sprint %q tiene las historias %v, se esperaba que no tuviera ninguna", nombreSprint, sp.Historias)
	}
	return nil
}

// elSprintTieneSPPlanificados suma los Story Points de las historias del sprint con metrics.SPPlanificados (HU-24).
// Sprint.SPPlanificados recién queda fijo al cerrar el sprint (HU-13).
func (m *mundo) elSprintTieneSPPlanificados(nombreSprint string, esperados int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.buscarSprint(nombreSprint)
	if err != nil {
		return err
	}
	var historias []metrics.HistoriaSprint
	for _, id := range sp.Historias {
		h, err := m.backlog.Buscar(id)
		if err != nil {
			return fmt.Errorf("historia %q del sprint %q: %w", id, nombreSprint, err)
		}
		historias = append(historias, metrics.HistoriaSprint{StoryPoints: h.StoryPoints, Completada: h.Estado == backlog.Completada})
	}
	if got := metrics.SPPlanificados(historias); got != esperados {
		return fmt.Errorf("el sprint %q tiene %d Story Points planificados, se esperaban %d", nombreSprint, got, esperados)
	}
	return nil
}

// buscarSprint busca un sprint por su nombre del escenario.
func (m *mundo) buscarSprint(nombre string) (sprint.Sprint, error) {
	sp, err := m.sprints.Buscar(domain.SprintID(m.id("sprint", nombre)))
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("sprint %q: %w", nombre, err)
	}
	return sp, nil
}
