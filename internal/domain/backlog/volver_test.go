package backlog

import "testing"

func TestVolverAlBacklog(t *testing.T) {
	for _, estado := range []Estado{Pendiente, EnProgreso} {
		h := Historia{ID: "HU-1", Estado: estado, SprintID: "S1"}
		h.VolverAlBacklog()
		if h.Estado != Pendiente || h.SprintID != "" {
			t.Errorf("desde %s: historia = %+v, se esperaba Pendiente y sin sprint (RN-13.3)", estado, h)
		}
	}
}
