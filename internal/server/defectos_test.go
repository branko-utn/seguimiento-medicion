package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// escenarioDefectos arma el proyecto "Demo" con la historia HU-1 y tres sprints: S1 (21 al 23/10), S2 (24 al
// 26/10) y S3 (27 al 29/10). Si conDefecto, registra D1, detectado en S2.
func escenarioDefectos(t *testing.T, conDefecto bool) (http.Handler, Servicios) {
	t.Helper()
	proyectos, _ := memory.NuevoProyectos("")
	historias, _ := memory.NuevoHistorias("")
	sprints, _ := memory.NuevoSprints("")
	defectos, _ := memory.NuevoDefectos("")
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(proyectos),
		Backlog:   service.NuevoBacklog(historias, proyectos, sprintsSvc, nil, nil),
		Sprints:   sprintsSvc,
		Defectos:  service.NuevoDefectos(defectos, historias, sprints),
	}
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	p, err := svc.Proyectos.Crear("Demo", "", dia(21), dia(31))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Backlog.Crear(p.ID, "Crear proyecto", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	for _, inicio := range []int{21, 24, 27} {
		if _, err := svc.Sprints.Crear(p.ID, "Sprint de prueba", dia(inicio), dia(inicio+2)); err != nil {
			t.Fatal(err)
		}
	}
	if conDefecto {
		if _, err := svc.Defectos.Registrar("Acepta fechas invertidas", defect.Alta, "HU-1", "S2"); err != nil {
			t.Fatal(err)
		}
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}
	return router, svc
}

