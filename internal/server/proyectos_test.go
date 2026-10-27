package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServidorProyectos arma el router con la sección Proyectos sobre servicios en memoria, sin datos. Devuelve
// también los servicios para preparar casos y verificar qué quedó guardado.
func nuevoServidorProyectos(t *testing.T) (http.Handler, Servicios) {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoIntegrantes, err := memory.NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}
	svc := Servicios{
		Proyectos:   service.NuevoProyectos(repoProyectos),
		Integrantes: service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return router, svc
}

// crearDemo da de alta "Demo" (P1) del 21/10 al 02/11 directamente con el servicio.
func crearDemo(t *testing.T, svc Servicios) {
	t.Helper()
	inicio := time.Date(2026, time.October, 21, 0, 0, 0, 0, time.Local)
	fin := time.Date(2026, time.November, 2, 0, 0, 0, 0, time.Local)
	if _, err := svc.Proyectos.Crear("Demo", "Proyecto de prueba", inicio, fin); err != nil {
		t.Fatal(err)
	}
}

// servir manda el pedido al router y devuelve la respuesta grabada.
func servir(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// contieneTodo hace fallar el test por cada texto que falte en el cuerpo de la respuesta.
func contieneTodo(t *testing.T, cuerpo string, textos ...string) {
	t.Helper()
	for _, texto := range textos {
		if !strings.Contains(cuerpo, texto) {
			t.Errorf("la página no contiene %q", texto)
		}
	}
}

func TestProyectos_lista(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemo(t, svc)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		`<a href="/proyectos/P1">Demo</a>`, "21/10/2026", "02/11/2026", "Planificado", `href="/proyectos/nuevo"`)
	if strings.Contains(rec.Body.String(), "Esta sección se construye") {
		t.Error("la sección Proyectos sigue marcada como en construcción")
	}
}

func TestProyectos_listaVacia(t *testing.T) {
	router, _ := nuevoServidorProyectos(t)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "Todavía no hay proyectos", `href="/proyectos/nuevo"`)
}

func TestProyectos_formularioNuevo(t *testing.T) {
	router, _ := nuevoServidorProyectos(t)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos/nuevo", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		`action="/proyectos"`, `name="nombre"`, `name="descripcion"`, `type="date" name="inicio"`, `type="date" name="fin"`)
}

func TestProyectos_crear(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	form := url.Values{"nombre": {"  Software Metrics  "}, "descripcion": {"TPI"}, "inicio": {"2026-10-21"}, "fin": {"2026-11-02"}}

	rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos", form))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/proyectos/P1" {
		t.Errorf("redirige a %q, se esperaba /proyectos/P1", lugar)
	}
	p, err := svc.Proyectos.Buscar("P1")
	if err != nil || p.Nombre != "Software Metrics" || p.Descripcion != "TPI" {
		t.Errorf("quedó guardado %+v (error %v), se esperaba Software Metrics · TPI", p, err)
	}
}

func TestProyectos_crear_errores(t *testing.T) {
	casos := []struct {
		titulo, nombre, inicio, fin, quiereMensaje string
	}{
		{"nombre de 2 caracteres", "AB", "2026-10-21", "2026-11-02", "el nombre debe tener entre 3 y 100 caracteres"},
		{"fin anterior al inicio", "Nuevo", "2026-11-02", "2026-10-21", "la fecha de fin no puede ser anterior a la de inicio"},
		{"fecha vacía", "Nuevo", "", "2026-11-02", "la fecha de fin no puede ser anterior a la de inicio"},
		{"nombre repetido", "demo", "2026-10-21", "2026-11-02", "ya existe un proyecto con ese nombre"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			router, svc := nuevoServidorProyectos(t)
			crearDemo(t, svc)
			form := url.Values{"nombre": {caso.nombre}, "descripcion": {"Mi descripción"}, "inicio": {caso.inicio}, "fin": {caso.fin}}

			rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos", form))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// El formulario vuelve con el mensaje de la spec y con lo que se escribió.
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, `action="/proyectos"`,
				`value="`+caso.nombre+`"`, "Mi descripción", `value="`+caso.inicio+`"`, `value="`+caso.fin+`"`)
			if lista, _ := svc.Proyectos.Listar(); len(lista) != 1 {
				t.Errorf("hay %d proyectos, se esperaba 1: no se debe guardar uno inválido", len(lista))
			}
		})
	}
}
