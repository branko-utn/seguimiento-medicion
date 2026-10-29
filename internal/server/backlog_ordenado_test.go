package server

import (
	"net/http"
	"strings"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
)

// backlogConPrioridades carga en Demo HU-1 Baja, HU-2 Alta, HU-3 Media y HU-4 Alta (los Antecedentes de HU-09) y en
// Otro HU-5 Alta. HU-3 queda En progreso, asignada al sprint S1 de Demo.
func backlogConPrioridades(t *testing.T) entornoBacklog {
	t.Helper()
	e := nuevoServidorBacklog(t)
	for _, h := range []struct {
		proyecto  domain.ProyectoID
		prioridad backlog.Prioridad
	}{
		{"P1", backlog.PrioridadBaja}, {"P1", backlog.PrioridadAlta}, {"P1", backlog.PrioridadMedia},
		{"P1", backlog.PrioridadAlta}, {"P2", backlog.PrioridadAlta},
	} {
		if _, err := e.svc.Backlog.Crear(h.proyecto, "Historia", "", h.prioridad); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.svc.Backlog.AgregarCriterio("HU-3", "Funciona"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Backlog.Estimar("HU-3", 3); err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.Sprints.Crear("P1", "Primer MVP", fechaDeFormulario("2026-10-23"), fechaDeFormulario("2026-10-25")); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Sprints.Asignar("S1", "HU-3"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Backlog.CambiarEstado("HU-3", backlog.EnProgreso); err != nil {
		t.Fatal(err)
	}
	return e
}

// backlogDeProyecto devuelve la sección del backlog de un proyecto: desde su id hasta el cierre de la sección.
func backlogDeProyecto(t *testing.T, cuerpo, proyecto string) string {
	t.Helper()
	inicio := strings.Index(cuerpo, `id="backlog-`+proyecto+`"`)
	if inicio < 0 {
		t.Fatalf("la página no tiene el backlog de %s", proyecto)
	}
	fin := strings.Index(cuerpo[inicio:], "</section>")
	if fin < 0 {
		t.Fatalf("el backlog de %s no se cierra", proyecto)
	}
	return cuerpo[inicio : inicio+fin]
}

// enOrden verifica que las historias aparezcan en ese orden y que no aparezca ninguna otra de las indicadas en
// ausentes.
func enOrden(t *testing.T, seccion string, orden []string, ausentes ...string) {
	t.Helper()
	anterior := -1
	for _, id := range orden {
		i := strings.Index(seccion, `href="/backlog/`+id+`"`)
		if i < 0 {
			t.Errorf("no aparece %s", id)
			continue
		}
		if i < anterior {
			t.Errorf("%s aparece fuera de orden, se esperaba %v", id, orden)
		}
		anterior = i
	}
	for _, id := range ausentes {
		if strings.Contains(seccion, `href="/backlog/`+id+`"`) {
			t.Errorf("%s no debería aparecer", id)
		}
	}
}

func TestBacklog_agrupadoPorProyectoYOrdenado(t *testing.T) {
	e := backlogConPrioridades(t)

	rec := ver(e.router, "/backlog")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	// RN-09.1 y RN-09.3: cada proyecto con su backlog, por prioridad e identificador.
	demo := backlogDeProyecto(t, cuerpo, "P1")
	contiene(t, demo, "Demo")
	enOrden(t, demo, []string{"HU-2", "HU-4", "HU-3", "HU-1"}, "HU-5")
	otro := backlogDeProyecto(t, cuerpo, "P2")
	contiene(t, otro, "Otro")
	enOrden(t, otro, []string{"HU-5"}, "HU-1", "HU-2", "HU-3", "HU-4")
	// El desplegable para filtrar por estado, con "Todos" elegido.
	contiene(t, cuerpo, `<select name="estado"`, `<option value="" selected>Todos</option>`,
		"<option>Pendiente</option>", "<option>En progreso</option>", "<option>Completada</option>")
}

func TestBacklog_filtradoPorEstado(t *testing.T) {
	e := backlogConPrioridades(t)

	cuerpo := ver(e.router, "/backlog?estado=Pendiente").Body.String()

	// RN-09.2: solo las pendientes, en el mismo orden; el desplegable muestra el estado elegido.
	enOrden(t, backlogDeProyecto(t, cuerpo, "P1"), []string{"HU-2", "HU-4", "HU-1"}, "HU-3")
	contiene(t, cuerpo, "<option selected>Pendiente</option>")
}

func TestBacklog_filtradoSinHistorias(t *testing.T) {
	e := backlogConPrioridades(t)

	cuerpo := ver(e.router, "/backlog?estado=Completada").Body.String()

	// Ninguna historia en el estado pedido: no hay secciones y se ve el mensaje de backlog vacío.
	if strings.Contains(cuerpo, `id="backlog-`) {
		t.Error("no debería mostrar el backlog de ningún proyecto")
	}
	contiene(t, cuerpo, "No hay historias en el backlog.")
}

func TestBacklog_estadoInvalido(t *testing.T) {
	e := backlogConPrioridades(t)

	rec := ver(e.router, "/backlog?estado=Cualquiera")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	contiene(t, rec.Body.String(), "el estado debe ser Pendiente, En progreso o Completada")
}
