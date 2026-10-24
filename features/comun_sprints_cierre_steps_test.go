package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/cucumber/godog"
)

// registrarPasosCierreSprints registra los pasos compartidos de sprints que le tocan a Branko (features/README.md).
func registrarPasosCierreSprints(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que el sprint "([^"]*)" está cerrado$`, m.queElSprintEstaCerrado)
}

// queElSprintEstaCerrado deja el sprint en estado Cerrado directamente en el repositorio ("Dado"). Es preparación
// del escenario: cerrar un sprint de verdad (historias de vuelta al backlog, SP fijos) es HU-13, del Sprint 2.
func (m *mundo) queElSprintEstaCerrado(nombre string) error {
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", nombre)))
	if err != nil {
		return fmt.Errorf("sprint %q: %w", nombre, err)
	}
	sp.Estado = sprint.Cerrado
	return m.repoSprints.Actualizar(sp)
}

func (m *mundo) queLaHistoriaEstaCompletada(historia string) error { return nil }

func (m *mundo) queLaHistoriaYaEstaCompletadaEnElSprint(historia, nombreSprint string) error {
	return nil
}