func ver(router http.Handler, ruta string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

func enviarFormulario(router http.Handler, ruta string, campos url.Values, htmx bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(campos.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func contiene(t *testing.T, cuerpo string, textos ...string) {
	t.Helper()
	for _, s := range textos {
		if !strings.Contains(cuerpo, s) {
			t.Errorf("la página no contiene %q", s)
		}
	}
}

func TestDefectos_seccionImplementada(t *testing.T) {
	for _, s := range secciones {
		if s.Ruta == "/defectos" && s.EnConstruccion {
			t.Fatal("la sección /defectos sigue marcada en construcción")
		}
	}
}

func TestDefectos_listaVacia(t *testing.T) {
	router, _ := escenarioDefectos(t, false)

	rec := ver(router, "/defectos")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contiene(t, rec.Body.String(), "<h1>Defectos</h1>", "Todavía no hay defectos registrados.")
}

func TestDefectos_listaConAcciones(t *testing.T) {
	router, _ := escenarioDefectos(t, true)

	rec := ver(router, "/defectos")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	contiene(t, cuerpo,
		"D1", "Acepta fechas invertidas", "Alta", "HU-1", "S2", "Abierto",
		`action="/defectos/D1/estado"`, `value="En progreso"`,
		`action="/defectos/D1/resolver"`, `<option value="S2">`, `<option value="S3">`,
	)
	// RN-23.4: no se ofrece un sprint anterior al de detección (en el formulario de registro sí aparece S1).
	if strings.Contains(formularioDe(t, cuerpo, "/defectos/D1/resolver"), `<option value="S1">`) {
		t.Error("se ofrece resolver en S1, que es anterior al sprint de detección")
	}
}

func TestDefectos_resueltoSinAcciones(t *testing.T) {
	router, svc := escenarioDefectos(t, true)
	if err := svc.Defectos.Resolver("D1", "S3"); err != nil {
		t.Fatal(err)
	}

	cuerpo := ver(router, "/defectos").Body.String()

	contiene(t, cuerpo, "Resuelto", "S3")
	for _, accion := range []string{`action="/defectos/D1/estado"`, `action="/defectos/D1/resolver"`} {
		if strings.Contains(cuerpo, accion) {
			t.Errorf("un defecto resuelto no debe ofrecer %s", accion)
		}
	}
}

func TestDefectos_enProgresoOfreceVolverAAbrir(t *testing.T) {
	router, svc := escenarioDefectos(t, true)
	if err := svc.Defectos.CambiarEstado("D1", defect.EnProgreso); err != nil {
		t.Fatal(err)
	}

	contiene(t, ver(router, "/defectos").Body.String(), "En progreso", `value="Abierto"`)
}

func TestDefectos_cambiarEstado(t *testing.T) {
	router, svc := escenarioDefectos(t, true)

	rec := enviarFormulario(router, "/defectos/D1/estado", url.Values{"estado": {"En progreso"}}, false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/defectos" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /defectos", rec.Code, rec.Header().Get("Location"))
	}
	if d, _ := svc.Defectos.Buscar("D1"); d.Estado != defect.EnProgreso {
		t.Errorf("D1 quedó en %s, se esperaba En progreso", d.Estado)
	}
}

func TestDefectos_cambiarEstado_conHTMX(t *testing.T) {
	router, _ := escenarioDefectos(t, true)

	rec := enviarFormulario(router, "/defectos/D1/estado", url.Values{"estado": {"En progreso"}}, true)

	if rec.Code != http.StatusNoContent || rec.Header().Get("HX-Redirect") != "/defectos" {
		t.Fatalf("respuesta = %d (HX-Redirect %q), se esperaba 204 con HX-Redirect a /defectos", rec.Code, rec.Header().Get("HX-Redirect"))
	}
}

func TestDefectos_resolver(t *testing.T) {
	router, svc := escenarioDefectos(t, true)

	rec := enviarFormulario(router, "/defectos/D1/resolver", url.Values{"sprint": {"S3"}}, false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/defectos" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /defectos", rec.Code, rec.Header().Get("Location"))
	}
	if d, _ := svc.Defectos.Buscar("D1"); d.Estado != defect.Resuelto || d.SprintResolucion != "S3" {
		t.Errorf("D1 = %+v, se esperaba Resuelto en S3", d)
	}
}

func TestDefectos_accionesConError(t *testing.T) {
	casos := []struct {
		nombre  string
		ruta    string
		campos  url.Values
		codigo  int
		mensaje string
	}{
		{"defecto inexistente", "/defectos/D9/estado", url.Values{"estado": {"En progreso"}}, http.StatusNotFound, "el defecto no existe"},
		{"al mismo estado", "/defectos/D1/estado", url.Values{"estado": {"Abierto"}}, http.StatusBadRequest, "no se puede pasar de Abierto a Abierto"},
		{"a Resuelto sin sprint", "/defectos/D1/estado", url.Values{"estado": {"Resuelto"}}, http.StatusBadRequest, "para resolver un defecto hay que indicar el sprint de resolución"},
		{"resolver sin sprint", "/defectos/D1/resolver", url.Values{}, http.StatusBadRequest, "para resolver un defecto hay que indicar el sprint de resolución"},
		{"sprint inexistente", "/defectos/D1/resolver", url.Values{"sprint": {"S99"}}, http.StatusNotFound, "el sprint no existe"},
		{"sprint anterior al de detección", "/defectos/D1/resolver", url.Values{"sprint": {"S1"}}, http.StatusBadRequest, "el sprint de resolución no puede ser anterior al de detección"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			router, _ := escenarioDefectos(t, true)
			rec := enviarFormulario(router, c.ruta, c.campos, false)
			if rec.Code != c.codigo {
				t.Fatalf("código = %d, se esperaba %d", rec.Code, c.codigo)
			}
			if !strings.Contains(rec.Body.String(), c.mensaje) {
				t.Errorf("la página de error no muestra %q", c.mensaje)
			}
		})
	}
}

