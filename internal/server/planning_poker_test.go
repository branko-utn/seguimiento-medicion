package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// Fechas de los proyectos de prueba (como los tests de servicio).
var (
	inicioPoker = time.Date(2026, time.October, 21, 0, 0, 0, 0, time.Local)
	finPoker    = time.Date(2026, time.November, 2, 0, 0, 0, 0, time.Local)
)

// nuevoServidorPoker arma el router con la sección de Planning Poker sobre servicios en memoria: un
// proyecto "Demo" (P1) con "Ana Pérez" (I1) y "Beto Gómez" (I2), y la historia "Login" (HU-1) en su
// backlog. Devuelve también el repositorio de historias para preparar casos (como completar una historia
// directamente, igual que los pasos "Dado" de los escenarios).
func nuevoServidorPoker(t *testing.T) (http.Handler, *memory.Historias) {
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

	svc := Servicios{
		Proyectos:   service.NuevoProyectos(repoProyectos),
		Integrantes: service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		Backlog:     service.NuevoBacklog(repoHistorias, repoProyectos, service.NuevoSprints(repoSprints, repoProyectos, repoHistorias), nil, nil),
		Poker:       service.NuevoPlanningPoker(repoSesiones, repoHistorias, repoIntegrantes),
	}
	if _, err := svc.Proyectos.Crear("Demo", "", inicioPoker, finPoker); err != nil {
		t.Fatal(err)
	}
	for _, nombre := range []string{"Ana Pérez", "Beto Gómez"} {
		if _, err := svc.Integrantes.Registrar("P1", nombre, project.ProductBuilder); err != nil {
			t.Fatal(err)
		}
	}
	// El Agile Enabler del proyecto: solo él revela (regla de interfaz, ADR 0004).
	if _, err := svc.Integrantes.Registrar("P1", "Branko Almeira", project.AgileEnabler); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Backlog.Crear("P1", "Login", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}

	router, err := NewRouter(svc)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return router, repoHistorias
}

// pedido arma un pedido con un formulario POST (application/x-www-form-urlencoded).
func pedidoFormulario(metodo, ruta string, valores url.Values) *http.Request {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(valores.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestPlanningPoker_pantalla(t *testing.T) {
	router, _ := nuevoServidorPoker(t)

	// El backlog aparece con el botón para abrir la votación de cada historia.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	if cuerpo := rec.Body.String(); !strings.Contains(cuerpo, "Login") || !strings.Contains(cuerpo, "Abrir votación") {
		t.Errorf("la página no muestra la historia con su botón para abrir la votación")
	}

	// Se abre la votación: el botón desaparece y aparece el formulario para votar.
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303 (redirección a la lista)", rec.Code)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, "Votás como") {
		t.Error("la página no ofrece el formulario para votar")
	}
	if !strings.Contains(cuerpo, "Votaron: todavía nadie") {
		t.Error("la votación recién abierta no tiene votos")
	}
	if strings.Contains(cuerpo, "Abrir votación") {
		t.Error("el botón de abrir sigue visible con la votación ya abierta")
	}

	// Ana vota: se ve su nombre entre los que votaron (RN-16.4: solo quién, no los valores).
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303 (redirección a la lista)", rec.Code)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	if cuerpo := rec.Body.String(); !strings.Contains(cuerpo, "Votaron: Ana Pérez") {
		t.Error("la página no muestra quién votó")
	}
}

