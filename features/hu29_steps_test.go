package features

import (
	"fmt"
	"math"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
	"github.com/cucumber/godog"
)

// registrarPasosHU29 registra los pasos del reporte de sprint.
func registrarPasosHU29(sc *godog.ScenarioContext, m *mundo) {
	var r report.ReporteSprint

	sc.Step(`^genero el reporte del sprint "([^"]*)"$`, func(nombre string) error {
		var err error
		r, err = m.reportes.Sprint(domain.SprintID(m.id("sprint", nombre)))
		return m.resultado(err)
	})

	sc.Step(`^el reporte tiene las historias:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas := tabla.Rows[1:]
		if len(r.Historias) != len(filas) {
			return fmt.Errorf("el reporte tiene %d historias, se esperaban %d", len(r.Historias), len(filas))
		}
		for i, fila := range filas {
			h := r.Historias[i]
			c := fila.Cells
			sp, _ := strconv.Atoi(c[1].Value)
			estimadas, _ := strconv.ParseFloat(c[2].Value, 64)
			reales, _ := strconv.ParseFloat(c[3].Value, 64)
			completada := c[4].Value == "sí"
			if string(h.ID) != m.id("historia", c[0].Value) || h.StoryPoints != sp || !igual(h.HorasEstimadas, estimadas) ||
				!igual(h.HorasReales, reales) || h.Completada != completada {
				return fmt.Errorf("fila %d del reporte = %+v, se esperaba %s %s %s %s %s", i+1, h, c[0].Value, c[1].Value, c[2].Value, c[3].Value, c[4].Value)
			}
		}
		return nil
	})

	sc.Step(`^el reporte muestra (\d+) Story Points planificados y (\d+) completados$`, func(planificados, completados int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if r.SPPlanificados != planificados || r.SPCompletados != completados {
			return fmt.Errorf("el reporte muestra %d/%d SP, se esperaba %d/%d", r.SPPlanificados, r.SPCompletados, planificados, completados)
		}
		return nil
	})

	sc.Step(`^el reporte muestra ([\d.]+) % de historias completadas$`, func(porcentaje float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !igual(redondear(r.Porcentaje), porcentaje) {
			return fmt.Errorf("el reporte muestra %v %%, se esperaba %v %%", r.Porcentaje, porcentaje)
		}
		return nil
	})

	sc.Step(`^el reporte muestra ([\d.]+) horas estimadas, ([\d.]+) reales y una desviación de (-?[\d.]+) %$`, func(estimadas, reales, desviacion float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !igual(r.HorasEstimadas, estimadas) || !igual(r.HorasReales, reales) || !r.DesviacionCalculable || !igual(redondear(r.Desviacion), desviacion) {
			return fmt.Errorf("el reporte muestra %v h estimadas, %v reales y desviación %v (calculable %v)", r.HorasEstimadas, r.HorasReales, r.Desviacion, r.DesviacionCalculable)
		}
		return nil
	})

	sc.Step(`^el reporte muestra ([\d.]+) horas de "([^"]*)"$`, func(horas float64, integrante string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		for _, t := range r.PorIntegrante {
			if t.Integrante == integrante {
				if !igual(t.Horas, horas) {
					return fmt.Errorf("el reporte muestra %v horas de %q, se esperaban %v", t.Horas, integrante, horas)
				}
				return nil
			}
		}
		return fmt.Errorf("el reporte no muestra horas de %q: %+v", integrante, r.PorIntegrante)
	})

	sc.Step(`^el reporte muestra (\d+) defectos? detectados? y (\d+) resueltos?$`, func(detectados, resueltos int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if r.Detectados != detectados || r.Resueltos != resueltos {
			return fmt.Errorf("el reporte muestra %d detectados y %d resueltos, se esperaba %d y %d", r.Detectados, r.Resueltos, detectados, resueltos)
		}
		return nil
	})

	sc.Step(`^el reporte muestra ([\d.]+) horas reales$`, func(reales float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !igual(r.HorasReales, reales) {
			return fmt.Errorf("el reporte muestra %v horas reales, se esperaban %v", r.HorasReales, reales)
		}
		return nil
	})

	sc.Step(`^el reporte es parcial$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !r.Parcial {
			return fmt.Errorf("el reporte del sprint %s (%s) no está marcado como parcial", r.SprintID, r.Estado)
		}
		return nil
	})
}

// igual compara horas y porcentajes sin errores de redondeo de punto flotante.
func igual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// redondear deja el número con un decimal, como se muestra en la pantalla.
func redondear(v float64) float64 { return math.Round(v*10) / 10 }
