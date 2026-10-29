package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// entornoBacklog es el router de la sección Backlog sobre servicios en memoria, con "Demo" (P1) y "Otro" (P2).
type entornoBacklog struct {
	router http.Handler
	svc    Servicios
}

func nuevoServidorBacklog(t *testing.T) entornoBacklog {
	t.Helper()
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoSprints, _ := memory.NuevoSprints("")
	repoIntegrantes, _ := memory.NuevoIntegrantes("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	sprints := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	registros := service.NuevoRegistros(repoRegistros, repoIntegrantes, repoHistorias, repoProyectos)
	defectos := service.NuevoDefectos(repoDefectos, repoHistorias, repoSprints)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(repoProyectos),
		Sprints:   sprints,
		Backlog:   service.NuevoBacklog(repoHistorias, repoProyectos, sprints, registros, defectos),
		Registros: registros,
		Defectos:  defectos,
	}
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := svc.Proyectos.Crear(nombre, "", fechaDeFormulario("2026-10-21"), fechaDeFormulario("2026-11-02")); err != nil {
			t.Fatal(err)
		}
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return entornoBacklog{router: router, svc: svc}
}

// crearHistoria crea HU-1 "Crear proyecto" en Demo, con prioridad Alta.
func (e entornoBacklog) crearHistoria(t *testing.T) {
	t.Helper()
	if _, err := e.svc.Backlog.Crear("P1", "Crear proyecto", "Alta de proyectos", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
}

func (e entornoBacklog) historia(t *testing.T) backlog.Historia {
	t.Helper()
	h, err := e.svc.Backlog.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestBacklog_listaEnlazaAlDetalleYAlAlta(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog", nil)).Body.String()

	contieneTodo(t, cuerpo, `href="/backlog/nueva"`, `<a href="/backlog/HU-1">Crear proyecto</a>`)
}

func TestBacklog_formularioNuevo(t *testing.T) {
	e := nuevoServidorBacklog(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/nueva", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), `action="/backlog"`, `name="titulo"`, `name="descripcion"`,
		`<option value="P1">Demo</option>`, `<option value="P2">Otro</option>`,
		"<option>Alta</option>", "<option selected>Media</option>", "<option>Baja</option>")
}

func TestBacklog_crear(t *testing.T) {
	e := nuevoServidorBacklog(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog",
		url.Values{"proyecto": {"P1"}, "titulo": {"  Crear sprint  "}, "descripcion": {"Con Sprint Goal"}, "prioridad": {"Media"}}))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/backlog/HU-1" {
		t.Errorf("redirige a %q, se esperaba /backlog/HU-1", lugar)
	}
	h := e.historia(t)
	if h.Titulo != "Crear sprint" || h.ProyectoID != "P1" || h.Prioridad != backlog.PrioridadMedia {
		t.Errorf("historia = %+v", h)
	}
}

func TestBacklog_crear_errores(t *testing.T) {
	casos := []struct {
		nombre  string
		valores url.Values
		mensaje string
	}{
		{"título vacío", url.Values{"proyecto": {"P1"}, "titulo": {"  "}, "prioridad": {"Alta"}}, "el título es obligatorio"},
		{"prioridad inválida", url.Values{"proyecto": {"P1"}, "titulo": {"Algo"}, "prioridad": {"Urgente"}}, "la prioridad debe ser Alta, Media o Baja"},
		{"proyecto inexistente", url.Values{"proyecto": {"P9"}, "titulo": {"Algo"}, "prioridad": {"Alta"}}, "el proyecto no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoServidorBacklog(t)

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog", c.valores))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// Vuelve el formulario con el mensaje de la spec y con lo que se escribió.
			contieneTodo(t, rec.Body.String(), c.mensaje, `action="/backlog"`)
			if c.valores.Get("titulo") == "Algo" {
				contieneTodo(t, rec.Body.String(), `value="Algo"`)
			}
		})
	}
}

