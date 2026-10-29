package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/cucumber/godog"
)

// registrarPasosHU21 registra los pasos de "Consultar el esfuerzo registrado" (HU-21). El resumen consultado vive en
// una variable local: inicializarEscenario se ejecuta para cada escenario, así que no se comparte entre escenarios
// (igual que en HU-24 y HU-25).
func registrarPasosHU21(sc *godog.ScenarioContext, m *mundo) {
	var resumen service.ResumenEsfuerzo

	// Cuando: guarda el resumen y el error (puede ser el esperado: proyecto inexistente).
	sc.Step(`^consulto el esfuerzo registrado en "([^"]*)"$`, func(proyecto string) error {
		r, err := m.registros.DelProyecto(domain.ProyectoID(m.id("proyecto", proyecto)))
		resumen = r
		return m.resultado(err)
	})

	sc.Step(`^los totales por historia son:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas, err := filasDeTabla(tabla, "historia", "horas")
		if err != nil {
			return err
		}
		if len(filas) != len(resumen.PorHistoria) {
			return fmt.Errorf("hay %d totales por historia, se esperaban %d: %+v", len(resumen.PorHistoria), len(filas), resumen.PorHistoria)
		}
		for i, f := range filas {
			got := resumen.PorHistoria[i]
			if string(got.Historia) != m.id("historia", f["historia"]) || horasConUnDecimal(got.Horas) != f["horas"] {
				return fmt.Errorf("total %d = %s %s, se esperaba %s %s", i+1, got.Historia, horasConUnDecimal(got.Horas), f["historia"], f["horas"])
			}
		}
		return nil
	})

	sc.Step(`^los totales por integrante son:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas, err := filasDeTabla(tabla, "integrante", "horas")
		if err != nil {
			return err
		}
		if len(filas) != len(resumen.PorIntegrante) {
			return fmt.Errorf("hay %d totales por integrante, se esperaban %d: %+v", len(resumen.PorIntegrante), len(filas), resumen.PorIntegrante)
		}
		for i, f := range filas {
			got := resumen.PorIntegrante[i]
			if got.Nombre != f["integrante"] || horasConUnDecimal(got.Horas) != f["horas"] {
				return fmt.Errorf("total %d = %s %s, se esperaba %s %s", i+1, got.Nombre, horasConUnDecimal(got.Horas), f["integrante"], f["horas"])
			}
		}
		return nil
	})

	sc.Step(`^el esfuerzo total del proyecto es (\d+(?:\.\d+)?) horas$`, func(horas string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if got := horasConUnDecimal(resumen.Total); got != horas {
			return fmt.Errorf("el esfuerzo total es %s horas, se esperaban %s", got, horas)
		}
		return nil
	})

	sc.Step(`^los registros son:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		filas, err := filasDeTabla(tabla, "fecha", "integrante", "historia", "actividad", "horas")
		if err != nil {
			return err
		}
		if len(filas) != len(resumen.Registros) {
			return fmt.Errorf("hay %d registros, se esperaban %d", len(resumen.Registros), len(filas))
		}
		for i, f := range filas {
			r := resumen.Registros[i]
			got := fmt.Sprintf("%s %s %s %s %s", r.Fecha.Format("2006-01-02"), r.Integrante, r.HistoriaID, r.Actividad, horasConUnDecimal(r.Horas))
			want := fmt.Sprintf("%s %s %s %s %s", f["fecha"], f["integrante"], m.id("historia", f["historia"]), f["actividad"], f["horas"])
			if got != want {
				return fmt.Errorf("registro %d = %q, se esperaba %q", i+1, got, want)
			}
		}
		return nil
	})

	sc.Step(`^no hay esfuerzo registrado$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if len(resumen.Registros) != 0 || len(resumen.PorHistoria) != 0 || len(resumen.PorIntegrante) != 0 {
			return fmt.Errorf("se esperaba un resumen vacío, se obtuvo %+v", resumen)
		}
		return nil
	})
}

// filasDeTabla lee una tabla del escenario como una lista de filas columna → valor y verifica que tenga las columnas
// pedidas.
func filasDeTabla(tabla *godog.Table, columnas ...string) ([]map[string]string, error) {
	if len(tabla.Rows) == 0 {
		return nil, fmt.Errorf("la tabla está vacía")
	}
	var encabezado []string
	for _, celda := range tabla.Rows[0].Cells {
		encabezado = append(encabezado, celda.Value)
	}
	for _, c := range columnas {
		encontrada := false
		for _, e := range encabezado {
			encontrada = encontrada || e == c
		}
		if !encontrada {
			return nil, fmt.Errorf("a la tabla le falta la columna %q", c)
		}
	}
	var filas []map[string]string
	for _, fila := range tabla.Rows[1:] {
		valores := map[string]string{}
		for i, celda := range fila.Cells {
			valores[encabezado[i]] = celda.Value
		}
		filas = append(filas, valores)
	}
	return filas, nil
}

// horasConUnDecimal escribe las horas como en los escenarios: "5.5", "4.0".
func horasConUnDecimal(h float64) string {
	return fmt.Sprintf("%.1f", h)
}
