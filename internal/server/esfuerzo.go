package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

// hoy devuelve la fecha de hoy, sin hora. El dominio la recibe como parámetro (RN-20.4); los tests la fijan.
var hoy = func() time.Time {
	a, m, d := time.Now().Date()
	return time.Date(a, m, d, 0, 0, 0, 0, time.Local)
}

// rutasEsfuerzo registra la pantalla de esfuerzo (épica E5, TEC-15): el registro de horas trabajadas con su lista
// (HU-20) y las horas estimadas de cada historia (HU-19).
func (a *app) rutasEsfuerzo(mux *http.ServeMux) {
	mux.HandleFunc("GET /esfuerzo", a.verEsfuerzo)
	mux.HandleFunc("POST /esfuerzo", a.registrarEsfuerzo)
	mux.HandleFunc("POST /esfuerzo/historias/{historia}/estimacion", a.estimarHistoria)
}

// datosEsfuerzo es lo que recibe esfuerzo.html.
type datosEsfuerzo struct {
	Registros   []filaRegistro
	Integrantes []opcionIntegrante
	Historias   []filaHistoria
	Form        formRegistro
}

// filaRegistro es un registro de esfuerzo de la lista, con el nombre del integrante.
type filaRegistro struct {
	effort.Registro
	Integrante string
}

// opcionIntegrante es un integrante del formulario, con el nombre de su proyecto para distinguirlo.
type opcionIntegrante struct {
	ID       domain.IntegranteID
	Nombre   string
	Proyecto string
}

// filaHistoria es una historia con su estimación en horas (HU-19) y el formulario para cargarla.
type filaHistoria struct {
	ID       domain.HistoriaID
	Titulo   string
	Estimada bool
	Horas    float64
	// Valor es lo que muestra el campo de horas: lo que se escribió si hubo un error, o la estimación guardada.
	Valor string
	Error string
}

// formRegistro es el formulario de registro de esfuerzo (HU-20), con lo que se escribió para volver a mostrarlo si
// hay un error.
type formRegistro struct {
	Integrante              domain.IntegranteID
	Historia                domain.HistoriaID
	Fecha, Actividad, Horas string
	Error                   string
}

// formEstimacion es la estimación que se intentó cargar y no se pudo, para mostrar el error en su historia.
type formEstimacion struct {
	Historia     domain.HistoriaID
	Horas, Error string
}

func (a *app) verEsfuerzo(w http.ResponseWriter, r *http.Request) {
	form := formRegistro{Fecha: hoy().Format(formatoFecha)}
	a.mostrarEsfuerzo(w, r, http.StatusOK, form, formEstimacion{})
}

// registrarEsfuerzo registra las horas trabajadas (HU-20) y vuelve a la pantalla. Si el servicio las rechaza, vuelve
// a mostrar la pantalla con el mensaje de la spec en el formulario y con lo que se escribió.
func (a *app) registrarEsfuerzo(w http.ResponseWriter, r *http.Request) {
	form := formRegistro{
		Integrante: domain.IntegranteID(r.FormValue("integrante")),
		Historia:   domain.HistoriaID(r.FormValue("historia")),
		Fecha:      r.FormValue("fecha"),
		Actividad:  r.FormValue("actividad"),
		Horas:      r.FormValue("horas"),
	}
	_, err := a.svc.Registros.Registrar(form.Integrante, form.Historia, fechaDeFormulario(form.Fecha), form.Actividad, horasDeFormulario(form.Horas), hoy())
	if err != nil {
		form.Error = err.Error()
		a.mostrarEsfuerzo(w, r, http.StatusBadRequest, form, formEstimacion{})
		return
	}
	http.Redirect(w, r, "/esfuerzo", http.StatusSeeOther)
}

// estimarHistoria guarda las horas estimadas de la historia de la ruta (HU-19) y vuelve a la pantalla. Si la historia
// no existe muestra el error 404; si las horas no son válidas, muestra el mensaje de la spec en esa historia.
func (a *app) estimarHistoria(w http.ResponseWriter, r *http.Request) {
	est := formEstimacion{Historia: domain.HistoriaID(r.PathValue("historia")), Horas: r.FormValue("horas")}
	_, err := a.svc.Estimaciones.Estimar(est.Historia, horasDeFormulario(est.Horas))
	switch {
	case errors.Is(err, backlog.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case err != nil:
		est.Error = err.Error()
		a.mostrarEsfuerzo(w, r, http.StatusBadRequest, formRegistro{Fecha: hoy().Format(formatoFecha)}, est)
	default:
		http.Redirect(w, r, "/esfuerzo", http.StatusSeeOther)
	}
}

// mostrarEsfuerzo arma la pantalla: el formulario de registro, la lista de registros y las horas estimadas.
func (a *app) mostrarEsfuerzo(w http.ResponseWriter, r *http.Request, codigo int, form formRegistro, est formEstimacion) {
	integrantes, err := a.opcionesIntegrante()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los integrantes.")
		return
	}
	historias, err := a.svc.Backlog.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
		return
	}
	nombres := make(map[domain.IntegranteID]string, len(integrantes))
	for _, i := range integrantes {
		nombres[i.ID] = i.Nombre
	}

	datos := datosEsfuerzo{Integrantes: integrantes, Form: form}
	for _, h := range historias {
		for _, reg := range a.svc.Registros.DeHistoria(h.ID) {
			datos.Registros = append(datos.Registros, filaRegistro{Registro: reg, Integrante: nombres[reg.IntegranteID]})
		}
		fila := filaHistoria{ID: h.ID, Titulo: h.Titulo}
		e, ok, err := a.svc.Estimaciones.Buscar(h.ID)
		if err != nil {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las horas estimadas.")
			return
		}
		if ok {
			fila.Estimada, fila.Horas = true, e.Horas
			fila.Valor = strconv.FormatFloat(e.Horas, 'f', -1, 64)
		}
		if h.ID == est.Historia {
			fila.Valor, fila.Error = est.Horas, est.Error
		}
		datos.Historias = append(datos.Historias, fila)
	}
	a.renderConCodigo(w, r, codigo, "esfuerzo.html", datos)
}

// opcionesIntegrante devuelve los integrantes de todos los proyectos, proyecto por proyecto.
func (a *app) opcionesIntegrante() ([]opcionIntegrante, error) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		return nil, err
	}
	var opciones []opcionIntegrante
	for _, p := range proyectos {
		integrantes, err := a.svc.Integrantes.DelProyecto(p.ID)
		if err != nil {
			return nil, err
		}
		for _, i := range integrantes {
			opciones = append(opciones, opcionIntegrante{ID: i.ID, Nombre: i.Nombre, Proyecto: p.Nombre})
		}
	}
	return opciones, nil
}

// horasDeFormulario convierte las horas escritas en el formulario; acepta coma o punto decimal ("1,5" o "1.5"). Si no
// es un número devuelve 0, que el dominio rechaza con ErrHorasInvalidas: el borde no inventa mensajes propios.
func horasDeFormulario(valor string) float64 {
	h, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(valor), ",", "."), 64)
	if err != nil {
		return 0
	}
	return h
}
