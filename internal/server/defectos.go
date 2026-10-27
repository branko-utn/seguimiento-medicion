package server

import (
	"net/http"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
)

// rutasDefectos registra la pantalla de defectos (épica E6). HU-23: lista de defectos con las acciones para cambiar el
// estado y resolverlos. El formulario para registrar defectos (HU-22) lo agrega TEC-15.
func (a *app) rutasDefectos(mux *http.ServeMux) {
	mux.HandleFunc("GET /defectos", a.listarDefectos)
	mux.HandleFunc("POST /defectos/{defecto}/estado", a.cambiarEstadoDefecto)
	mux.HandleFunc("POST /defectos/{defecto}/resolver", a.resolverDefecto)
}

// filaDefecto es un defecto de la lista con lo que la pantalla puede ofrecer para él.
type filaDefecto struct {
	defect.Defecto
	// Siguiente es el estado al que pasa con CambiarEstado (Abierto ↔ En progreso); el cero si está resuelto.
	Siguiente defect.Estado
	// Sprints son los sprints donde se puede resolver: del mismo proyecto y que no empiezan antes que el de detección
	// (RN-23.3 y RN-23.4). El servicio lo vuelve a validar.
	Sprints []domain.SprintID
}

func (a *app) listarDefectos(w http.ResponseWriter, r *http.Request) {
	var filas []filaDefecto
	for _, d := range a.svc.Defectos.Listar() {
		fila := filaDefecto{Defecto: d}
		switch d.Estado {
		case defect.Abierto:
			fila.Siguiente = defect.EnProgreso
		case defect.EnProgreso:
			fila.Siguiente = defect.Abierto
		}
		if d.Estado != defect.Resuelto {
			sprints, err := a.sprintsParaResolver(d)
			if err != nil {
				a.renderError(w, r, http.StatusInternalServerError, "No se pudieron cargar los sprints del defecto "+string(d.ID)+".")
				return
			}
			fila.Sprints = sprints
		}
		filas = append(filas, fila)
	}
	a.render(w, r, "defectos.html", filas)
}

// sprintsParaResolver devuelve los sprints del proyecto del defecto que no empiezan antes que el de detección.
func (a *app) sprintsParaResolver(d defect.Defecto) ([]domain.SprintID, error) {
	deteccion, err := a.svc.Sprints.Buscar(d.SprintDeteccion)
	if err != nil {
		return nil, err
	}
	sprints, err := a.svc.Sprints.DelProyecto(deteccion.ProyectoID)
	if err != nil {
		return nil, err
	}
	var ids []domain.SprintID
	for _, s := range sprints {
		if !s.Inicio.Before(deteccion.Inicio) {
			ids = append(ids, s.ID)
		}
	}
	return ids, nil
}

func (a *app) cambiarEstadoDefecto(w http.ResponseWriter, r *http.Request) {
	id := domain.DefectoID(r.PathValue("defecto"))
	err := a.svc.Defectos.CambiarEstado(id, estadoDefectoDesdeTexto(r.FormValue("estado")))
	a.responderAccion(w, r, err, "/defectos")
}

func (a *app) resolverDefecto(w http.ResponseWriter, r *http.Request) {
	id := domain.DefectoID(r.PathValue("defecto"))
	err := a.svc.Defectos.Resolver(id, domain.SprintID(r.FormValue("sprint")))
	a.responderAccion(w, r, err, "/defectos")
}

// estadoDefectoDesdeTexto convierte el texto del formulario en un Estado (ADR 0003: la conversión va en el borde).
// Si el texto no es un estado, devuelve el cero, que el dominio rechaza.
func estadoDefectoDesdeTexto(texto string) defect.Estado {
	for _, e := range []defect.Estado{defect.Abierto, defect.EnProgreso, defect.Resuelto} {
		if e.String() == texto {
			return e
		}
	}
	return 0
}
