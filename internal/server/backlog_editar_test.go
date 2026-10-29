package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
)

func TestBacklog_detalleTieneEditarYEliminar(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/HU-1", nil)).Body.String()

	// El formulario "Editar" viene cargado con los datos actuales; "Eliminar" es un botón aparte (HU-06).
	contieneTodo(t, cuerpo,
		`action="/backlog/HU-1/editar"`, `name="titulo" value="Crear proyecto"`, ">Alta de proyectos</textarea>",
		"<option selected>Alta</option>", "<option>Media</option>", "<option>Baja</option>",
		`action="/backlog/HU-1/eliminar"`, ">Eliminar</button>")
}

func TestBacklog_editar(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/editar",
		url.Values{"titulo": {"  Crear y listar proyectos  "}, "descripcion": {"Alta y listado"}, "prioridad": {"Baja"}}))

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/backlog/HU-1" {
		t.Fatalf("código = %d, Location = %q; se esperaba 303 a /backlog/HU-1", rec.Code, rec.Header().Get("Location"))
	}
	h := e.historia(t)
	if h.Titulo != "Crear y listar proyectos" || h.Descripcion != "Alta y listado" || h.Prioridad != backlog.PrioridadBaja {
		t.Errorf("historia = %+v", h)
	}
}

func TestBacklog_editar_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		completar bool
		valores   url.Values
		mensaje   string
	}{
		{"título vacío", false, url.Values{"titulo": {"  "}, "prioridad": {"Alta"}}, "el título es obligatorio"},
		{"prioridad inválida", false, url.Values{"titulo": {"Algo"}, "prioridad": {"Urgente"}}, "la prioridad debe ser Alta, Media o Baja"},
		{"historia completada", true, url.Values{"titulo": {"Algo"}, "prioridad": {"Alta"}}, "una historia completada no se puede editar"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoServidorBacklog(t)
			e.crearHistoria(t)
			if c.completar {
				e.completarHistoria(t)
			}

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/editar", c.valores))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// Vuelve al detalle con el mensaje de la spec, y la historia no cambia.
			contieneTodo(t, rec.Body.String(), c.mensaje, "<h1>HU-1 · Crear proyecto</h1>")
			if h := e.historia(t); h.Titulo != "Crear proyecto" {
				t.Errorf("título = %q, no debería cambiar", h.Titulo)
			}
		})
	}
}

func TestBacklog_eliminar(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/eliminar", url.Values{}))

	// Al eliminar vuelve a la lista del backlog (spec HU-06, restricciones).
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/backlog" {
		t.Fatalf("código = %d, Location = %q; se esperaba 303 a /backlog", rec.Code, rec.Header().Get("Location"))
	}
	if _, err := e.svc.Backlog.Buscar("HU-1"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Buscar HU-1: error = %v, se esperaba ErrNoExiste", err)
	}
}

func TestBacklog_eliminar_errores(t *testing.T) {
	casos := []struct {
		nombre   string
		preparar func(e entornoBacklog, t *testing.T)
		mensaje  string
	}{
		{"asignada a un sprint", entornoBacklog.asignarASprint, "no se puede eliminar una historia asignada a un sprint"},
		{"completada", entornoBacklog.completarHistoria, "no se puede eliminar una historia completada"},
		{"con defectos", entornoBacklog.registrarDefecto, "no se puede eliminar una historia con esfuerzo o defectos registrados"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoServidorBacklog(t)
			e.crearHistoria(t)
			c.preparar(e, t)

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/backlog/HU-1/eliminar", url.Values{}))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			contieneTodo(t, rec.Body.String(), c.mensaje, "<h1>HU-1 · Crear proyecto</h1>")
			e.historia(t) // sigue existiendo
		})
	}
}

func TestBacklog_editarYEliminarInexistente(t *testing.T) {
	e := nuevoServidorBacklog(t)
	for _, ruta := range []string{"/backlog/HU-9/editar", "/backlog/HU-9/eliminar"} {
		rec := servir(e.router, pedidoFormulario(http.MethodPost, ruta, url.Values{"titulo": {"x"}, "prioridad": {"Alta"}}))
		if rec.Code != http.StatusNotFound {
			t.Errorf("POST %s: código = %d, se esperaba 404", ruta, rec.Code)
		}
		contieneTodo(t, rec.Body.String(), "la historia no existe")
	}
}

// asignarASprint crea el sprint S1 en Demo y le asigna HU-1.
func (e entornoBacklog) asignarASprint(t *testing.T) {
	t.Helper()
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Backlog.Estimar("HU-1", 3); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Backlog.AgregarCriterio("HU-1", "Funciona"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
}

// completarHistoria asigna HU-1 a S1, inicia el sprint y la completa (HU-12).
func (e entornoBacklog) completarHistoria(t *testing.T) {
	t.Helper()
	e.asignarASprint(t)
	if err := e.svc.Sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.CompletarHistoria("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
}

// registrarDefecto crea el sprint S1 y registra un defecto en HU-1 detectado en ese sprint (HU-22).
func (e entornoBacklog) registrarDefecto(t *testing.T) {
	t.Helper()
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.Defectos.Registrar("Falla al guardar", defect.Critica, "HU-1", "S1"); err != nil {
		t.Fatal(err)
	}
}
