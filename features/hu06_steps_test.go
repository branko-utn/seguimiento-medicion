package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/cucumber/godog"
)

// registrarPasosHU06 registra los pasos de "Editar y eliminar ítem del backlog" (HU-06).
func registrarPasosHU06(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^edito la historia "([^"]*)" con título "([^"]*)", descripción "([^"]*)" y prioridad "([^"]*)"$`, m.editoLaHistoria)
	sc.Step(`^la historia "([^"]*)" tiene el título "([^"]*)" y prioridad "([^"]*)"$`, m.laHistoriaTieneElTituloYPrioridad)
	sc.Step(`^elimino la historia "([^"]*)"$`, m.eliminoLaHistoria)
	sc.Step(`^el backlog de "([^"]*)" no tiene la historia "([^"]*)"$`, m.elBacklogNoTieneLaHistoria)
	sc.Step(`^el backlog de "([^"]*)" tiene (\d+) historias$`, m.elBacklogTieneNHistorias)
}

// editoLaHistoria ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) editoLaHistoria(historia, titulo, descripcion, prioridad string) error {
	err := m.backlog.Editar(domain.HistoriaID(m.id("historia", historia)), titulo, descripcion, prioridadDesdeTexto(prioridad))
	return m.resultado(err)
}

// laHistoriaTieneElTituloYPrioridad verifica el título y la prioridad guardados ("Entonces").
func (m *mundo) laHistoriaTieneElTituloYPrioridad(historia, titulo, prioridad string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if h.Titulo != titulo || h.Prioridad.String() != prioridad {
		return fmt.Errorf("la historia %q tiene título %q y prioridad %q, se esperaba %q y %q",
			historia, h.Titulo, h.Prioridad, titulo, prioridad)
	}
	return nil
}

// eliminoLaHistoria ejecuta el caso de uso y guarda su resultado ("Cuando").
func (m *mundo) eliminoLaHistoria(historia string) error {
	return m.resultado(m.backlog.Eliminar(domain.HistoriaID(m.id("historia", historia))))
}

// elBacklogNoTieneLaHistoria verifica que la historia ya no esté en el backlog del proyecto ("Entonces").
func (m *mundo) elBacklogNoTieneLaHistoria(proyecto, historia string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	lista, err := m.backlogDe(proyecto)
	if err != nil {
		return err
	}
	id := domain.HistoriaID(m.id("historia", historia))
	for _, h := range lista {
		if h.ID == id {
			return fmt.Errorf("el backlog de %q todavía tiene la historia %q", proyecto, historia)
		}
	}
	return nil
}

// elBacklogTieneNHistorias verifica cuántas historias tiene el backlog del proyecto ("Entonces").
func (m *mundo) elBacklogTieneNHistorias(proyecto string, n int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	lista, err := m.backlogDe(proyecto)
	if err != nil {
		return err
	}
	if len(lista) != n {
		return fmt.Errorf("el backlog de %q tiene %d historias, se esperaban %d", proyecto, len(lista), n)
	}
	return nil
}

// backlogDe devuelve las historias del backlog de un proyecto (Listar devuelve las de todos los proyectos).
func (m *mundo) backlogDe(proyecto string) ([]backlog.Historia, error) {
	p, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return nil, err
	}
	todas, err := m.backlog.Listar()
	if err != nil {
		return nil, err
	}
	var lista []backlog.Historia
	for _, h := range todas {
		if h.ProyectoID == p.ID {
			lista = append(lista, h)
		}
	}
	return lista, nil
}
