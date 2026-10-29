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

// TestSesion_Revelar_Consenso: RN-17.1 (se ven todos los votos) y RN-17.2 (todos iguales → consenso).
func TestSesion_Revelar_Consenso(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)
	s.Votar(domain.IntegranteID("i2"), 5)

	res, err := s.Revelar()
	if err != nil {
		t.Fatalf("Revelar error = %v", err)
	}
	if len(res.Votos) != 2 || res.Votos[0].IntegranteID != "i1" || res.Votos[0].Carta != 5 ||
		res.Votos[1].IntegranteID != "i2" || res.Votos[1].Carta != 5 {
		t.Fatalf("los votos revelados = %v, se esperaban i1:5 e i2:5 en orden", res.Votos)
	}
	if !res.Consenso {
		t.Error("con todos los votos iguales hay consenso (RN-17.2)")
	}
	if res.Maximo != 5 || res.Minimo != 5 {
		t.Errorf("consenso: máximo = %d, mínimo = %d, se esperaba 5 y 5", res.Maximo, res.Minimo)
	}
}

// TestSesion_Revelar_Diferencias: RN-17.3, sin consenso se informan el valor más alto y el más bajo con
// quiénes los votaron.
func TestSesion_Revelar_Diferencias(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Votar(domain.IntegranteID("i2"), 8)
	s.Votar(domain.IntegranteID("i3"), 5)

	res, err := s.Revelar()
	if err != nil {
		t.Fatalf("Revelar error = %v", err)
	}
	if res.Consenso {
		t.Error("con votos distintos no hay consenso")
	}
	if res.Maximo != 8 || res.Minimo != 3 {
		t.Fatalf("máximo = %d, mínimo = %d, se esperaban 8 y 3", res.Maximo, res.Minimo)
	}
	if len(res.VotaronMaximo) != 1 || res.VotaronMaximo[0] != domain.IntegranteID("i2") {
		t.Errorf("votaron el máximo %v, se esperaba [i2]", res.VotaronMaximo)
	}
	if len(res.VotaronMinimo) != 1 || res.VotaronMinimo[0] != domain.IntegranteID("i1") {
		t.Errorf("votaron el mínimo %v, se esperaba [i1]", res.VotaronMinimo)
	}
}

// TestSesion_Revelar_UnSoloVoto: caso límite de la spec, un solo voto es consenso (RN-17.2).
func TestSesion_Revelar_UnSoloVoto(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)

	res, err := s.Revelar()
	if err != nil {
		t.Fatalf("Revelar error = %v", err)
	}
	if !res.Consenso {
		t.Error("un solo voto también es consenso (RN-17.2)")
	}
	if res.Maximo != 3 || res.Minimo != 3 {
		t.Errorf("máximo = %d, mínimo = %d, se esperaba 3 y 3", res.Maximo, res.Minimo)
	}
}

// TestSesion_Revelar_VariosVotaronElMaximo: caso límite de la spec, todos los que votaron el valor más
// alto aparecen en VotaronMaximo (RN-17.3).
func TestSesion_Revelar_VariosVotaronElMaximo(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 8)
	s.Votar(domain.IntegranteID("i2"), 8)
	s.Votar(domain.IntegranteID("i3"), 2)

	res, err := s.Revelar()
	if err != nil {
		t.Fatalf("Revelar error = %v", err)
	}
	if res.Consenso {
		t.Error("con votos distintos no hay consenso")
	}
	quiereMax := []domain.IntegranteID{"i1", "i2"}
	if len(res.VotaronMaximo) != len(quiereMax) || res.VotaronMaximo[0] != quiereMax[0] || res.VotaronMaximo[1] != quiereMax[1] {
		t.Errorf("votaron el máximo %v, se esperaba %v (en orden de voto)", res.VotaronMaximo, quiereMax)
	}
	if len(res.VotaronMinimo) != 1 || res.VotaronMinimo[0] != domain.IntegranteID("i3") {
		t.Errorf("votaron el mínimo %v, se esperaba [i3]", res.VotaronMinimo)
	}
}

// TestSesion_Revelar_SinVotos: RN-17.4, solo se revela con al menos un voto.
func TestSesion_Revelar_SinVotos(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	if _, err := s.Revelar(); !errors.Is(err, ErrSinVotos) {
		t.Fatalf("se esperaba ErrSinVotos, se obtuvo %v", err)
	}
}

