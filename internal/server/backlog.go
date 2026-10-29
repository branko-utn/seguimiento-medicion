package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// rutasBacklog registra las rutas de la sección Backlog (épica E2, TEC-16): la lista y el alta de historias (HU-05),
// la edición y la eliminación (HU-06) y el detalle con sus criterios (HU-07), su estado (HU-08) y su estimación en
// Story Points (HU-15).
func (a *app) rutasBacklog(mux *http.ServeMux) {
	mux.HandleFunc("GET /backlog", a.listarBacklog)
	mux.HandleFunc("GET /backlog/nueva", a.nuevaHistoria)
	mux.HandleFunc("POST /backlog", a.crearHistoria)
	mux.HandleFunc("GET /backlog/{historia}", a.verHistoria)
	mux.HandleFunc("POST /backlog/{historia}/editar", a.editarHistoria)
	mux.HandleFunc("POST /backlog/{historia}/eliminar", a.eliminarHistoria)
	mux.HandleFunc("POST /backlog/{historia}/criterios", a.agregarCriterio)
	mux.HandleFunc("POST /backlog/{historia}/criterios/{numero}/quitar", a.quitarCriterio)
	mux.HandleFunc("POST /backlog/{historia}/estimacion", a.estimarStoryPoints)
	mux.HandleFunc("POST /backlog/{historia}/estado", a.cambiarEstadoHistoria)
}

// datosBacklog es lo que recibe la plantilla del backlog.
type datosBacklog struct {
	Historias []backlog.Historia
	Proyectos map[string]string
}

func (a *app) listarBacklog(w http.ResponseWriter, r *http.Request) {
	var historias []backlog.Historia
	if a.svc.Backlog != nil {
		if h, err := a.svc.Backlog.Listar(); err == nil {
			historias = h
		} else {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
			return
		}
	}
	nombres := make(map[string]string)
	if a.svc.Proyectos != nil {
		if ps, err := a.svc.Proyectos.Listar(); err == nil {
			for _, p := range ps {
				nombres[string(p.ID)] = p.Nombre
			}
		}
	}
	a.render(w, r, "backlog.html", datosBacklog{Historias: historias, Proyectos: nombres})
}

// prioridades son las prioridades que se ofrecen en el formulario, en el orden de la spec HU-05.
var prioridades = []backlog.Prioridad{backlog.PrioridadAlta, backlog.PrioridadMedia, backlog.PrioridadBaja}

// estadosHistoria son los estados que se ofrecen en el detalle, en el orden de la spec HU-08.
var estadosHistoria = []backlog.Estado{backlog.Pendiente, backlog.EnProgreso, backlog.Completada}

// formHistoria es el formulario de alta de historia, con el texto tal como se escribió.
type formHistoria struct {
	Proyecto, Titulo, Descripcion, Prioridad, Error string
	Proyectos                                       []project.Proyecto
	Prioridades                                     []string
}

func (a *app) nuevaHistoria(w http.ResponseWriter, r *http.Request) {
	a.mostrarFormHistoria(w, r, http.StatusOK, formHistoria{Prioridad: backlog.PrioridadMedia.String()})
}

// crearHistoria agrega la historia al backlog del proyecto elegido (HU-05) y redirige a su detalle. Si el servicio
// la rechaza (también si el proyecto no existe: es un dato del formulario), vuelve a mostrar el formulario con el
// mensaje de la spec y con lo que se escribió.
func (a *app) crearHistoria(w http.ResponseWriter, r *http.Request) {
	form := formHistoria{
		Proyecto:    r.FormValue("proyecto"),
		Titulo:      r.FormValue("titulo"),
		Descripcion: r.FormValue("descripcion"),
		Prioridad:   r.FormValue("prioridad"),
	}
	h, err := a.svc.Backlog.Crear(domain.ProyectoID(form.Proyecto), form.Titulo, form.Descripcion, prioridadDesdeTexto(form.Prioridad))
	if err != nil {
		form.Error = err.Error()
		a.mostrarFormHistoria(w, r, http.StatusBadRequest, form)
		return
	}
	http.Redirect(w, r, rutaHistoria(h.ID), http.StatusSeeOther)
}

func (a *app) mostrarFormHistoria(w http.ResponseWriter, r *http.Request, codigo int, form formHistoria) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los proyectos.")
		return
	}
	form.Proyectos = proyectos
	for _, p := range prioridades {
		form.Prioridades = append(form.Prioridades, p.String())
	}
	a.renderConCodigo(w, r, codigo, "historia_form.html", form)
}

// datosHistoria es lo que recibe historia.html. Los criterios van numerados desde 1, como los quita HU-07, y
// EstadoActual es el texto del estado para marcarlo en el desplegable; Prioridades y PrioridadActual, lo mismo para
// el formulario de edición (HU-06). Error es el mensaje de la última acción que falló.
type datosHistoria struct {
	Historia        backlog.Historia
	Proyecto        string
	Criterios       []criterioNumerado
	Escala          []int
	Estados         []string
	EstadoActual    string
	Prioridades     []string
	PrioridadActual string
	Error           string
}

type criterioNumerado struct {
	Numero int
	Texto  string
}

func (a *app) verHistoria(w http.ResponseWriter, r *http.Request) {
	a.mostrarHistoria(w, r, http.StatusOK, "")
}