// TestPlanningPoker_revelar: con votos aparece el formulario "Revelo como" (solo Agile Enablers) y al
// revelar se ven todos los votos con sus valores y las diferencias (HU-17).
func TestPlanningPoker_revelar(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}}))
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I2"}, "carta": {"8"}}))

	// Con votos, la pantalla ofrece revelar eligiendo entre los Agile Enablers del proyecto.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, "Revelar votos") || !strings.Contains(cuerpo, "Revelo como") {
		t.Error("la página no ofrece revelar los votos")
	}
	if !strings.Contains(cuerpo, "Branko Almeira") {
		t.Error("el selector de revelar no incluye al Agile Enabler")
	}

	// El Agile Enabler revela: se ven los valores, el máximo y el mínimo con quiénes votaron cada uno.
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/revelar",
		url.Values{"historia": {"HU-1"}, "revelo": {"I3"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303 (redirección a la lista)", rec.Code)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	cuerpo = rec.Body.String()
	if !strings.Contains(cuerpo, "No hay consenso") {
		t.Error("la página no informa la falta de consenso")
	}
	if !strings.Contains(cuerpo, "el voto más alto es 8 (Beto Gómez)") {
		t.Error("la página no muestra el voto más alto con quién lo votó")
	}
	if !strings.Contains(cuerpo, "el voto más bajo es 5 (Ana Pérez)") {
		t.Error("la página no muestra el voto más bajo con quién lo votó")
	}
	if !strings.Contains(cuerpo, "<td>5</td>") || !strings.Contains(cuerpo, "<td>8</td>") {
		t.Error("la página no muestra los votos revelados con sus valores")
	}
	// La ronda revelada está cerrada: ya no se vota ni se revela (RN-17.4).
	if strings.Contains(cuerpo, "Votás como") || strings.Contains(cuerpo, "Revelar votos") {
		t.Error("la ronda revelada sigue ofreciendo votar o revelar")
	}
}

// TestPlanningPoker_revelarConsenso: con todos los votos iguales, la página anuncia el consenso.
func TestPlanningPoker_revelarConsenso(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))
	for _, voto := range []url.Values{
		{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}},
		{"historia": {"HU-1"}, "integrante": {"I2"}, "carta": {"5"}},
	} {
		router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/votar", voto))
	}
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/revelar",
		url.Values{"historia": {"HU-1"}, "revelo": {"I3"}}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	if cuerpo := rec.Body.String(); !strings.Contains(cuerpo, "Hay consenso en 5.") {
		t.Errorf("la página no anuncia el consenso: %q", cuerpo)
	}
}

// TestPlanningPoker_revelar_noEsAgileEnabler: un Product Builder no puede revelar (regla de interfaz,
// ADR 0004): el pedido se rechaza con un mensaje claro.
func TestPlanningPoker_revelar_noEsAgileEnabler(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/revelar",
		url.Values{"historia": {"HU-1"}, "revelo": {"I1"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "solo un Agile Enabler puede revelar los votos") {
		t.Errorf("la página de error no muestra el mensaje de la regla: %q", rec.Body.String())
	}
}

// TestPlanningPoker_revelar_sinVotos: revelar una ronda sin votos muestra el error de la spec (RN-17.4),
// aunque la pantalla no ofrezca el botón hasta que haya votos.
func TestPlanningPoker_revelar_sinVotos(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	router.ServeHTTP(httptest.NewRecorder(), pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/revelar",
		url.Values{"historia": {"HU-1"}, "revelo": {"I3"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "no hay votos para revelar") {
		t.Errorf("la página de error no muestra el mensaje de la spec: %q", rec.Body.String())
	}
}

func TestPlanningPoker_votarDosVeces_muestraElError(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	abrir := pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}})
	router.ServeHTTP(httptest.NewRecorder(), abrir)
	voto := pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"5"}})
	router.ServeHTTP(httptest.NewRecorder(), voto)

	// RN-16.2: cada integrante vota una vez por ronda; el error de la spec se muestra en la pantalla.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, pedidoFormulario(http.MethodPost, "/planning-poker/votar",
		url.Values{"historia": {"HU-1"}, "integrante": {"I1"}, "carta": {"8"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "el integrante ya votó en esta ronda") {
		t.Errorf("la página de error no muestra el mensaje de la spec: %q", rec.Body.String())
	}
}

func TestPlanningPoker_abrirDosVeces_muestraElError(t *testing.T) {
	router, _ := nuevoServidorPoker(t)
	abrir := pedidoFormulario(http.MethodPost, "/planning-poker/abrir", url.Values{"historia": {"HU-1"}})
	router.ServeHTTP(httptest.NewRecorder(), abrir)

	// RN-16.5: una historia tiene como máximo una votación abierta.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, abrir)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "la historia ya tiene una votación abierta") {
		t.Errorf("la página de error no muestra el mensaje de la spec: %q", rec.Body.String())
	}
}

func TestPlanningPoker_historiaCompletadaNoSeVota(t *testing.T) {
	router, repoHistorias := nuevoServidorPoker(t)

	// Preparación: la historia se completa directamente en el repositorio (como los pasos "Dado"):
	// completarla por la aplicación exige sprint (HU-08).
	h, err := repoHistorias.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	h.Estado = backlog.Completada
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/planning-poker", nil))
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, "La historia está completada: no se vota.") {
		t.Error("la página no avisa que una historia completada no se vota")
	}
	if strings.Contains(cuerpo, "Abrir votación") {
		t.Error("la historia completada ofrece abrir una votación (RN-16.5)")
	}
}
