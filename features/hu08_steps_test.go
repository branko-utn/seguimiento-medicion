package features

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU08 registra los pasos de "Cambiar el estado de una historia" (HU-08). El "Dado" de asignación
// es compartido de Sosa, los de estado y completada son compartidos de backlog y de Branko (features/README.md).
func registrarPasosHU08(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^cambio el estado de "([^"]*)" a "([^"]*)"$`, m.cambioElEstado)
}

// cambioElEstado ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado). El
// texto del escenario se convierte en Estado en el borde (ADR 0003); un texto desconocido falla el paso.
func (m *mundo) cambioElEstado(historia, estado string) error {
	e, err := estadoDesdeTexto(estado)
	if err != nil {
		return err
	}
	err = m.backlog.CambiarEstado(domain.HistoriaID(m.id("historia", historia)), e)
	return m.resultado(err)
}
