package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/cucumber/godog"
)

// Fechas de un proyecto cuando el escenario no las indica: el período del TP.
const (
	inicioProyectoPorDefecto = "2026-10-21"
	finProyectoPorDefecto    = "2026-11-02"
)

// registrarPasosComunesProyectos registra los pasos de proyectos que usan varias historias (features/README.md).
// Responsable: Federico Sosa.
func registrarPasosComunesProyectos(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^que existe un proyecto "([^"]*)"$`, m.queExisteUnProyecto)
	sc.Step(`^que existe un proyecto "([^"]*)" del "([^"]*)" al "([^"]*)"$`, m.queExisteUnProyectoDelAl)
	sc.Step(`^que "([^"]*)" es integrante del proyecto "([^"]*)"$`, m.queEsIntegranteDelProyecto)
}

// queEsIntegranteDelProyecto registra a la persona como Product Builder: el escenario no indica el rol.
func (m *mundo) queEsIntegranteDelProyecto(nombre, proyecto string) error {
	return m.registrarIntegrante(nombre, project.ProductBuilder, proyecto)
}

// registrarIntegrante registra un integrante y recuerda su ID con el nombre guardado (sin espacios sobrantes).
func (m *mundo) registrarIntegrante(nombre string, rol project.Rol, proyecto string) error {
	i, err := m.integrantes.Registrar(domain.ProyectoID(m.id("proyecto", proyecto)), nombre, rol)
	if err != nil {
		return err
	}
	m.recordar("integrante", i.Nombre, string(i.ID))
	return nil
}

// rolDesdeTexto convierte el texto del escenario en un Rol (ADR 0003: la conversión va en el borde).
// Si el texto no es un rol, devuelve el cero, que project.NuevoIntegrante rechaza con ErrRolInvalido.
func rolDesdeTexto(texto string) project.Rol {
	for _, rol := range []project.Rol{project.AgileEnabler, project.ProductBuilder, project.ProductArchitect} {
		if rol.String() == texto {
			return rol
		}
	}
	return 0
}

func (m *mundo) queExisteUnProyecto(nombre string) error {
	return m.crearProyecto(nombre, inicioProyectoPorDefecto, finProyectoPorDefecto)
}

func (m *mundo) queExisteUnProyectoDelAl(nombre, inicio, fin string) error {
	return m.crearProyecto(nombre, inicio, fin)
}

// crearProyecto crea un proyecto y recuerda su ID con el nombre guardado (sin espacios sobrantes).
func (m *mundo) crearProyecto(nombre, inicio, fin string) error {
	i, err := parseFecha(inicio)
	if err != nil {
		return err
	}
	f, err := parseFecha(fin)
	if err != nil {
		return err
	}
	p, err := m.proyectos.Crear(nombre, "", i, f)
	if err != nil {
		return err
	}
	m.recordar("proyecto", p.Nombre, string(p.ID))
	return nil
}

// proyectoLlamado busca un proyecto guardado por su nombre.
func (m *mundo) proyectoLlamado(nombre string) (project.Proyecto, error) {
	lista, err := m.proyectos.Listar()
	if err != nil {
		return project.Proyecto{}, err
	}
	for _, p := range lista {
		if p.Nombre == nombre {
			return p, nil
		}
	}
	return project.Proyecto{}, fmt.Errorf("no existe un proyecto llamado %q", nombre)
}
