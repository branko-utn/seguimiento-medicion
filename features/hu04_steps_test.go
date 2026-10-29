package features

import "github.com/cucumber/godog"

// registrarPasosHU04 registra los pasos de "Consultar estado del proyecto" (HU-04).
func registrarPasosHU04(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^consulto el estado del proyecto "([^"]*)"$`, m.consultoElEstadoDelProyecto)
	sc.Step(`^el estado del proyecto es "([^"]*)"$`, m.elEstadoDelProyectoEs)
	sc.Step(`^quedan (\d+) días para el fin del proyecto$`, m.quedanDiasParaElFin)
	sc.Step(`^el proyecto va del "([^"]*)" al "([^"]*)"$`, m.elProyectoVaDelAl)
	sc.Step(`^el sprint activo del proyecto es "([^"]*)"$`, m.elSprintActivoDelProyectoEs)
	sc.Step(`^el proyecto no tiene sprint activo$`, m.elProyectoNoTieneSprintActivo)
	sc.Step(`^el porcentaje de historias completadas del proyecto es (\d+(?:\.\d+)?)$`, m.elPorcentajeDelProyectoEs)
}

func (m *mundo) consultoElEstadoDelProyecto(proyecto string) error { return godog.ErrPending }

func (m *mundo) elEstadoDelProyectoEs(estado string) error { return godog.ErrPending }

func (m *mundo) quedanDiasParaElFin(dias int) error { return godog.ErrPending }

func (m *mundo) elProyectoVaDelAl(inicio, fin string) error { return godog.ErrPending }

func (m *mundo) elSprintActivoDelProyectoEs(nombre string) error { return godog.ErrPending }

func (m *mundo) elProyectoNoTieneSprintActivo() error { return godog.ErrPending }

func (m *mundo) elPorcentajeDelProyectoEs(texto string) error { return godog.ErrPending }
