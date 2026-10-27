package server

import (
	"errors"
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

// rutasAccionesSprint registra las acciones sobre un sprint de HU-12 (completar una historia) y HU-13 (cerrar el
// sprint). La pantalla del sprint (TEC-14) solo agrega los botones: un formulario o un hx-post a estas rutas.
func (a *app) rutasAccionesSprint(mux *http.ServeMux) {
	mux.HandleFunc("POST /sprints/{sprint}/historias/{historia}/completar", a.completarHistoria)
	mux.HandleFunc("POST /sprints/{sprint}/cerrar", a.cerrarSprint)
}

func (a *app) completarHistoria(w http.ResponseWriter, r *http.Request) {
	sprintID := domain.SprintID(r.PathValue("sprint"))
	err := a.svc.Sprints.CompletarHistoria(sprintID, domain.HistoriaID(r.PathValue("historia")))
	a.responderAccion(w, r, err, "/sprints/"+string(sprintID))
}

func (a *app) cerrarSprint(w http.ResponseWriter, r *http.Request) {
	sprintID := domain.SprintID(r.PathValue("sprint"))
	err := a.svc.Sprints.CerrarSprint(sprintID)
	a.responderAccion(w, r, err, "/sprints/"+string(sprintID))
}

// responderAccion vuelve a la pantalla del sprint si la acción salió bien (redirección normal, o HX-Redirect si el
// pedido viene de HTMX) y, si no, muestra el error con su mensaje: 404 si algo no existe, 400 si se rompió una regla.
func (a *app) responderAccion(w http.ResponseWriter, r *http.Request, err error, destino string) {
	if err != nil {
		codigo := http.StatusBadRequest
		if errors.Is(err, sprint.ErrNoExiste) || errors.Is(err, backlog.ErrNoExiste) {
			codigo = http.StatusNotFound
		}
		a.renderError(w, r, codigo, err.Error())
		return
	}
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", destino)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, destino, http.StatusSeeOther)
}
