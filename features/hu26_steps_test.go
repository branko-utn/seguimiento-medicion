package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
	"github.com/cucumber/godog"
)

func registrarPasosHU26(sc *godog.ScenarioContext, m *mundo) {
	var horasReales float64
	var desviacion float64

	sc.Step(`^registros de esfuerzo de ([\d.]+), ([\d.]+) y ([\d.]+) horas$`, func(h1, h2, h3 string) error {
		v1, err := strconv.ParseFloat(h1, 64)
		if err != nil {
			return err
		}
		v2, err := strconv.ParseFloat(h2, 64)
		if err != nil {
			return err
		}
		v3, err := strconv.ParseFloat(h3, 64)
		if err != nil {
			return err
		}
		horasReales = metrics.HorasReales([]float64{v1, v2, v3})
		return nil
	})

	sc.Step(`^que no hay registros de esfuerzo$`, func() error {
		horasReales = metrics.HorasReales(nil)
		return nil
	})

	sc.Step(`^las horas reales son ([\d.]+)$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		esperado, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("horas reales esperadas %q inválidas: %w", texto, err)
		}
		if fmt.Sprintf("%.1f", horasReales) != fmt.Sprintf("%.1f", esperado) {
			return fmt.Errorf("horas reales = %.1f, se esperaban %.1f", horasReales, esperado)
		}
		return nil
	})

	sc.Step(`^calculo la desviación con ([\d.]+) horas estimadas y ([\d.]+) reales$`, func(estimadasTxt, realesTxt string) error {
		estimadas, err := strconv.ParseFloat(estimadasTxt, 64)
		if err != nil {
			return m.resultado(err)
		}
		reales, err := strconv.ParseFloat(realesTxt, 64)
		if err != nil {
			return m.resultado(err)
		}
		d, err := metrics.Desviacion(estimadas, reales)
		if err != nil {
			return m.resultado(err)
		}
		desviacion = d
		return m.resultado(nil)
	})

	sc.Step(`^la desviación es (-?\d+(?:\.\d+)?) \%$`, func(texto string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		esperado, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			return fmt.Errorf("desviación esperada %q inválida: %w", texto, err)
		}
		if fmt.Sprintf("%.1f", desviacion) != fmt.Sprintf("%.1f", esperado) {
			return fmt.Errorf("desviación = %.1f%%, se esperaba %.1f%%", desviacion, esperado)
		}
		return nil
	})
}
