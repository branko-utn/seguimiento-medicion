package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/cucumber/godog"
)

// registrarPasosHU17 registra los pasos de "Revelar votos y detectar diferencias" (HU-17). Reusa los
// "Dado" de HU-16 (votos) y los compartidos de proyectos, integrantes e historias (features/README.md).
// resultado guarda el del último "revelo": los "Entonces" que no mencionan la historia (por ejemplo
// "hay consenso en 5") verifican contra él (patrón de HU-24: estado del escenario en variables locales
// capturadas por los pasos, sin tocar el mundo compartido).
func registrarPasosHU17(sc *godog.ScenarioContext, m *mundo) {
	var resultado *estimation.Resultado

	revelar := func(historia string) (estimation.Resultado, error) {
		return m.poker.Revelar(domain.HistoriaID(m.id("historia", historia)))
	}

	sc.Step(`^revelo los votos de la historia "([^"]*)"$`, func(historia string) error {
		res, err := revelar(historia)
		if err == nil {
			resultado = &res
		}
		return m.resultado(err)
	})

	// "Dado": si falla, el escenario falla.
	sc.Step(`^que se revelaron los votos de la historia "([^"]*)"$`, func(historia string) error {
		res, err := revelar(historia)
		if err != nil {
			return err
		}
		resultado = &res
		return nil
	})

	// La tabla consulta el servicio y no el resultado guardado: ejercita la misma lectura que hace la
	// pantalla tras revelar.
	sc.Step(`^los votos de la historia "([^"]*)" son:$`, func(historia string, tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		res, err := m.poker.Resultado(domain.HistoriaID(m.id("historia", historia)))
		if err != nil {
			return err
		}
		return verificarTablaVotos(res, tabla, m)
	})

	sc.Step(`^hay consenso en (\d+)$`, func(valor int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		res, err := rondaRevelada(resultado)
		if err != nil {
			return err
		}
		if !res.Consenso || res.Maximo != valor {
			return fmt.Errorf("no hay consenso en %d: consenso = %v, valor = %d", valor, res.Consenso, res.Maximo)
		}
		return nil
	})

	sc.Step(`^no hay consenso$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		res, err := rondaRevelada(resultado)
		if err != nil {
			return err
		}
		if res.Consenso {
			return fmt.Errorf("hay consenso en %d, se esperaban diferencias", res.Maximo)
		}
		return nil
	})

	sc.Step(`^el voto más alto es (\d+), de "([^"]*)"$`, func(valor int, de string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		res, err := rondaRevelada(resultado)
		if err != nil {
			return err
		}
		return verificarExtremo("alto", res.Maximo, res.VotaronMaximo, valor, de, m)
	})

	sc.Step(`^el voto más bajo es (\d+), de "([^"]*)"$`, func(valor int, de string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		res, err := rondaRevelada(resultado)
		if err != nil {
			return err
		}
		return verificarExtremo("bajo", res.Minimo, res.VotaronMinimo, valor, de, m)
	})
}

// rondaRevelada devuelve el resultado del último "revelo" del escenario.
func rondaRevelada(resultado *estimation.Resultado) (estimation.Resultado, error) {
	if resultado == nil {
		return estimation.Resultado{}, fmt.Errorf("este escenario no reveló ninguna votación")
	}
	return *resultado, nil
}

// verificarTablaVotos compara los votos revelados con la tabla "integrante | voto" del escenario,
// en el orden en que votaron (RN-17.1).
func verificarTablaVotos(res estimation.Resultado, tabla *godog.Table, m *mundo) error {
	if len(tabla.Rows) < 2 {
		return fmt.Errorf("la tabla de votos no tiene filas")
	}
	if len(res.Votos) != len(tabla.Rows)-1 {
		return fmt.Errorf("se revelaron %d votos, la tabla tiene %d", len(res.Votos), len(tabla.Rows)-1)
	}
	for n, fila := range tabla.Rows[1:] {
		nombre := fila.Cells[0].Value
		carta, err := strconv.Atoi(fila.Cells[1].Value)
		if err != nil {
			return fmt.Errorf("voto de %s: %w", nombre, err)
		}
		voto := res.Votos[n]
		if voto.IntegranteID != domain.IntegranteID(m.id("integrante", nombre)) || voto.Carta != carta {
			return fmt.Errorf("voto %d = %s:%d, se esperaba %s:%d", n+1, voto.IntegranteID, voto.Carta, nombre, carta)
		}
	}
	return nil
}

// verificarExtremo compara el valor más alto o más bajo de la ronda y todos los que lo votaron con los
// del escenario, en orden de voto (RN-17.3).
func verificarExtremo(lado string, valor int, votaron []domain.IntegranteID, quiereValor int, quiereQuienes string, m *mundo) error {
	if valor != quiereValor {
		return fmt.Errorf("el voto más %s es %d, se esperaba %d", lado, valor, quiereValor)
	}
	quiere := lista(quiereQuienes)
	if len(votaron) != len(quiere) {
		return fmt.Errorf("votaron el más %s %v, se esperaba %q", lado, votaron, quiereQuienes)
	}
	for n, quien := range quiere {
		if votaron[n] != domain.IntegranteID(m.id("integrante", quien)) {
			return fmt.Errorf("votaron el más %s %v, se esperaba %q", lado, votaron, quiereQuienes)
		}
	}
	return nil
}
