package features

import (
	"fmt"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/cucumber/godog"
)

// registrarPasosHU03 registra los pasos de "Registrar integrantes" (HU-03).
func registrarPasosHU03(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^registro a "([^"]*)" como "([^"]*)" en el proyecto "([^"]*)"$`, m.registroA)
	sc.Step(`^que el proyecto "([^"]*)" tiene los integrantes:$`, m.queElProyectoTieneLosIntegrantes)
	sc.Step(`^el proyecto "([^"]*)" tiene los integrantes:$`, m.elProyectoTieneLosIntegrantes)
	sc.Step(`^el proyecto "([^"]*)" tiene (\d+) integrantes$`, m.elProyectoTieneNIntegrantes)
}

func (m *mundo) registroA(nombre, rol, proyecto string) error {
	return m.resultado(m.registrarIntegrante(nombre, rolDesdeTexto(rol), proyecto))
}

func (m *mundo) queElProyectoTieneLosIntegrantes(proyecto string, tabla *godog.Table) error {
	filas, err := filasDeIntegrantes(tabla)
	if err != nil {
		return err
	}
	for _, f := range filas {
		if err := m.registrarIntegrante(f.nombre, rolDesdeTexto(f.rol), proyecto); err != nil {
			return fmt.Errorf("registrar a %q: %w", f.nombre, err)
		}
	}
	return nil
}

func (m *mundo) elProyectoTieneLosIntegrantes(proyecto string, tabla *godog.Table) error {
	if err := m.sinError(); err != nil {
		return err
	}
	quiere, err := filasDeIntegrantes(tabla)
	if err != nil {
		return err
	}
	lista, err := m.integrantesDe(proyecto)
	if err != nil {
		return err
	}
	if len(lista) != len(quiere) {
		return fmt.Errorf("el proyecto %q tiene %d integrantes, se esperaban %d", proyecto, len(lista), len(quiere))
	}
	for n, f := range quiere {
		if lista[n].Nombre != f.nombre || lista[n].Rol.String() != f.rol {
			return fmt.Errorf("integrante %d: %q (%s), se esperaba %q (%s)", n+1, lista[n].Nombre, lista[n].Rol, f.nombre, f.rol)
		}
	}
	return nil
}

func (m *mundo) elProyectoTieneNIntegrantes(proyecto string, n int) error {
	if err := m.sinError(); err != nil {
		return err
	}
	lista, err := m.integrantesDe(proyecto)
	if err != nil {
		return err
	}
	if len(lista) != n {
		return fmt.Errorf("el proyecto %q tiene %d integrantes, se esperaban %d", proyecto, len(lista), n)
	}
	return nil
}

func (m *mundo) integrantesDe(proyecto string) ([]project.Integrante, error) {
	return m.integrantes.DelProyecto(domain.ProyectoID(m.id("proyecto", proyecto)))
}

type filaIntegrante struct {
	nombre, rol string
}

// filasDeIntegrantes lee una tabla con las columnas "nombre" y "rol" (la primera fila es el encabezado).
func filasDeIntegrantes(tabla *godog.Table) ([]filaIntegrante, error) {
	if len(tabla.Rows) == 0 {
		return nil, fmt.Errorf("la tabla de integrantes está vacía")
	}
	columna := map[string]int{}
	for n, celda := range tabla.Rows[0].Cells {
		columna[celda.Value] = n
	}
	colNombre, okNombre := columna["nombre"]
	colRol, okRol := columna["rol"]
	if !okNombre || !okRol {
		return nil, fmt.Errorf("la tabla de integrantes debe tener las columnas \"nombre\" y \"rol\"")
	}

	var filas []filaIntegrante
	for _, fila := range tabla.Rows[1:] {
		filas = append(filas, filaIntegrante{
			nombre: fila.Cells[colNombre].Value,
			rol:    fila.Cells[colRol].Value,
		})
	}
	return filas, nil
}
