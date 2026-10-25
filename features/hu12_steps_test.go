package features

import (
	"fmt"
	"slices"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU12 registra los pasos de "Registrar historias completadas" (HU-12). Los de historia completada y
// sprint cerrado son compartidos (comun_sprints_cierre_steps_test.go).
func registrarPasosHU12(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^existe un sprint activo "([^"]*)" con las historias "([^"]*)"$`, m.existeUnSprintActivoConLasHistorias)
	sc.Step(`^completo la historia "([^"]*)" en el sprint "([^"]*)"$`, m.completoLaHistoria)
	sc.Step(`^el sprint "([^"]*)" tiene completadas las historias "([^"]*)"$`, m.elSprintTieneCompletadas)
	sc.Step(`^el sprint "([^"]*)" tiene (\d+) Story Points completados de (\d+) planificados$`, m.elSprintTieneSPCompletadosDe)
}

// existeUnSprintActivoConLasHistorias crea un sprint activo en el último proyecto y le asigna las historias con el
// caso de uso de HU-11 ("Dado").
func (m *mundo) existeUnSprintActivoConLasHistorias(nombre, historias string) error {
	proyectos, err := m.proyectos.Listar()
	if err != nil {
		return err
	}
	if len(proyectos) == 0 {
		return fmt.Errorf("no hay ningún proyecto creado para el sprint %q", nombre)
	}
	if err := m.existeUnSprintActivoEn(nombre, proyectos[len(proyectos)-1].Nombre); err != nil {
		return err
	}
	for _, h := range lista(historias) {
		if err := m.sprints.Asignar(domain.SprintID(m.id("sprint", nombre)), domain.HistoriaID(m.id("historia", h))); err != nil {
			return fmt.Errorf("asignar %s al sprint %s: %w", h, nombre, err)
		}
	}
	return nil
}

// completoLaHistoria ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) completoLaHistoria(historia, nombreSprint string) error {
	err := m.sprints.CompletarHistoria(domain.SprintID(m.id("sprint", nombreSprint)), domain.HistoriaID(m.id("historia", historia)))
	return m.resultado(err)
}

// elSprintTieneCompletadas verifica la lista de historias completadas del sprint, en orden ("Entonces").
func (m *mundo) elSprintTieneCompletadas(nombreSprint, historias string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", nombreSprint)))
	if err != nil {
		return err
	}
	var quiere []domain.HistoriaID
	for _, h := range lista(historias) {
		quiere = append(quiere, domain.HistoriaID(m.id("historia", h)))
	}
	if !slices.Equal(sp.Completadas, quiere) {
		return fmt.Errorf("el sprint %q tiene completadas %v, se esperaba %v", nombreSprint, sp.Completadas, quiere)
	}
	return nil
}

// elSprintTieneSPCompletadosDe verifica los SP completados y planificados del sprint con las métricas de HU-24
// ("Entonces"). Se usa también después de un error esperado, para verificar que los SP no cambiaron.
func (m *mundo) elSprintTieneSPCompletadosDe(nombreSprint string, completados, planificados int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.repoSprints.Buscar(domain.SprintID(m.id("sprint", nombreSprint)))
	if err != nil {
		return err
	}
	var hs []metrics.HistoriaSprint
	for _, id := range sp.Historias {
		h, err := m.backlog.Buscar(id)
		if err != nil {
			return err
		}
		hs = append(hs, metrics.HistoriaSprint{StoryPoints: h.StoryPoints, Completada: slices.Contains(sp.Completadas, id)})
	}
	if got := metrics.SPCompletados(hs); got != completados {
		return fmt.Errorf("el sprint %q tiene %d SP completados, se esperaban %d", nombreSprint, got, completados)
	}
	if got := metrics.SPPlanificados(hs); got != planificados {
		return fmt.Errorf("el sprint %q tiene %d SP planificados, se esperaban %d", nombreSprint, got, planificados)
	}
	return nil
}