// formularioDe devuelve el HTML del formulario cuya acción es accion, para revisar solo sus opciones.
func formularioDe(t *testing.T, cuerpo, accion string) string {
	t.Helper()
	inicio := strings.Index(cuerpo, `action="`+accion+`"`)
	if inicio < 0 {
		t.Fatalf("la página no tiene el formulario %s", accion)
	}
	fin := strings.Index(cuerpo[inicio:], "</form>")
	if fin < 0 {
		t.Fatalf("el formulario %s no se cierra", accion)
	}
	return cuerpo[inicio : inicio+fin]
}

func TestDefectos_formularioDeRegistro(t *testing.T) {
	router, _ := escenarioDefectos(t, false)

	cuerpo := ver(router, "/defectos").Body.String()

	form := formularioDe(t, cuerpo, "/defectos")
	contiene(t, form,
		`name="descripcion"`,
		"<option>Crítica</option>", "<option>Alta</option>", "<option>Media</option>", "<option>Baja</option>",
		`<option value="HU-1">HU-1 · Crear proyecto</option>`,
		`<option value="S1">S1 · Demo</option>`, `<option value="S2">S2 · Demo</option>`, `<option value="S3">S3 · Demo</option>`,
	)
}

func TestDefectos_registrar(t *testing.T) {
	router, svc := escenarioDefectos(t, false)
	campos := url.Values{"descripcion": {"Acepta fechas invertidas"}, "severidad": {"Crítica"}, "historia": {"HU-1"}, "sprint": {"S2"}}

	rec := enviarFormulario(router, "/defectos", campos, false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/defectos" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /defectos", rec.Code, rec.Header().Get("Location"))
	}
	d, err := svc.Defectos.Buscar("D1")
	if err != nil {
		t.Fatal(err)
	}
	if d.Descripcion != "Acepta fechas invertidas" || d.Severidad != defect.Critica || d.HistoriaID != "HU-1" || d.SprintDeteccion != "S2" || d.Estado != defect.Abierto {
		t.Errorf("defecto registrado = %+v", d)
	}
}

func TestDefectos_registrarConError(t *testing.T) {
	casos := []struct {
		nombre  string
		campos  url.Values
		mensaje string
	}{
		{"descripción vacía", url.Values{"descripcion": {"  "}, "severidad": {"Alta"}, "historia": {"HU-1"}, "sprint": {"S2"}}, "la descripción es obligatoria"},
		{"severidad inválida", url.Values{"descripcion": {"Se cae"}, "severidad": {"Urgente"}, "historia": {"HU-1"}, "sprint": {"S2"}}, "la severidad debe ser Crítica, Alta, Media o Baja"},
		{"historia inexistente", url.Values{"descripcion": {"Se cae"}, "severidad": {"Alta"}, "historia": {"HU-99"}, "sprint": {"S2"}}, "la historia no existe"},
		{"sprint inexistente", url.Values{"descripcion": {"Se cae"}, "severidad": {"Alta"}, "historia": {"HU-1"}, "sprint": {"S99"}}, "el sprint no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			router, svc := escenarioDefectos(t, false)

			rec := enviarFormulario(router, "/defectos", c.campos, false)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			form := formularioDe(t, rec.Body.String(), "/defectos")
			contiene(t, form, c.mensaje)
			// Vuelve a mostrar lo que se escribió.
			if d := c.campos.Get("descripcion"); strings.TrimSpace(d) != "" {
				contiene(t, form, `value="`+d+`"`)
			}
			if len(svc.Defectos.Listar()) != 0 {
				t.Error("el defecto rechazado no debe guardarse")
			}
		})
	}
}

func TestDefectos_registrarConErrorConservaLaSeveridad(t *testing.T) {
	router, _ := escenarioDefectos(t, false)
	campos := url.Values{"descripcion": {""}, "severidad": {"Baja"}, "historia": {"HU-1"}, "sprint": {"S3"}}

	form := formularioDe(t, enviarFormulario(router, "/defectos", campos, false).Body.String(), "/defectos")

	contiene(t, form, "<option selected>Baja</option>", `<option value="S3" selected>S3 · Demo</option>`)
}
