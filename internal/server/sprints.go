package server

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// rutasSprints registra las rutas de la sección Sprints (épica E3, TEC-14): la lista por proyecto, el alta, el detalle
// y el inicio (HU-10), y la asignación de historias (HU-11).
func (a *app) rutasSprints(mux *http.ServeMux) {
	mux.HandleFunc("GET /sprints", a.listarSprints)
	mux.HandleFunc("GET /sprints/nuevo", a.nuevoSprint)
	mux.HandleFunc("POST /sprints", a.crearSprint)
	mux.HandleFunc("GET /sprints/{sprint}", a.verSprint)
	mux.HandleFunc("POST /sprints/{sprint}/iniciar", a.iniciarSprint)
	mux.HandleFunc("POST /sprints/{sprint}/historias", a.asignarHistoria)
}

// datosSprint es lo que recibe sprint.html. Historias son las del Sprint Backlog y Disponibles las del proyecto que
// todavía no tienen sprint (se ofrecen todas: las reglas de HU-11 las valida el dominio). PuedeIniciar,
// PuedeAsignar y PuedeCerrar los decide el handler con el enum (ADR 0003): la plantilla no compara textos de estado.
// Completadas marca las historias completadas en el sprint (HU-12). Metricas son las del sprint (TEC-17). Error es el mensaje de la última acción que
// falló y Elegida la historia que se había elegido, para volver a seleccionarla.
type datosSprint struct {
	Sprint         sprint.Sprint
	Tono           string // tono de la insignia del estado (TEC-21, HU-32.1)
	Proyecto       project.Proyecto
	Historias      []backlog.Historia
	Completadas    map[domain.HistoriaID]bool
	SPPlanificados int
	SPCompletados  int
	Metricas       metricasSprint
	Disponibles    []backlog.Historia
	PuedeIniciar   bool
	PuedeAsignar   bool
	PuedeCerrar    bool // también habilita completar historias (HU-12): las dos acciones son de un sprint activo
	Cerrado        bool // muestra "No completada" en las historias que volvieron al backlog (HU-14)
	Error          string
	Elegida        string
}

func (a *app) verSprint(w http.ResponseWriter, r *http.Request) {
	a.mostrarSprint(w, r, http.StatusOK, "", "")
}

// asignarHistoria agrega la historia elegida al Sprint Backlog (HU-11) y vuelve al detalle. Si el servicio la
// rechaza, vuelve a mostrar el detalle con el mensaje de la spec y la historia elegida.
func (a *app) asignarHistoria(w http.ResponseWriter, r *http.Request) {
	id := domain.SprintID(r.PathValue("sprint"))
	historia := r.FormValue("historia")
	err := a.svc.Sprints.Asignar(id, domain.HistoriaID(historia))
	switch {
	case errors.Is(err, sprint.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case err != nil:
		a.mostrarSprint(w, r, http.StatusBadRequest, err.Error(), historia)
	default:
		http.Redirect(w, r, rutaSprint(id), http.StatusSeeOther)
	}
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
		a.mostrarSprint(w, r, http.StatusBadRequest, err.Error(), "")
	default:
		http.Redirect(w, r, rutaSprint(id), http.StatusSeeOther)
	}
}

