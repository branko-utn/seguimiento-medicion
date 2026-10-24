package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
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

func TestBacklog_AgregarCriterio(t *testing.T) {
	s := nuevoServicioBacklog(t)
	h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}

	// El criterio queda al final y sin espacios sobrantes (spec HU-07, salidas esperadas).
	if err := s.AgregarCriterio(h.ID, "El nombre es obligatorio"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if err := s.AgregarCriterio(h.ID, "  La fecha de fin no es anterior a la de inicio  "); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	guardada, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	quiere := []string{"El nombre es obligatorio", "La fecha de fin no es anterior a la de inicio"}
	if len(guardada.Criterios) != len(quiere) {
		t.Fatalf("hay %d criterios, se esperaban %d", len(guardada.Criterios), len(quiere))
	}
	for n, texto := range quiere {
		if guardada.Criterios[n].Texto != texto {
			t.Errorf("criterio %d = %q, se esperaba %q", n+1, guardada.Criterios[n].Texto, texto)
		}
	}
	if !guardada.TieneCriterios() {
		t.Error("la historia debería tener criterios (RN-07.2)")
	}
}

func TestBacklog_AgregarCriterio_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		historia  domain.HistoriaID
		texto     string
		quiereErr error
	}{
		{"historia inexistente", "HU-99", "El nombre es obligatorio", backlog.ErrNoExiste},
		{"criterio vacío", "HU-1", "   ", backlog.ErrCriterioVacio},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s := nuevoServicioBacklog(t)
			if _, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta); err != nil {
				t.Fatal(err)
			}

			if err := s.AgregarCriterio(caso.historia, caso.texto); !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}

			// El criterio inválido no queda guardado.
			h, _ := s.Buscar("HU-1")
			if len(h.Criterios) != 0 {
				t.Errorf("hay %d criterios guardados, se esperaba 0", len(h.Criterios))
			}
		})
	}
}

func TestBacklog_QuitarCriterio(t *testing.T) {
	s := nuevoServicioBacklog(t)
	h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}
	for _, texto := range []string{"A", "B", "C"} {
		if err := s.AgregarCriterio(h.ID, texto); err != nil {
			t.Fatal(err)
		}
	}

	// Los siguientes se renumeran: quitar el 2 deja "A, C" (spec HU-07, salidas esperadas).
	if err := s.QuitarCriterio(h.ID, 2); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	guardada, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	quiere := []string{"A", "C"}
	if len(guardada.Criterios) != len(quiere) {
		t.Fatalf("quedan %d criterios, se esperaban %d", len(guardada.Criterios), len(quiere))
	}
	for n, texto := range quiere {
		if guardada.Criterios[n].Texto != texto {
			t.Errorf("criterio %d = %q, se esperaba %q", n+1, guardada.Criterios[n].Texto, texto)
		}
	}
}

func TestBacklog_QuitarCriterio_elUnico(t *testing.T) {
	// Caso límite: quitar el único criterio deja la historia sin criterios, y es válido.
	s := nuevoServicioBacklog(t)
	h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AgregarCriterio(h.ID, "A"); err != nil {
		t.Fatal(err)
	}
	// Preparación: la historia debe tener exactamente el criterio "A" antes de quitarlo.
	conCriterio, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(conCriterio.Criterios) != 1 {
		t.Fatalf("preparación: hay %d criterios, se esperaba 1", len(conCriterio.Criterios))
	}

	if err := s.QuitarCriterio(h.ID, 1); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	guardada, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(guardada.Criterios) != 0 || guardada.TieneCriterios() {
		t.Errorf("se esperaba la historia sin criterios, se obtuvo %+v", guardada.Criterios)
	}
}

