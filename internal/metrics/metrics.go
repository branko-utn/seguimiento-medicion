package metrics

// HistoriaSprint es una historia del sprint vista desde las métricas: solo lo que hace falta para calcular los
// Story Points planificados y completados (ADR 0002, sección 3).
type HistoriaSprint struct {
	StoryPoints int
	Completada  bool
}

// SPPlanificados suma los Story Points de las historias asignadas al sprint (RN-24.1).
func SPPlanificados(hs []HistoriaSprint) int {
	return 0
}

// SPCompletados suma los Story Points de las historias completadas del sprint (RN-24.2).
func SPCompletados(hs []HistoriaSprint) int {
	return 0
}

// Velocidad devuelve el promedio de Story Points completados de los sprints cerrados (RN-24.3).
func Velocidad(completados []int) float64 {
	return 0
}
