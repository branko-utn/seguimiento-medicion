package server

import (
	"errors"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// rutasSprints registra las rutas de la sección Sprints (épica E3, TEC-14): la lista por proyecto, el alta, el detalle
// y el inicio (HU-10).
func (a *app) rutasSprints(mux *http.ServeMux) {
	mux.HandleFunc("GET /sprints", a.listarSprints)
	mux.HandleFunc("GET /sprints/nuevo", a.nuevoSprint)
	mux.HandleFunc("POST /sprints", a.crearSprint)
	mux.HandleFunc("GET /sprints/{sprint}", a.verSprint)
	mux.HandleFunc("POST /sprints/{sprint}/iniciar", a.iniciarSprint)
	mux.HandleFunc("POST /sprints/{sprint}/historias", a.asignarHistoria)
}

// asignarHistoria: asignación de historias al sprint (HU-11). Todavía sin implementar (RED).
func (a *app) asignarHistoria(w http.ResponseWriter, r *http.Request) {}

// datosSprint es lo que recibe sprint.html. Error es el mensaje de la última acción que falló, si hubo.
// PuedeIniciar lo decide el handler con el enum (ADR 0003): la plantilla no compara textos de estado.
type datosSprint struct {
	Sprint       sprint.Sprint
	Proyecto     project.Proyecto
	PuedeIniciar bool
	Error        string
}

func (a *app) verSprint(w http.ResponseWriter, r *http.Request) {
	a.mostrarSprint(w, r, http.StatusOK, "")
}

// iniciarSprint pasa el sprint a Activo (HU-10) y vuelve a su detalle. Si el servicio lo rechaza, vuelve a mostrar
// el detalle con el mensaje de la spec.
func (a *app) iniciarSprint(w http.ResponseWriter, r *http.Request) {
	id := domain.SprintID(r.PathValue("sprint"))
	err := a.svc.Sprints.Iniciar(id)
	switch {
	case errors.Is(err, sprint.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case err != nil:
		a.mostrarSprint(w, r, http.StatusBadRequest, err.Error())
	default:
		http.Redirect(w, r, rutaSprint(id), http.StatusSeeOther)
	}
}

// mostrarSprint muestra el sprint de la ruta con el nombre de su proyecto, con el código y el mensaje de error que se
// le pidan. Si el sprint no existe, muestra el error 404.
func (a *app) mostrarSprint(w http.ResponseWriter, r *http.Request, codigo int, mensaje string) {
	sp, err := a.svc.Sprints.Buscar(domain.SprintID(r.PathValue("sprint")))
	if errors.Is(err, sprint.ErrNoExiste) {
		a.renderError(w, r, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar el sprint.")
		return
	}
	p, err := a.svc.Proyectos.Buscar(sp.ProyectoID)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar el proyecto del sprint.")
		return
	}
	a.renderConCodigo(w, r, codigo, "sprint.html", datosSprint{
		Sprint:       sp,
		Proyecto:     p,
		PuedeIniciar: sp.Estado == sprint.Planificado,
		Error:        mensaje,
	})
}

// grupoSprints es un proyecto con sus sprints, en el orden en que se crearon.
type grupoSprints struct {
	Proyecto project.Proyecto
	Sprints  []sprint.Sprint
}

// formSprint es lo que recibe sprint_form.html. Los campos guardan el texto tal como se escribió, para volver a
// mostrarlo si hay un error; Proyectos son las opciones del desplegable.
type formSprint struct {
	Proyecto, Goal, Inicio, Fin string
	Error                       string
	Proyectos                   []project.Proyecto
}

func (a *app) listarSprints(w http.ResponseWriter, r *http.Request) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los proyectos.")
		return
	}
	grupos := make([]grupoSprints, 0, len(proyectos))
	for _, p := range proyectos {
		sprints, err := a.svc.Sprints.DelProyecto(p.ID)
		if err != nil {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints.")
			return
		}
		grupos = append(grupos, grupoSprints{Proyecto: p, Sprints: sprints})
	}
	a.render(w, r, "sprints.html", grupos)
}

// nuevoSprint muestra el formulario de alta. El proyecto del enlace (?proyecto=P1) viene elegido.
func (a *app) nuevoSprint(w http.ResponseWriter, r *http.Request) {
	a.formularioSprint(w, r, http.StatusOK, formSprint{Proyecto: r.URL.Query().Get("proyecto")})
}

// crearSprint da de alta el sprint (HU-10) y redirige a su detalle. Si el servicio lo rechaza, vuelve a mostrar el
// formulario con el mensaje de la spec y con lo que se escribió.
func (a *app) crearSprint(w http.ResponseWriter, r *http.Request) {
	form := formSprint{
		Proyecto: r.FormValue("proyecto"),
		Goal:     r.FormValue("goal"),
		Inicio:   r.FormValue("inicio"),
		Fin:      r.FormValue("fin"),
	}
	sp, err := a.svc.Sprints.Crear(domain.ProyectoID(form.Proyecto), form.Goal, fechaDeFormulario(form.Inicio), fechaDeFormulario(form.Fin))
	if err != nil {
		form.Error = err.Error()
		a.formularioSprint(w, r, http.StatusBadRequest, form)
		return
	}
	http.Redirect(w, r, rutaSprint(sp.ID), http.StatusSeeOther)
}

// formularioSprint completa las opciones del desplegable de proyectos y muestra el formulario de alta.
func (a *app) formularioSprint(w http.ResponseWriter, r *http.Request, codigo int, form formSprint) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los proyectos.")
		return
	}
	form.Proyectos = proyectos
	a.renderConCodigo(w, r, codigo, "sprint_form.html", form)
}

func rutaSprint(id domain.SprintID) string {
	return "/sprints/" + string(id)
}
