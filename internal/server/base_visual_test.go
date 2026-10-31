package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TEC-21: los componentes visuales comunes que usan las cuatro partes de HU-32.

func TestAppCSS_componentes(t *testing.T) {
	rec := pedir(t, http.MethodGet, "/static/app.css")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	contieneTodo(t, rec.Body.String(),
		"--sme-marca", "--sme-exito", "--sme-alerta", "--sme-peligro", // paleta
		"header.container", ".encabezado", // cabecera y encabezado de página
		".tarjetas", ".tarjeta-metrica", ".valor", ".etiqueta", // tarjetas de métricas
		".insignia", ".tono-neutro", ".tono-info", ".tono-exito", ".tono-alerta", ".tono-peligro", // insignias
		".vacio", "table.tabla", ".num", ".acciones", "form.en-linea", // estado vacío, tablas y acciones
		"@media print") // impresión (HU-29)
}

func TestInicio_seccionesComoTarjetasConIcono(t *testing.T) {
	rec := pedir(t, http.MethodGet, "/")

	cuerpo := rec.Body.String()
	contieneTodo(t, cuerpo, `class="encabezado"`, `class="tarjetas-secciones"`)
	if n := strings.Count(cuerpo, `<svg class="icono"`); n != len(secciones) {
		t.Errorf("la portada tiene %d íconos, se esperaba uno por sección (%d)", n, len(secciones))
	}
	for _, s := range secciones {
		if s.Icono == "" {
			t.Errorf("la sección %s no tiene ícono", s.Nombre)
		}
	}
}

func TestPaginasBase_sinEstilosEnLinea(t *testing.T) {
	a, err := nuevaApp(Servicios{})
	if err != nil {
		t.Fatal(err)
	}
	paginas := map[string]*httptest.ResponseRecorder{"inicio": pedir(t, http.MethodGet, "/")}
	paginas["error"] = httptest.NewRecorder()
	a.renderError(paginas["error"], httptest.NewRequest(http.MethodGet, "/", nil), http.StatusNotFound, "el proyecto no existe")
	paginas["en construcción"] = httptest.NewRecorder()
	a.enConstruccion(secciones[0])(paginas["en construcción"], httptest.NewRequest(http.MethodGet, secciones[0].Ruta, nil))

	for nombre, rec := range paginas {
		if strings.Contains(rec.Body.String(), `style="`) {
			t.Errorf("la página %s tiene estilos en línea: van en app.css (TEC-21)", nombre)
		}
	}
}

func TestError_conInsigniaDePeligro(t *testing.T) {
	a, err := nuevaApp(Servicios{})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	a.renderError(rec, httptest.NewRequest(http.MethodGet, "/", nil), http.StatusNotFound, "el proyecto no existe")

	contieneTodo(t, rec.Body.String(), `<span class="insignia tono-peligro">404</span>`, "el proyecto no existe", `href="/"`)
}

func TestEnConstruccion_comoEstadoVacio(t *testing.T) {
	a, err := nuevaApp(Servicios{})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	s := seccion{Ruta: "/demo", Nombre: "Demo", Descripcion: "Sección de prueba.", Epica: "E0", Sprint: 4}
	a.enConstruccion(s)(rec, httptest.NewRequest(http.MethodGet, s.Ruta, nil))

	contieneTodo(t, rec.Body.String(), `class="vacio"`, "Demo", "Sprint 4", `href="/"`)
}
