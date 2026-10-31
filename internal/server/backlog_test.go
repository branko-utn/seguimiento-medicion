package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

func TestBacklog_vacio(t *testing.T) {
	rec := pedir(t, http.MethodGet, "/backlog")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, "Backlog") {
		t.Errorf("no contiene 'Backlog'")
	}
	// RN-32.4: sin historias se ve el estado vacío, con un enlace a la acción que llena el backlog.
	contieneTodo(t, rec.Body.String(), "No hay historias en el backlog.", `class="vacio"`, `href="/backlog/nueva"`)
}

func TestBacklog_conHistorias(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoSprints, _ := memory.NuevoSprints("")

	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	proyectosSvc := service.NuevoProyectos(repoProyectos)
	backlogSvc := service.NuevoBacklog(repoHistorias, repoProyectos, sprintsSvc, nil, nil)

	p, err := proyectosSvc.Crear("Proyecto Alpha", "Desc", time.Now(), time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	_, err = backlogSvc.Crear(p.ID, "HU-01", "Primera historia", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}

	svc := Servicios{
		Proyectos: proyectosSvc,
		Backlog:   backlogSvc,
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/backlog", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, "HU-01") {
		t.Errorf("no contiene el título 'HU-01'")
	}
	if !strings.Contains(cuerpo, "Proyecto Alpha") {
		t.Errorf("no contiene el nombre del proyecto 'Proyecto Alpha'")
	}
}

// RN-32.1: la prioridad y el estado de cada historia se ven con insignias de tono, tanto en la lista
// como en el detalle. Los tonos los decide el handler en Go (ADR 0003), no la plantilla.
func TestTonoPrioridad(t *testing.T) {
	// De más urgente (alerta) a menos (neutro), igual que tonoSeveridad; el cero no es una prioridad.
	casos := []struct {
		prioridad backlog.Prioridad
		quiere    string
	}{
		{backlog.PrioridadAlta, tonoAlerta},
		{backlog.PrioridadMedia, tonoInfo},
		{backlog.PrioridadBaja, tonoNeutro},
		{backlog.Prioridad(0), tonoNeutro},
	}
	for _, caso := range casos {
		if got := tonoPrioridad(caso.prioridad); got != caso.quiere {
			t.Errorf("tonoPrioridad(%v) = %q, se esperaba %q", caso.prioridad, got, caso.quiere)
		}
	}
}

func TestBacklog_insignias(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t) // HU-1 "Crear proyecto", con prioridad Alta y estado Pendiente

	cuerpo := ver(e.router, "/backlog").Body.String()

	contieneTodo(t, cuerpo,
		`class="insignia tono-alerta">Alta</span>`,
		`class="insignia tono-neutro">Pendiente</span>`)
}

func TestDetalleHistoria_insignias(t *testing.T) {
	e := nuevoServidorBacklog(t)
	e.crearHistoria(t)

	cuerpo := servir(e.router, httptest.NewRequest(http.MethodGet, "/backlog/HU-1", nil)).Body.String()

	contieneTodo(t, cuerpo,
		`class="insignia tono-alerta">Alta</span>`,
		`class="insignia tono-neutro">Pendiente</span>`)
}
