package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/cucumber/godog"
)

// registrarPasosHU22 registra los pasos de "Registrar defecto" (HU-22).
func registrarPasosHU22(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" detectado en "([^"]*)"$`, m.registroElDefecto)
	sc.Step(`^existe el defecto "([^"]*)" en estado "([^"]*)"$`, m.existeElDefectoEnEstado)
	sc.Step(`^el defecto "([^"]*)" no tiene sprint de resolución$`, m.elDefectoNoTieneSprintDeResolucion)
}

// severidadDesdeTexto convierte el texto del escenario en una Severidad (ADR 0003: la conversión va en el borde).
// Si el texto no es una severidad, devuelve el cero, que defect.Nuevo rechaza con ErrSeveridadInvalida.
func severidadDesdeTexto(texto string) defect.Severidad {
	for _, s := range []defect.Severidad{defect.Critica, defect.Alta, defect.Media, defect.Baja} {
		if s.String() == texto {
			return s
		}
	}
	return 0
}

// registroElDefecto ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) registroElDefecto(descripcion, severidad, historia, sprint string) error {
	_, err := m.defectos.Registrar(
		descripcion,
		severidadDesdeTexto(severidad),
		domain.HistoriaID(m.id("historia", historia)),
		domain.SprintID(m.id("sprint", sprint)),
	)
	return m.resultado(err)
}

// buscarDefecto devuelve el defecto guardado con ese ID ("D1"), o defect.ErrNoExiste.
func (m *mundo) buscarDefecto(defecto string) (defect.Defecto, error) {
	d, err := m.defectos.Buscar(domain.DefectoID(m.id("defecto", defecto)))
	if err != nil {
		return defect.Defecto{}, fmt.Errorf("defecto %q: %w", defecto, err)
	}
	return d, nil
}

// existeElDefectoEnEstado verifica que el defecto exista y esté en ese estado ("Entonces").
func (m *mundo) existeElDefectoEnEstado(defecto, estado string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	d, err := m.buscarDefecto(defecto)
	if err != nil {
		return err
	}
	if d.Estado.String() != estado {
		return fmt.Errorf("el defecto %q está en estado %q, se esperaba %q", defecto, d.Estado, estado)
	}
	return nil
}

// elDefectoNoTieneSprintDeResolucion verifica que el defecto no tenga sprint de resolución ("Entonces").
func (m *mundo) elDefectoNoTieneSprintDeResolucion(defecto string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	d, err := m.buscarDefecto(defecto)
	if err != nil {
		return err
	}
	if d.SprintResolucion != "" {
		return fmt.Errorf("el defecto %q tiene sprint de resolución %q, se esperaba ninguno", defecto, d.SprintResolucion)
	}
	return nil
}
