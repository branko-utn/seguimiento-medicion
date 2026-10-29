package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServidorProyectos arma el router con la sección Proyectos sobre servicios en memoria, sin datos. Devuelve
// también los servicios para preparar casos y verificar qué quedó guardado.
func nuevoServidorProyectos(t *testing.T) (http.Handler, Servicios) {
	t.Helper()
	router, svc, _ := armarServidorProyectos(t)
	return router, svc
}

// armarServidorProyectos es nuevoServidorProyectos, pero devuelve también el repositorio de proyectos para preparar
// estados que todavía no tienen caso de uso.
func armarServidorProyectos(t *testing.T) (http.Handler, Servicios, *memory.Proyectos) {
	t.Helper()
	e := nuevoEntornoProyectos(t)
	return e.router, e.svc, e.repoProyectos
}

// entornoProyectos es el router de la sección Proyectos con sus servicios y los repositorios para preparar casos
// (por ejemplo, el estado de las historias).
type entornoProyectos struct {
	router        http.Handler
	svc           Servicios
	repoProyectos *memory.Proyectos
	repoHistorias *memory.Historias
}

// nuevoEntornoProyectos arma el entorno sin datos. Hoy es el 27/10/2026, dentro de los proyectos de prueba: el estado
// (HU-04) y la edición (HU-02) se calculan con las fechas, y el test no debe depender de la fecha real.
func nuevoEntornoProyectos(t *testing.T) entornoProyectos {
	t.Helper()
	original := hoy
	hoy = func() time.Time { return fechaDeFormulario("2026-10-27") }
	t.Cleanup(func() { hoy = original })
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoIntegrantes, err := memory.NuevoIntegrantes("")
	if err != nil {
		t.Fatal(err)
	}
	repoSprints, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	repoHistorias, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	svc := Servicios{
		Proyectos:   service.NuevoProyectos(repoProyectos).ConResumen(repoSprints, repoHistorias),
		Integrantes: service.NuevoIntegrantes(repoIntegrantes, repoProyectos),
		Sprints:     service.NuevoSprints(repoSprints, repoProyectos, repoHistorias),
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return entornoProyectos{router: router, svc: svc, repoProyectos: repoProyectos, repoHistorias: repoHistorias}
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
		`<a href="/proyectos/P1">Demo</a>`, "21/10/2026", "02/11/2026", "<td>En curso</td>", `href="/proyectos/nuevo"`)
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

// crearDemoConAna es crearDemo más "Ana Pérez" como Agile Enabler.
func crearDemoConAna(t *testing.T, svc Servicios) {
	t.Helper()
	crearDemo(t, svc)
	if _, err := svc.Integrantes.Registrar("P1", "Ana Pérez", project.AgileEnabler); err != nil {
		t.Fatal(err)
	}
}

func TestProyectos_detalle(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemoConAna(t, svc)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		"<h1>Demo</h1>", "Proyecto de prueba", "21/10/2026", "02/11/2026", "Planificado",
		"<td>Ana Pérez</td>", "<td>Agile Enabler</td>",
		`action="/proyectos/P1/integrantes"`, `name="nombre"`,
		"<option>Agile Enabler</option>", "<option>Product Builder</option>", "<option>Product Architect</option>")
}

func TestProyectos_detalleSinIntegrantes(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemo(t, svc)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "Todavía no hay integrantes", `action="/proyectos/P1/integrantes"`)
}

func TestProyectos_proyectoInexistente(t *testing.T) {
	pedidos := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/proyectos/P9", nil),
		pedidoFormulario(http.MethodPost, "/proyectos/P9/integrantes", url.Values{"nombre": {"Ana"}, "rol": {"Product Builder"}}),
	}
	for _, req := range pedidos {
		t.Run(req.Method+" "+req.URL.Path, func(t *testing.T) {
			router, _ := nuevoServidorProyectos(t)

			rec := servir(router, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba 404", rec.Code)
			}
			contieneTodo(t, rec.Body.String(), "el proyecto no existe")
		})
	}
}

func TestProyectos_registrarIntegrante(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemoConAna(t, svc)
	form := url.Values{"nombre": {"  Beto Gómez  "}, "rol": {"Product Architect"}}

	rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos/P1/integrantes", form))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/proyectos/P1" {
		t.Errorf("redirige a %q, se esperaba /proyectos/P1", lugar)
	}
	lista, _ := svc.Integrantes.DelProyecto("P1")
	if len(lista) != 2 || lista[1].Nombre != "Beto Gómez" || lista[1].Rol != project.ProductArchitect {
		t.Errorf("integrantes = %+v, se esperaba que Beto Gómez quedara como Product Architect", lista)
	}
}

