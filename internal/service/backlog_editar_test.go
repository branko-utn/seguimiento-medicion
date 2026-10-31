package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// esfuerzoFalso cumple EsfuerzoDeHistorias con registros cargados a mano por historia.
type esfuerzoFalso map[domain.HistoriaID][]effort.Registro

func (f esfuerzoFalso) DeHistoria(id domain.HistoriaID) []effort.Registro { return f[id] }

// defectosFalsos cumple ListadorDefectos con una lista fija.
type defectosFalsos []defect.Defecto

func (f defectosFalsos) Listar() []defect.Defecto { return f }

// backlogParaEditar arma el servicio con el proyecto P1 y las historias HU-1 "Crear proyecto" (Media, 3 SP, con un
// criterio) y HU-2 "Crear sprint". Devuelve también el repositorio, para preparar estados que no son de HU-06.
func backlogParaEditar(t *testing.T, esfuerzo EsfuerzoDeHistorias, defectos ListadorDefectos) (*Backlog, *memory.Historias) {
	t.Helper()
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoSprints, _ := memory.NuevoSprints("")
	if _, err := NuevoProyectos(repoProyectos).Crear("Demo", "", inicio, fin); err != nil {
		t.Fatal(err)
	}
	s := NuevoBacklog(repoHistorias, repoProyectos, NuevoSprints(repoSprints, repoProyectos, repoHistorias), esfuerzo, defectos)
	for _, titulo := range []string{"Crear proyecto", "Crear sprint"} {
		if _, err := s.Crear("P1", titulo, "", backlog.PrioridadMedia); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Estimar("HU-1", 3); err != nil {
		t.Fatal(err)
	}
	if err := s.AgregarCriterio("HU-1", "Valida el nombre"); err != nil {
		t.Fatal(err)
	}
	return s, repoHistorias
}

// prepararHistoria cambia la historia directo en el repositorio (estado y sprint).
func prepararHistoria(t *testing.T, repo *memory.Historias, id domain.HistoriaID, estado backlog.Estado, sprint domain.SprintID) {
	t.Helper()
	h, err := repo.Buscar(id)
	if err != nil {
		t.Fatal(err)
	}
	h.Estado, h.SprintID = estado, sprint
	if err := repo.Guardar(h); err != nil {
		t.Fatal(err)
	}
}

func TestBacklog_Editar(t *testing.T) {
	s, repo := backlogParaEditar(t, esfuerzoFalso{}, defectosFalsos{})
	prepararHistoria(t, repo, "HU-1", backlog.EnProgreso, "S1")
	antes, _ := s.Buscar("HU-1")

	// Una historia asignada a un sprint (no completada) se puede editar.
	if err := s.Editar("HU-1", "  Crear y listar proyectos  ", "Alta y listado", backlog.PrioridadAlta); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	got, _ := s.Buscar("HU-1")
	want := antes
	want.Titulo, want.Descripcion, want.Prioridad = "Crear y listar proyectos", "Alta y listado", backlog.PrioridadAlta
	// El ID, el estado, los Story Points, los criterios y el sprint no cambian.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("historia guardada:\nquedó    %+v\nesperada %+v", got, want)
	}
}

func TestBacklog_EditarErrores(t *testing.T) {
	casos := []struct {
		nombre    string
		id        domain.HistoriaID
		completar bool
		titulo    string
		prioridad backlog.Prioridad
		wantErr   error
	}{
		{"historia inexistente", "HU-99", false, "Login", backlog.PrioridadAlta, backlog.ErrNoExiste},
		{"título vacío (RN-06.1)", "HU-1", false, "   ", backlog.PrioridadAlta, backlog.ErrTituloVacio},
		{"prioridad inválida (RN-06.1)", "HU-1", false, "Login", backlog.Prioridad(0), backlog.ErrPrioridadInvalida},
		{"historia completada (RN-06.2)", "HU-1", true, "Otro título", backlog.PrioridadAlta, backlog.ErrCompletadaNoSeEdita},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s, repo := backlogParaEditar(t, esfuerzoFalso{}, defectosFalsos{})
			if c.completar {
				prepararHistoria(t, repo, "HU-1", backlog.Completada, "S1")
			}
			antes, _ := s.Buscar("HU-1")

			if err := s.Editar(c.id, c.titulo, "", c.prioridad); !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, c.wantErr)
			}
			// Si falla, no se guarda nada.
			if got, _ := s.Buscar("HU-1"); !reflect.DeepEqual(got, antes) {
				t.Errorf("con error la historia no debería cambiar:\nquedó %+v\nantes %+v", got, antes)
			}
		})
	}
}

func TestBacklog_Eliminar(t *testing.T) {
	// Registros y defectos de otra historia no impiden eliminar HU-2.
	esfuerzo := esfuerzoFalso{"HU-1": {{ID: "E1", HistoriaID: "HU-1", Horas: 2}}}
	defectos := defectosFalsos{{ID: "D1", HistoriaID: "HU-1"}}
	s, _ := backlogParaEditar(t, esfuerzo, defectos)

	if err := s.Eliminar("HU-2"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if _, err := s.Buscar("HU-2"); !errors.Is(err, backlog.ErrNoExiste) {
		t.Errorf("Buscar HU-2 después de eliminar: error = %v, se esperaba ErrNoExiste", err)
	}
	lista, _ := s.Listar()
	if len(lista) != 1 || lista[0].ID != "HU-1" {
		t.Errorf("Listar() = %+v, se esperaba solo HU-1", lista)
	}

	// RN-06.4: la próxima historia no reutiliza el ID eliminado.
	h, err := s.Crear("P1", "Cerrar sprint", "", backlog.PrioridadBaja)
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != "HU-3" {
		t.Errorf("ID = %q, se esperaba HU-3", h.ID)
	}
}

func TestBacklog_EliminarErrores(t *testing.T) {
	casos := []struct {
		nombre   string
		id       domain.HistoriaID
		estado   backlog.Estado
		sprint   domain.SprintID
		esfuerzo esfuerzoFalso
		defectos defectosFalsos
		wantErr  error
	}{
		{"historia inexistente", "HU-99", backlog.Pendiente, "", nil, nil, backlog.ErrNoExiste},
		{"asignada a un sprint (RN-06.3)", "HU-2", backlog.Pendiente, "S1", nil, nil, backlog.ErrAsignadaNoSeElimina},
		{"completada (RN-06.3)", "HU-2", backlog.Completada, "S1", nil, nil, backlog.ErrCompletadaNoSeElimina},
		{"con esfuerzo registrado (RN-06.5)", "HU-2", backlog.Pendiente, "",
			esfuerzoFalso{"HU-2": {{ID: "E1", HistoriaID: "HU-2", Horas: 2}}}, nil, ErrHistoriaConRegistros},
		{"con defectos (RN-06.5)", "HU-2", backlog.Pendiente, "",
			nil, defectosFalsos{{ID: "D1", HistoriaID: "HU-2"}}, ErrHistoriaConRegistros},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s, repo := backlogParaEditar(t, c.esfuerzo, c.defectos)
			prepararHistoria(t, repo, "HU-2", c.estado, c.sprint)

			if err := s.Eliminar(c.id); !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, c.wantErr)
			}
			// Si falla, la historia sigue en el backlog.
			if lista, _ := s.Listar(); len(lista) != 2 {
				t.Errorf("el backlog tiene %d historias, se esperaban 2", len(lista))
			}
		})
	}
}