// TestSesion_Revelar_DosVeces: una ronda revelada no se vuelve a revelar (RN-17.4).
func TestSesion_Revelar_DosVeces(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	if _, err := s.Revelar(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Revelar(); !errors.Is(err, ErrYaRevelada) {
		t.Fatalf("se esperaba ErrYaRevelada, se obtuvo %v", err)
	}
}

// TestSesion_Votar_TrasRevelar: una ronda revelada no acepta votos nuevos (RN-17.4).
func TestSesion_Votar_TrasRevelar(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	if _, err := s.Revelar(); err != nil {
		t.Fatal(err)
	}
	if err := s.Votar(domain.IntegranteID("i2"), 5); !errors.Is(err, ErrYaRevelada) {
		t.Fatalf("se esperaba ErrYaRevelada, se obtuvo %v", err)
	}
}

// TestSesion_Votos_TrasRevelar: los valores se pueden consultar recién al revelar (spec HU-16: el error
// de Votos aplica "antes de revelar").
func TestSesion_Votos_TrasRevelar(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)

	if _, err := s.Votos(); !errors.Is(err, ErrVotosOcultos) {
		t.Fatalf("antes de revelar se esperaba ErrVotosOcultos, se obtuvo %v", err)
	}
	if _, err := s.Revelar(); err != nil {
		t.Fatal(err)
	}
	votos, err := s.Votos()
	if err != nil {
		t.Fatalf("tras revelar, Votos error = %v", err)
	}
	if len(votos) != 1 || votos[0].IntegranteID != domain.IntegranteID("i1") || votos[0].Carta != 5 {
		t.Fatalf("votos = %v, se esperaba [i1:5]", votos)
	}
}

// TestSesion_Resultado: la pantalla consulta el resultado de una ronda revelada sin volver a revelar
// (Revelar ya no puede repetirse): mientras está abierta da ErrVotosOcultos y después es idempotente.
func TestSesion_Resultado(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)
	s.Votar(domain.IntegranteID("i2"), 8)

	if _, err := s.Resultado(); !errors.Is(err, ErrVotosOcultos) {
		t.Fatalf("mientras está abierta se esperaba ErrVotosOcultos, se obtuvo %v", err)
	}
	if _, err := s.Revelar(); err != nil {
		t.Fatal(err)
	}
	primera, err := s.Resultado()
	if err != nil {
		t.Fatalf("tras revelar, Resultado error = %v", err)
	}
	segunda, err := s.Resultado()
	if err != nil {
		t.Fatalf("Resultado es idempotente, error = %v", err)
	}
	if primera.Consenso != segunda.Consenso || primera.Maximo != segunda.Maximo || primera.Minimo != segunda.Minimo {
		t.Errorf("dos consultas dieron resultados distintos: %+v y %+v", primera, segunda)
	}
	if primera.Consenso || primera.Maximo != 8 || primera.Minimo != 5 {
		t.Errorf("resultado = %+v, se esperaba sin consenso, máximo 8, mínimo 5", primera)
	}
}