// editarHistoria cambia el título, la descripción y la prioridad de la historia (HU-06) y vuelve a su detalle.
func (a *app) editarHistoria(w http.ResponseWriter, r *http.Request) {
	id := domain.HistoriaID(r.PathValue("historia"))
	err := a.svc.Backlog.Editar(id, r.FormValue("titulo"), r.FormValue("descripcion"), prioridadDesdeTexto(r.FormValue("prioridad")))
	a.responderHistoria(w, r, err)
}

// eliminarHistoria saca la historia del backlog (HU-06) y vuelve a la lista. Si no se puede eliminar, muestra el
// detalle con el mensaje de la spec (o el 404 si no existe), igual que las demás acciones del detalle.
func (a *app) eliminarHistoria(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.Backlog.Eliminar(domain.HistoriaID(r.PathValue("historia"))); err != nil {
		a.responderHistoria(w, r, err)
		return
	}
	http.Redirect(w, r, "/backlog", http.StatusSeeOther)
}

// agregarCriterio agrega un criterio de aceptación a la historia (HU-07).
func (a *app) agregarCriterio(w http.ResponseWriter, r *http.Request) {
	id := domain.HistoriaID(r.PathValue("historia"))
	a.responderHistoria(w, r, a.svc.Backlog.AgregarCriterio(id, r.FormValue("texto")))
}

// quitarCriterio quita el criterio de la posición {numero}, empezando en 1 (HU-07). Un número que no es un número
// llega al dominio como 0, que responde "el criterio no existe".
func (a *app) quitarCriterio(w http.ResponseWriter, r *http.Request) {
	id := domain.HistoriaID(r.PathValue("historia"))
	numero, _ := strconv.Atoi(r.PathValue("numero"))
	a.responderHistoria(w, r, a.svc.Backlog.QuitarCriterio(id, numero))
}

// estimarStoryPoints asigna los Story Points elegidos (HU-15). Un valor que no es un número llega al dominio como
// 0, que no es de la escala.
func (a *app) estimarStoryPoints(w http.ResponseWriter, r *http.Request) {
	id := domain.HistoriaID(r.PathValue("historia"))
	sp, _ := strconv.Atoi(r.FormValue("story_points"))
	a.responderHistoria(w, r, a.svc.Backlog.Estimar(id, sp))
}

// cambiarEstadoHistoria deja la historia en el estado elegido (HU-08). Completada delega en HU-12.
func (a *app) cambiarEstadoHistoria(w http.ResponseWriter, r *http.Request) {
	id := domain.HistoriaID(r.PathValue("historia"))
	a.responderHistoria(w, r, a.svc.Backlog.CambiarEstado(id, estadoDesdeTexto(r.FormValue("estado"))))
}

// responderHistoria vuelve al detalle de la historia si la acción salió bien. Si la historia no existe muestra el
// error 404 y, si se rompió una regla, el detalle con el mensaje de la spec (400).
func (a *app) responderHistoria(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, backlog.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case err != nil:
		a.mostrarHistoria(w, r, http.StatusBadRequest, err.Error())
	default:
		http.Redirect(w, r, rutaHistoria(domain.HistoriaID(r.PathValue("historia"))), http.StatusSeeOther)
	}
}

// mostrarHistoria muestra la historia de la ruta con su proyecto, la escala de Story Points y los estados, con el
// código y el mensaje de error que se le pidan. Si la historia no existe, muestra el error 404.
func (a *app) mostrarHistoria(w http.ResponseWriter, r *http.Request, codigo int, mensaje string) {
	h, err := a.svc.Backlog.Buscar(domain.HistoriaID(r.PathValue("historia")))
	if errors.Is(err, backlog.ErrNoExiste) {
		a.renderError(w, r, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar la historia.")
		return
	}
	datos := datosHistoria{Historia: h, Escala: estimation.Escala, EstadoActual: h.Estado.String(),
		PrioridadActual: h.Prioridad.String(), Error: mensaje}
	for i, c := range h.Criterios {
		datos.Criterios = append(datos.Criterios, criterioNumerado{Numero: i + 1, Texto: c.Texto})
	}
	if p, err := a.svc.Proyectos.Buscar(h.ProyectoID); err == nil {
		datos.Proyecto = p.Nombre
	}
	for _, e := range estadosHistoria {
		datos.Estados = append(datos.Estados, e.String())
	}
	for _, p := range prioridades {
		datos.Prioridades = append(datos.Prioridades, p.String())
	}
	a.renderConCodigo(w, r, codigo, "historia.html", datos)
}

func rutaHistoria(id domain.HistoriaID) string {
	return "/backlog/" + string(id)
}

// prioridadDesdeTexto convierte la prioridad elegida en el formulario (ADR 0003: la conversión va en el borde). Si
// el texto no es una prioridad, devuelve el cero, que el dominio rechaza con backlog.ErrPrioridadInvalida.
func prioridadDesdeTexto(texto string) backlog.Prioridad {
	for _, p := range prioridades {
		if p.String() == texto {
			return p
		}
	}
	return 0
}

// estadoDesdeTexto convierte el estado elegido en el detalle. Si el texto no es un estado, devuelve el cero, que el
// dominio rechaza con backlog.ErrEstadoInvalido.
func estadoDesdeTexto(texto string) backlog.Estado {
	for _, e := range estadosHistoria {
		if e.String() == texto {
			return e
		}
	}
	return 0
}
