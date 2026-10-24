package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/cucumber/godog"
)

// registrarPasosComunesBacklog registra los pasos de historias que usan varias épicas (features/README.md).
// Responsable: Ignacio Bianchi.
func registrarPasosComunesBacklog(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que el backlog de "([^"]*)" tiene la historia "([^"]*)" "([^"]*)"$`, m.queElBacklogTieneLaHistoria)
	sc.Step(`^la historia "([^"]*)" está en estado "([^"]*)"$`, m.laHistoriaEstaEnEstado)
}

// prioridadDesdeTexto convierte el texto del escenario en una Prioridad (ADR 0003: la conversión va en el borde).
// Si el texto no es una prioridad, devuelve el cero, que backlog.Nueva rechaza con ErrPrioridadInvalida.
func prioridadDesdeTexto(texto string) backlog.Prioridad {
	for _, p := range []backlog.Prioridad{backlog.PrioridadAlta, backlog.PrioridadMedia, backlog.PrioridadBaja} {
		if p.String() == texto {
			return p
		}
	}
	return 0
}

// queElBacklogTieneLaHistoria crea en el backlog del proyecto la historia con ese título ("Dado"). Como el paso
// no indica prioridad, usa Media: los escenarios que lo usan no dependen de ella.
func (m *mundo) queElBacklogTieneLaHistoria(proyecto, historia, titulo string) error {
	h, err := m.backlog.Crear(domain.ProyectoID(m.id("proyecto", proyecto)), titulo, "", prioridadDesdeTexto("Media"))
	if err != nil {
		return err
	}
	m.recordar("historia", historia, string(h.ID))
	return nil
}

// laHistoriaEstaEnEstado verifica el estado de la historia ("Entonces").
func (m *mundo) laHistoriaEstaEnEstado(historia, estado string) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if h.Estado.String() != estado {
		return fmt.Errorf("la historia %q está en estado %q, se esperaba %q", historia, h.Estado, estado)
	}
	return nil
}

// buscarHistoria busca una historia por su nombre del escenario. Si el mismo ID existe en más de un proyecto
// (cada proyecto numera desde HU-1, RN-05.1), el paso debe indicar el proyecto.
func (m *mundo) buscarHistoria(nombre string) (backlog.Historia, error) {
	id := domain.HistoriaID(m.id("historia", nombre))
	lista, err := m.backlog.Listar()
	if err != nil {
		return backlog.Historia{}, err
	}
	var encontradas []backlog.Historia
	for _, h := range lista {
		if h.ID == id {
			encontradas = append(encontradas, h)
		}
	}
	switch len(encontradas) {
	case 1:
		return encontradas[0], nil
	case 0:
		return backlog.Historia{}, fmt.Errorf("no existe la historia %q", nombre)
	default:
		return backlog.Historia{}, fmt.Errorf("la historia %q está en más de un proyecto, indicá el proyecto en el paso", nombre)
	}
}
