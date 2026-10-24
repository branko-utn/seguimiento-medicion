package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU20 registra los pasos de "Registrar esfuerzo trabajado" (HU-20).
func registrarPasosHU20(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^"([^"]*)" registra (-?\d+(?:\.\d+)?) horas de "([^"]*)" en "([^"]*)" el "([^"]*)"$`, m.registraHoras)
	sc.Step(`^que "([^"]*)" registró (\d+(?:\.\d+)?) horas de "([^"]*)" en "([^"]*)" el "([^"]*)"$`, m.queRegistroHoras)
	sc.Step(`^la historia "([^"]*)" tiene (\d+(?:\.\d+)?) horas reales$`, m.laHistoriaTieneHorasReales)
	sc.Step(`^"([^"]*)" tiene (\d+(?:\.\d+)?) horas registradas el "([^"]*)"$`, m.tieneHorasRegistradasEl)
}

// registrar ejecuta el caso de uso con la fecha de hoy del escenario (paso "hoy es …").
func (m *mundo) registrar(integrante string, horas float64, actividad, historia, fecha string) error {
	f, err := parseFecha(fecha)
	if err != nil {
		return err
	}
	_, err = m.registros.Registrar(
		domain.IntegranteID(m.id("integrante", integrante)),
		domain.HistoriaID(m.id("historia", historia)),
		f, actividad, horas, m.hoy,
	)
	return err
}

// registraHoras ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) registraHoras(integrante string, horas float64, actividad, historia, fecha string) error {
	return m.resultado(m.registrar(integrante, horas, actividad, historia, fecha))
}

// queRegistroHoras deja cargado un registro previo ("Dado").
func (m *mundo) queRegistroHoras(integrante string, horas float64, actividad, historia, fecha string) error {
	return m.registrar(integrante, horas, actividad, historia, fecha)
}

// laHistoriaTieneHorasReales verifica la suma de los registros de la historia ("Entonces").
func (m *mundo) laHistoriaTieneHorasReales(historia string, horas float64) error {
	if err := m.sinError(); err != nil {
		return err
	}
	var total float64
	for _, r := range m.registros.DeHistoria(domain.HistoriaID(m.id("historia", historia))) {
		total += r.Horas
	}
	if total != horas {
		return fmt.Errorf("la historia %q tiene %v horas reales, se esperaban %v", historia, total, horas)
	}
	return nil
}

// tieneHorasRegistradasEl verifica la suma de los registros del integrante en ese día ("Entonces").
func (m *mundo) tieneHorasRegistradasEl(integrante string, horas float64, fecha string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	f, err := parseFecha(fecha)
	if err != nil {
		return err
	}
	var total float64
	for _, r := range m.registros.DelIntegrante(domain.IntegranteID(m.id("integrante", integrante))) {
		if r.Fecha.Equal(f) {
			total += r.Horas
		}
	}
	if total != horas {
		return fmt.Errorf("%q tiene %v horas registradas el %s, se esperaban %v", integrante, total, fecha, horas)
	}
	return nil
}
