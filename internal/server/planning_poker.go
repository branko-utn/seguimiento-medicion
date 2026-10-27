package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/estimation"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// rutasPlanningPoker registra las rutas de la sección Planning Poker (épica E4): la lista de historias
// con sus votaciones, abrir una votación y votar (HU-16).
func (a *app) rutasPlanningPoker(mux *http.ServeMux) {
	mux.HandleFunc("GET /planning-poker", a.listarPlanningPoker)
	mux.HandleFunc("POST /planning-poker/abrir", a.abrirVotacion)
	mux.HandleFunc("POST /planning-poker/votar", a.votar)
}

// filaPoker es lo que se muestra de cada historia: si tiene la votación abierta, quiénes votaron
// (solo los nombres: RN-16.4) y los integrantes habilitados para votar (RN-16.1).
type filaPoker struct {
	Historia    backlog.Historia
	Abierta     bool
	Completada  bool
	Votaron     string
	Integrantes []project.Integrante
}

// datosPoker es lo que recibe la plantilla de la sección.
type datosPoker struct {
	Filas  []filaPoker
	Escala []int
}

// listarPlanningPoker muestra el backlog con el estado de la votación de cada historia (HU-16).
// Los valores de los votos nunca llegan a la plantilla: la pantalla solo muestra quién votó (RN-16.4).
func (a *app) listarPlanningPoker(w http.ResponseWriter, r *http.Request) {
	historias, err := a.svc.Backlog.Listar()
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar las historias.")
		return
	}
	filas := make([]filaPoker, 0, len(historias))
	for _, h := range historias {
		// RN-16.1: en la pantalla solo se elige entre los integrantes del proyecto de la historia.
		integrantes, err := a.svc.Integrantes.DelProyecto(h.ProyectoID)
		if err != nil {
			a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los integrantes.")
			return
		}
		fila := filaPoker{Historia: h, Completada: h.Estado == backlog.Completada, Integrantes: integrantes}

		votantes, err := a.svc.Poker.QuienesVotaron(h.ID)
		switch {
		case err == nil:
			fila.Abierta = true
			if nombres := nombresDe(integrantes, votantes); len(nombres) > 0 {
				fila.Votaron = strings.Join(nombres, ", ")
			} else {
				fila.Votaron = "todavía nadie"
			}
		case errors.Is(err, estimation.ErrSinVotacionAbierta):
			// Sin votación abierta: la fila ofrece abrirla (si la historia no está completada, RN-16.5).
		default:
			a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar la votación.")
			return
		}
		filas = append(filas, fila)
	}
	a.render(w, r, "planning_poker.html", datosPoker{Filas: filas, Escala: estimation.Escala})
}

// abrirVotacion inicia la votación de la historia (HU-16) y vuelve a la lista.
func (a *app) abrirVotacion(w http.ResponseWriter, r *http.Request) {
	if _, err := a.svc.Poker.IniciarVotacion(domain.HistoriaID(r.FormValue("historia"))); err != nil {
		a.renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, "/planning-poker", http.StatusSeeOther)
}

// votar registra el voto del integrante en la votación abierta de la historia (HU-16) y vuelve a la
// lista. Una carta que no es un número cae en el 0 y la rechaza el dominio con el error de la escala
// (RN-16.3): no hace falta un mensaje propio del borde.
func (a *app) votar(w http.ResponseWriter, r *http.Request) {
	carta, _ := strconv.Atoi(r.FormValue("carta"))
	err := a.svc.Poker.Votar(
		domain.HistoriaID(r.FormValue("historia")),
		domain.IntegranteID(r.FormValue("integrante")),
		carta,
	)
	if err != nil {
		a.renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, "/planning-poker", http.StatusSeeOther)
}

// nombresDe traduce los IDs de los votantes a sus nombres, en el orden en que votaron (RN-16.4:
// la pantalla muestra quién votó, nunca los valores).
func nombresDe(integrantes []project.Integrante, ids []domain.IntegranteID) []string {
	nombre := make(map[domain.IntegranteID]string, len(integrantes))
	for _, i := range integrantes {
		nombre[i.ID] = i.Nombre
	}
	nombres := make([]string, 0, len(ids))
	for _, id := range ids {
		nombres = append(nombres, nombre[id])
	}
	return nombres
}