// TestSesion_JSONRoundTrip_Revelada: el estado revelado se persiste: una sesión revelada y recargada
// sigue cerrada para votos nuevos y conserva el resultado.
func TestSesion_JSONRoundTrip_Revelada(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	if _, err := s.Revelar(); err != nil {
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

	if err := otra.Votar(domain.IntegranteID("i2"), 5); !errors.Is(err, ErrYaRevelada) {
		t.Fatalf("la sesión recargada sigue cerrada: se esperaba ErrYaRevelada, se obtuvo %v", err)
	}
	res, err := otra.Resultado()
	if err != nil {
		t.Fatalf("la sesión recargada conserva el resultado: %v", err)
	}
	if !res.Consenso || res.Maximo != 3 {
		t.Errorf("resultado tras recargar = %+v, se esperaba consenso en 3", res)
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

// ===== HU-18: Nueva ronda y estimación acordada =====

// TestSesion_NuevaRonda_SinRevelar: RN-18.1, solo se abre una ronda nueva después de revelar.
func TestSesion_NuevaRonda_SinRevelar(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)
	// No revelar, intentar nueva ronda
	if err := s.NuevaRonda(); !errors.Is(err, ErrRondaSinRevelar) {
		t.Fatalf("se esperaba ErrRondaSinRevelar, se obtuvo %v", err)
	}
}

// TestSesion_NuevaRonda: RN-18.1, abre una nueva ronda vacía y archiva la anterior.
func TestSesion_NuevaRonda(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Votar(domain.IntegranteID("i2"), 8)
	if _, err := s.Revelar(); err != nil {
		t.Fatal(err)
	}
	// Abrir nueva ronda
	if err := s.NuevaRonda(); err != nil {
		t.Fatalf("NuevaRonda error = %v", err)
	}
	// La ronda actual debe ser 2, sin votos
	if s.Ronda != 2 {
		t.Fatalf("Ronda = %d, se esperaba 2", s.Ronda)
	}
	if len(s.QuienesVotaron()) != 0 {
		t.Fatalf("nueva ronda debe empezar sin votos, got %v", s.QuienesVotaron())
	}
	// Historial: debe tener la ronda 1 con sus votos
	rondas := s.Rondas()
	if len(rondas) != 1 {
		t.Fatalf("historial debe tener 1 ronda, got %d", len(rondas))
	}
	if rondas[0].Numero != 1 || len(rondas[0].Votos) != 2 {
		t.Fatalf("ronda archivada = %+v, se esperaba ronda 1 con 2 votos", rondas[0])
	}
}

// TestSesion_NuevaRonda_Varias: varias rondas se archivan en orden.
func TestSesion_NuevaRonda_Varias(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	// Ronda 1
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Votar(domain.IntegranteID("i2"), 8)
	s.Revelar()
	s.NuevaRonda()
	// Ronda 2
	s.Votar(domain.IntegranteID("i1"), 5)
	s.Votar(domain.IntegranteID("i2"), 5)
	s.Revelar()
	s.NuevaRonda()
	// Ronda 3
	if s.Ronda != 3 {
		t.Fatalf("Ronda = %d, se esperaba 3", s.Ronda)
	}
	rondas := s.Rondas()
	if len(rondas) != 2 {
		t.Fatalf("historial debe tener 2 rondas, got %d", len(rondas))
	}
	if rondas[0].Numero != 1 || rondas[1].Numero != 2 {
		t.Fatalf("números de ronda = %v, se esperaba [1, 2]", []int{rondas[0].Numero, rondas[1].Numero})
	}
}

// TestSesion_Acordar: RN-18.3/RN-18.4, registra estimación acordada y cierra la votación.
func TestSesion_Acordar(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Votar(domain.IntegranteID("i2"), 8)
	s.Revelar()
	// Acordar valor (puede ser distinto a lo votado, RN-18.3)
	if err := s.Acordar(5); err != nil {
		t.Fatalf("Acordar error = %v", err)
	}
	// El dominio valida (revealed + escala); el cierre (Ronda=0 / eliminación) lo hace el servicio.
	// Historial completo preservado
	rondas := s.Rondas()
	if len(rondas) != 1 {
		t.Fatalf("historial debe tener 1 ronda, got %d", len(rondas))
	}
}

// TestSesion_Acordar_ValorFueraDeEscala: RN-18.3, el valor debe ser de la escala.
func TestSesion_Acordar_ValorFueraDeEscala(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)
	s.Revelar()
	if err := s.Acordar(4); !errors.Is(err, ErrStoryPointsInvalidos) {
		t.Fatalf("se esperaba ErrStoryPointsInvalidos, se obtuvo %v", err)
	}
}

// TestSesion_Acordar_SinRevelar: RN-18.4, solo se acuerda después de revelar.
func TestSesion_Acordar_SinRevelar(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 5)
	if err := s.Acordar(5); !errors.Is(err, ErrRondaSinRevelar) {
		t.Fatalf("se esperaba ErrRondaSinRevelar, se obtuvo %v", err)
	}
}

// TestSesion_Rondas: el historial incluye la ronda actual si está revelada.
func TestSesion_Rondas(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Revelar()
	rondas := s.Rondas()
	if len(rondas) != 1 || rondas[0].Numero != 1 {
		t.Fatalf("rondas = %+v, se esperaba 1 ronda revelada", s.Rondas())
	}
}

// TestSesion_JSONRoundTrip_ConHistorial: el historial se persiste y recupera.
func TestSesion_JSONRoundTrip_ConHistorial(t *testing.T) {
	s, _ := NuevaSesion(domain.HistoriaID("h1"))
	s.Votar(domain.IntegranteID("i1"), 3)
	s.Votar(domain.IntegranteID("i2"), 8)
	s.Revelar()
	s.NuevaRonda()
	s.Votar(domain.IntegranteID("i1"), 5)
	s.Revelar()
	s.Acordar(5)

	contenido, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}
	var otra Sesion
	if err := json.Unmarshal(contenido, &otra); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	// Tras Acordar, la ronda actual sigue siendo 2 (el cierre lo hace el servicio eliminando la sesión).
	rondas := otra.Rondas()
	if len(rondas) != 2 {
		t.Fatalf("historial debe tener 2 rondas, got %d", len(rondas))
	}
	if rondas[0].Numero != 1 || len(rondas[0].Votos) != 2 {
		t.Fatalf("ronda 1 = %+v", rondas[0])
	}
	if rondas[1].Numero != 2 || len(rondas[1].Votos) != 1 {
		t.Fatalf("ronda 2 = %+v", rondas[1])
	}
}
