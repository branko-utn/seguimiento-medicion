package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInicio_muestraLayoutYNavegacion(t *testing.T) {
	rec := pedir(t, http.MethodGet, "/")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, se esperaba text/html", ct)
	}
	cuerpo := rec.Body.String()
	for _, s := range secciones {
		if !strings.Contains(cuerpo, `href="`+s.Ruta+`"`) {
			t.Errorf("falta el enlace a %s en la navegación", s.Ruta)
		}
	}
	for _, recurso := range []string{"/static/vendor/pico.min.css", "/static/vendor/htmx.min.js", "/static/vendor/chart.umd.min.js"} {
		if !strings.Contains(cuerpo, recurso) {
			t.Errorf("el layout no incluye %s", recurso)
		}
	}
}

func TestSeccionEnConstruccion_avisaElSprint(t *testing.T) {
	for _, s := range secciones {
		if !s.EnConstruccion {
			continue
		}
		t.Run(s.Ruta, func(t *testing.T) {
			rec := pedir(t, http.MethodGet, s.Ruta)
			if rec.Code != http.StatusOK {
				t.Fatalf("código = %d, se esperaba 200", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), s.Nombre) {
				t.Errorf("la página no menciona la sección %q", s.Nombre)
			}
		})
	}
}

func TestStatic_sirveRecursosEmbebidos(t *testing.T) {
	rec := pedir(t, http.MethodGet, "/static/vendor/htmx.min.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	if rec.Body.Len() == 0 {
		t.Error("el archivo está vacío")
	}
}

func TestRenderError_muestraMensajeEscapado(t *testing.T) {
	a, err := nuevaApp(Servicios{})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	a.renderError(rec, httptest.NewRequest(http.MethodGet, "/", nil), http.StatusBadRequest, "<script>mal</script>")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("código = %d, se esperaba 400", rec.Code)
	}
	cuerpo := rec.Body.String()
	if strings.Contains(cuerpo, "<script>mal") {
		t.Error("el mensaje no se escapó")
	}
	if !strings.Contains(cuerpo, "&lt;script&gt;mal") {
		t.Error("la página no muestra el mensaje de error")
	}
}

func TestRender_paginaInexistenteDevuelve500(t *testing.T) {
	a, err := nuevaApp(Servicios{})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	a.render(rec, httptest.NewRequest(http.MethodGet, "/", nil), "no-existe.html", nil)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("código = %d, se esperaba 500", rec.Code)
	}
}
