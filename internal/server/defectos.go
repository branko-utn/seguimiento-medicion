package server

import (
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
)

// rutasDefectos registra la pantalla de defectos (épica E6): la lista con el alta (HU-22, TEC-15) y las acciones para
// cambiar el estado y resolverlos (HU-23).
func (a *app) rutasDefectos(mux *http.ServeMux) {
	mux.HandleFunc("GET /defectos", a.listarDefectos)
	mux.HandleFunc("POST /defectos", a.registrarDefecto)
	mux.HandleFunc("POST /defectos/{defecto}/estado", a.cambiarEstadoDefecto)
	mux.HandleFunc("POST /defectos/{defecto}/resolver", a.resolverDefecto)
}

// datosDefectos es lo que recibe defectos.html.
type datosDefectos struct {
	Filas       []filaDefecto
	Severidades []string
	Historias   []backlog.Historia
	Sprints     []opcionSprint
	Form        formDefecto
}

// filaDefecto es un defecto de la lista con lo que la pantalla puede ofrecer para él.
type filaDefecto struct {
	defect.Defecto
	// TonoSeveridad y TonoEstado son las clases de las insignias (HU-32.3, RN-32.1): las decide el handler.
	TonoSeveridad, TonoEstado string
	// Siguiente es el estado al que pasa con CambiarEstado (Abierto ↔ En progreso); el cero si está resuelto.
	Siguiente defect.Estado
	// Sprints son los sprints donde se puede resolver: del mismo proyecto y que no empiezan antes que el de detección
	// (RN-23.3 y RN-23.4). El servicio lo vuelve a validar.
	Sprints []domain.SprintID
}

// opcionSprint es un sprint del formulario de alta, con el nombre de su proyecto para distinguirlo.
type opcionSprint struct {
	ID       domain.SprintID
	Proyecto string
}

// formDefecto es el formulario de alta de defecto (HU-22), con lo que se escribió para volver a mostrarlo si hay un
// error.
type formDefecto struct {
	Descripcion, Severidad string
	Historia               domain.HistoriaID
	Sprint                 domain.SprintID
	Error                  string
}

// severidadesDefecto son las severidades que se ofrecen en el formulario, en el orden de la spec HU-22.
var severidadesDefecto = []defect.Severidad{defect.Critica, defect.Alta, defect.Media, defect.Baja}

func (a *app) listarDefectos(w http.ResponseWriter, r *http.Request) {
	a.mostrarDefectos(w, r, http.StatusOK, formDefecto{})
}

// registrarDefecto da de alta el defecto (HU-22) y vuelve a la lista. Si el servicio lo rechaza, vuelve a mostrar la
// lista con el mensaje de la spec en el formulario y con lo que se escribió.
func (a *app) registrarDefecto(w http.ResponseWriter, r *http.Request) {
	form := formDefecto{
		Descripcion: r.FormValue("descripcion"),
		Severidad:   r.FormValue("severidad"),
		Historia:    domain.HistoriaID(r.FormValue("historia")),
		Sprint:      domain.SprintID(r.FormValue("sprint")),
	}
	_, err := a.svc.Defectos.Registrar(form.Descripcion, severidadDesdeTexto(form.Severidad), form.Historia, form.Sprint)
	if err != nil {
		form.Error = err.Error()
		a.mostrarDefectos(w, r, http.StatusBadRequest, form)
		return
	}
	http.Redirect(w, r, "/defectos", http.StatusSeeOther)
}

// mostrarDefectos muestra la lista de defectos con sus acciones y el formulario de alta.
func (a *app) mostrarDefectos(w http.ResponseWriter, r *http.Request, codigo int, form formDefecto) {
	datos := datosDefectos{Form: form}
	for _, s := range severidadesDefecto {
		datos.Severidades = append(datos.Severidades, s.String())
	}
	for _, d := range a.svc.Defectos.Listar() {
		fila := filaDefecto{Defecto: d, TonoSeveridad: tonoSeveridad(d.Severidad), TonoEstado: tonoEstado(d.Estado)}
		switch d.Estado {
		case defect.Abierto:
			fila.Siguiente = defect.EnProgreso
		case defect.EnProgreso:
			fila.Siguiente = defect.Abierto
		}
		if d.Estado != defect.Resuelto {
			sprints, err := a.sprintsParaResolver(d)
			if err != nil {
				a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints del defecto "+string(d.ID)+".")
				return
			}
			fila.Sprints = sprints
		}
		datos.Filas = append(datos.Filas, fila)
	}
	historias, err := a.svc.Backlog.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
		return
	}
	datos.Historias = historias
	sprints, err := a.opcionesSprint()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints.")
		return
	}
	datos.Sprints = sprints
	a.renderConCodigo(w, r, codigo, "defectos.html", datos)
}

// opcionesSprint devuelve los sprints de todos los proyectos, proyecto por proyecto.
func (a *app) opcionesSprint() ([]opcionSprint, error) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		return nil, err
	}
	var opciones []opcionSprint
	for _, p := range proyectos {
		sprints, err := a.svc.Sprints.DelProyecto(p.ID)
		if err != nil {
			return nil, err
		}
		for _, s := range sprints {
			opciones = append(opciones, opcionSprint{ID: s.ID, Proyecto: p.Nombre})
		}
	}
	return opciones, nil
}

// sprintsParaResolver devuelve los sprints del proyecto del defecto que no empiezan antes que el de detección.
func (a *app) sprintsParaResolver(d defect.Defecto) ([]domain.SprintID, error) {
	deteccion, err := a.svc.Sprints.Buscar(d.SprintDeteccion)
	if err != nil {
		return nil, err
	}
	sprints, err := a.svc.Sprints.DelProyecto(deteccion.ProyectoID)
	if err != nil {
		return nil, err
	}
	var ids []domain.SprintID
	for _, s := range sprints {
		if !s.Inicio.Before(deteccion.Inicio) {
			ids = append(ids, s.ID)
		}
	}
	return ids, nil
}

func (a *app) cambiarEstadoDefecto(w http.ResponseWriter, r *http.Request) {
	id := domain.DefectoID(r.PathValue("defecto"))
	err := a.svc.Defectos.CambiarEstado(id, estadoDefectoDesdeTexto(r.FormValue("estado")))
	a.responderAccion(w, r, err, "/defectos")
}

func (a *app) resolverDefecto(w http.ResponseWriter, r *http.Request) {
	id := domain.DefectoID(r.PathValue("defecto"))
	err := a.svc.Defectos.Resolver(id, domain.SprintID(r.FormValue("sprint")))
	a.responderAccion(w, r, err, "/defectos")
}

// estadoDefectoDesdeTexto convierte el texto del formulario en un Estado (ADR 0003: la conversión va en el borde).
// Si el texto no es un estado, devuelve el cero, que el dominio rechaza.
func estadoDefectoDesdeTexto(texto string) defect.Estado {
	for _, e := range []defect.Estado{defect.Abierto, defect.EnProgreso, defect.Resuelto} {
		if e.String() == texto {
			return e
		}
	}
	return 0
}

// severidadDesdeTexto convierte la severidad elegida en el formulario (ADR 0003). Si el texto no es una severidad,
// devuelve el cero, que el dominio rechaza con defect.ErrSeveridadInvalida.
func severidadDesdeTexto(texto string) defect.Severidad {
	for _, s := range severidadesDefecto {
		if s.String() == texto {
			return s
		}
	}
	return 0
}
