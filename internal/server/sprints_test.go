package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// entornoSprints es el router de la sección Sprints sobre servicios en memoria, con los repositorios para preparar
// estados que no tienen caso de uso en la pantalla (Story Points, criterios, sprint cerrado).
type entornoSprints struct {
	router        http.Handler
	svc           Servicios
	repoSprints   *memory.Sprints
	repoHistorias *memory.Historias
}

// nuevoServidorSprints arma el entorno con "Demo" (P1) y "Otro" (P2), los dos del 21/10 al 02/11, sin sprints.
func nuevoServidorSprints(t *testing.T) entornoSprints {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
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
	sprints := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(repoProyectos),
		Sprints:   sprints,
		Backlog:   service.NuevoBacklog(repoHistorias, repoProyectos, sprints),
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
	return entornoSprints{router: router, svc: svc, repoSprints: repoSprints, repoHistorias: repoHistorias}
}

// crearPrimerMVP crea el sprint "Primer MVP" (S1) en Demo, del 23/10 al 25/10.
func (e entornoSprints) crearPrimerMVP(t *testing.T) {
	t.Helper()
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
}

// cargarHistorias crea el backlog del escenario de HU-11. En Demo: HU-1 "Crear proyecto" (3 SP, con criterio),
// HU-2 "Crear sprint" (sin estimar) y HU-3 "Cerrar sprint" (5 SP, sin criterios). En Otro: HU-4 "Ver métricas"
// (2 SP, con criterio). Story Points y criterios se cargan directo en el repositorio, como los pasos "Dado".
func (e entornoSprints) cargarHistorias(t *testing.T) {
	t.Helper()
	historias := []struct {
		proyecto domain.ProyectoID
		titulo   string
		sp       int
		criterio bool
	}{
		{"P1", "Crear proyecto", 3, true},
		{"P1", "Crear sprint", 0, true},
		{"P1", "Cerrar sprint", 5, false},
		{"P2", "Ver métricas", 2, true},
	}
	for _, d := range historias {
		h, err := e.svc.Backlog.Crear(d.proyecto, d.titulo, "", backlog.PrioridadMedia)
		if err != nil {
			t.Fatal(err)
		}
		h.StoryPoints = d.sp
		if d.criterio {
			h.Criterios = []backlog.Criterio{{Texto: "Funciona"}}
		}
		if err := e.repoHistorias.Guardar(h); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSprints_lista(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	// Cada proyecto con sus sprints y su botón de alta; un proyecto sin sprints lo avisa.
	contieneTodo(t, cuerpo,
		"<h2>Demo</h2>", `<a href="/sprints/S1">Primer MVP</a>`, "23/10/2026", "25/10/2026", "Planificado",
		`href="/sprints/nuevo?proyecto=P1"`, "<h2>Otro</h2>", "Todavía no hay sprints", `href="/sprints/nuevo?proyecto=P2"`)
	if strings.Contains(cuerpo, "Esta sección se construye") {
		t.Error("la sección Sprints sigue marcada como en construcción")
	}
}

func TestSprints_listaSinProyectos(t *testing.T) {
	router, _ := nuevoServidorProyectos(t)

	rec := servir(router, httptest.NewRequest(http.MethodGet, "/sprints", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "Todavía no hay proyectos", `href="/proyectos/nuevo"`)
}

func TestSprints_formularioNuevo(t *testing.T) {
	e := nuevoServidorSprints(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/nuevo?proyecto=P2", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	// El proyecto del enlace viene elegido en el desplegable.
	contieneTodo(t, rec.Body.String(), `action="/sprints"`, `name="proyecto"`,
		`<option value="P1">Demo</option>`, `<option value="P2" selected>Otro</option>`,
		`name="goal"`, `type="date" name="inicio"`, `type="date" name="fin"`)
}

func TestSprints_crear(t *testing.T) {
	e := nuevoServidorSprints(t)
	form := url.Values{"proyecto": {"P1"}, "goal": {"  Primer MVP  "}, "inicio": {"2026-10-23"}, "fin": {"2026-10-25"}}

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints", form))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/sprints/S1" {
		t.Errorf("redirige a %q, se esperaba /sprints/S1", lugar)
	}
	sp, err := e.svc.Sprints.Buscar("S1")
	if err != nil || sp.ProyectoID != "P1" || sp.Goal != "Primer MVP" {
		t.Errorf("quedó guardado %+v (error %v), se esperaba Primer MVP en P1", sp, err)
	}
}

func TestSprints_crear_errores(t *testing.T) {
	casos := []struct {
		titulo, proyecto, goal, inicio, fin, quiereMensaje string
	}{
		{"Sprint Goal vacío", "P1", "   ", "2026-10-26", "2026-10-27", "el Sprint Goal es obligatorio"},
		{"fin anterior al inicio", "P1", "Interfaz", "2026-10-27", "2026-10-26", "la fecha de fin no puede ser anterior a la de inicio"},
		{"fuera del proyecto", "P1", "Cierre", "2026-11-01", "2026-11-05", "las fechas del sprint deben estar dentro de las del proyecto"},
		{"se superpone con otro sprint", "P1", "Interfaz", "2026-10-25", "2026-10-27", "las fechas se superponen con otro sprint del proyecto"},
		{"proyecto inexistente", "P9", "Interfaz", "2026-10-26", "2026-10-27", "el proyecto no existe"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			e := nuevoServidorSprints(t)
			e.crearPrimerMVP(t)
			form := url.Values{"proyecto": {caso.proyecto}, "goal": {caso.goal}, "inicio": {caso.inicio}, "fin": {caso.fin}}

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints", form))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// El formulario vuelve con el mensaje de la spec y con lo que se escribió.
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, `action="/sprints"`,
				`value="`+caso.goal+`"`, `value="`+caso.inicio+`"`, `value="`+caso.fin+`"`)
			if caso.proyecto == "P1" {
				contieneTodo(t, rec.Body.String(), `<option value="P1" selected>Demo</option>`)
			}
			if lista, _ := e.svc.Sprints.DelProyecto("P1"); len(lista) != 1 {
				t.Errorf("hay %d sprints en Demo, se esperaba 1: no se debe guardar uno inválido", len(lista))
			}
		})
	}
}

func TestSprints_detalle(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		"<h1>Primer MVP</h1>", `<a href="/proyectos/P1">Demo</a>`, "Planificado", "del 23/10/2026 al 25/10/2026",
		`action="/sprints/S1/iniciar"`, "Todavía no hay historias en este sprint")
}

func TestSprints_detalleDeSprintActivo(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	if err := e.svc.Sprints.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil))

	contieneTodo(t, rec.Body.String(), "Activo")
	if strings.Contains(rec.Body.String(), `action="/sprints/S1/iniciar"`) {
		t.Error("un sprint activo no debe ofrecer el botón para iniciarlo")
	}
}

