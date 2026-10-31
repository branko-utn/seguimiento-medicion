// Comando demo carga nuestro propio proyecto (TEC-06): integrantes, el Product Backlog con sus criterios y
// estimaciones, los Sprints 1 a 3 cerrados y el 4 en curso, horas y defectos. Usa los mismos servicios que la
// aplicación, así los datos pasan por todas las validaciones.
//
// Uso: go run ./cmd/demo            (carga en data/, la carpeta de la aplicación)
//
//	DATA_DIR=otra go run ./cmd/demo
//
// Si el proyecto ya está cargado, no lo duplica: termina con el error del alta.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/app"
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/server"
)

func main() {
	dir := "data"
	if d := os.Getenv("DATA_DIR"); d != "" {
		dir = d
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("crear %s: %v", dir, err)
	}
	svc, err := app.NuevosServicios(dir)
	if err != nil {
		log.Fatal(err)
	}
	p, err := cargar(svc, time.Now())
	if err != nil {
		log.Fatalf("cargar la demo: %v", err)
	}
	fmt.Printf("Listo: %q (%s) cargado en %s. Levantá la app con `make run` y abrí /dashboard?proyecto=%s\n", p.Nombre, p.ID, dir, p.ID)
}

// cargar da de alta el proyecto de la demo con todos sus datos y lo devuelve. hoy es la fecha contra la que se
// validan los registros de esfuerzo (no pueden ser del futuro).
func cargar(svc server.Servicios, hoy time.Time) (project.Proyecto, error) {
	p, err := svc.Proyectos.Crear(proyectoDemo.nombre, proyectoDemo.descripcion, proyectoDemo.inicio, proyectoDemo.fin)
	if err != nil {
		return project.Proyecto{}, err
	}

	integrantes := make(map[string]domain.IntegranteID, len(integrantesDemo))
	for _, clave := range ordenIntegrantes {
		i := integrantesDemo[clave]
		alta, err := svc.Integrantes.Registrar(p.ID, i.nombre, i.rol)
		if err != nil {
			return p, fmt.Errorf("integrante %s: %w", i.nombre, err)
		}
		integrantes[clave] = alta.ID
	}

	// El backlog: cada historia con su criterio, sus Story Points y sus horas estimadas.
	ids := make([]domain.HistoriaID, len(historiasDemo))
	for n, h := range historiasDemo {
		alta, err := svc.Backlog.Crear(p.ID, h.titulo, "", prioridadDelSprint(h.sprint))
		if err != nil {
			return p, fmt.Errorf("HU-%02d: %w", n+1, err)
		}
		ids[n] = alta.ID
		if err := svc.Backlog.AgregarCriterio(alta.ID, h.criterio); err != nil {
			return p, fmt.Errorf("HU-%02d: %w", n+1, err)
		}
		if err := svc.Backlog.Estimar(alta.ID, h.sp); err != nil {
			return p, fmt.Errorf("HU-%02d: %w", n+1, err)
		}
		if _, err := svc.Estimaciones.Estimar(alta.ID, float64(h.sp)*horasPorSP); err != nil {
			return p, fmt.Errorf("HU-%02d: %w", n+1, err)
		}
	}

	// Los sprints, uno detrás del otro: se crean, se inician, se les asignan sus historias, se registran las horas
	// en sus días, se completan las historias y se cierran (el 4 queda en curso).
	sprints := make([]domain.SprintID, len(sprintsDemo))
	for s, sd := range sprintsDemo {
		sp, err := svc.Sprints.Crear(p.ID, sd.goal, sd.inicio, sd.fin)
		if err != nil {
			return p, fmt.Errorf("sprint %d: %w", s+1, err)
		}
		sprints[s] = sp.ID
		if err := svc.Sprints.Iniciar(sp.ID); err != nil {
			return p, fmt.Errorf("sprint %d: %w", s+1, err)
		}
		for n, h := range historiasDemo {
			if h.sprint != s+1 {
				continue
			}
			if err := svc.Sprints.Asignar(sp.ID, ids[n]); err != nil {
				return p, fmt.Errorf("HU-%02d al sprint %d: %w", n+1, s+1, err)
			}
			if err := registrarHoras(svc, integrantes[h.responsable], ids[n], h, sd, hoy); err != nil {
				return p, fmt.Errorf("horas de HU-%02d: %w", n+1, err)
			}
			if sd.cerrar {
				if err := svc.Sprints.CompletarHistoria(sp.ID, ids[n]); err != nil {
					return p, fmt.Errorf("completar HU-%02d: %w", n+1, err)
				}
			}
		}
		if sd.cerrar {
			if err := svc.Sprints.CerrarSprint(sp.ID); err != nil {
				return p, fmt.Errorf("cerrar el sprint %d: %w", s+1, err)
			}
		}
	}

	for _, d := range defectosDemo {
		alta, err := svc.Defectos.Registrar(d.descripcion, d.severidad, ids[d.historia-1], sprints[d.deteccion-1])
		if err != nil {
			return p, fmt.Errorf("defecto %q: %w", d.descripcion, err)
		}
		if d.resolucion > 0 {
			if err := svc.Defectos.Resolver(alta.ID, sprints[d.resolucion-1]); err != nil {
				return p, fmt.Errorf("resolver %q: %w", d.descripcion, err)
			}
		}
	}
	return p, nil
}

// registrarHoras reparte las horas reales de la historia entre los días de su sprint: 70 % de desarrollo y 30 % de
// pruebas. Las horas reales son las estimadas por el factor de la historia.
func registrarHoras(svc server.Servicios, integrante domain.IntegranteID, historia domain.HistoriaID, h historiaDemo, sd sprintDemo, hoy time.Time) error {
	reales := float64(h.sp) * horasPorSP * h.factor
	dias := int(sd.fin.Sub(sd.inicio).Hours()/24) + 1
	for _, tramo := range []struct {
		actividad string
		parte     float64
	}{{"Desarrollo", 0.7}, {"Pruebas", 0.3}} {
		porDia := redondear(reales * tramo.parte / float64(dias))
		if porDia == 0 {
			continue
		}
		for d := 0; d < dias; d++ {
			if _, err := svc.Registros.Registrar(integrante, historia, sd.inicio.AddDate(0, 0, d), tramo.actividad, porDia, hoy); err != nil {
				return err
			}
		}
	}
	return nil
}

// redondear deja las horas con un decimal.
func redondear(h float64) float64 {
	return float64(int(h*10+0.5)) / 10
}
