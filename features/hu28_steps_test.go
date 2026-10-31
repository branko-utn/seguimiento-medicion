package features

import (
	"fmt"
	"math"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/dashboard"
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

func registrarPasosHU28(sc *godog.ScenarioContext, m *mundo) {
	var d dashboard.Dashboard

	sc.Step(`^veo el dashboard de "([^"]*)"$`, func(nombre string) error {
		var err error
		d, err = m.dashboard.Del(domain.ProyectoID(m.id("proyecto", nombre)), m.hoy)
		return m.resultado(err)
	})

	sc.Step(`^el dashboard muestra el estado "([^"]*)" con (\d+) días restantes$`, func(estado string, dias int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if d.Estado != estado || d.DiasRestantes != dias {
			return fmt.Errorf("estado=%s días=%d, se esperaba %s %d", d.Estado, d.DiasRestantes, estado, dias)
		}
		return nil
	})

	sc.Step(`^el dashboard muestra el sprint activo "([^"]*)"$`, func(sprint string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if string(d.SprintActivo) != sprint {
			return fmt.Errorf("sprint activo=%s, se esperaba %s", d.SprintActivo, sprint)
		}
		return nil
	})

	sc.Step(`^el dashboard muestra una velocidad de ([\d.]+)$`, func(v float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if math.Abs(d.Velocidad-v) > 1e-9 {
			return fmt.Errorf("velocidad=%v, se esperaba %v", d.Velocidad, v)
		}
		return nil
	})

	sc.Step(`^el dashboard muestra ([\d.]+) % de historias completadas$`, func(p float64) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if math.Abs(math.Round(d.PorcentajeCompletadas*10)/10-p) > 1e-9 && math.Abs(d.PorcentajeCompletadas-p) > 1e-9 {
			return fmt.Errorf("porcentaje=%v, se esperaba %v", d.PorcentajeCompletadas, p)
		}
		return nil
	})

	sc.Step(`^el dashboard tiene los puntos por sprint:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas := tabla.Rows[1:]
		if len(d.Puntos) != len(filas) {
			return fmt.Errorf("hay %d puntos, se esperaban %d", len(d.Puntos), len(filas))
		}
		for i, fila := range filas {
			c := fila.Cells
			p := d.Puntos[i]
			if string(p.Sprint) != c[0].Value {
				return fmt.Errorf("fila %d sprint=%s esperado %s", i+1, p.Sprint, c[0].Value)
			}
			spPlan, _ := strconv.Atoi(c[1].Value)
			spComp, _ := strconv.Atoi(c[2].Value)
			he, _ := strconv.ParseFloat(c[3].Value, 64)
			hr, _ := strconv.ParseFloat(c[4].Value, 64)
			dd, _ := strconv.Atoi(c[5].Value)
			dr, _ := strconv.Atoi(c[6].Value)
			if p.SPPlanificados != spPlan || p.SPCompletados != spComp || math.Abs(p.HorasEstimadas-he) > 1e-9 || math.Abs(p.HorasReales-hr) > 1e-9 || p.DefectosDetectados != dd || p.DefectosResueltos != dr {
				return fmt.Errorf("fila %d = %+v", i+1, p)
			}
		}
		return nil
	})

	sc.Step(`^los gráficos muestran el mensaje "([^"]*)"$`, func(msg string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if d.Mensajes.Velocidad != msg || d.Mensajes.SP != msg || d.Mensajes.Horas != msg || d.Mensajes.Defectos != msg {
			return fmt.Errorf("mensajes = %+v, se esperaba todos %q", d.Mensajes, msg)
		}
		return nil
	})

	sc.Step(`^el dashboard muestra la desviación como "([^"]*)"$`, func(msg string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if msg == "no calculable" {
			if d.DesviacionCalculable {
				return fmt.Errorf("se esperaba desviación no calculable")
			}
			return nil
		}
		return fmt.Errorf("mensaje de desviación no soportado: %s", msg)
	})

	sc.Step(`^el gráfico de defectos muestra el mensaje "([^"]*)"$`, func(msg string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if d.Mensajes.Defectos != msg {
			return fmt.Errorf("mensaje defectos = %q, se esperaba %q", d.Mensajes.Defectos, msg)
		}
		return nil
	})

	sc.Step(`^el gráfico de horas se dibuja$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if d.Mensajes.Horas != "" {
			return fmt.Errorf("se esperaba que el gráfico de horas se dibujara, mensaje=%q", d.Mensajes.Horas)
		}
		return nil
	})

}

// ensure project error type is referenced if needed
