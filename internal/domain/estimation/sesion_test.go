package estimation

import (
	"encoding/json"
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

// TestNuevaSesion_AbertaEnRonda1 cubre la salida esperada de la spec HU-16:
// "Iniciar: una estimation.Sesion abierta en la ronda 1, sin votos".
func TestNuevaSesion_AbiertaEnRonda1(t *testing.T) {
	s, err := NuevaSesion(domain.HistoriaID("h1"))
	if err != nil {
		t.Fatalf("NuevaSesion error = %v", err)
	}
	if s.Ronda != 1 {
		t.Fatalf("la votación nueva debe abrirse en la ronda 1, got %d", s.Ronda)
	}
}

// TestSesion_JSONRoundTrip: la sesión se persiste completa (Store usa encoding/json y los votos no están
// exportados) y al deserializar conserva las reglas (RN-16.2 y RN-16.4 siguen vigentes).
func TestSesion_JSONRoundTrip(t *testing.T) {
	s, err := NuevaSesion(domain.HistoriaID("h1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Votar(domain.IntegranteID("i1"), 5); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar(domain.IntegranteID("i2"), 8); err != nil {
		t.Fatal(err)
	}

	contenido, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}
	var otra Sesion
	if err := json.Unmarshal(contenido, &otra); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if otra.HistoriaID != domain.HistoriaID("h1") {
		t.Fatalf("HistoriaID no se preservó: %s", otra.HistoriaID)
	}
	if otra.Ronda != 1 {
		t.Fatalf("Ronda no se preservó: %d", otra.Ronda)
	}
	got := otra.QuienesVotaron()
	if len(got) != 2 || got[0] != domain.IntegranteID("i1") || got[1] != domain.IntegranteID("i2") {
		t.Fatalf("los votos no se preservaron: %v", got)
	}
	// La sesión deserializada sigue sin exponer los valores (RN-16.4)…
	if _, err := otra.Votos(); !errors.Is(err, ErrVotosOcultos) {
		t.Fatalf("se esperaba ErrVotosOcultos, got %v", err)
	}
	// …y sigue rechazando el doble voto de la misma ronda (RN-16.2).
	if err := otra.Votar(domain.IntegranteID("i1"), 3); !errors.Is(err, ErrYaVoto) {
		t.Fatalf("se esperaba ErrYaVoto, got %v", err)
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
