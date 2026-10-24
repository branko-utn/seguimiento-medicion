package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

// registrarPasosHU25 registra los pasos de Porcentaje de historias completadas. El porcentaje calculado vive en una
// variable local: inicializarEscenario se ejecuta para cada escenario, así que no se comparte entre escenarios
// (igual que en HU-24).
func registrarPasosHU25(sc *godog.ScenarioContext, m *mundo) {
	var porcentaje float64

	sc.Step(`^calculo el porcentaje con (\d+) historias completadas de (\d+)$`, func(completadas, total int) error {
		p, err := metrics.PorcentajeCompletadas(completadas, total)
		porcentaje = p
		return m.resultado(err)
	})
	sc.Step(`^el porcentaje de historias completadas es (-?\d+(?:\.\d+)?)$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		esperado, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("porcentaje esperado %q inválido: %w", texto, err)
		}
		if fmt.Sprintf("%.1f", porcentaje) != fmt.Sprintf("%.1f", esperado) {
			return fmt.Errorf("porcentaje = %.1f, se esperaba %.1f", porcentaje, esperado)
		}
		return nil
	})
}
