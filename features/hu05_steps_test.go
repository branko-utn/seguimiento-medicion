package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU05 registra los pasos de "Crear ítem del Product Backlog" (HU-05).
func registrarPasosHU05(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^agrego al proyecto "([^"]*)" la historia "([^"]*)" con prioridad "([^"]*)"$`, m.agregoLaHistoria)
	sc.Step(`^que el backlog de "([^"]*)" tiene (\d+) historias$`, m.queElBacklogTieneNHistorias)
	sc.Step(`^el backlog de "([^"]*)" tiene la historia "([^"]*)" "([^"]*)"$`, m.elBacklogTieneLaHistoria)
	sc.Step(`^la historia "([^"]*)" no tiene Story Points$`, m.laHistoriaNoTieneStoryPoints)
}

// agregoLaHistoria ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) agregoLaHistoria(proyecto, titulo, prioridad string) error {
	h, err := m.backlog.Crear(domain.ProyectoID(m.id("proyecto", proyecto)), titulo, "", prioridadDesdeTexto(prioridad))
	if err == nil {
		m.recordar("historia", titulo, string(h.ID))
	}
	return m.resultado(err)
}

// queElBacklogTieneNHistorias crea n historias genéricas en el backlog del proyecto ("Dado").
func (m *mundo) queElBacklogTieneNHistorias(proyecto string, n int) error {
	for i := 1; i <= n; i++ {
		titulo := fmt.Sprintf("Historia %d", i)
		if _, err := m.backlog.Crear(domain.ProyectoID(m.id("proyecto", proyecto)), titulo, "", prioridadDesdeTexto("Media")); err != nil {
			return err
		}
	}
	return nil
}

// elBacklogTieneLaHistoria verifica que el backlog del proyecto tenga la historia con ese ID y título ("Entonces").
// Filtra por proyecto: el mismo ID ("HU-1") puede existir en dos proyectos (RN-05.1).
func (m *mundo) elBacklogTieneLaHistoria(proyecto, historia, titulo string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	p, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	lista, err := m.backlog.Listar()
	if err != nil {
		return err
	}
	for _, h := range lista {
		if h.ProyectoID == p.ID && h.ID == domain.HistoriaID(m.id("historia", historia)) {
			if h.Titulo != titulo {
				return fmt.Errorf("la historia %s de %q se llama %q, se esperaba %q", historia, proyecto, h.Titulo, titulo)
			}
			return nil
		}
	}
	return fmt.Errorf("el backlog de %q no tiene la historia %s", proyecto, historia)
}

// laHistoriaNoTieneStoryPoints verifica que la historia esté sin estimar (RN-05.4, "Entonces").
func (m *mundo) laHistoriaNoTieneStoryPoints(historia string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if h.StoryPoints != 0 {
		return fmt.Errorf("la historia %q tiene %d Story Points, se esperaba 0", historia, h.StoryPoints)
	}
	return nil
}
