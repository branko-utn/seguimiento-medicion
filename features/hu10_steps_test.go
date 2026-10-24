package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU10 registra los pasos de "Crear sprint con Sprint Goal" (HU-10).
func registrarPasosHU10(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^creo en "([^"]*)" un sprint del "([^"]*)" al "([^"]*)" con el objetivo "([^"]*)"$`, m.creoUnSprint)
	sc.Step(`^inicio el sprint "([^"]*)"$`, m.inicioElSprint)
	sc.Step(`^que en "([^"]*)" existe un sprint del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintDelAl)
	sc.Step(`^que en "([^"]*)" existe un sprint "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintNombradoDelAl)
	sc.Step(`^que en "([^"]*)" existe un sprint activo "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queEnExisteUnSprintActivoDelAl)
	sc.Step(`^el proyecto "([^"]*)" tiene un sprint "([^"]*)" en estado "([^"]*)" con el objetivo "([^"]*)"$`, m.elProyectoTieneUnSprint)
	sc.Step(`^el proyecto "([^"]*)" tiene (\d+) sprints$`, m.elProyectoTieneNSprints)
	sc.Step(`^el sprint "([^"]*)" está en estado "([^"]*)"$`, m.elSprintEstaEnEstado)
}

func (m *mundo) creoUnSprint(proyecto, inicio, fin, goal string) error {
	i, err := parseFecha(inicio)
	if err != nil {
		return err
	}
	f, err := parseFecha(fin)
	if err != nil {
		return err
	}
	_, err = m.sprints.Crear(domain.ProyectoID(m.id("proyecto", proyecto)), goal, i, f)
	return m.resultado(err)
}

func (m *mundo) inicioElSprint(nombre string) error {
	return m.resultado(m.sprints.Iniciar(domain.SprintID(m.id("sprint", nombre))))
}

func (m *mundo) queEnExisteUnSprintDelAl(proyecto, inicio, fin string) error {
	return m.crearSprintConFechas(proyecto, "", inicio, fin, false)
}

func (m *mundo) queEnExisteUnSprintNombradoDelAl(proyecto, nombre, inicio, fin string) error {
	return m.crearSprintConFechas(proyecto, nombre, inicio, fin, false)
}

func (m *mundo) queEnExisteUnSprintActivoDelAl(proyecto, nombre, inicio, fin string) error {
	return m.crearSprintConFechas(proyecto, nombre, inicio, fin, true)
}

func (m *mundo) crearSprintConFechas(proyecto, nombre, inicio, fin string, activo bool) error {
	i, err := parseFecha(inicio)
	if err != nil {
		return err
	}
	f, err := parseFecha(fin)
	if err != nil {
		return err
	}
	return m.crearSprint(proyecto, nombre, i, f, activo)
}

func (m *mundo) elProyectoTieneUnSprint(proyecto, nombre, estado, goal string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.sprints.Buscar(domain.SprintID(m.id("sprint", nombre)))
	if err != nil {
		return err
	}
	if string(sp.ProyectoID) != m.id("proyecto", proyecto) {
		return fmt.Errorf("el sprint %q es del proyecto %q, se esperaba %q", nombre, sp.ProyectoID, proyecto)
	}
	if sp.Estado.String() != estado || sp.Goal != goal {
		return fmt.Errorf("el sprint %q está en estado %q con el objetivo %q, se esperaba %q con %q",
			nombre, sp.Estado, sp.Goal, estado, goal)
	}
	return nil
}

func (m *mundo) elProyectoTieneNSprints(proyecto string, n int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	lista, err := m.sprints.DelProyecto(domain.ProyectoID(m.id("proyecto", proyecto)))
	if err != nil {
		return err
	}
	if len(lista) != n {
		return fmt.Errorf("el proyecto %q tiene %d sprints, se esperaban %d", proyecto, len(lista), n)
	}
	return nil
}

func (m *mundo) elSprintEstaEnEstado(nombre, estado string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	sp, err := m.sprints.Buscar(domain.SprintID(m.id("sprint", nombre)))
	if err != nil {
		return err
	}
	if sp.Estado.String() != estado {
		return fmt.Errorf("el sprint %q está en estado %q, se esperaba %q", nombre, sp.Estado, estado)
	}
	return nil
}