func TestSprints_iniciar(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/iniciar", url.Values{}))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/sprints/S1" {
		t.Errorf("redirige a %q, se esperaba /sprints/S1", lugar)
	}
	if sp, _ := e.svc.Sprints.Buscar("S1"); sp.Estado != sprint.Activo {
		t.Errorf("estado = %v, se esperaba Activo", sp.Estado)
	}
}

func TestSprints_iniciar_errores(t *testing.T) {
	// En Demo: S1 (23–25) activo y S2 (26–27) planificado.
	casos := []struct {
		titulo, sprint, quiereMensaje, quiereTitulo string
	}{
		{"ya hay otro sprint activo", "S2", "ya hay un sprint activo en el proyecto", "<h1>Interfaz</h1>"},
		{"el sprint ya está activo", "S1", "solo se puede iniciar un sprint planificado", "<h1>Primer MVP</h1>"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			e := nuevoServidorSprints(t)
			e.crearPrimerMVP(t)
			if _, err := e.svc.Sprints.Crear("P1", "Interfaz", fechaDeFormulario("2026-10-26"), fechaDeFormulario("2026-10-27")); err != nil {
				t.Fatal(err)
			}
			if err := e.svc.Sprints.Iniciar("S1"); err != nil {
				t.Fatal(err)
			}

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/"+caso.sprint+"/iniciar", url.Values{}))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// Vuelve el detalle del sprint con el mensaje de la spec.
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, caso.quiereTitulo)
			if s2, _ := e.svc.Sprints.Buscar("S2"); s2.Estado != sprint.Planificado {
				t.Errorf("S2 quedó %v, se esperaba Planificado: un error no debe cambiar el sprint", s2.Estado)
			}
		})
	}
}

func TestSprints_sprintInexistente(t *testing.T) {
	pedidos := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/sprints/S9", nil),
		pedidoFormulario(http.MethodPost, "/sprints/S9/iniciar", url.Values{}),
	}
	for _, req := range pedidos {
		t.Run(req.Method+" "+req.URL.Path, func(t *testing.T) {
			e := nuevoServidorSprints(t)

			rec := servir(e.router, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba 404", rec.Code)
			}
			contieneTodo(t, rec.Body.String(), "el sprint no existe")
		})
	}
}

