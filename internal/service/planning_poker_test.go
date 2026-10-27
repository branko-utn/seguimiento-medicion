package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServicioPoker arma el Planning Poker con el proyecto "Demo" (P1) y sus integrantes "Ana Pérez" (I1)
// y "Beto Gómez" (I2), el proyecto "Otro" (P2) con "Carla Ruiz" (I3) y la historia "Login" (HU-1) en el
// backlog de "Demo". Devuelve también el repositorio de historias para preparar casos (como completar la
// historia directamente, igual que los pasos "Dado" de los escenarios).
func nuevoServicioPoker(t *testing.T) (*PlanningPoker, *memory.Historias) {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoIntegrantes, err := memory.NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}
	repoHistorias, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	repoSprints, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	repoSesiones, err := memory.NuevoSesiones("")
	if err != nil {
		t.Fatal(err)
	}

	proyectos := NuevoProyectos(repoProyectos)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	integrantes := NuevoIntegrantes(repoIntegrantes, repoProyectos)
	for _, nombre := range []string{"Ana Pérez", "Beto Gómez"} {
		if _, err := integrantes.Registrar("P1", nombre, project.ProductBuilder); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := integrantes.Registrar("P2", "Carla Ruiz", project.ProductBuilder); err != nil {
		t.Fatal(err)
	}

	back := NuevoBacklog(repoHistorias, repoProyectos, NuevoSprints(repoSprints, repoProyectos, repoHistorias))
	if _, err := back.Crear("P1", "Login", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	return NuevoPlanningPoker(repoSesiones, repoHistorias, repoIntegrantes), repoHistorias
}

func TestPlanningPoker_IniciarVotacion(t *testing.T) {
	s, _ := nuevoServicioPoker(t)

	sesion, err := s.IniciarVotacion("HU-1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// Salidas esperadas de la spec HU-16: una sesión abierta en la ronda 1, sin votos.
	if sesion.HistoriaID != "HU-1" || sesion.Ronda != 1 {
		t.Errorf("IniciarVotacion() = %s · ronda %d, se esperaba HU-1 · ronda 1", sesion.HistoriaID, sesion.Ronda)
	}
	if votantes, err := s.QuienesVotaron("HU-1"); err != nil || len(votantes) != 0 {
		t.Errorf("la votación nueva no tiene votos, se obtuvieron %v · %v", votantes, err)
	}
}

func TestPlanningPoker_IniciarVotacion_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		historia  domain.HistoriaID
		quiereErr error
	}{
		{"historia inexistente", "HU-99", backlog.ErrNoExiste},
		{"ya hay una votación abierta", "HU-1", estimation.ErrVotacionAbierta},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s, _ := nuevoServicioPoker(t)
			// Preparación: HU-1 ya tiene una votación abierta (RN-16.5).
			if _, err := s.IniciarVotacion("HU-1"); err != nil {
				t.Fatal(err)
			}

			if _, err := s.IniciarVotacion(caso.historia); !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}
		})
	}
}

