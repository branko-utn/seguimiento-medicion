package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU13 registra los pasos de "Cerrar sprint" (HU-13).
func registrarPasosHU13(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^cierro el sprint "([^"]*)"$`, m.cierroElSprint)
	sc.Step(`^el sprint "([^"]*)" cerró con (\d+) Story Points planificados y (\d+) completados$`, m.elSprintCerroCon)
	sc.Step(`^la historia "([^"]*)" no está asignada a ningún sprint$`, m.laHistoriaNoEstaAsignada)
}

// cierroElSprint ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) cierroElSprint(nombre string) error {
	return m.resultado(m.sprints.CerrarSprint(domain.SprintID(m.id("sprint", nombre))))
}

// elSprintCerroCon verifica los SP que quedaron fijos al cerrar el sprint ("Entonces").
func (m *mundo) elSprintCerroCon(nombre string, planificados, completados int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", nombre)))
	if err != nil {
		return err
	}
	if sp.SPPlanificados != planificados || sp.SPCompletados != completados {
		return fmt.Errorf("el sprint %q cerró con %d planificados y %d completados, se esperaban %d y %d",
			nombre, sp.SPPlanificados, sp.SPCompletados, planificados, completados)
	}
	return nil
}

// laHistoriaNoEstaAsignada verifica que la historia volvió al backlog sin sprint ("Entonces").
func (m *mundo) laHistoriaNoEstaAsignada(historia string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if h.SprintID != "" {
		return fmt.Errorf("la historia %q sigue asignada al sprint %s", historia, h.SprintID)
	}
	return nil
}
