package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/web"
)

// app tiene lo que comparten todos los handlers: los servicios y las plantillas.
type app struct {
	svc     Servicios
	paginas map[string]*template.Template
}

// vista es el dato que recibe toda plantilla. Los datos propios de cada página van en Datos.
type vista struct {
	Ruta      string
	Secciones []seccion
	Datos     any
}

// Activa indica si la ruta corresponde a la sección que se está viendo (para resaltarla en la navegación).
func (v vista) Activa(ruta string) bool {
	return v.Ruta == ruta || strings.HasPrefix(v.Ruta, ruta+"/")
}

type datosError struct {
	Codigo  int
	Texto   string
	Mensaje string
}

func nuevaApp(svc Servicios) (*app, error) {
	paginas, err := cargarPaginas(web.FS)
	if err != nil {
		return nil, err
	}
	return &app{svc: svc, paginas: paginas}, nil
}

// cargarPaginas arma una plantilla por cada archivo de templates/paginas/, combinada con el layout base.
func cargarPaginas(fsys fs.FS) (map[string]*template.Template, error) {
	archivos, err := fs.Glob(fsys, "templates/paginas/*.html")
	if err != nil {
		return nil, err
	}
	if len(archivos) == 0 {
		return nil, fmt.Errorf("no hay plantillas en templates/paginas")
	}
	paginas := make(map[string]*template.Template, len(archivos))
	for _, archivo := range archivos {
		t, err := template.ParseFS(fsys, "templates/base.html", archivo)
		if err != nil {
			return nil, fmt.Errorf("cargar plantilla %s: %w", archivo, err)
		}
		paginas[path.Base(archivo)] = t
	}
	return paginas, nil
}

// render muestra una página de templates/paginas/ con código 200.
func (a *app) render(w http.ResponseWriter, r *http.Request, pagina string, datos any) {
	a.renderConCodigo(w, r, http.StatusOK, pagina, datos)
}

// renderError muestra la página de error con el código y el mensaje indicados.
func (a *app) renderError(w http.ResponseWriter, r *http.Request, codigo int, mensaje string) {
	a.renderConCodigo(w, r, codigo, "error.html", datosError{Codigo: codigo, Texto: http.StatusText(codigo), Mensaje: mensaje})
}

func (a *app) renderConCodigo(w http.ResponseWriter, r *http.Request, codigo int, pagina string, datos any) {
	t, ok := a.paginas[pagina]
	if !ok {
		log.Printf("render: no existe la página %q", pagina)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	// Se renderiza primero en memoria para no mandar una página a medias si la plantilla falla.
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", vista{Ruta: r.URL.Path, Secciones: secciones, Datos: datos}); err != nil {
		log.Printf("render %s: %v", pagina, err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(codigo)
	_, _ = buf.WriteTo(w)
}
