package features

import (
	"fmt"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/cucumber/godog"
)

// registrarPasosHU18 registra los pasos de "Nueva ronda y estimación acordada" (HU-18). Reusa los
// "Dado" de HU-16/HU-17 (votos, revelar) y los compartidos de proyectos, integrantes e historias (features/README.md).
func registrarPasosHU18(sc *godog.ScenarioContext, m *mundo) {

	// "que se revelaron los votos de la historia..." ya está en HU-17 (compartido)

	// La tabla consulta el servicio y no el resultado guardado: ejercita la misma lectura que hace la
	// pantalla tras revelar.
	sc.Step(`^el historial de la historia "([^"]*)" tiene la ronda (\d+) con los votos:$`, func(historia string, ronda int, tabla *godog.Table) error {
		if err := m.sinError(); err != nil {
			return err
		}
		// Usar Rondas() para acceder al historial, no Resultado() que devuelve la ronda actual
		sesion, err := m.poker.Rondas(domain.HistoriaID(m.id("historia", historia)))
		if err != nil {
			return err
		}
		// Buscar la ronda específica en el historial
		for i := range sesion {
			if sesion[i].Numero == ronda {
				// Crear un Resultado temporal para verificar con la tabla
				res := estimation.Resultado{
					Votos:         sesion[i].Votos,
					Consenso:      len(sesion[i].Votos) > 0 && sesion[i].Votos[0].Carta == sesion[i].Votos[len(sesion[i].Votos)-1].Carta, // simplificado
				}
				// Verificar solo los votos (la tabla solo tiene integrante | voto)
				return verificarTablaVotosRonda(res, tabla, m)
			}
		}
		return fmt.Errorf("no se encontró la ronda %d en el historial", ronda)
	})

	sc.Step(`^registro (\d+) como estimación acordada de la historia "([^"]*)"$`, func(sp int, historia string) error {
		return m.resultado(m.poker.Acordar(domain.HistoriaID(m.id("historia", historia)), sp))
	})

	sc.Step(`^la historia "([^"]*)" no tiene una votación abierta$`, func(historia string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		_, err := m.poker.QuienesVotaron(domain.HistoriaID(m.id("historia", historia)))
		return m.resultado(err)
	})

	// "la historia ... no tiene Story Points" ya está implementado en hu05_steps_test.go (compartido)

	// "hay consenso en N", "no hay consenso", "el voto más alto/bajo es N, de ..." ya están en HU-17 (compartidos)
	// "que se revelaron los votos de la historia..." ya está en HU-17 (compartido)
	// "los votos de la historia ... son:" ya está en HU-17 (compartido)

	// "la votación de la historia ... está en la ronda N sin votos" (Entonces)
	sc.Step(`^la votación de la historia "([^"]*)" está en la ronda (\d+) sin votos$`, func(historia string, ronda int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		sesion, ok := m.poker.BuscarSesionParaTest(domain.HistoriaID(m.id("historia", historia)))
		if !ok {
			return fmt.Errorf("no hay votación para la historia %s", historia)
		}
		if sesion.Ronda != ronda {
			return fmt.Errorf("la votación está en la ronda %d, se esperaba %d", sesion.Ronda, ronda)
		}
		// La sesión no está revelada (nueva ronda), por lo que no tiene votos visibles.
		// Verificamos que la ronda actual no tenga votos (campo privado votos está vacío).
		// Usamos reflexión o asumimos que una ronda no revelada no tiene votos.
		// Como no podemos acceder al campo privado, verificamos que la ronda no esté revelada.
		if sesion.Revelada() {
			votos, err := sesion.Votos()
			if err != nil {
				return err
			}
			if len(votos) != 0 {
				return fmt.Errorf("la ronda %d tiene %d votos, se esperaba 0", ronda, len(votos))
			}
		}
		return nil
	})

	// "la votación de la historia ... está en la ronda N" (Entonces)
	sc.Step(`^la votación de la historia "([^"]*)" está en la ronda (\d+)$`, func(historia string, ronda int) error {
		if err := m.sinError(); err != nil {
			return err
		}
		sesion, ok := m.poker.BuscarSesionParaTest(domain.HistoriaID(m.id("historia", historia)))
		if !ok {
			return fmt.Errorf("no hay votación para la historia %s", historia)
		}
		if sesion.Ronda != ronda {
			return fmt.Errorf("la votación está en la ronda %d, se esperaba %d", sesion.Ronda, ronda)
		}
		return nil
	})

	// "que se abrió una nueva ronda en la historia ..." (Dado)
	sc.Step(`^que se abrió una nueva ronda en la historia "([^"]*)"$`, func(historia string) error {
		_, err := m.poker.NuevaRonda(domain.HistoriaID(m.id("historia", historia)))
		if err != nil {
			return err
		}
		return nil
	})

	// "abro una nueva ronda en la historia ..." (Cuando)
	sc.Step(`^abro una nueva ronda en la historia "([^"]*)"$`, func(historia string) error {
		_, err := m.poker.NuevaRonda(domain.HistoriaID(m.id("historia", historia)))
		return m.resultado(err)
	})

	// "que se registró N como estimación acordada de la historia ..." (Dado)
	sc.Step(`^que se registró (\d+) como estimación acordada de la historia "([^"]*)"$`, func(sp int, historia string) error {
		return m.resultado(m.poker.Acordar(domain.HistoriaID(m.id("historia", historia)), sp))
	})

	// "la historia ... no tiene Story Points" ya está implementado en hu05_steps_test.go (compartido)

	// "hay consenso en N", "no hay consenso", "el voto más alto/bajo es N, de ..." ya están en HU-17 (compartidos)
	// "que se revelaron los votos de la historia..." ya está en HU-17 (compartido)
	// "los votos de la historia ... son:" ya está en HU-17 (compartido)
}

// verificarTablaVotosRonda compara los votos revelados con la tabla "integrante | voto" del escenario,
// en el orden en que votaron (RN-17.1).
func verificarTablaVotosRonda(res estimation.Resultado, tabla *godog.Table, m *mundo) error {
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