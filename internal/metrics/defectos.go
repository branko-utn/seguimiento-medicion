package metrics

// DefectoSprint es un defecto visto desde las métricas (HU-27).
type DefectoSprint struct {
	SprintDeteccion  string
	SprintResolucion string // "" mientras no esté resuelto
}

func DefectosDetectados(ds []DefectoSprint, sprint string) int {
	return -1
}

func DefectosResueltos(ds []DefectoSprint, sprint string) int {
	return -1
}

func TotalDetectados(ds []DefectoSprint) int {
	return -1
}

func TotalResueltos(ds []DefectoSprint) int {
	return -1
}
