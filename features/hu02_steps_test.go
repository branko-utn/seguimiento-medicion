package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/cucumber/godog"
)

// registrarPasosHU02 registra los pasos de "Modificar proyecto" (HU-02).
func registrarPasosHU02(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^modifico el proyecto "([^"]*)" con el nombre "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.modificoElProyecto)
	sc.Step(`^que el proyecto "([^"]*)" está finalizado$`, m.queElProyectoEstaFinalizado)
}

// modificoElProyecto cambia nombre y fechas. El paso no menciona la descripción, así que se conserva la que tiene.
// Si sale bien, recuerda el nombre nuevo para que los pasos siguientes encuentren el proyecto.
func (m *mundo) modificoElProyecto(proyecto, nombre, inicio, fin string) error {
	i, err := parseFecha(inicio)
	if err != nil {
		return err
	}
	f, err := parseFecha(fin)
	if err != nil {
		return err
	}
	id := domain.ProyectoID(m.id("proyecto", proyecto))
	descripcion := ""
	if actual, err := m.repoProyectos.Buscar(id); err == nil {
		descripcion = actual.Descripcion
	}
	p, err := m.proyectos.Modificar(id, nombre, descripcion, i, f)
	if err == nil {
		m.recordar("proyecto", p.Nombre, string(p.ID))
	}
	return m.resultado(err)
}

// queElProyectoEstaFinalizado deja el proyecto en estado Finalizado directamente en el repositorio ("Dado"). Es
// preparación del escenario: todavía no hay un caso de uso que finalice un proyecto.
func (m *mundo) queElProyectoEstaFinalizado(proyecto string) error {
	p, err := m.repoProyectos.Buscar(domain.ProyectoID(m.id("proyecto", proyecto)))
	if err != nil {
		return fmt.Errorf("proyecto %q: %w", proyecto, err)
	}
	p.Estado = project.Finalizado
	return m.repoProyectos.Actualizar(p)
}
