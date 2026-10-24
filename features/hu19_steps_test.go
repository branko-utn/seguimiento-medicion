package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU19 registra los pasos de "Registrar horas estimadas de una historia" (HU-19).
func registrarPasosHU19(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^estimo la historia "([^"]*)" en (-?\d+(?:\.\d+)?) horas$`, m.estimoLaHistoria)
	sc.Step(`^que la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas estimadas$`, m.queLaHistoriaTieneHorasEstimadas)
	sc.Step(`^la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas estimadas$`, m.laHistoriaTieneHorasEstimadas)
}

// estimoLaHistoria ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) estimoLaHistoria(historia string, horas float64) error {
	_, err := m.estimaciones.Estimar(domain.HistoriaID(m.id("historia", historia)), horas)
	return m.resultado(err)
}

// queLaHistoriaTieneHorasEstimadas deja cargada una estimación previa ("Dado").
func (m *mundo) queLaHistoriaTieneHorasEstimadas(historia string, horas float64) error {
	_, err := m.estimaciones.Estimar(domain.HistoriaID(m.id("historia", historia)), horas)
	return err
}

// laHistoriaTieneHorasEstimadas verifica las horas guardadas para la historia ("Entonces").
func (m *mundo) laHistoriaTieneHorasEstimadas(historia string, horas float64) error {
	if err := m.sinError(); err != nil {
		return err
	}
	e, ok, err := m.estimaciones.Buscar(domain.HistoriaID(m.id("historia", historia)))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("la historia %q no tiene horas estimadas", historia)
	}
	if e.Horas != horas {
		return fmt.Errorf("la historia %q tiene %v horas estimadas, se esperaban %v", historia, e.Horas, horas)
	}
	return nil
}
