package features

import (
	"fmt"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosComunesSprints registra los pasos de sprints que usan varias historias (features/README.md).
// Responsable: Federico Sosa. Los escenarios los escriben con y sin "que": "Y existe un sprint …" y
// "Dado que existe un sprint …".
func registrarPasosComunesSprints(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^(?:que )?existe un sprint "([^"]*)" en "([^"]*)"$`, m.existeUnSprintEn)
	sc.Step(`^(?:que )?existe un sprint activo "([^"]*)" en "([^"]*)"$`, m.existeUnSprintActivoEn)
	sc.Step(`^que la historia "([^"]*)" está asignada al sprint "([^"]*)"$`, m.queLaHistoriaEstaAsignadaAlSprint)
}

func (m *mundo) queLaHistoriaEstaAsignadaAlSprint(historia, sprint string) error {
	return godog.ErrPending
}

func (m *mundo) existeUnSprintEn(nombre, proyecto string) error {
	return m.crearSprintSinFechas(proyecto, nombre, false)
}

func (m *mundo) existeUnSprintActivoEn(nombre, proyecto string) error {
	return m.crearSprintSinFechas(proyecto, nombre, true)
}

// diasPorSprint es el largo de los sprints que se crean sin fechas en el escenario.
const diasPorSprint = 3

// crearSprintSinFechas crea un sprint en los diasPorSprint días siguientes al último sprint del proyecto (el
// primero arranca el día de inicio del proyecto), sin pasar del fin. Así no se superpone ni queda fuera.
func (m *mundo) crearSprintSinFechas(proyecto, nombre string, activo bool) error {
	p, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	sprints, err := m.sprints.DelProyecto(p.ID)
	if err != nil {
		return err
	}
	inicio := p.Inicio
	for _, s := range sprints {
		if !s.Fin.Before(inicio) {
			inicio = s.Fin.AddDate(0, 0, 1)
		}
	}
	if inicio.After(p.Fin) {
		return fmt.Errorf("no quedan días libres en el proyecto %q para el sprint %q", proyecto, nombre)
	}
	fin := inicio.AddDate(0, 0, diasPorSprint-1)
	if fin.After(p.Fin) {
		fin = p.Fin
	}
	return m.crearSprint(proyecto, nombre, inicio, fin, activo)
}

// crearSprint crea un sprint (y lo inicia si activo es true). Si el escenario le da un nombre ("S1"), recuerda su ID.
func (m *mundo) crearSprint(proyecto, nombre string, inicio, fin time.Time, activo bool) error {
	sp, err := m.sprints.Crear(domain.ProyectoID(m.id("proyecto", proyecto)), "Sprint de prueba", inicio, fin)
	if err != nil {
		return err
	}
	if nombre != "" {
		m.recordar("sprint", nombre, string(sp.ID))
	}
	if activo {
		return m.sprints.Iniciar(sp.ID)
	}
	return nil
}
