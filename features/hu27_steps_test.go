package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU27 registra los pasos de "Defectos detectados y resueltos" (HU-27). Los defectos viven en una
// variable local: inicializarEscenario se ejecuta para cada escenario, así que no se comparten (igual que en HU-25).
// Son funciones puras: los escenarios no pasan por el servicio de defectos.
func registrarPasosHU27(sc *godog.ScenarioContext, m *mundo) {
	var defectos []metrics.DefectoSprint

	// Tabla con las columnas "defecto", "detectado en" y "resuelto en" (vacía si no está resuelto).
	sc.Step(`^los defectos:$`, func(tabla *godog.Table) error {
		if len(tabla.Rows) < 2 {
			return fmt.Errorf("la tabla de defectos no tiene filas")
		}
		columnas := map[string]int{}
		for i, celda := range tabla.Rows[0].Cells {
			columnas[celda.Value] = i
		}
		for _, c := range []string{"detectado en", "resuelto en"} {
			if _, ok := columnas[c]; !ok {
				return fmt.Errorf("a la tabla de defectos le falta la columna %q", c)
			}
		}
		for _, fila := range tabla.Rows[1:] {
			defectos = append(defectos, metrics.DefectoSprint{
				SprintDeteccion:  fila.Cells[columnas["detectado en"]].Value,
				SprintResolucion: fila.Cells[columnas["resuelto en"]].Value,
			})
		}
		return nil
	})

	sc.Step(`^en el sprint "([^"]*)" se detectaron (\d+) defectos y se resolvieron (\d+)$`, func(sprint string, detectados, resueltos int) error {
		d, r := metrics.DefectosDetectados(defectos, sprint), metrics.DefectosResueltos(defectos, sprint)
		if d != detectados || r != resueltos {
			return fmt.Errorf("en el sprint %q se detectaron %d y se resolvieron %d, se esperaba %d y %d", sprint, d, r, detectados, resueltos)
		}
		return nil
	})

	sc.Step(`^en total se detectaron (\d+) defectos y se resolvieron (\d+)$`, func(detectados, resueltos int) error {
		d, r := metrics.TotalDetectados(defectos), metrics.TotalResueltos(defectos)
		if d != detectados || r != resueltos {
			return fmt.Errorf("en total se detectaron %d y se resolvieron %d, se esperaba %d y %d", d, r, detectados, resueltos)
		}
		return nil
	})
}
