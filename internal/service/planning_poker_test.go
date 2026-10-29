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

// TestPlanningPoker_Revelar: la revelada devuelve el resultado con las diferencias y queda persistida:
// la ronda no acepta votos nuevos (RN-17.4) y el resultado se puede volver a consultar.
func TestPlanningPoker_Revelar(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 5); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}

	res, err := s.Revelar("HU-1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// RN-17.1 y RN-17.3: todos los votos y las diferencias.
	if res.Consenso || res.Maximo != 8 || res.Minimo != 5 {
		t.Errorf("resultado = %+v, se esperaba sin consenso, máximo 8, mínimo 5", res)
	}
	if len(res.VotaronMaximo) != 1 || res.VotaronMaximo[0] != "I2" {
		t.Errorf("votaron el máximo %v, se esperaba [I2]", res.VotaronMaximo)
	}
	if len(res.VotaronMinimo) != 1 || res.VotaronMinimo[0] != "I1" {
		t.Errorf("votaron el mínimo %v, se esperaba [I1]", res.VotaronMinimo)
	}

	// El estado revelado queda guardado: la ronda no acepta votos nuevos (RN-17.4)…
	// (I1 ya votó: la ronda revelada se rechaza antes que el doble voto.)
	if err := s.Votar("HU-1", "I1", 13); !errors.Is(err, estimation.ErrYaRevelada) {
		t.Errorf("votar tras revelar = %v, se esperaba ErrYaRevelada", err)
	}
	// …y el resultado se puede volver a consultar sin repetir la revelada.
	otra, err := s.Resultado("HU-1")
	if err != nil {
		t.Fatalf("consultar el resultado: %v", err)
	}
	if otra.Consenso != res.Consenso || otra.Maximo != res.Maximo || otra.Minimo != res.Minimo {
		t.Errorf("Resultado() = %+v, se esperaba el mismo que Revelar(): %+v", otra, res)
	}
}

// TestPlanningPoker_Revelar_errores: los tres casos de error de la spec HU-17.
func TestPlanningPoker_Revelar_errores(t *testing.T) {
	t.Run("sin votación abierta", func(t *testing.T) {
		s, _ := nuevoServicioPoker(t)
		if _, err := s.Revelar("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
			t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
		}
	})

	t.Run("sin votos", func(t *testing.T) {
		// RN-17.4: solo se revela una ronda con al menos un voto.
		s, _ := nuevoServicioPoker(t)
		if _, err := s.IniciarVotacion("HU-1"); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Revelar("HU-1"); !errors.Is(err, estimation.ErrSinVotos) {
			t.Fatalf("se esperaba ErrSinVotos, se obtuvo %v", err)
		}
	})

	t.Run("ya revelada", func(t *testing.T) {
		s, _ := nuevoServicioPoker(t)
		if _, err := s.IniciarVotacion("HU-1"); err != nil {
			t.Fatal(err)
		}
		if err := s.Votar("HU-1", "I1", 3); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Revelar("HU-1"); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Revelar("HU-1"); !errors.Is(err, estimation.ErrYaRevelada) {
			t.Fatalf("se esperaba ErrYaRevelada, se obtuvo %v", err)
		}
	})
}

