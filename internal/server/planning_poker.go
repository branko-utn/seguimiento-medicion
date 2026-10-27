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
// con sus votaciones, abrir una votación, votar y revelar (HU-16 y HU-17).
func (a *app) rutasPlanningPoker(mux *http.ServeMux) {
	mux.HandleFunc("GET /planning-poker", a.listarPlanningPoker)
	mux.HandleFunc("POST /planning-poker/abrir", a.abrirVotacion)
	mux.HandleFunc("POST /planning-poker/votar", a.votar)
	mux.HandleFunc("POST /planning-poker/revelar", a.revelar)
}

// filaPoker es lo que se muestra de cada historia: si tiene la votación abierta o revelada, quiénes
// votaron (solo los nombres: RN-16.4), los integrantes habilitados para votar (RN-16.1) y los Agile
// Enablers habilitados para revelar (regla de interfaz, ADR 0004).
type filaPoker struct {
	Historia    backlog.Historia
	Abierta     bool
	Revelada    bool
	Completada  bool
	ConVotos    bool
	Votaron     string
	Integrantes []project.Integrante
	Enablers    []project.Integrante
	Resultado   *vistaResultado // nil mientras la ronda no está revelada
}

// vistaResultado es el resultado de la ronda con nombres, listo para la plantilla (HU-17).
type vistaResultado struct {
	Votos         []vistaVoto
	Consenso      bool
	Maximo        int
	Minimo        int
	VotaronMaximo string
	VotaronMinimo string
}

// vistaVoto es un voto revelado: el integrante por su nombre y su carta.
type vistaVoto struct {
	Integrante string
	Carta      int
}

// datosPoker es lo que recibe la plantilla de la sección.
type datosPoker struct {
	Filas  []filaPoker
	Escala []int
}

// listarPlanningPoker muestra el backlog con el estado de la votación de cada historia (HU-16 y HU-17).
// Los valores de los votos llegan a la plantilla solo cuando la ronda está revelada (RN-16.4).
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
		fila := filaPoker{
			Historia:    h,
			Completada:  h.Estado == backlog.Completada,
			Integrantes: integrantes,
			Enablers:    enablers(integrantes),
		}

		votantes, err := a.svc.Poker.QuienesVotaron(h.ID)
		switch {
		case errors.Is(err, estimation.ErrSinVotacionAbierta):
			// Sin votación abierta: la fila ofrece abrirla (si la historia no está completada, RN-16.5).
		case err != nil:
			a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar la votación.")
			return
		default:
			fila.Abierta = true
			fila.ConVotos = len(votantes) > 0
			if nombres := nombresDe(integrantes, votantes); len(nombres) > 0 {
				fila.Votaron = strings.Join(nombres, ", ")
			} else {
				fila.Votaron = "todavía nadie"
			}
			// Con la ronda revelada, el resultado queda a la vista y la votación se cierra (RN-17.4).
			res, err := a.svc.Poker.Resultado(h.ID)
			switch {
			case err == nil:
				fila.Abierta, fila.Revelada = false, true
				fila.Resultado = nuevaVistaResultado(integrantes, res)
			case errors.Is(err, estimation.ErrVotosOcultos):
				// Abierta: los valores siguen ocultos (RN-16.4).
			default:
				a.renderError(w, r, http.StatusInternalServerError, "No se pudo cargar el resultado.")
				return
			}
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

// revelar cierra la ronda mostrando todos los votos y sus diferencias (HU-17). La regla de quién revela
// es de interfaz (ADR 0004): solo aparece en la pantalla si hay Agile Enablers, y el pedido se valida acá
// — el caso de uso no distingue quién revela.
func (a *app) revelar(w http.ResponseWriter, r *http.Request) {
	historiaID := domain.HistoriaID(r.FormValue("historia"))
	h, err := a.svc.Backlog.Buscar(historiaID)
	if err != nil {
		a.renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	integrantes, err := a.svc.Integrantes.DelProyecto(h.ProyectoID)
	if err != nil {
		a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los integrantes.")
		return
	}
	if !esAgileEnabler(integrantes, domain.IntegranteID(r.FormValue("revelo"))) {
		a.renderError(w, r, http.StatusBadRequest, "solo un Agile Enabler puede revelar los votos")
		return
	}
	if _, err := a.svc.Poker.Revelar(historiaID); err != nil {
		a.renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, "/planning-poker", http.StatusSeeOther)
}

// enablers filtra los Agile Enablers del proyecto: solo ellos revelan (regla de interfaz, ADR 0004).
func enablers(integrantes []project.Integrante) []project.Integrante {
	var aes []project.Integrante
	for _, i := range integrantes {
		if i.Rol == project.AgileEnabler {
			aes = append(aes, i)
		}
	}
	return aes
}

// esAgileEnabler informa si el integrante elegido es uno de los Agile Enablers del proyecto.
func esAgileEnabler(integrantes []project.Integrante, id domain.IntegranteID) bool {
	for _, i := range enablers(integrantes) {
		if i.ID == id {
			return true
		}
	}
	return false
}

// nuevaVistaResultado traduce los IDs del resultado a nombres, listo para la plantilla (HU-17).
func nuevaVistaResultado(integrantes []project.Integrante, res estimation.Resultado) *vistaResultado {
	v := &vistaResultado{Consenso: res.Consenso, Maximo: res.Maximo, Minimo: res.Minimo}
	for _, voto := range res.Votos {
		v.Votos = append(v.Votos, vistaVoto{
			Integrante: strings.Join(nombresDe(integrantes, []domain.IntegranteID{voto.IntegranteID}), ""),
			Carta:      voto.Carta,
		})
	}
	v.VotaronMaximo = strings.Join(nombresDe(integrantes, res.VotaronMaximo), ", ")
	v.VotaronMinimo = strings.Join(nombresDe(integrantes, res.VotaronMinimo), ", ")
	return v
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
