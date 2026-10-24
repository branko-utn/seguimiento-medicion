package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/cucumber/godog"
)

// registrarPasosHU07 registra los pasos de "Criterios de aceptación de una historia" (HU-07).
func registrarPasosHU07(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^agrego a "([^"]*)" el criterio "([^"]*)"$`, m.agregoElCriterio)
	sc.Step(`^quito de "([^"]*)" el criterio (\d+)$`, m.quitoElCriterio)
	sc.Step(`^que la historia "([^"]*)" tiene los criterios "([^"]*)"$`, m.queLaHistoriaTieneLosCriterios)
	sc.Step(`^la historia "([^"]*)" tiene los criterios:$`, m.laHistoriaTieneLosCriterios)
	sc.Step(`^la historia "([^"]*)" no tiene criterios$`, m.laHistoriaNoTieneCriterios)
}

// agregoElCriterio ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
func (m *mundo) agregoElCriterio(historia, criterio string) error {
	err := m.backlog.AgregarCriterio(domain.HistoriaID(m.id("historia", historia)), criterio)
	return m.resultado(err)
}

// quitoElCriterio ejecuta el caso de uso y guarda su resultado ("Cuando").
func (m *mundo) quitoElCriterio(historia string, numero int) error {
	err := m.backlog.QuitarCriterio(domain.HistoriaID(m.id("historia", historia)), numero)
	return m.resultado(err)
}

// queLaHistoriaTieneLosCriterios agrega a la historia los criterios de la lista "A, B, C" ("Dado").
func (m *mundo) queLaHistoriaTieneLosCriterios(historia, criterios string) error {
	for _, texto := range lista(criterios) {
		if err := m.backlog.AgregarCriterio(domain.HistoriaID(m.id("historia", historia)), texto); err != nil {
			return err
		}
	}
	return nil
}

// laHistoriaTieneLosCriterios verifica la tabla número/criterio ("Entonces"). El número de la tabla debe coincidir
// con la posición del criterio: se identifican por su posición, empezando en 1 (RN-07.3).
func (m *mundo) laHistoriaTieneLosCriterios(historia string, tabla *godog.Table) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if len(tabla.Rows) < 2 {
		return fmt.Errorf("la tabla de criterios no tiene filas")
	}
	columnas := map[string]int{}
	for i, celda := range tabla.Rows[0].Cells {
		columnas[celda.Value] = i
	}
	for _, c := range []string{"número", "criterio"} {
		if _, ok := columnas[c]; !ok {
			return fmt.Errorf("a la tabla de criterios le falta la columna %q", c)
		}
	}
	filas := tabla.Rows[1:]
	if len(filas) != len(h.Criterios) {
		return fmt.Errorf("la historia %q tiene %d criterios, la tabla tiene %d filas", historia, len(h.Criterios), len(filas))
	}
	for i, fila := range filas {
		numero := i + 1
		if fila.Cells[columnas["número"]].Value != fmt.Sprintf("%d", numero) {
			return fmt.Errorf("en la tabla, el criterio %d aparece como %q", numero, fila.Cells[columnas["número"]].Value)
		}
		if texto := fila.Cells[columnas["criterio"]].Value; h.Criterios[i].Texto != texto {
			return fmt.Errorf("el criterio %d es %q, se esperaba %q", numero, h.Criterios[i].Texto, texto)
		}
	}
	return nil
}

// laHistoriaNoTieneCriterios verifica que la historia no tenga ninguno ("Entonces").
func (m *mundo) laHistoriaNoTieneCriterios(historia string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if len(h.Criterios) != 0 {
		return fmt.Errorf("la historia %q tiene %d criterios, se esperaba 0", historia, len(h.Criterios))
	}
	return nil
}
