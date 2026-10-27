package metrics

// DefectoSprint es un defecto visto desde las métricas (HU-27). Los sprints se identifican con su ID como texto.
type DefectoSprint struct {
	SprintDeteccion  string
	SprintResolucion string // "" mientras no esté resuelto
}

// DefectosDetectados cuenta los defectos cuyo sprint de detección es sprint (RN-27.1).
func DefectosDetectados(ds []DefectoSprint, sprint string) int {
	return contar(ds, func(d DefectoSprint) bool { return d.SprintDeteccion == sprint })
}

// DefectosResueltos cuenta los defectos cuyo sprint de resolución es sprint. Un defecto sin sprint de resolución no
// está resuelto en ningún sprint (RN-27.2).
func DefectosResueltos(ds []DefectoSprint, sprint string) int {
	return contar(ds, func(d DefectoSprint) bool { return d.resuelto() && d.SprintResolucion == sprint })
}

// TotalDetectados es la cantidad de defectos del proyecto (RN-27.3).
func TotalDetectados(ds []DefectoSprint) int {
	return len(ds)
}

// TotalResueltos es la cantidad de defectos resueltos del proyecto (RN-27.3).
func TotalResueltos(ds []DefectoSprint) int {
	return contar(ds, DefectoSprint.resuelto)
}

func (d DefectoSprint) resuelto() bool {
	return d.SprintResolucion != ""
}

func contar(ds []DefectoSprint, incluir func(DefectoSprint) bool) int {
	n := 0
	for _, d := range ds {
		if incluir(d) {
			n++
		}
	}
	return n
}
