package metrics

import "testing"

// defectosEjemplo son los defectos de los escenarios de HU-27: D1 detectado y resuelto en S1, D2 detectado en S1 y
// resuelto en S2, D3 detectado en S2 y sin resolver.
var defectosEjemplo = []DefectoSprint{
	{SprintDeteccion: "S1", SprintResolucion: "S1"},
	{SprintDeteccion: "S1", SprintResolucion: "S2"},
	{SprintDeteccion: "S2", SprintResolucion: ""},
}

func TestDefectosPorSprint(t *testing.T) {
	casos := []struct {
		nombre                      string
		defectos                    []DefectoSprint
		sprint                      string
		quiereDetectados, quiereRes int
	}{
		{"S1: detecta D1 y D2, resuelve D1", defectosEjemplo, "S1", 2, 1},
		{"S2: detecta D3, resuelve D2 (detectado en S1)", defectosEjemplo, "S2", 1, 1},
		{"sprint sin defectos", defectosEjemplo, "S3", 0, 0},
		{"sin defectos", nil, "S1", 0, 0},
		{"un defecto sin resolver no cuenta como resuelto en el sprint vacío", defectosEjemplo, "", 0, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := DefectosDetectados(c.defectos, c.sprint); got != c.quiereDetectados {
				t.Errorf("DefectosDetectados(%q) = %d, se esperaba %d", c.sprint, got, c.quiereDetectados)
			}
			if got := DefectosResueltos(c.defectos, c.sprint); got != c.quiereRes {
				t.Errorf("DefectosResueltos(%q) = %d, se esperaba %d", c.sprint, got, c.quiereRes)
			}
		})
	}
}

func TestDefectosTotales(t *testing.T) {
	casos := []struct {
		nombre                      string
		defectos                    []DefectoSprint
		quiereDetectados, quiereRes int
	}{
		{"tres defectos, dos resueltos", defectosEjemplo, 3, 2},
		{"sin defectos", nil, 0, 0},
		{"ninguno resuelto", []DefectoSprint{{SprintDeteccion: "S1"}, {SprintDeteccion: "S2"}}, 2, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := TotalDetectados(c.defectos); got != c.quiereDetectados {
				t.Errorf("TotalDetectados = %d, se esperaba %d", got, c.quiereDetectados)
			}
			if got := TotalResueltos(c.defectos); got != c.quiereRes {
				t.Errorf("TotalResueltos = %d, se esperaba %d", got, c.quiereRes)
			}
		})
	}
}
