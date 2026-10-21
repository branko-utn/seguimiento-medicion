package features

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/server"
	"github.com/cucumber/godog"
)

// Pasos genéricos para escenarios que se verifican contra HTTP.

type claveRespuesta struct{}

func registrarPasosHTTP(sc *godog.ScenarioContext) {
	sc.Step(`^hago un pedido "([^"]*)" a "([^"]*)"$`, hacerPedido)
	sc.Step(`^la respuesta tiene código (\d+)$`, respuestaTieneCodigo)
	sc.Step(`^el cuerpo es el JSON (.+)$`, cuerpoEsJSON)
}

func hacerPedido(ctx context.Context, metodo, ruta string) (context.Context, error) {
	rec := httptest.NewRecorder()
	server.NewRouter().ServeHTTP(rec, httptest.NewRequest(metodo, ruta, nil))
	return context.WithValue(ctx, claveRespuesta{}, rec), nil
}

func respuestaTieneCodigo(ctx context.Context, codigo int) error {
	rec, ok := ctx.Value(claveRespuesta{}).(*httptest.ResponseRecorder)
	if !ok {
		return fmt.Errorf("no se hizo ningún pedido")
	}
	if rec.Code != codigo {
		return fmt.Errorf("código = %d (%s), se esperaba %d", rec.Code, http.StatusText(rec.Code), codigo)
	}
	return nil
}

func cuerpoEsJSON(ctx context.Context, esperado string) error {
	rec, ok := ctx.Value(claveRespuesta{}).(*httptest.ResponseRecorder)
	if !ok {
		return fmt.Errorf("no se hizo ningún pedido")
	}
	if got := strings.TrimSpace(rec.Body.String()); got != esperado {
		return fmt.Errorf("cuerpo = %s, se esperaba %s", got, esperado)
	}
	return nil
}
