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
	if !strings.Contains(rec.Body.String(), "Backlog") {
		t.Errorf("no contiene 'Backlog'")
	}
	if !strings.Contains(rec.Body.String(), "No hay historias en el backlog.") {
		t.Errorf("no contiene mensaje de backlog vacío")
	}
}

func TestBacklog_conHistorias(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoSprints, _ := memory.NuevoSprints("")

	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	proyectosSvc := service.NuevoProyectos(repoProyectos)
	backlogSvc := service.NuevoBacklog(repoHistorias, repoProyectos, sprintsSvc)

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
