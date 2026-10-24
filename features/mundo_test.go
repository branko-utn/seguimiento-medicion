package features

import (
	"fmt"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
	"github.com/cucumber/godog"
)

// mundo es el estado de un escenario (TEC-10). Se crea uno nuevo para cada escenario, así los escenarios no se
// afectan entre sí.
//
// Convención de los pasos:
//   - "Dado": prepara datos. Si falla, el escenario falla (devuelve el error).
//   - "Cuando": ejecuta la acción y guarda su error con m.resultado(err), sin fallar: puede ser el error esperado.
//   - "Entonces": verifica. "veo el error" compara el error guardado; los demás llaman primero a m.sinError().
type mundo struct {
	hoy   time.Time
	err   error             // error de la última acción ("Cuando")
	alias map[string]string // "tipo:nombre en el escenario" → ID real

	// Cada épica agrega acá su servicio y lo crea en nuevoMundo con repositorios en memoria (ruta ""),
	// por ejemplo: proyectos *service.Proyectos
	proyectos   *service.Proyectos
	integrantes *service.Integrantes
	backlog     *service.Backlog
	historias   *memory.Historias // repositorio, para preparar escenarios (pasos "Dado" de backlog)
	sprints     *service.Sprints
	repoSprints *memory.Sprints // repositorio, para preparar escenarios (paso "que el sprint … está cerrado")
}

func nuevoMundo() *mundo {
	// Con ruta "" los repositorios funcionan solo en memoria y Open no puede fallar.
	repoProyectos, _ := memory.NuevoProyectos("")
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoSprints, _ := memory.NuevoSprints("")

	return &mundo{
		hoy:         time.Now().Truncate(24 * time.Hour),
		alias:       map[string]string{},
		proyectos:   service.NuevoProyectos(repoProyectos),
		integrantes: service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		backlog:     service.NuevoBacklog(repoHistorias, repoProyectos),
		historias:   repoHistorias,
		sprints:     service.NuevoSprints(repoSprints, repoProyectos, repoHistorias),
		repoSprints: repoSprints,
	}
}

// recordar asocia el nombre que usa el escenario ("S9", "Demo") con el ID real que generó el repositorio.
// tipo separa los espacios de nombres: "proyecto", "integrante", "historia", "sprint", "defecto".
func (m *mundo) recordar(tipo, nombre, id string) {
	m.alias[tipo+":"+nombre] = id
}

// id devuelve el ID real de un nombre del escenario. Si no hay alias, devuelve el nombre tal cual: así un escenario
// puede referirse a algo que no existe ("HU-99") y obtener el error de "no existe".
func (m *mundo) id(tipo, nombre string) string {
	if id, ok := m.alias[tipo+":"+nombre]; ok {
		return id
	}
	return nombre
}

// resultado guarda el error de una acción ("Cuando") para que lo verifique un paso "Entonces".
func (m *mundo) resultado(err error) error {
	m.err = err
	return nil
}

// sinError hace fallar un paso "Entonces" si la acción anterior falló.
func (m *mundo) sinError() error {
	if m.err != nil {
		return fmt.Errorf("la acción anterior falló: %v", m.err)
	}
	return nil
}

func (m *mundo) veoElError(mensaje string) error {
	if m.err == nil {
		return fmt.Errorf("se esperaba el error %q, pero la acción no falló", mensaje)
	}
	if m.err.Error() != mensaje {
		return fmt.Errorf("se esperaba el error %q, pero fue %q", mensaje, m.err.Error())
	}
	return nil
}

func (m *mundo) hoyEs(fecha string) error {
	f, err := parseFecha(fecha)
	if err != nil {
		return err
	}
	m.hoy = f
	return nil
}

// parseFecha lee una fecha de los escenarios ("2026-10-24") como un día sin hora.
func parseFecha(s string) (time.Time, error) {
	f, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("fecha %q inválida: usar el formato AAAA-MM-DD", s)
	}
	return f, nil
}

// lista separa una lista escrita en un paso: "HU-1, HU-2" → ["HU-1", "HU-2"].
func lista(s string) []string {
	var items []string
	for _, item := range strings.Split(s, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, item)
		}
	}
	return items
}

// registrarPasosComunes registra los pasos que no pertenecen a ninguna épica (features/README.md).
func registrarPasosComunes(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^veo el error "([^"]*)"$`, m.veoElError)
	sc.Step(`^hoy es "([^"]*)"$`, m.hoyEs)
}