// TestPlanningPoker_Resultado: la pantalla consulta el resultado sin repetir la revelada; mientras la
// votación está abierta los valores siguen ocultos (RN-16.4).
func TestPlanningPoker_Resultado(t *testing.T) {
	t.Run("sin votación abierta", func(t *testing.T) {
		s, _ := nuevoServicioPoker(t)
		if _, err := s.Resultado("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
			t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
		}
	})

	t.Run("abierta, los valores ocultos", func(t *testing.T) {
		s, _ := nuevoServicioPoker(t)
		if _, err := s.IniciarVotacion("HU-1"); err != nil {
			t.Fatal(err)
		}
		if err := s.Votar("HU-1", "I1", 5); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Resultado("HU-1"); !errors.Is(err, estimation.ErrVotosOcultos) {
			t.Fatalf("se esperaba ErrVotosOcultos, se obtuvo %v", err)
		}
	})

	t.Run("revelada, el resultado queda a la vista", func(t *testing.T) {
		s, _ := nuevoServicioPoker(t)
		if _, err := s.IniciarVotacion("HU-1"); err != nil {
			t.Fatal(err)
		}
		if err := s.Votar("HU-1", "I1", 5); err != nil {
			t.Fatal(err)
		}
		if err := s.Votar("HU-1", "I2", 5); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Revelar("HU-1"); err != nil {
			t.Fatal(err)
		}

		res, err := s.Resultado("HU-1")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if !res.Consenso || res.Maximo != 5 {
			t.Errorf("resultado = %+v, se esperaba consenso en 5", res)
		}
	})
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

// ===== HU-18: Nueva ronda y estimación acordada =====

func TestPlanningPoker_NuevaRonda(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}

	// Abrir nueva ronda
	sesion, err := s.NuevaRonda("HU-1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// La ronda actual debe ser 2, sin votos
	if sesion.Ronda != 2 {
		t.Fatalf("Ronda = %d, se esperaba 2", sesion.Ronda)
	}
	if len(sesion.QuienesVotaron()) != 0 {
		t.Fatalf("nueva ronda debe empezar sin votos")
	}
}

func TestPlanningPoker_NuevaRonda_SinRevelar(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	// Intentar nueva ronda sin revelar
	if _, err := s.NuevaRonda("HU-1"); !errors.Is(err, estimation.ErrRondaSinRevelar) {
		t.Fatalf("se esperaba ErrRondaSinRevelar, se obtuvo %v", err)
	}
}

func TestPlanningPoker_NuevaRonda_SinVotacionAbierta(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.NuevaRonda("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Acordar(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}

	// Acordar valor (puede ser distinto a lo votado, RN-18.3)
	if err := s.Acordar("HU-1", 5); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	// La votación se cierra: no hay votación abierta
	if _, err := s.QuienesVotaron("HU-1"); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Errorf("votación cerrada no debe estar abierta: %v", err)
	}
	// La historia debe tener los Story Points acordados (via HU-15/Estimar)
}

func TestPlanningPoker_Acordar_SinRevelar(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Acordar("HU-1", 5); !errors.Is(err, estimation.ErrRondaSinRevelar) {
		t.Fatalf("se esperaba ErrRondaSinRevelar, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Acordar_ValorFueraDeEscala(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 5); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Acordar("HU-1", 4); !errors.Is(err, estimation.ErrStoryPointsInvalidos) {
		t.Fatalf("se esperaba ErrStoryPointsInvalidos, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Acordar_HistoriaCompletada(t *testing.T) {
	s, repoHistorias := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 5); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}
	// Completar la historia
	h, err := repoHistorias.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	h.Estado = backlog.Completada
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}
	// Intentar acordar
	if err := s.Acordar("HU-1", 5); !errors.Is(err, backlog.ErrHistoriaCompletada) {
		t.Fatalf("se esperaba ErrHistoriaCompletada, se obtuvo %v", err)
	}
}

func TestPlanningPoker_Rondas(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}

	rondas, err := s.Rondas("HU-1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rondas) != 1 || rondas[0].Numero != 1 || len(rondas[0].Votos) != 2 {
		t.Fatalf("rondas = %+v, se esperaba 1 ronda con 2 votos", rondas)
	}
}

// TestPlanningPoker_Acordar_IniciarVotacionSinError: RN-18.5, tras acordar se puede abrir otra votación.
func TestPlanningPoker_Acordar_IniciarVotacionSinError(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}
	// Acordar cierra la votación
	if err := s.Acordar("HU-1", 5); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// Tras acordar se puede abrir otra votación (RN-18.5)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatalf("tras acordar, IniciarVotacion debería ser válido: %v", err)
	}
}

// TestPlanningPoker_Acordar_AcordarConError: tras acordar, un segundo Acordar devuelve ErrSinVotacionAbierta.
func TestPlanningPoker_Acordar_AcordarConError(t *testing.T) {
	s, _ := nuevoServicioPoker(t)
	if _, err := s.IniciarVotacion("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar("HU-1", "I2", 8); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar("HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Acordar("HU-1", 5); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// Segundo acordar: sin votación abierta → ErrSinVotacionAbierta
	if err := s.Acordar("HU-1", 8); !errors.Is(err, estimation.ErrSinVotacionAbierta) {
		t.Fatalf("se esperaba ErrSinVotacionAbierta, se obtuvo %v", err)
	}
}
