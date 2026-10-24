package features

import (
	"fmt"
	"slices"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/cucumber/godog"
)

// registrarPasosCierreSprints registra los pasos compartidos de sprints que le tocan a Branko (features/README.md).
func registrarPasosCierreSprints(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que el sprint "([^"]*)" está cerrado$`, m.queElSprintEstaCerrado)
	sc.Step(`^que la historia "([^"]*)" está completada$`, m.queLaHistoriaEstaCompletada)
	sc.Step(`^que la historia "([^"]*)" ya está completada en el sprint "([^"]*)"$`, m.queLaHistoriaYaEstaCompletadaEnElSprint)
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

// queLaHistoriaEstaCompletada deja la historia en estado Completada directamente en el repositorio ("Dado").
// Es preparación del escenario (HU-15, HU-22): el caso de uso real de completar es HU-12, dentro de un sprint.
func (m *mundo) queLaHistoriaEstaCompletada(historia string) error {
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	h.Estado = backlog.Completada
	return m.historias.Guardar(h)
}

// queLaHistoriaYaEstaCompletadaEnElSprint deja la historia completada dentro del sprint ("Dado"; HU-08, HU-12): si no
// estaba asignada, la asigna; queda en Sprint.Historias y en Sprint.Completadas (sin repetirse) y en estado Completada.
// Es preparación del escenario: no pasa por los casos de uso de HU-11 ni de HU-12.
func (m *mundo) queLaHistoriaYaEstaCompletadaEnElSprint(historia, nombreSprint string) error {
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", nombreSprint)))
	if err != nil {
		return fmt.Errorf("sprint %q: %w", nombreSprint, err)
	}
	if h.SprintID != "" && h.SprintID != sp.ID {
		return fmt.Errorf("la historia %q está asignada a otro sprint (%s)", historia, h.SprintID)
	}
	h.SprintID = sp.ID
	h.Estado = backlog.Completada
	if !slices.Contains(sp.Historias, h.ID) {
		sp.Historias = append(sp.Historias, h.ID)
	}
	if !slices.Contains(sp.Completadas, h.ID) {
		sp.Completadas = append(sp.Completadas, h.ID)
	}
	if err := m.repoSprints.Actualizar(sp); err != nil {
		return err
	}
	return m.historias.Guardar(h)
}
