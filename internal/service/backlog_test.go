package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// nuevoServicioBacklog arma el servicio con dos proyectos ya creados: "Demo" (P1) y "Otro" (P2).
func nuevoServicioBacklog(t *testing.T) *Backlog {
	t.Helper()
	repoProyectos, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoHistorias, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	proyectos := NuevoProyectos(repoProyectos)
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := proyectos.Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	return NuevoBacklog(repoHistorias, repoProyectos)
}

func TestBacklog_Crear(t *testing.T) {
	s := nuevoServicioBacklog(t)

	h, err := s.Crear("P1", "  Crear proyecto  ", "con detalle", backlog.PrioridadAlta)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// RN-05.1: el ID lo genera el repositorio. RN-05.2: el título se guarda sin espacios sobrantes.
	// RN-05.4: nueva = Pendiente, sin Story Points, sin criterios y sin sprint.
	if h.ID != "HU-1" || h.ProyectoID != "P1" || h.Titulo != "Crear proyecto" || h.Descripcion != "con detalle" ||
		h.Prioridad != backlog.PrioridadAlta || h.Estado != backlog.Pendiente || h.StoryPoints != 0 ||
		len(h.Criterios) != 0 || h.SprintID != "" {
		t.Errorf("Crear() = %+v, se esperaba HU-1 · Crear proyecto · Alta · Pendiente · sin estimar", h)
	}

	segunda, err := s.Crear("P1", "Crear sprint", "", backlog.PrioridadMedia)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if segunda.ID != "HU-2" {
		t.Errorf("RN-05.1: los IDs son correlativos, se obtuvo %q", segunda.ID)
	}
}

func TestBacklog_Crear_losIDsNoSeRepitenEntreProyectos(t *testing.T) {
	// Caso límite de la spec HU-05: con 3 historias en otro proyecto, la siguiente es HU-4.
	s := nuevoServicioBacklog(t)
	for i := 0; i < 3; i++ {
		if _, err := s.Crear("P1", "Historia", "", backlog.PrioridadBaja); err != nil {
			t.Fatal(err)
		}
	}

	otra, err := s.Crear("P2", "Login", "", backlog.PrioridadBaja)
	if err != nil {
		t.Fatal(err)
	}
	if otra.ID != "HU-4" {
		t.Errorf("los IDs no se repiten entre proyectos: se esperaba HU-4, se obtuvo %q", otra.ID)
	}
}

func TestBacklog_Buscar(t *testing.T) {
	s := nuevoServicioBacklog(t)
	h, _ := s.Crear("P1", "Login", "", backlog.PrioridadAlta)

	got, err := s.Buscar(h.ID)
	if err != nil || got.Titulo != "Login" {
		t.Errorf("Buscar(%s) = %+v, %v", h.ID, got, err)
	}
	if _, err := s.Buscar("HU-99"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Buscar(HU-99) = %v, se esperaba backlog.ErrNoExiste", err)
	}
	if backlog.ErrNoExiste.Error() != "la historia no existe" {
		t.Errorf("mensaje de ErrNoExiste = %q", backlog.ErrNoExiste.Error())
	}
}

func TestBacklog_Crear_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		proyecto  domain.ProyectoID
		titulo    string
		prioridad backlog.Prioridad
		quiereErr error
	}{
		{"proyecto inexistente", "P9", "Login", backlog.PrioridadAlta, project.ErrNoExiste},
		{"primero se verifica el proyecto", "P9", "", backlog.Prioridad(99), project.ErrNoExiste},
		{"título vacío", "P1", "   ", backlog.PrioridadAlta, backlog.ErrTituloVacio},
		{"prioridad inválida", "P1", "Login", backlog.Prioridad(99), backlog.ErrPrioridadInvalida},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s := nuevoServicioBacklog(t)

			_, err := s.Crear(caso.proyecto, caso.titulo, "", caso.prioridad)
			if !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}

			lista, _ := s.Listar()
			if len(lista) != 0 {
				t.Errorf("hay %d historias guardadas, se esperaba 0: no se debe guardar una historia inválida", len(lista))
			}
		})
	}
}

func TestBacklog_Listar(t *testing.T) {
	s := nuevoServicioBacklog(t)

	vacio, err := s.Listar()
	if err != nil {
		t.Fatalf("backlog sin historias: error inesperado: %v", err)
	}
	if len(vacio) != 0 {
		t.Errorf("backlog sin historias: se obtuvieron %d, se esperaba 0", len(vacio))
	}

	altas := []struct {
		proyecto domain.ProyectoID
		titulo   string
	}{
		{"P1", "Crear proyecto"},
		{"P2", "Login"},
		{"P1", "Crear sprint"},
	}
	for _, a := range altas {
		if _, err := s.Crear(a.proyecto, a.titulo, "", backlog.PrioridadMedia); err != nil {
			t.Fatalf("Crear %s: %v", a.titulo, err)
		}
	}

	lista, err := s.Listar()
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 3 {
		t.Fatalf("Listar() tiene %d historias, se esperaban 3", len(lista))
	}
	for n, a := range altas {
		if lista[n].Titulo != a.titulo || lista[n].ProyectoID != a.proyecto {
			t.Errorf("historia %d = %s · %s, se esperaba %s · %s (en orden de creación)", n, lista[n].ProyectoID, lista[n].Titulo, a.proyecto, a.titulo)
		}
	}
}