func TestBacklog_detalle(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/HU-1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		"<h1>HU-1 · Crear proyecto</h1>", "Demo", "Alta de proyectos", "Pendiente", "Sin estimar",
		"Todavía no tiene criterios de aceptación",
		`action="/backlog/HU-1/criterios"`, `action="/backlog/HU-1/estimacion"`, `action="/backlog/HU-1/estado"`,
		`<option value="21">21</option>`)
}

func TestBacklog_historiaInexistente(t *testing.T) {
	e := nuevoServidorBacklog(t)
	for _, pedido := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/backlog/HU-9", nil),
		pedidoFormulario(http.MethodPost, "/backlog/HU-9/criterios", url.Values{"texto": {"x"}}),
		pedidoFormulario(http.MethodPost, "/backlog/HU-9/estimacion", url.Values{"story_points": {"3"}}),
		pedidoFormulario(http.MethodPost, "/backlog/HU-9/estado", url.Values{"estado": {"Pendiente"}}),
		pedidoFormulario(http.MethodPost, "/backlog/HU-9/criterios/1/quitar", url.Values{}),
	} {
		rec := servir(e.router, pedido)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s %s: código = %d, se esperaba 404", pedido.Method, pedido.URL.Path, rec.Code)
		}
		contieneTodo(t, rec.Body.String(), "la historia no existe")
	}
}

func TestBacklog_criterios(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	for _, texto := range []string{"Se guarda el nombre", "Valida las fechas"} {
		rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/criterios", url.Values{"texto": {texto}}))
		if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/backlog/HU-1" {
			t.Fatalf("agregar %q: código = %d, Location = %q", texto, rec.Code, rec.Header().Get("Location"))
		}
	}
	cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/HU-1", nil)).Body.String()
	contieneTodo(t, cuerpo, "Se guarda el nombre", "Valida las fechas", `action="/backlog/HU-1/criterios/2/quitar"`)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/criterios/1/quitar", url.Values{}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("quitar: código = %d, se esperaba 303", rec.Code)
	}
	if h := e.historia(t); len(h.Criterios) != 1 || h.Criterios[0].Texto != "Valida las fechas" {
		t.Errorf("criterios = %+v, se esperaba solo \"Valida las fechas\"", h.Criterios)
	}
}

func TestBacklog_criterios_errores(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/criterios", url.Values{"texto": {"   "}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("criterio vacío: código = %d, se esperaba 400", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "el criterio de aceptación no puede estar vacío", "<h1>HU-1 · Crear proyecto</h1>")

	rec = servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/criterios/5/quitar", url.Values{}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("criterio inexistente: código = %d, se esperaba 400", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "el criterio no existe")
}

func TestBacklog_estimar(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/estimacion", url.Values{"story_points": {"5"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if h := e.historia(t); h.StoryPoints != 5 {
		t.Errorf("Story Points = %d, se esperaba 5", h.StoryPoints)
	}

	rec = servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/estimacion", url.Values{"story_points": {"4"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("fuera de escala: código = %d, se esperaba 400", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21")
}

func TestBacklog_cambiarEstado(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	// Sin sprint, la historia no puede pasar a En progreso (RN-08.2): se ve el mensaje de la spec.
	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/estado", url.Values{"estado": {"En progreso"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("sin sprint: código = %d, se esperaba 400", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "la historia debe estar asignada a un sprint")

	h := e.historia(t)
	if err := e.svc.Backlog.AgregarCriterio(h.ID, "Funciona"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Backlog.Estimar(h.ID, 3); err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}

	rec = servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/estado", url.Values{"estado": {"En progreso"}}))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("con sprint: código = %d, se esperaba 303", rec.Code)
	}
	if h := e.historia(t); h.Estado != backlog.EnProgreso {
		t.Errorf("estado = %v, se esperaba En progreso", h.Estado)
	}
	cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/HU-1", nil)).Body.String()
	contieneTodo(t, cuerpo, `<a href="/sprints/S1">S1</a>`)
	if strings.Contains(cuerpo, "Sin estimar") {
		t.Error("una historia estimada no debe decir \"Sin estimar\"")
	}
}
