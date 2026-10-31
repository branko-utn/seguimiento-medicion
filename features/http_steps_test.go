package features

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/server"
	"github.com/cucumber/godog"
)

// Pasos genéricos para escenarios que se verifican contra HTTP.

type claveRespuesta struct{}

func registrarPasosHTTP(sc *godog.ScenarioContext, m *mundo) {
	sc.Step(`^hago un pedido "([^"]*)" a "([^"]*)"$`, m.hacerPedido)
	sc.Step(`^la respuesta tiene código (\d+)$`, respuestaTieneCodigo)
	sc.Step(`^el cuerpo es el JSON (.+)$`, cuerpoEsJSON)
	sc.Step(`^la página contiene "([^"]*)"$`, paginaContiene)
	sc.Step(`^la página no contiene "([^"]*)"$`, paginaNoContiene) // TEC-21: lo usan los escenarios de HU-32
	sc.Step(`^la página contiene un enlace a "([^"]*)" con el texto "([^"]*)"$`, paginaContieneEnlace)
}

// hacerPedido construye el router con los servicios del mundo, así los datos creados en los pasos "Dado" son
// visibles para el handler HTTP que renderiza la página.
func (m *mundo) hacerPedido(ctx context.Context, metodo, ruta string) (context.Context, error) {
	svc := server.Servicios{
		Proyectos:    m.proyectos,
		Integrantes:  m.integrantes,
		Backlog:      m.backlog,
		Sprints:      m.sprints,
		Estimaciones: m.estimaciones,
		Registros:    m.registros,
		Defectos:     m.defectos,
		Poker:        m.poker,
		Reportes:     m.reportes,
		Dashboard:    m.dashboard,
	}
	router, err := server.NewRouter(svc)
	if err != nil {
		return ctx, err
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(metodo, ruta, nil))
	return context.WithValue(ctx, claveRespuesta{}, rec), nil
}

func respuesta(ctx context.Context) (*httptest.ResponseRecorder, error) {
	rec, ok := ctx.Value(claveRespuesta{}).(*httptest.ResponseRecorder)
	if !ok {
		return nil, fmt.Errorf("no se hizo ningún pedido")
	}
	return rec, nil
}

func paginaContiene(ctx context.Context, texto string) error {
	rec, err := respuesta(ctx)
	if err != nil {
		return err
	}
	if !strings.Contains(rec.Body.String(), texto) {
		return fmt.Errorf("la página no contiene %q", texto)
	}
	return nil
}

// paginaNoContiene verifica que un texto no aparezca en la página (TEC-21), por ejemplo `style="` (RN-32.5). Sin un
// pedido previo falla, para no dar por buena una página que nunca se pidió.
func paginaNoContiene(ctx context.Context, texto string) error {
	rec, err := respuesta(ctx)
	if err != nil {
		return err
	}
	if strings.Contains(rec.Body.String(), texto) {
		return fmt.Errorf("la página contiene %q y no debería", texto)
	}
	return nil
}

func paginaContieneEnlace(ctx context.Context, ruta, texto string) error {
	rec, err := respuesta(ctx)
	if err != nil {
		return err
	}
	enlace := regexp.MustCompile(`<a[^>]*href="` + regexp.QuoteMeta(ruta) + `"[^>]*>\s*` + regexp.QuoteMeta(texto) + `\s*</a>`)
	if !enlace.MatchString(rec.Body.String()) {
		return fmt.Errorf("la página no tiene un enlace a %s con el texto %q", ruta, texto)
	}
	return nil
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
