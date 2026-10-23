package features

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU24 registra los pasos de Story Points planificados, completados y velocidad.
// Las métricas son funciones puras, así que el estado del escenario son solo sus entradas.
func registrarPasosHU24(sc *godog.ScenarioContext, m *mundo) {
	var historias []metrics.HistoriaSprint
	var sprintsCerrados []int

	sc.Step(`^un sprint con las historias:$`, func(tabla *godog.Table) error {
		hs, err := historiasDesdeTabla(tabla)
		if err != nil {
			return err
		}
		historias = hs
		return nil
	})
	sc.Step(`^un sprint sin historias$`, func() error {
		historias = nil
		return nil
	})
	sc.Step(`^los Story Points planificados son (\d+)$`, func(esperados int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if got := metrics.SPPlanificados(historias); got != esperados {
			return fmt.Errorf("SP planificados = %d, se esperaban %d", got, esperados)
		}
		return nil
	})
	sc.Step(`^los Story Points completados son (\d+)$`, func(esperados int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		if got := metrics.SPCompletados(historias); got != esperados {
			return fmt.Errorf("SP completados = %d, se esperaban %d", got, esperados)
		}
		return nil
	})
	sc.Step(`^sprints cerrados con (.+) Story Points completados$`, func(texto string) error {
		completados, err := sprintsDesdeTexto(texto)
		if err != nil {
			return err
		}
		sprintsCerrados = completados
		return nil
	})
	sc.Step(`^que no hay sprints cerrados$`, func() error {
		sprintsCerrados = nil
		return nil
	})
	sc.Step(`^la velocidad del equipo es (\d+(?:\.\d+)?)$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		esperada, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("velocidad esperada %q inválida: %w", texto, err)
		}
		got := metrics.Velocidad(sprintsCerrados)
		if fmt.Sprintf("%.1f", got) != fmt.Sprintf("%.1f", esperada) {
			return fmt.Errorf("velocidad = %.1f, se esperaba %.1f", got, esperada)
		}
		return nil
	})
}

// historiasDesdeTabla lee la tabla "story points | completada" de los escenarios.
func historiasDesdeTabla(tabla *godog.Table) ([]metrics.HistoriaSprint, error) {
	if len(tabla.Rows) == 0 {
		return nil, fmt.Errorf("la tabla de historias está vacía")
	}
	colSP, colCompletada := -1, -1
	for i, celda := range tabla.Rows[0].Cells {
		switch strings.ToLower(strings.TrimSpace(celda.Value)) {
		case "story points":
			colSP = i
		case "completada":
			colCompletada = i
		}
	}
	if colSP < 0 || colCompletada < 0 {
		return nil, fmt.Errorf("la tabla debe tener las columnas \"story points\" y \"completada\"")
	}

	historias := make([]metrics.HistoriaSprint, 0, len(tabla.Rows)-1)
	for _, fila := range tabla.Rows[1:] {
		sp, err := strconv.Atoi(strings.TrimSpace(fila.Cells[colSP].Value))
		if err != nil {
			return nil, fmt.Errorf("story points %q inválidos: %w", fila.Cells[colSP].Value, err)
		}
		completada, err := completadaDesdeTexto(fila.Cells[colCompletada].Value)
		if err != nil {
			return nil, err
		}
		historias = append(historias, metrics.HistoriaSprint{StoryPoints: sp, Completada: completada})
	}
	return historias, nil
}

// completadaDesdeTexto interpreta la columna "completada" ("sí" o "no").
func completadaDesdeTexto(texto string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(texto)) {
	case "sí", "si":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("valor de completada %q inválido: usar \"sí\" o \"no\"", texto)
	}
}

// sprintsDesdeTexto lee una lista escrita como "10, 13 y 8".
func sprintsDesdeTexto(texto string) ([]int, error) {
	texto = strings.ReplaceAll(texto, " y ", ",")
	var completados []int
	for _, parte := range strings.Split(texto, ",") {
		parte = strings.TrimSpace(parte)
		if parte == "" {
			continue
		}
		sp, err := strconv.Atoi(parte)
		if err != nil {
			return nil, fmt.Errorf("Story Points completados %q inválidos: %w", parte, err)
		}
		completados = append(completados, sp)
	}
	return completados, nil
}
