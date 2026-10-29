package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// escenarioSprint arma un proyecto con un sprint activo S1 que tiene HU-1 asignada, y HU-2 sin asignar.
func escenarioSprint(t *testing.T) (http.Handler, Servicios) {
	t.Helper()
	proyectos, _ := memory.NuevoProyectos("")
	historias, _ := memory.NuevoHistorias("")
	sprints, _ := memory.NuevoSprints("")
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(proyectos),
		Backlog:   service.NuevoBacklog(historias, proyectos, sprintsSvc, nil, nil),
		Sprints:   sprintsSvc,
	}
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	p, err := svc.Proyectos.Crear("Demo", "", dia(21), dia(31))
	if err != nil {
		t.Fatal(err)
	}
	for _, titulo := range []string{"Crear proyecto", "Crear sprint"} {
		h, err := svc.Backlog.Crear(p.ID, titulo, "", backlog.PrioridadAlta)
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.Backlog.AgregarCriterio(h.ID, "funciona"); err != nil {
			t.Fatal(err)
		}
		if err := svc.Backlog.Estimar(h.ID, 3); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.Sprints.Crear(p.ID, "MVP", dia(23), dia(25)); err != nil {
		t.Fatal(err)
	}
	if err := svc.Sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}
	return router, svc
}

func enviar(router http.Handler, ruta string, htmx bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, ruta, nil)
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCompletarHistoria(t *testing.T) {
	router, svc := escenarioSprint(t)

	rec := enviar(router, "/sprints/S1/historias/HU-1/completar", false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/sprints/S1" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /sprints/S1", rec.Code, rec.Header().Get("Location"))
	}
	if h, _ := svc.Backlog.Buscar("HU-1"); h.Estado != backlog.Completada {
		t.Errorf("HU-1 quedó en %s, se esperaba Completada", h.Estado)
	}
}

func TestCompletarHistoria_conHTMX(t *testing.T) {
	router, _ := escenarioSprint(t)

	rec := enviar(router, "/sprints/S1/historias/HU-1/completar", true)

	if rec.Code != http.StatusNoContent || rec.Header().Get("HX-Redirect") != "/sprints/S1" {
		t.Fatalf("respuesta = %d (HX-Redirect %q), se esperaba 204 con HX-Redirect a /sprints/S1", rec.Code, rec.Header().Get("HX-Redirect"))
	}
}

func TestCerrarSprint(t *testing.T) {
	router, svc := escenarioSprint(t)

	rec := enviar(router, "/sprints/S1/cerrar", false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/sprints/S1" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /sprints/S1", rec.Code, rec.Header().Get("Location"))
	}
	if h, _ := svc.Backlog.Buscar("HU-1"); h.SprintID != "" {
		t.Errorf("HU-1 no volvió al backlog: sprint %q", h.SprintID)
	}
}

func TestAccionesSprint_errores(t *testing.T) {
	casos := []struct {
		nombre  string
		ruta    string
		codigo  int
		mensaje string
	}{
		{"historia no asignada", "/sprints/S1/historias/HU-2/completar", http.StatusBadRequest, "la historia no está asignada a este sprint"},
		{"historia inexistente", "/sprints/S1/historias/HU-9/completar", http.StatusNotFound, "la historia no existe"},
		{"sprint inexistente al completar", "/sprints/S9/historias/HU-1/completar", http.StatusNotFound, "el sprint no existe"},
		{"sprint inexistente al cerrar", "/sprints/S9/cerrar", http.StatusNotFound, "el sprint no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			router, _ := escenarioSprint(t)
			rec := enviar(router, c.ruta, false)
			if rec.Code != c.codigo {
				t.Fatalf("código = %d, se esperaba %d", rec.Code, c.codigo)
			}
			if !strings.Contains(rec.Body.String(), c.mensaje) {
				t.Errorf("la página de error no muestra %q", c.mensaje)
			}
		})
	}
}

func TestCerrarSprint_yaCerrado(t *testing.T) {
	router, svc := escenarioSprint(t)
	if err := svc.Sprints.CerrarSprint(domain.SprintID("S1")); err != nil {
		t.Fatal(err)
	}
	rec := enviar(router, "/sprints/S1/cerrar", false)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "el sprint está cerrado y no se puede modificar") {
		t.Errorf("respuesta = %d, se esperaba 400 con el mensaje de sprint cerrado", rec.Code)
	}
}

// La pantalla del sprint (TEC-14) ofrece las acciones de HU-12 y HU-13 solo mientras el sprint está activo.
func TestSprints_botonesDeCompletarYCerrar(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)
	if err := e.svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
	completar := `action="/sprints/S1/historias/HU-1/completar"`
	cerrar := `action="/sprints/S1/cerrar"`

	t.Run("planificado: sin acciones", func(t *testing.T) {
		cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil)).Body.String()
		for _, accion := range []string{completar, cerrar} {
			if strings.Contains(cuerpo, accion) {
				t.Errorf("un sprint planificado no debe ofrecer %s", accion)
			}
		}
	})

	if err := e.svc.Sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	t.Run("activo: completar cada historia y cerrar", func(t *testing.T) {
		cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil)).Body.String()
		contieneTodo(t, cuerpo, completar, cerrar, "0 Story Points completados")
	})

	if rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/historias/HU-1/completar", url.Values{})); rec.Code != http.StatusSeeOther {
		t.Fatalf("completar: código = %d, se esperaba 303", rec.Code)
	}
	t.Run("activo con la historia completada", func(t *testing.T) {
		cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil)).Body.String()
		contieneTodo(t, cuerpo, "✔ Completada", cerrar, "3 Story Points completados")
		if strings.Contains(cuerpo, completar) {
			t.Error("una historia completada no debe ofrecer el botón para completarla")
		}
	})

	if rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/cerrar", url.Values{})); rec.Code != http.StatusSeeOther {
		t.Fatalf("cerrar: código = %d, se esperaba 303", rec.Code)
	}
	t.Run("cerrado: sin acciones y con el resultado", func(t *testing.T) {
		cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil)).Body.String()
		contieneTodo(t, cuerpo, "Cerrado", "3 Story Points planificados", "3 Story Points completados")
		for _, accion := range []string{completar, cerrar} {
			if strings.Contains(cuerpo, accion) {
				t.Errorf("un sprint cerrado no debe ofrecer %s", accion)
			}
		}
	})
}
