package features

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/cucumber/godog"
)

// registrarPasosHU14 registra los pasos de "Consultar sprints anteriores" (HU-14). El estado del escenario son el
// listado y el detalle devueltos por el servicio, como los lee la pantalla.
func registrarPasosHU14(sc *godog.ScenarioContext, m *mundo) {
	var cerrados []sprint.Sprint
	var detalle service.DetalleSprint

	sc.Step(`^consulto los sprints cerrados de "([^"]*)"$`, func(proyecto string) error {
		lista, err := m.sprints.Cerrados(domain.ProyectoID(m.id("proyecto", proyecto)))
		if err == nil {
			cerrados = lista
		}
		return m.resultado(err)
	})

	sc.Step(`^veo los sprints cerrados:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		columnas, err := columnasDeTabla(tabla, "sprint", "sprint goal", "story points planificados", "story points completados")
		if err != nil {
			return err
		}
		filas := tabla.Rows[1:]
		if len(cerrados) != len(filas) {
			return fmt.Errorf("hay %d sprints cerrados, se esperaban %d", len(cerrados), len(filas))
		}
		for i, fila := range filas {
			valor := func(c string) string { return fila.Cells[columnas[c]].Value }
			planificados, err := strconv.Atoi(strings.TrimSpace(valor("story points planificados")))
			if err != nil {
				return fmt.Errorf("story points planificados de %s: %w", valor("sprint"), err)
			}
			completados, err := strconv.Atoi(strings.TrimSpace(valor("story points completados")))
			if err != nil {
				return fmt.Errorf("story points completados de %s: %w", valor("sprint"), err)
			}
			sp := cerrados[i]
			if string(sp.ID) != m.id("sprint", valor("sprint")) || sp.Goal != valor("sprint goal") ||
				sp.SPPlanificados != planificados || sp.SPCompletados != completados {
				return fmt.Errorf("sprint %d = %s %q %d/%d, se esperaba %s %q %d/%d", i+1, sp.ID, sp.Goal,
					sp.SPPlanificados, sp.SPCompletados, valor("sprint"), valor("sprint goal"), planificados, completados)
			}
		}
		return nil
	})

	sc.Step(`^no hay sprints cerrados$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		if len(cerrados) != 0 {
			return fmt.Errorf("hay %d sprints cerrados, se esperaba ninguno", len(cerrados))
		}
		return nil
	})

	sc.Step(`^veo el detalle del sprint "([^"]*)"$`, func(nombre string) error {
		d, err := m.sprints.Detalle(domain.SprintID(m.id("sprint", nombre)))
		if err == nil {
			detalle = d
		}
		return m.resultado(err)
	})

	sc.Step(`^el detalle tiene las historias:$`, func(tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		columnas, err := columnasDeTabla(tabla, "historia", "story points", "completada")
		if err != nil {
			return err
		}
		filas := tabla.Rows[1:]
		if len(detalle.Historias) != len(filas) {
			return fmt.Errorf("el detalle tiene %d historias, se esperaban %d", len(detalle.Historias), len(filas))
		}
		for i, fila := range filas {
			valor := func(c string) string { return fila.Cells[columnas[c]].Value }
			storyPoints, err := strconv.Atoi(strings.TrimSpace(valor("story points")))
			if err != nil {
				return fmt.Errorf("story points de %s: %w", valor("historia"), err)
			}
			completada, err := completadaDesdeTexto(valor("completada"))
			if err != nil {
				return err
			}
			h := detalle.Historias[i]
			if string(h.ID) != m.id("historia", valor("historia")) || h.StoryPoints != storyPoints || h.Completada != completada {
				return fmt.Errorf("historia %d = %s %d SP completada=%v, se esperaba %s %d SP completada=%v", i+1,
					h.ID, h.StoryPoints, h.Completada, valor("historia"), storyPoints, completada)
			}
		}
		return nil
	})
}

// columnasDeTabla indexa las columnas de una tabla de godog por su encabezado y controla que estén las requeridas.
func columnasDeTabla(tabla *godog.Table, requeridas ...string) (map[string]int, error) {
	if len(tabla.Rows) == 0 {
		return nil, fmt.Errorf("la tabla está vacía")
	}
	columnas := make(map[string]int, len(tabla.Rows[0].Cells))
	for i, celda := range tabla.Rows[0].Cells {
		columnas[strings.ToLower(strings.TrimSpace(celda.Value))] = i
	}
	for _, c := range requeridas {
		if _, ok := columnas[c]; !ok {
			return nil, fmt.Errorf("a la tabla le falta la columna %q", c)
		}
	}
	return columnas, nil
}
