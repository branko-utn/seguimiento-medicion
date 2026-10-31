package features

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/report"
	"github.com/cucumber/godog"
)

// registrarPasosHU30 registra los pasos del reporte de proyecto.
func registrarPasosHU30(sc *godog.ScenarioContext, m *mundo) {
	var r report.ReporteProyecto

	sc.Step(`^genero el reporte del proyecto "([^"]*)"$`, func(nombre string) error {
		var err error
		r, err = m.reportes.Proyecto(domain.ProyectoID(m.id("proyecto", nombre)), m.hoy)
		return m.resultado(err)
	})

	sc.Step(`^el reporte del proyecto muestra el estado "([^"]*)"$`, func(estado string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if r.Estado != estado {
			return fmt.Errorf("el reporte del proyecto muestra el estado %q, se esperaba %q", r.Estado, estado)
		}
		return nil
	})

	sc.Step(`^el resumen de sprints es:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas := tabla.Rows[1:]
		if len(r.Sprints) != len(filas) {
			return fmt.Errorf("el resumen tiene %d sprints, se esperaban %d", len(r.Sprints), len(filas))
		}
		for i, fila := range filas {
			s, c := r.Sprints[i], fila.Cells
			planificados, _ := strconv.Atoi(c[2].Value)
			completados, _ := strconv.Atoi(c[3].Value)
			if string(s.SprintID) != m.id("sprint", c[0].Value) || s.Estado != c[1].Value || s.SPPlanificados != planificados || s.SPCompletados != completados {
				return fmt.Errorf("fila %d del resumen = %s %s %d/%d, se esperaba %s %s %s/%s", i+1,
					s.SprintID, s.Estado, s.SPPlanificados, s.SPCompletados, c[0].Value, c[1].Value, c[2].Value, c[3].Value)
			}
		}
		return nil
	})

	sc.Step(`^el resumen de sprints está vacío$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if len(r.Sprints) != 0 {
			return fmt.Errorf("el resumen tiene %d sprints, se esperaba vacío", len(r.Sprints))
		}
		return nil
	})

	sc.Step(`^el reporte del proyecto muestra una velocidad de ([\d.]+)$`, func(velocidad float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !igual(redondear(r.Velocidad), velocidad) {
			return fmt.Errorf("el reporte del proyecto muestra una velocidad de %v, se esperaba %v", r.Velocidad, velocidad)
		}
		return nil
	})

	sc.Step(`^el reporte del proyecto muestra ([\d.]+) % de historias completadas$`, func(porcentaje float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if !igual(redondear(r.Porcentaje), porcentaje) {
			return fmt.Errorf("el reporte del proyecto muestra %v %%, se esperaba %v %%", r.Porcentaje, porcentaje)
		}
		return nil
	})

	sc.Step(`^el reporte del proyecto muestra la desviación como "no calculable"$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if r.DesviacionCalculable {
			return fmt.Errorf("el reporte del proyecto muestra una desviación de %v, se esperaba \"no calculable\"", r.Desviacion)
		}
		return nil
	})

	sc.Step(`^el reporte del proyecto muestra (\d+) defectos detectados, (\d+) resueltos y (\d+) abiertos$`, func(detectados, resueltos, abiertos int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if r.Detectados != detectados || r.Resueltos != resueltos || r.Abiertos != abiertos {
			return fmt.Errorf("el reporte del proyecto muestra %d detectados, %d resueltos y %d abiertos, se esperaba %d, %d y %d",
				r.Detectados, r.Resueltos, r.Abiertos, detectados, resueltos, abiertos)
		}
		return nil
	})

	sc.Step(`^los defectos abiertos son "([^"]*)"$`, func(lista string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		var got []string
		for _, d := range r.DefectosAbiertos {
			got = append(got, d.Descripcion)
		}
		if strings.Join(got, ", ") != lista {
			return fmt.Errorf("los defectos abiertos son %q, se esperaba %q", strings.Join(got, ", "), lista)
		}
		return nil
	})
}
