package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/cucumber/godog"
)

// registrarPasosHU04 registra los pasos de "Consultar estado del proyecto" (HU-04). El resumen consultado se guarda
// en una variable local de los pasos (como en HU-25), sin agregar campos al mundo compartido.
func registrarPasosHU04(sc *godog.ScenarioContext, m *mundo) {
	var resumen service.ResumenProyecto

	sc.Step(`^consulto el estado del proyecto "([^"]*)"$`, func(proyecto string) error {
		r, err := m.proyectos.Resumen(domain.ProyectoID(m.id("proyecto", proyecto)), m.hoy)
		resumen = r
		return m.resultado(err)
	})

	sc.Step(`^el estado del proyecto es "([^"]*)"$`, func(estado string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if resumen.Estado.String() != estado {
			return fmt.Errorf("el estado del proyecto es %q, se esperaba %q", resumen.Estado, estado)
		}
		return nil
	})

	sc.Step(`^quedan (\d+) días para el fin del proyecto$`, func(dias int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if resumen.DiasRestantes != dias {
			return fmt.Errorf("quedan %d días, se esperaban %d", resumen.DiasRestantes, dias)
		}
		return nil
	})

	sc.Step(`^el proyecto va del "([^"]*)" al "([^"]*)"$`, func(inicio, fin string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		i, err := parseFecha(inicio)
		if err != nil {
			return err
		}
		f, err := parseFecha(fin)
		if err != nil {
			return err
		}
		p := resumen.Proyecto
		if !p.Inicio.Equal(i) || !p.Fin.Equal(f) {
			return fmt.Errorf("el proyecto va del %s al %s, se esperaba del %s al %s",
				p.Inicio.Format("2006-01-02"), p.Fin.Format("2006-01-02"), inicio, fin)
		}
		return nil
	})

	sc.Step(`^el sprint activo del proyecto es "([^"]*)"$`, func(nombre string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if resumen.SprintActivo == nil {
			return fmt.Errorf("el proyecto no tiene sprint activo, se esperaba %q", nombre)
		}
		if quiere := m.id("sprint", nombre); string(resumen.SprintActivo.ID) != quiere {
			return fmt.Errorf("el sprint activo es %q, se esperaba %q", resumen.SprintActivo.ID, quiere)
		}
		return nil
	})

	sc.Step(`^el proyecto no tiene sprint activo$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if resumen.SprintActivo != nil {
			return fmt.Errorf("el sprint activo es %q, se esperaba que no hubiera", resumen.SprintActivo.ID)
		}
		return nil
	})

	sc.Step(`^el porcentaje de historias completadas del proyecto es (\d+(?:\.\d+)?)$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		esperado, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("porcentaje esperado %q inválido: %w", texto, err)
		}
		// Se compara con un decimal, como en HU-25.
		if fmt.Sprintf("%.1f", resumen.PorcentajeCompletadas) != fmt.Sprintf("%.1f", esperado) {
			return fmt.Errorf("el porcentaje es %.1f, se esperaba %.1f", resumen.PorcentajeCompletadas, esperado)
		}
		return nil
	})
}
