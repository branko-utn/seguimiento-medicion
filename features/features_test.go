// Package features ejecuta los escenarios BDD de esta carpeta con godog.
// Cada historia agrega sus pasos en un archivo hu<xx>_steps_test.go y los registra en inicializarEscenario.
// Los escenarios con la etiqueta @pendiente no se ejecutan (ver Options.Tags).
package features

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: inicializarEscenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  archivosFeature(t),
			Strict: true, // un paso sin implementar hace fallar el test
			// Los escenarios marcados @pendiente se escribieron antes que sus pasos (primer commit de la historia).
			// Se quita la etiqueta al implementar los pasos: los escenarios pasan a fallar (RED) hasta tener el código.
			Tags:     "~@pendiente",
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron escenarios BDD")
	}
}

// inicializarEscenario se ejecuta antes de cada escenario: crea un mundo nuevo y registra los pasos.
func inicializarEscenario(sc *godog.ScenarioContext) {
	m := nuevoMundo()
	registrarPasosComunes(sc, m)
	registrarPasosHTTP(sc, m)
	registrarPasosComunesProyectos(sc, m)
	registrarPasosComunesBacklog(sc, m)
	registrarPasosHU01(sc, m)
	registrarPasosHU02(sc, m)
	registrarPasosHU03(sc, m)
	registrarPasosComunesSprints(sc, m)
	registrarPasosCierreSprints(sc, m)
	registrarPasosHU12(sc, m)
	registrarPasosHU13(sc, m)
	registrarPasosHU10(sc, m)
	registrarPasosHU11(sc, m)
	registrarPasosHU05(sc, m)
	registrarPasosHU06(sc, m)
	registrarPasosHU07(sc, m)
	registrarPasosHU15(sc, m)
	registrarPasosHU08(sc, m)
	registrarPasosHU24(sc, m)
	registrarPasosHU25(sc, m)
	registrarPasosHU26(sc, m)
	registrarPasosHU19(sc, m)
	registrarPasosHU20(sc, m)
	registrarPasosHU22(sc, m)
	registrarPasosHU27(sc, m)
	registrarPasosHU16(sc, m)
	registrarPasosHU23(sc, m)
	registrarPasosHU17(sc, m)
	registrarPasosTEC16(sc, m)
	registrarPasosHU29(sc, m)
	// Cada historia registra acá sus pasos, por ejemplo: registrarPasosHU01(sc, m)
	// Los pasos compartidos entre historias los registra una sola vez su responsable (features/README.md).
}

// archivosFeature devuelve los .feature de la carpeta, salvo las plantillas (empiezan con "_").
func archivosFeature(t *testing.T) []string {
	t.Helper()
	todos, err := filepath.Glob("*.feature")
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, p := range todos {
		if !strings.HasPrefix(filepath.Base(p), "_") {
			paths = append(paths, p)
		}
	}
	if len(paths) == 0 {
		t.Fatal("no hay archivos .feature")
	}
	return paths
}
