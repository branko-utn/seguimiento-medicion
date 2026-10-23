package metrics

// HistoriaSprint es una historia del sprint vista desde las métricas: solo lo que hace falta para calcular los
// Story Points planificados y completados (ADR 0002, sección 3).
type HistoriaSprint struct {
	StoryPoints int
	Completada  bool
}

// SPPlanificados suma los Story Points de las historias asignadas al sprint (RN-24.1).
func SPPlanificados(hs []HistoriaSprint) int {
	total := 0
	for _, h := range hs {
		total += h.StoryPoints
	}
	return total
}

// SPCompletados suma los Story Points de las historias completadas del sprint (RN-24.2).
func SPCompletados(hs []HistoriaSprint) int {
	total := 0
	for _, h := range hs {
		if h.Completada {
			total += h.StoryPoints
		}
	}
	return total
}

// Velocidad devuelve el promedio de Story Points completados de los sprints cerrados (RN-24.3).
// Sin sprints cerrados devuelve 0 (RN-24.4).
func Velocidad(completados []int) float64 {
	if len(completados) == 0 {
		return 0
	}
	total := 0
	for _, sp := range completados {
		total += sp
	}
	return float64(total) / float64(len(completados))
}
