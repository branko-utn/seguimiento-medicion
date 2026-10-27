package server

import (
	"net/http"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// rutasProyectos registra las rutas de la sección Proyectos (épica E1, TEC-13): la lista y el alta (HU-01).
func (a *app) rutasProyectos(mux *http.ServeMux) {
	mux.HandleFunc("GET /proyectos", a.listarProyectos)
	mux.HandleFunc("GET /proyectos/nuevo", a.nuevoProyecto)
	mux.HandleFunc("POST /proyectos", a.crearProyecto)
	mux.HandleFunc("GET /proyectos/{id}", a.verProyecto)
	mux.HandleFunc("POST /proyectos/{id}/integrantes", a.registrarIntegrante)
}

// verProyecto y registrarIntegrante: detalle del proyecto e integrantes (HU-03). Todavía sin implementar (RED).
func (a *app) verProyecto(w http.ResponseWriter, r *http.Request) {}

func (a *app) registrarIntegrante(w http.ResponseWriter, r *http.Request) {}

// formatoFecha es el formato de los campos <input type="date">.
const formatoFecha = "2006-01-02"

// formProyecto es lo que recibe proyecto_form.html, el formulario que comparten el alta (HU-01) y la edición (HU-02).
// Los campos guardan el texto tal como se escribió, para volver a mostrarlo si hay un error.
type formProyecto struct {
	Titulo, Accion, Boton string
	Nombre, Descripcion   string
	Inicio, Fin           string
	Error                 string
}

func (a *app) listarProyectos(w http.ResponseWriter, r *http.Request) {
	proyectos, err := a.svc.Proyectos.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los proyectos.")
		return
	}
	a.render(w, r, "proyectos.html", proyectos)
}

func (a *app) nuevoProyecto(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "proyecto_form.html", formAlta())
}

// crearProyecto da de alta el proyecto (HU-01) y redirige a su detalle. Si el servicio lo rechaza, vuelve a mostrar
// el formulario con el mensaje de la spec y con lo que se escribió.
func (a *app) crearProyecto(w http.ResponseWriter, r *http.Request) {
	form := formAlta()
	form.leer(r)

	p, err := a.svc.Proyectos.Crear(form.Nombre, form.Descripcion, fechaDeFormulario(form.Inicio), fechaDeFormulario(form.Fin))
	if err != nil {
		form.Error = err.Error()
		a.renderConCodigo(w, r, http.StatusBadRequest, "proyecto_form.html", form)
		return
	}
	http.Redirect(w, r, rutaProyecto(p.ID), http.StatusSeeOther)
}

func formAlta() formProyecto {
	return formProyecto{Titulo: "Nuevo proyecto", Accion: "/proyectos", Boton: "Crear proyecto"}
}

// leer copia los campos enviados en el formulario.
func (f *formProyecto) leer(r *http.Request) {
	f.Nombre = r.FormValue("nombre")
	f.Descripcion = r.FormValue("descripcion")
	f.Inicio = r.FormValue("inicio")
	f.Fin = r.FormValue("fin")
}

// fechaDeFormulario convierte el valor de un <input type="date">. Si está vacío o mal formado devuelve la fecha cero,
// que el dominio rechaza con su error de fechas (RN-01.2): el borde no inventa mensajes propios.
func fechaDeFormulario(valor string) time.Time {
	f, err := time.ParseInLocation(formatoFecha, valor, time.Local)
	if err != nil {
		return time.Time{}
	}
	return f
}

func rutaProyecto(id domain.ProyectoID) string {
	return "/proyectos/" + string(id)
}