func TestPlanningPoker_IniciarVotacion_historiaCompletada(t *testing.T) {
	// RN-16.5: no se vota una historia completada. La historia se completa directamente en el
	// repositorio (preparación, como los pasos "Dado"): completarla por la aplicación exige sprint (HU-08).
	s, repoHistorias := nuevoServicioPoker(t)
	h, err := repoHistorias.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	h.Estado = backlog.Completada
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}

	if _, err := s.IniciarVotacion("HU-1"); !errors.Is(err, estimation.ErrHistoriaCompletada) {
		t.Fatalf("se esperaba ErrHistoriaCompletada, se obtuvo %v", err)
	}
	// La votación rechazada no queda guardada.
	if _, err := s.QuienesVotaron("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Errorf("no debería quedar una votación abierta, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Votar(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}

	// El voto queda registrado en la ronda actual (spec HU-16, salidas esperadas).
	if err := s.Votar("HU-1", "I1", 5); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	// Solo se puede consultar quién votó, en el orden en que votaron (RN-16.4).
	votantes, err := s.QuienesVotaron("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	quiere := []domain.IntegranteID{"I1", "I2"}
	if len(votantes) != len(quiere) {
		t.Fatalf("hay %d votantes, se esperaban %d", len(votantes), len(quiere))
	}
	for n, id := range quiere {
		if votantes[n] != id {
			t.Errorf("votante %d = %s, se esperaba %s", n+1, votantes[n], id)
		}
	}
}

func TestPlanningPoker_Votar_errores(t *testing.T) {
	// Ana (I1) ya votó antes de cada caso: así "ya votó" tiene un voto previo que encontrar, y los
	// otros casos prueban su condición sin chocar con el doble voto de I1.
	casos := []struct {
		nombre     string
		integrante domain.IntegranteID
		carta      int
		quiereErr  error
	}{
		{"carta fuera de la escala", "I2", 4, estimation.ErrStoryPointsInvalidos},
		{"integrante ya votó en esta ronda", "I1", 8, estimation.ErrYaVoto},
		{"integrante inexistente", "I99", 5, project.ErrIntegranteNoExiste},
		{"integrante de otro proyecto", "I3", 5, estimation.ErrNoEsDelProyecto},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s, _ := nuevoServicioPoker(t)
			if _, err := s.IniciarVotacion("HU-1"); err != nil {
				t.Fatal(err)
			}
			if err := s.Votar("HU-1", "I1", 5); err != nil {
				t.Fatal(err)
			}

			if err := s.Votar("HU-1", caso.integrante, caso.carta); !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}

			// El voto inválido no queda registrado: sigue habiendo un solo votante.
			votantes, err := s.QuienesVotaron("HU-1")
			if err != nil {
				t.Fatal(err)
			}
			if len(votantes) != 1 {
				t.Errorf("hay %d votantes, se esperaba 1: el voto inválido no debe quedar registrado", len(votantes))
			}
		})
	}
}

func TestPlanningPoker_Votar_historiaInexistente(t *testing.T) {
	s, _ := nuevoServicioPoker(t)

	if err := s.Votar("HU-99", "I1", 5); !errors.Is(err, backlog.ErrNoExiste) {
		t.Fatalf("se esperaba backlog.ErrNoExiste, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Votar_sinVotacionAbierta(t *testing.T) {
	// Caso no cubierto por la spec (la entrada de "Votar" es la votación): error acordado por el equipo
	// en la planning del Sprint 2. El mensaje es el que verá la pantalla.
	s, _ := nuevoServicioPoker(t)

	if err := s.Votar("HU-1", "I1", 5); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
	}
	if estimation.ErrSinVotacionAbierta.Error() != "la historia no tiene una votación abierta" {
		t.Errorf("mensaje de ErrSinVotacionAbierta = %q", estimation.ErrSinVotacionAbierta.Error())
	}
}

func TestPlanningPoker_QuienesVotaron_sinVotacionAbierta(t *testing.T) {
	s, _ := nuevoServicioPoker(t)

	if _, err := s.QuienesVotaron("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
	}
}

func TestPlanningPoker_ConsultarVotos_ocultos(t *testing.T) {
	// RN-16.4: mientras la votación está abierta, los valores de los votos no se pueden consultar.
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 5); err != nil {
		t.Fatal(err)
	}

	if _, err := s.ConsultarVotos("HU-1"); !errors.Is(err, estimation.ErrVotosOcultos) {
		t.Fatalf("se esperaba ErrVotosOcultos, se obtuvo %v", err)
	}
}

func TestPlanningPoker_ConsultarVotos_sinVotacionAbierta(t *testing.T) {
	s, _ := nuevoServicioPoker(t)

	if _, err := s.ConsultarVotos("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
	}
}
