package features

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// registrarPasosHU01 registra los pasos de "Crear proyecto" (HU-01).
func registrarPasosHU01(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^creo el proyecto "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.creoElProyecto)
	sc.Step(`^creo un proyecto con un nombre de (\d+) caracteres del "([^"]*)" al "([^"]*)"$`, m.creoUnProyectoConNombreDeLargo)
	sc.Step(`^el proyecto "([^"]*)" existe$`, m.elProyectoExiste)
	sc.Step(`^el proyecto "([^"]*)" está en estado "([^"]*)"$`, m.elProyectoEstaEnEstado)
	sc.Step(`^el proyecto se crea$`, m.elProyectoSeCrea)
}

func (m *mundo) creoElProyecto(nombre, inicio, fin string) error {
	i, err := parseFecha(inicio)
	if err != nil {
		return err
	}
	f, err := parseFecha(fin)
	if err != nil {
		return err
	}
	p, err := m.proyectos.Crear(nombre, "", i, f)
	if err == nil {
		m.recordar("proyecto", p.Nombre, string(p.ID))
	}
	return m.resultado(err)
}

func (m *mundo) creoUnProyectoConNombreDeLargo(largo int, inicio, fin string) error {
	return m.creoElProyecto(strings.Repeat("a", largo), inicio, fin)
}

func (m *mundo) elProyectoExiste(nombre string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	_, err := m.proyectoLlamado(nombre)
	return err
}

func (m *mundo) elProyectoEstaEnEstado(nombre, estado string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	p, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	if p.Estado.String() != estado {
		return fmt.Errorf("el proyecto %q está en estado %q, se esperaba %q", nombre, p.Estado, estado)
	}
	return nil
}

func (m *mundo) elProyectoSeCrea() error {
	if err := m.sinError(); err != nil {
		return err
	}
	lista, err := m.proyectos.Listar()
	if err != nil {
		return err
	}
	if len(lista) != 1 {
		return fmt.Errorf("hay %d proyectos, se esperaba 1", len(lista))
	}
	return nil
}