func TestProyectos_registrarIntegrante_errores(t *testing.T) {
	casos := []struct {
		titulo, nombre, rol, quiereMensaje, quiereOpcion string
	}{
		{"nombre vacío", "   ", "Product Builder", "el nombre del integrante es obligatorio", "<option selected>Product Builder</option>"},
		{"rol que no existe", "Beto Gómez", "Tester", "el rol debe ser Agile Enabler, Product Builder o Product Architect", ""},
		{"integrante repetido", "ana pérez", "Product Architect", "el integrante ya está registrado en el proyecto", "<option selected>Product Architect</option>"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			router, svc := nuevoServidorProyectos(t)
			crearDemoConAna(t, svc)
			form := url.Values{"nombre": {caso.nombre}, "rol": {caso.rol}}

			rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos/P1/integrantes", form))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// Vuelve el detalle con el error, lo que se escribió y el rol elegido.
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, "<h1>Demo</h1>", "<td>Ana Pérez</td>", `value="`+caso.nombre+`"`)
			if caso.quiereOpcion != "" {
				contieneTodo(t, rec.Body.String(), caso.quiereOpcion)
			}
			if lista, _ := svc.Integrantes.DelProyecto("P1"); len(lista) != 1 {
				t.Errorf("hay %d integrantes, se esperaba 1: no se debe registrar uno inválido", len(lista))
			}
		})
	}
}

func TestProyectos_detalleEnlazaLaEdicion(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemo(t, svc)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	contieneTodo(t, rec.Body.String(), `href="/proyectos/P1/editar"`)
}

func TestProyectos_formularioEdicion(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemo(t, svc)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/proyectos/P1/editar", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	// El mismo formulario del alta, con los datos actuales y apuntando a la modificación.
	contieneTodo(t, rec.Body.String(), "Editar proyecto", `action="/proyectos/P1"`,
		`value="Demo"`, "Proyecto de prueba", `type="date" name="inicio" value="2026-10-21"`, `type="date" name="fin" value="2026-11-02"`)
}

func TestProyectos_modificar(t *testing.T) {
	router, svc := nuevoServidorProyectos(t)
	crearDemo(t, svc)
	form := url.Values{"nombre": {"Software Metrics"}, "descripcion": {"TPI"}, "inicio": {"2026-10-22"}, "fin": {"2026-11-03"}}

	rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos/P1", form))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/proyectos/P1" {
		t.Errorf("redirige a %q, se esperaba /proyectos/P1", lugar)
	}
	p, _ := svc.Proyectos.Buscar("P1")
	if p.Nombre != "Software Metrics" || p.Descripcion != "TPI" || p.Inicio.Format(formatoFecha) != "2026-10-22" || p.Fin.Format(formatoFecha) != "2026-11-03" {
		t.Errorf("quedó guardado %+v, se esperaba Software Metrics · TPI · 22/10 al 03/11", p)
	}
}

func TestProyectos_modificar_errores(t *testing.T) {
	casos := []struct {
		titulo, proyecto, nombre, inicio, fin, quiereMensaje string
	}{
		{"nombre de 2 caracteres", "P1", "AB", "2026-10-21", "2026-11-02", "el nombre debe tener entre 3 y 100 caracteres"},
		{"fin anterior al inicio", "P1", "Demo", "2026-11-02", "2026-10-21", "la fecha de fin no puede ser anterior a la de inicio"},
		{"nombre de otro proyecto", "P1", "otro", "2026-10-21", "2026-11-02", "ya existe un proyecto con ese nombre"},
		{"proyecto finalizado", "P3", "Final 2", "2026-10-21", "2026-11-02", "un proyecto finalizado no se puede modificar"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			router, svc, repo := armarServidorProyectos(t)
			crearDemo(t, svc)
			otro, err := svc.Proyectos.Crear("Otro", "", fechaDeFormulario("2026-10-21"), fechaDeFormulario("2026-11-02"))
			if err != nil {
				t.Fatal(err)
			}
			// P3 "Final" terminó en septiembre: al 27/10 está finalizado por sus fechas, aunque el campo diga Planificado.
			final := otro
			final.ID, final.Nombre = "", "Final"
			final.Inicio, final.Fin = fechaDeFormulario("2026-09-01"), fechaDeFormulario("2026-09-30")
			if _, err := repo.Agregar(final); err != nil {
				t.Fatal(err)
			}
			antes, _ := svc.Proyectos.Listar()
			form := url.Values{"nombre": {caso.nombre}, "descripcion": {"Cambiada"}, "inicio": {caso.inicio}, "fin": {caso.fin}}

			rec := servir(router, pedidoFormulario(http.MethodPost, "/proyectos/"+caso.proyecto, form))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, "Editar proyecto", `action="/proyectos/`+caso.proyecto+`"`,
				`value="`+caso.nombre+`"`, "Cambiada", `value="`+caso.inicio+`"`, `value="`+caso.fin+`"`)
			if despues, _ := svc.Proyectos.Listar(); !slices.Equal(despues, antes) {
				t.Errorf("los proyectos quedaron %+v, antes %+v: un error no debe cambiar nada", despues, antes)
			}
		})
	}
}

