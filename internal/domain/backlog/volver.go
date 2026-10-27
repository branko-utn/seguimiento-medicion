package backlog

// VolverAlBacklog devuelve al Product Backlog una historia que no se completó en su sprint: queda Pendiente y sin
// sprint asignado (HU-13, RN-13.3). Así se puede planificar en otro sprint. El servicio solo la usa con historias no
// completadas.
func (h *Historia) VolverAlBacklog() {
	h.Estado = Pendiente
	h.SprintID = ""
}
