package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU25 registra los pasos de Porcentaje de historias completadas.
func registrarPasosHU25(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^calculo el porcentaje con (\d+) historias completadas de (\d+)$`, func(completadas, total int) error {
		p, err := metrics.PorcentajeCompletadas(completadas, total)
		return m.resultado(err, porcentajeGuardado{valor: p})
	})
	sc.Step(`^el porcentaje de historias completadas es (-?\d+(?:\.\d+)?)$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		pg, ok := m.dato.(porcentajeGuardado)
		if !ok {
			return fmt.Errorf("no hay porcentaje calculado")
		}
		esperado, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("porcentaje esperado %q inválida: %w", texto, err)
		}
		if fmt.Sprintf("%.1f", pg.valor) != fmt.Sprintf("%.1f", esperado) {
			return fmt.Errorf("porcentaje = %.1f, se esperaba %.1f", pg.valor, esperado)
		}
		return nil
	})
}

type porcentajeGuardado struct {
	valor float64
}
