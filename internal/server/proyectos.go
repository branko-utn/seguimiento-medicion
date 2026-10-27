package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// rutasProyectos registra las rutas de la sección Proyectos (épica E1, TEC-13): la lista y el alta (HU-01), y el
// detalle con sus integrantes (HU-03).
func (a *app) rutasProyectos(mux *http.ServeMux) {
	mux.HandleFunc("GET /proyectos", a.listarProyectos)
	mux.HandleFunc("GET /proyectos/nuevo", a.nuevoProyecto)
	mux.HandleFunc("POST /proyectos", a.crearProyecto)
	mux.HandleFunc("GET /proyectos/{id}", a.verProyecto)
	mux.HandleFunc("POST /proyectos/{id}/integrantes", a.registrarIntegrante)
}

// formIntegrante es el formulario de alta de integrante del detalle, con el texto tal como se escribió.
type formIntegrante struct {
	Nombre, Rol, Error string
}

// datosDetalle es lo que recibe proyecto.html.
type datosDetalle struct {
	Proyecto    project.Proyecto
	Integrantes []project.Integrante
	Roles       []string
	Form        formIntegrante
}

// rolesIntegrante son los roles que se ofrecen en el formulario, en el orden de la spec HU-03.
var rolesIntegrante = []project.Rol{project.AgileEnabler, project.ProductBuilder, project.ProductArchitect}

func (a *app) verProyecto(w http.ResponseWriter, r *http.Request) {
	a.mostrarDetalle(w, r, http.StatusOK, formIntegrante{})
}

// registrarIntegrante registra un integrante en el proyecto (HU-03) y vuelve al detalle. Si el servicio lo rechaza,
// vuelve a mostrar el detalle con el mensaje de la spec, el nombre escrito y el rol elegido.
func (a *app) registrarIntegrante(w http.ResponseWriter, r *http.Request) {
	id := domain.ProyectoID(r.PathValue("id"))
	form := formIntegrante{Nombre: r.FormValue("nombre"), Rol: r.FormValue("rol")}

	_, err := a.svc.Integrantes.Registrar(id, form.Nombre, rolDesdeTexto(form.Rol))
	switch {
	case errors.Is(err, project.ErrNoExiste):
		a.renderError(w, r, http.StatusNotFound, err.Error())
	case err != nil:
		form.Error = err.Error()
		a.mostrarDetalle(w, r, http.StatusBadRequest, form)
	default:
		http.Redirect(w, r, rutaProyecto(id), http.StatusSeeOther)
	}
}

// mostrarDetalle muestra el proyecto de la ruta con sus integrantes y el formulario de alta de integrante.
func (a *app) mostrarDetalle(w http.ResponseWriter, r *http.Request, codigo int, form formIntegrante) {
	id := domain.ProyectoID(r.PathValue("id"))
	p, err := a.svc.Proyectos.Buscar(id)
	if errors.Is(err, project.ErrNoExiste) {
		a.renderError(w, r, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar el proyecto.")
		return
	}
	integrantes, err := a.svc.Integrantes.DelProyecto(id)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los integrantes.")
		return
	}
	roles := make([]string, 0, len(rolesIntegrante))
	for _, rol := range rolesIntegrante {
		roles = append(roles, rol.String())
	}
	a.renderConCodigo(w, r, codigo, "proyecto.html", datosDetalle{Proyecto: p, Integrantes: integrantes, Roles: roles, Form: form})
}

// rolDesdeTexto convierte el rol elegido en el formulario (ADR 0003: la conversión va en el borde). Si el texto no
// es un rol, devuelve el cero, que el dominio rechaza con project.ErrRolInvalido.
func rolDesdeTexto(texto string) project.Rol {
	for _, rol := range rolesIntegrante {
		if rol.String() == texto {
			return rol
		}
	}
	return 0
}

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
