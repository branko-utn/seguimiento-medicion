package estimation

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

func TestNuevaSesion(t *testing.T) {
	s, err := NuevaSesion(domain.HistoriaID("h1"))
	if err != nil {
		t.Fatalf("NuevaSesion error = %v", err)
	}
	if s.HistoriaID != domain.HistoriaID("h1") {
		t.Fatalf("HistoriaID = %s", s.HistoriaID)
	}
	if got := s.QuienesVotaron(); len(got) != 0 {
		t.Fatalf("se esperaba lista vacía de votantes, got %v", got)
	}
}

func TestSesion_Votar_QuienesVotaron(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if err := s.Votar(domain.IntegranteID("i1"), 5); err != nil {
		t.Fatalf("Votar error = %v", err)
	}
	if err := s.Votar(domain.IntegranteID("i2"), 8); err != nil {
		t.Fatalf("Votar error = %v", err)
	}
	got := s.QuienesVotaron()
	if len(got) != 2 {
		t.Fatalf("se esperaban 2 votantes, got %d", len(got))
	}
}

func TestSesion_Votar_DosVecesMismaRonda(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if err := s.Votar(domain.IntegranteID("i1"), 5); err != nil {
		t.Fatalf("Votar error = %v", err)
	}
	if err := s.Votar(domain.IntegranteID("i1"), 8); !errors.Is(err, ErrYaVoto) {
		t.Fatalf("se esperaba ErrYaVoto, got %v", err)
	}
}

func TestSesion_Votar_CartaFueraDeEscala(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if err := s.Votar(domain.IntegranteID("i1"), 4); !errors.Is(err, ErrStoryPointsInvalidos) {
		t.Fatalf("se esperaba ErrStoryPointsInvalidos, got %v", err)
	}
}

func TestSesion_VotosOcultosMientrasAbierta(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if err := s.Votar(domain.IntegranteID("i1"), 5); err != nil {
		t.Fatalf("Votar error = %v", err)
	}
	if _, err := s.Votos(); err == nil || !errors.Is(err, ErrVotosOcultos) {
		t.Fatalf("se esperaba ErrVotosOcultos al consultar valores, got %v", err)
	}
}

func TestSesion_Votar_ExtremosEscala(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if err := s.Votar(domain.IntegranteID("i1"), 1); err != nil {
		t.Fatalf("Votar 1 error = %v", err)
	}
	s2, _ := NuevaSesion(domain.HistoriaID("h2"))
	if err := s2.Votar(domain.IntegranteID("i1"), 21); err != nil {
		t.Fatalf("Votar 21 error = %v", err)
	}
}