func TestProyectos_edicionDeProyectoInexistente(t *testing.T) {
	pedidos := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/proyectos/P9/editar", nil),
		pedidoFormulario(http.MethodPost, "/proyectos/P9", url.Values{"nombre": {"Nuevo"}, "inicio": {"2026-10-21"}, "fin": {"2026-11-02"}}),
	}
	for _, req := range pedidos {
		t.Run(req.Method+" "+req.URL.Path, func(t *testing.T) {
			router, _ := nuevoServidorProyectos(t)

			rec := servir(router, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba 404", rec.Code)
			}
			contieneTodo(t, rec.Body.String(), "el proyecto no existe")
		})
	}
}

func TestProyectos_listaMuestraElEstadoCalculado(t *testing.T) {
	e := nuevoEntornoProyectos(t)
	// Al 27/10: Antiguo ya terminó, Demo está en curso y Futuro todavía no empezó. Los tres tienen Planificado guardado.
	for _, p := range []struct{ nombre, inicio, fin string }{
		{"Antiguo", "2026-09-01", "2026-09-30"}, {"Demo", "2026-10-21", "2026-11-02"}, {"Futuro", "2026-11-10", "2026-11-20"},
	} {
		if _, err := e.svc.Proyectos.Crear(p.nombre, "", fechaDeFormulario(p.inicio), fechaDeFormulario(p.fin)); err != nil {
			t.Fatal(err)
		}
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/proyectos", nil))

	cuerpo := strings.Join(strings.Fields(rec.Body.String()), " ")
	contieneTodo(t, cuerpo,
		`Antiguo</a></td> <td>01/09/2026</td> <td>30/09/2026</td> <td>Finalizado</td>`,
		`Demo</a></td> <td>21/10/2026</td> <td>02/11/2026</td> <td>En curso</td>`,
		`Futuro</a></td> <td>10/11/2026</td> <td>20/11/2026</td> <td>Planificado</td>`)
}

func TestProyectos_detalleConResumen(t *testing.T) {
	e := nuevoEntornoProyectos(t)
	crearDemo(t, e.svc)
	// S1 "Primer MVP" activo y 1 de 4 historias completadas.
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	for _, estado := range []backlog.Estado{backlog.Completada, backlog.EnProgreso, backlog.Pendiente, backlog.Pendiente} {
		if _, err := e.repoHistorias.Agregar(backlog.Historia{ProyectoID: "P1", Titulo: "Historia", Prioridad: backlog.PrioridadMedia, Estado: estado}); err != nil {
			t.Fatal(err)
		}
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	// El encabezado usa el estado calculado, y el resumen va arriba de los integrantes.
	contieneTodo(t, cuerpo, "En curso · del 21/10/2026 al 02/11/2026", "6 días restantes",
		`Sprint activo: <a href="/sprints/S1">Primer MVP</a>`, "25.0 % de historias completadas")
	if strings.Contains(cuerpo, "Planificado") {
		t.Error("el detalle no debe mostrar el estado guardado (Planificado): se calcula con las fechas")
	}
	if strings.Index(cuerpo, "6 días restantes") > strings.Index(cuerpo, "<h2>Integrantes</h2>") {
		t.Error("el resumen debe ir arriba de los integrantes")
	}
}

func TestProyectos_detalleSinSprintActivo(t *testing.T) {
	e := nuevoEntornoProyectos(t)
	crearDemo(t, e.svc)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	contieneTodo(t, rec.Body.String(), "Sin sprint activo", "0.0 % de historias completadas")
}

func TestProyectos_detalleDeProyectoFinalizado(t *testing.T) {
	e := nuevoEntornoProyectos(t)
	if _, err := e.svc.Proyectos.Crear("Antiguo", "", fechaDeFormulario("2026-09-01"), fechaDeFormulario("2026-09-30")); err != nil {
		t.Fatal(err)
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/proyectos/P1", nil))

	contieneTodo(t, rec.Body.String(), "Finalizado · del 01/09/2026 al 30/09/2026", "0 días restantes")
}