// mostrarSprint muestra el sprint de la ruta con su proyecto, sus historias y las que se le pueden asignar, con el
// código, el mensaje de error y la historia elegida que se le pidan. Si el sprint no existe, muestra el error 404.
func (a *app) mostrarSprint(w http.ResponseWriter, r *http.Request, codigo int, mensaje, elegida string) {
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
	historias, err := a.svc.Backlog.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
		return
	}

	datos := datosSprint{
		Sprint:       sp,
		Tono:         tonoEstado(sp.Estado),
		Proyecto:     p,
		PuedeIniciar: sp.Estado == sprint.Planificado,
		PuedeAsignar: sp.Estado != sprint.Cerrado,
		PuedeCerrar:  sp.Estado == sprint.Activo,
		Cerrado:      sp.Estado == sprint.Cerrado,
		Completadas:  make(map[domain.HistoriaID]bool, len(sp.Completadas)),
		Error:        mensaje,
		Elegida:      elegida,
	}
	for _, id := range sp.Completadas {
		datos.Completadas[id] = true
	}
	delSprint := make(map[domain.HistoriaID]backlog.Historia, len(sp.Historias))
	for _, h := range historias {
		switch {
		case slices.Contains(sp.Historias, h.ID):
			delSprint[h.ID] = h
		case h.ProyectoID == sp.ProyectoID && h.SprintID == "":
			datos.Disponibles = append(datos.Disponibles, h)
		}
	}
	// Las historias se muestran en el orden en que se asignaron, y los SP planificados y completados los calcula HU-24.
	entradas := make([]metrics.HistoriaSprint, 0, len(sp.Historias))
	for _, id := range sp.Historias {
		h := delSprint[id]
		datos.Historias = append(datos.Historias, h)
		entradas = append(entradas, metrics.HistoriaSprint{StoryPoints: h.StoryPoints, Completada: slices.Contains(sp.Completadas, id)})
	}
	if sp.Estado == sprint.Cerrado {
		// RN-14.3: un sprint cerrado muestra los SP que quedaron fijos al cerrarlo.
		datos.SPPlanificados = sp.SPPlanificados
		datos.SPCompletados = sp.SPCompletados
	} else {
		datos.SPPlanificados = metrics.SPPlanificados(entradas)
		datos.SPCompletados = metrics.SPCompletados(entradas)
	}
	datos.Metricas, err = a.metricasSprint(sp)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron calcular las métricas del sprint.")
		return
	}
	a.renderConCodigo(w, r, codigo, "sprint.html", datos)
}

// metricasSprint son las métricas de la pantalla del sprint (TEC-17), con los números ya formateados con un decimal
// (HU-24 y HU-25). La velocidad es la del equipo en el proyecto: el promedio de SP completados de sus sprints
// cerrados (RN-24.3), así que no depende del sprint que se está mirando.
type metricasSprint struct {
	Porcentaje             string
	Completadas, Historias int
	Velocidad              string
	SprintsCerrados        int
}

func (a *app) metricasSprint(sp sprint.Sprint) (metricasSprint, error) {
	porcentaje, err := metrics.PorcentajeCompletadas(len(sp.Completadas), len(sp.Historias))
	if err != nil {
		return metricasSprint{}, err
	}
	sprints, err := a.svc.Sprints.DelProyecto(sp.ProyectoID)
	if err != nil {
		return metricasSprint{}, err
	}
	var completados []int
	for _, s := range sprints {
		if s.Estado == sprint.Cerrado {
			completados = append(completados, s.SPCompletados)
		}
	}
	return metricasSprint{
		Porcentaje:      unDecimal(porcentaje),
		Completadas:     len(sp.Completadas),
		Historias:       len(sp.Historias),
		Velocidad:       unDecimal(metrics.Velocidad(completados)),
		SprintsCerrados: len(completados),
	}, nil
}

// unDecimal muestra el número con un decimal y coma decimal, como piden HU-24 y HU-25 (10,333… se muestra 10,3).
func unDecimal(v float64) string {
	return strings.Replace(strconv.FormatFloat(v, 'f', 1, 64), ".", ",", 1)
}

// filaSprint es un sprint de la lista con el tono de la insignia de su estado (TEC-21, HU-32.1). Embebe el sprint, así
// la plantilla sigue usando .ID, .Goal y sus fechas.
type filaSprint struct {
	sprint.Sprint
	Tono string
}

// grupoSprints es un proyecto con todos sus sprints, en el orden en que se crearon, y sus cerrados para el
// historial (HU-14).
type grupoSprints struct {
	Proyecto project.Proyecto
	Sprints  []filaSprint
	Cerrados []sprint.Sprint
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
		cerrados, err := a.svc.Sprints.Cerrados(p.ID)
		if err != nil {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints cerrados.")
			return
		}
		filas := make([]filaSprint, 0, len(sprints))
		for _, sp := range sprints {
			filas = append(filas, filaSprint{Sprint: sp, Tono: tonoEstado(sp.Estado)})
		}
		grupos = append(grupos, grupoSprints{Proyecto: p, Sprints: filas, Cerrados: cerrados})
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
