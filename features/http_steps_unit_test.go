package features

import (
	"context"
	"net/http/httptest"
	"testing"
)

// Todavía ningún .feature activo usa "la página no contiene" (lo usan los escenarios de HU-32), así que se prueba acá.

// conPagina devuelve un contexto como el que deja el paso "hago un pedido", con esa página como respuesta.
func conPagina(cuerpo string) context.Context {
	rec := httptest.NewRecorder()
	rec.Body.WriteString(cuerpo)
	return context.WithValue(context.Background(), claveRespuesta{}, rec)
}

func TestPaginaNoContiene(t *testing.T) {
	ctx := conPagina(`<span class="insignia tono-info">Activo</span>`)

	if err := paginaNoContiene(ctx, `style="`); err != nil {
		t.Errorf("la página no tiene estilos en línea, pero el paso falló: %v", err)
	}
	if err := paginaNoContiene(ctx, "insignia"); err == nil {
		t.Error("la página contiene \"insignia\": el paso debería fallar")
	}
}

func TestPaginaNoContiene_sinPedido(t *testing.T) {
	if err := paginaNoContiene(context.Background(), "insignia"); err == nil {
		t.Error("sin un pedido previo el paso debería fallar, no dar por buena una página que no existe")
	}
}