func TestSprints_detalleConHistorias(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)
	if err := e.svc.Sprints.Asignar("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil))

	cuerpo := rec.Body.String()
	// Las historias del sprint con sus SP y el total planificado (HU-24).
	contieneTodo(t, cuerpo, "<td>HU-1</td>", "<td>Crear proyecto</td>", "<td>3</td>", "3 Story Points planificados")
	// El desplegable ofrece las historias de Demo sin sprint, también las que no están listas: las valida el dominio.
	contieneTodo(t, cuerpo, `action="/sprints/S1/historias"`, `name="historia"`,
		`<option value="HU-2">HU-2 · Crear sprint (0 SP)</option>`, `<option value="HU-3">HU-3 · Cerrar sprint (5 SP)</option>`)
	for _, ajena := range []string{`<option value="HU-1"`, `<option value="HU-4"`} {
		if strings.Contains(cuerpo, ajena) {
			t.Errorf("el desplegable no debe ofrecer %s: ya está en un sprint o es de otro proyecto", ajena)
		}
	}
}

func TestSprints_detalleSinHistoriasDisponibles(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.cargarHistorias(t)
	if _, err := e.svc.Sprints.Crear("P2", "Métricas", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Asignar("S1", "HU-4"); err != nil {
		t.Fatal(err)
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil))

	contieneTodo(t, rec.Body.String(), "<td>Ver métricas</td>", "2 Story Points planificados", "No hay historias del backlog sin sprint")
}

func TestSprints_asignar(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/historias", url.Values{"historia": {"HU-1"}}))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if lugar := rec.Header().Get("Location"); lugar != "/sprints/S1" {
		t.Errorf("redirige a %q, se esperaba /sprints/S1", lugar)
	}
	sp, _ := e.svc.Sprints.Buscar("S1")
	h, _ := e.svc.Backlog.Buscar("HU-1")
	if len(sp.Historias) != 1 || sp.Historias[0] != "HU-1" || h.SprintID != "S1" {
		t.Errorf("sprint %v e historia con sprint %q, se esperaba HU-1 asignada a S1", sp.Historias, h.SprintID)
	}
}

func TestSprints_asignar_errores(t *testing.T) {
	casos := []struct {
		titulo, historia, quiereMensaje string
	}{
		{"historia sin estimar", "HU-2", "la historia no está estimada"},
		{"historia sin criterios", "HU-3", "la historia no tiene criterios de aceptación"},
		{"historia de otro proyecto", "HU-4", "la historia pertenece a otro proyecto"},
		{"historia inexistente", "HU-9", "la historia no existe"},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			e := nuevoServidorSprints(t)
			e.crearPrimerMVP(t)
			e.cargarHistorias(t)

			rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/historias", url.Values{"historia": {caso.historia}}))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			// Vuelve el detalle del sprint con el mensaje de la spec.
			contieneTodo(t, rec.Body.String(), caso.quiereMensaje, "<h1>Primer MVP</h1>", `action="/sprints/S1/historias"`)
			if sp, _ := e.svc.Sprints.Buscar("S1"); len(sp.Historias) != 0 {
				t.Errorf("el sprint quedó con %v, se esperaba sin historias: un error no debe cambiarlo", sp.Historias)
			}
		})
	}
}

func TestSprints_asignar_conservaLaHistoriaElegida(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/historias", url.Values{"historia": {"HU-3"}}))

	contieneTodo(t, rec.Body.String(), `<option value="HU-3" selected>`)
}

func TestSprints_sprintCerrado(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.crearPrimerMVP(t)
	e.cargarHistorias(t)
	// Cerrar un sprint de verdad es HU-13: acá se carga el estado directo en el repositorio.
	sp, _ := e.repoSprints.Buscar("S1")
	sp.Estado = sprint.Cerrado
	if err := e.repoSprints.Actualizar(sp); err != nil {
		t.Fatal(err)
	}

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/sprints/S1", nil))
	if strings.Contains(rec.Body.String(), `action="/sprints/S1/historias"`) {
		t.Error("un sprint cerrado no debe ofrecer el formulario para asignar historias")
	}

	rec = servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S1/historias", url.Values{"historia": {"HU-1"}}))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "el sprint está cerrado y no se puede modificar")
}

func TestSprints_asignarEnSprintInexistente(t *testing.T) {
	e := nuevoServidorSprints(t)
	e.cargarHistorias(t)

	rec := servir(e.router, pedidoFormulario(http.MethodPost, "/sprints/S9/historias", url.Values{"historia": {"HU-1"}}))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("código = %d, se esperaba 404", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "el sprint no existe")
}
