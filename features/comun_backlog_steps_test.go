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
	sc.Step(`^que la historia "([^"]*)" está en estado "([^"]*)"$`, m.queLaHistoriaEstaEnEstado)
	sc.Step(`^(?:que )?el proyecto tiene las historias:$`, m.queElProyectoTieneLasHistorias)
	sc.Step(`^que la historia "([^"]*)" tiene (\d+) Story Points$`, m.queLaHistoriaTieneStoryPoints)
	sc.Step(`^la historia "([^"]*)" tiene (\d+) Story Points$`, m.laHistoriaTieneStoryPoints)
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

// buscarHistoria busca una historia por su nombre del escenario (los IDs son únicos en el sistema, RN-05.1).
func (m *mundo) buscarHistoria(nombre string) (backlog.Historia, error) {
	h, err := m.backlog.Buscar(domain.HistoriaID(m.id("historia", nombre)))
	if err != nil {
		return backlog.Historia{}, fmt.Errorf("historia %q: %w", nombre, err)
	}
	return h, nil
}

// estadoDesdeTexto convierte el texto del escenario en un Estado (ADR 0003: la conversión va en el borde).
func estadoDesdeTexto(texto string) (backlog.Estado, error) {
	for _, e := range []backlog.Estado{backlog.Pendiente, backlog.EnProgreso, backlog.Completada} {
		if e.String() == texto {
			return e, nil
		}
	}
	return 0, fmt.Errorf("estado %q desconocido", texto)
}

// filaHistoria es una fila de la tabla "el proyecto tiene las historias:".
type filaHistoria struct {
	historia, titulo       string
	storyPoints, criterios int
	prioridad              backlog.Prioridad // opcional (HU-09): el cero vale Media
}

// queElProyectoTieneLasHistorias carga las historias de la tabla en el último proyecto creado ("Dado").
// Columnas: historia, título y story points; criterios (cantidad) es opcional y vale 1 si no está, para que la
// historia se pueda asignar a un sprint (HU-11); prioridad (Alta, Media o Baja) es opcional y vale Media (HU-09).
func (m *mundo) queElProyectoTieneLasHistorias(tabla *godog.Table) error {
	if len(tabla.Rows) < 2 {
		return fmt.Errorf("la tabla de historias no tiene filas")
	}
	columnas := map[string]int{}
	for i, celda := range tabla.Rows[0].Cells {
		columnas[celda.Value] = i
	}
	for _, c := range []string{"historia", "título", "story points"} {
		if _, ok := columnas[c]; !ok {
			return fmt.Errorf("a la tabla de historias le falta la columna %q", c)
		}
	}
	var filas []filaHistoria
	for _, fila := range tabla.Rows[1:] {
		valor := func(c string) string { return fila.Cells[columnas[c]].Value }
		f := filaHistoria{historia: valor("historia"), titulo: valor("título"), criterios: 1}
		if _, err := fmt.Sscan(valor("story points"), &f.storyPoints); err != nil {
			return fmt.Errorf("story points de %s: %w", f.historia, err)
		}
		if _, ok := columnas["criterios"]; ok {
			if _, err := fmt.Sscan(valor("criterios"), &f.criterios); err != nil {
				return fmt.Errorf("criterios de %s: %w", f.historia, err)
			}
		}
		if _, ok := columnas["prioridad"]; ok {
			if f.prioridad = prioridadDesdeTexto(valor("prioridad")); f.prioridad == 0 {
				return fmt.Errorf("prioridad de %s: %q no es Alta, Media ni Baja", f.historia, valor("prioridad"))
			}
		}
		filas = append(filas, f)
	}
	return m.cargarHistorias(filas)
}

// cargarHistorias crea las historias en el último proyecto creado, con prioridad Media si la fila no la trae, y les
// fija Story Points y criterios directamente en el repositorio: es preparación del escenario, no el caso de uso de
// HU-07 ni de HU-15.
func (m *mundo) cargarHistorias(filas []filaHistoria) error {
	proyectos, err := m.proyectos.Listar()
	if err != nil {
		return err
	}
	if len(proyectos) == 0 {
		return fmt.Errorf("no hay ningún proyecto creado para cargarle historias")
	}
	proyecto := proyectos[len(proyectos)-1]
	for _, f := range filas {
		prioridad := f.prioridad
		if prioridad == 0 {
			prioridad = backlog.PrioridadMedia
		}
		h, err := m.backlog.Crear(proyecto.ID, f.titulo, "", prioridad)
		if err != nil {
			return err
		}
		m.recordar("historia", f.historia, string(h.ID))
		h.StoryPoints = f.storyPoints
		for i := 1; i <= f.criterios; i++ {
			h.Criterios = append(h.Criterios, backlog.Criterio{Texto: fmt.Sprintf("Criterio %d", i)})
		}
		if err := m.historias.Guardar(h); err != nil {
			return err
		}
	}
	return nil
}

// queLaHistoriaEstaEnEstado deja la historia en ese estado ("Dado"; preparación del escenario).
func (m *mundo) queLaHistoriaEstaEnEstado(historia, estado string) error {
	e, err := estadoDesdeTexto(estado)
	if err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	h.Estado = e
	return m.historias.Guardar(h)
}

// queLaHistoriaTieneStoryPoints fija los Story Points de la historia ("Dado"; preparación del escenario).
func (m *mundo) queLaHistoriaTieneStoryPoints(historia string, sp int) error {
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	h.StoryPoints = sp
	return m.historias.Guardar(h)
}

// laHistoriaTieneStoryPoints verifica los Story Points de la historia ("Entonces").
func (m *mundo) laHistoriaTieneStoryPoints(historia string, sp int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	h, err := m.buscarHistoria(historia)
	if err != nil {
		return err
	}
	if h.StoryPoints != sp {
		return fmt.Errorf("la historia %q tiene %d Story Points, se esperaba %d", historia, h.StoryPoints, sp)
	}
	return nil
}
