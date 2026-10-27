package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/cucumber/godog"
)

// registrarPasosHU23 registra los pasos de "Cambiar estado de un defecto y registrar su resolución" (HU-23).
func registrarPasosHU23(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que existe el defecto "([^"]*)" en "([^"]*)" detectado en "([^"]*)"$`, m.queExisteElDefectoDetectadoEn)
	sc.Step(`^que el defecto "([^"]*)" está en estado "([^"]*)"$`, m.queElDefectoEstaEnEstado)
	sc.Step(`^que el defecto "([^"]*)" se resolvió en el sprint "([^"]*)"$`, m.queElDefectoSeResolvioEn)
	sc.Step(`^cambio el estado del defecto "([^"]*)" a "([^"]*)"$`, m.cambioElEstadoDelDefecto)
	sc.Step(`^resuelvo el defecto "([^"]*)" en el sprint "([^"]*)"$`, m.resuelvoElDefecto)
	sc.Step(`^el defecto "([^"]*)" se resolvió en el sprint "([^"]*)"$`, m.elDefectoSeResolvioEn)
}

// estadoDefectoDesdeTexto convierte el texto del escenario en un Estado (ADR 0003: la conversión va en el borde).
// Si el texto no es un estado, devuelve el cero.
func estadoDefectoDesdeTexto(texto string) defect.Estado {
	for _, e := range []defect.Estado{defect.Abierto, defect.EnProgreso, defect.Resuelto} {
		if e.String() == texto {
			return e
		}
	}
	return 0
}

// queExisteElDefectoDetectadoEn registra un defecto de severidad Alta ("Dado"): el escenario no indica la descripción
// ni la severidad. El nombre ("D1") queda como alias del ID que generó el repositorio.
func (m *mundo) queExisteElDefectoDetectadoEn(defecto, historia, sprint string) error {
	d, err := m.defectos.Registrar(
		"Defecto de prueba",
		defect.Alta,
		domain.HistoriaID(m.id("historia", historia)),
		domain.SprintID(m.id("sprint", sprint)),
	)
	if err != nil {
		return err
	}
	m.recordar("defecto", defecto, string(d.ID))
	return nil
}

func (m *mundo) idDefecto(defecto string) domain.DefectoID {
	return domain.DefectoID(m.id("defecto", defecto))
}

// queElDefectoEstaEnEstado deja el defecto en ese estado con el caso de uso ("Dado").
func (m *mundo) queElDefectoEstaEnEstado(defecto, estado string) error {
	return m.defectos.CambiarEstado(m.idDefecto(defecto), estadoDefectoDesdeTexto(estado))
}

// queElDefectoSeResolvioEn deja el defecto resuelto en ese sprint con el caso de uso ("Dado").
func (m *mundo) queElDefectoSeResolvioEn(defecto, sprint string) error {
	return m.defectos.Resolver(m.idDefecto(defecto), domain.SprintID(m.id("sprint", sprint)))
}

// cambioElEstadoDelDefecto ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) cambioElEstadoDelDefecto(defecto, estado string) error {
	return m.resultado(m.defectos.CambiarEstado(m.idDefecto(defecto), estadoDefectoDesdeTexto(estado)))
}

// resuelvoElDefecto ejecuta el caso de uso y guarda su resultado ("Cuando").
func (m *mundo) resuelvoElDefecto(defecto, sprint string) error {
	return m.resultado(m.defectos.Resolver(m.idDefecto(defecto), domain.SprintID(m.id("sprint", sprint))))
}

// elDefectoSeResolvioEn verifica el sprint de resolución del defecto ("Entonces").
func (m *mundo) elDefectoSeResolvioEn(defecto, sprint string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	d, err := m.buscarDefecto(defecto)
	if err != nil {
		return err
	}
	if quiere := domain.SprintID(m.id("sprint", sprint)); d.SprintResolucion != quiere {
		return fmt.Errorf("el defecto %q se resolvió en %q, se esperaba %q", defecto, d.SprintResolucion, quiere)
	}
	return nil
}
