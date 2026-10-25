package features

import "github.com/cucumber/godog"

// registrarPasosHU15 registra los pasos de "Estimar una historia en Story Points" (HU-15). Los pasos "Dado" y
// "Entonces" de Story Points y el de historia completada son compartidos (features/README.md).
func registrarPasosHU15(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^estimo la historia "([^"]*)" en (\d+) Story Points$`, m.estimoStoryPoints)
}

// estimoStoryPoints ejecuta el caso de uso y guarda su resultado ("Cuando": el error puede ser el esperado).
// El método no se llama estimoLaHistoria como el texto del paso, porque HU-19 (horas) ya usa ese nombre en el mundo.
func (m *mundo) estimoStoryPoints(historia string, sp int) error {
	return godog.ErrPending
}