func TestBacklog_QuitarCriterio_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		historia  domain.HistoriaID
		numero    int
		quiereErr error
	}{
		{"historia inexistente", "HU-99", 1, backlog.ErrNoExiste},
		{"criterio inexistente", "HU-1", 2, backlog.ErrCriterioNoExiste},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s := nuevoServicioBacklog(t)
			h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.AgregarCriterio(h.ID, "A"); err != nil {
				t.Fatal(err)
			}

			if err := s.QuitarCriterio(caso.historia, caso.numero); !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}

			// El criterio que sí existe no debe haberse tocado.
			guardada, _ := s.Buscar("HU-1")
			if len(guardada.Criterios) != 1 || guardada.Criterios[0].Texto != "A" {
				t.Errorf("el criterio existente no debe cambiar; se obtuvo %+v", guardada.Criterios)
			}
		})
	}
}

// nuevoServicioConHistoriaCompletada arma el servicio con una historia "Login" (HU-1) estimada en 3 SP y
// completada. Fija el estado directamente en el repositorio: es preparación, como los pasos "Dado" de los
// escenarios (la historia se completa desde la aplicación con HU-12).
func nuevoServicioConHistoriaCompletada(t *testing.T) *Backlog {
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
	if _, err := proyectos.Crear("Demo", "", inicio, fin); err != nil {
		t.Fatal(err)
	}
	s := NuevoBacklog(repoHistorias, repoProyectos)
	if _, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}

	h, err := repoHistorias.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	h.StoryPoints = 3
	h.Estado = backlog.Completada
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}

	guardada, err := s.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	if guardada.StoryPoints != 3 || guardada.Estado != backlog.Completada {
		t.Fatalf("preparación: se obtuvo %+v", guardada)
	}
	return s
}

func TestBacklog_Estimar(t *testing.T) {
	s := nuevoServicioBacklog(t)
	h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}

	// La estimación queda persistida (spec HU-15, salidas esperadas).
	if err := s.Estimar(h.ID, 3); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	guardada, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	if guardada.StoryPoints != 3 {
		t.Errorf("StoryPoints = %d, se esperaba 3", guardada.StoryPoints)
	}
}

func TestBacklog_Estimar_reemplazaElValorAnterior(t *testing.T) {
	// RN-15.2: se puede re-estimar mientras no esté completada; el valor anterior se reemplaza.
	s := nuevoServicioBacklog(t)
	h, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Estimar(h.ID, 3); err != nil {
		t.Fatal(err)
	}

	if err := s.Estimar(h.ID, 5); err != nil {
		t.Fatalf("re-estimar debería ser válido: %v", err)
	}
	guardada, err := s.Buscar(h.ID)
	if err != nil {
		t.Fatal(err)
	}
	if guardada.StoryPoints != 5 {
		t.Errorf("StoryPoints = %d, se esperaba 5", guardada.StoryPoints)
	}
}

func TestBacklog_Estimar_errores(t *testing.T) {
	casos := []struct {
		nombre    string
		historia  domain.HistoriaID
		sp        int
		quiereErr error
	}{
		{"historia inexistente", "HU-99", 5, backlog.ErrNoExiste},
		{"fuera de la escala", "HU-1", 4, estimation.ErrStoryPointsInvalidos},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			s := nuevoServicioBacklog(t)
			if _, err := s.Crear("P1", "Login", "", backlog.PrioridadAlta); err != nil {
				t.Fatal(err)
			}

			if err := s.Estimar(caso.historia, caso.sp); !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}

			// Un valor inválido no queda guardado.
			guardada, _ := s.Buscar("HU-1")
			if guardada.StoryPoints != 0 {
				t.Errorf("StoryPoints = %d, se esperaba 0", guardada.StoryPoints)
			}
		})
	}
}

func TestBacklog_Estimar_historiaCompletada(t *testing.T) {
	// RN-15.2: no se puede re-estimar una historia completada, y la estimación anterior no cambia.
	s := nuevoServicioConHistoriaCompletada(t)

	if err := s.Estimar("HU-1", 5); !errors.Is(err, backlog.ErrHistoriaCompletada) {
		t.Fatalf("se esperaba ErrHistoriaCompletada, se obtuvo %v", err)
	}
	guardada, err := s.Buscar("HU-1")
	if err != nil {
		t.Fatal(err)
	}
	if guardada.StoryPoints != 3 {
		t.Errorf("StoryPoints = %d, se esperaba 3: la estimación anterior no cambia", guardada.StoryPoints)
	}
}
