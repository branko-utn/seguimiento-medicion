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

func inicializarEscenario(sc *godog.ScenarioContext) {
	registrarPasosHTTP(sc)
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
