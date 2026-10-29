package features

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/cucumber/godog"
)

// registrarPasosHU09 registra los pasos de "Ver el backlog ordenado por prioridad" (HU-09). El backlog que se ve
// queda en una variable del escenario, como el resumen de HU-21.
func registrarPasosHU09(sc *godog.ScenarioContext, m *mundo) {
	var visto []backlog.Historia

	veo := func(proyecto string, estado *backlog.Estado) error {
		hs, err := m.backlog.DelProyecto(domain.ProyectoID(m.id("proyecto", proyecto)), estado)
		visto = hs
		return m.resultado(err)
	}
	sc.Step(`^veo el backlog de "([^"]*)"$`, func(proyecto string) error {
		return veo(proyecto, nil)
	})
	sc.Step(`^veo el backlog de "([^"]*)" con estado "([^"]*)"$`, func(proyecto, estado string) error {
		e, err := estadoDesdeTexto(estado)
		if err != nil {
			return err
		}
		return veo(proyecto, &e)
	})
	sc.Step(`^el backlog se muestra en el orden "([^"]*)"$`, func(orden string) error {
		if err := m.sinError(); err != nil {
			return err
		}
		var want []domain.HistoriaID
		for _, nombre := range strings.Split(orden, ", ") {
			want = append(want, domain.HistoriaID(m.id("historia", nombre)))
		}
		return mismoOrden(visto, want)
	})
	sc.Step(`^el backlog se muestra vacío$`, func() error {
		if err := m.sinError(); err != nil {
			return err
		}
		return mismoOrden(visto, nil)
	})

	// Este escenario prueba la regla de orden sola, sin repositorio: las historias de la tabla no se guardan.
	var ordenadas []backlog.Historia
	sc.Step(`^ordeno las historias:$`, func(tabla *godog.Table) error {
		filas, err := filasDeTabla(tabla, "historia", "prioridad")
		if err != nil {
			return err
		}
		var hs []backlog.Historia
		for _, f := range filas {
			p := prioridadDesdeTexto(f["prioridad"])
			if p == 0 {
				return fmt.Errorf("prioridad de %s: %q no es Alta, Media ni Baja", f["historia"], f["prioridad"])
			}
			hs = append(hs, backlog.Historia{ID: domain.HistoriaID(f["historia"]), Prioridad: p})
		}
		ordenadas = backlog.Ordenar(hs)
		return nil
	})
	sc.Step(`^quedan en el orden "([^"]*)"$`, func(orden string) error {
		var want []domain.HistoriaID
		for _, id := range strings.Split(orden, ", ") {
			want = append(want, domain.HistoriaID(id))
		}
		return mismoOrden(ordenadas, want)
	})
}

// mismoOrden compara los IDs de las historias con los esperados, en orden.
func mismoOrden(hs []backlog.Historia, want []domain.HistoriaID) error {
	var got []domain.HistoriaID
	for _, h := range hs {
		got = append(got, h.ID)
	}
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("el backlog quedó en el orden %v, se esperaba %v", got, want)
	}
	return nil
}
