package metrics

import (
	"math"
	"testing"
)

func TestSPPlanificados(t *testing.T) {
	casos := []struct {
		nombre    string
		historias []HistoriaSprint
		quiere    int
	}{
		{"lista nula", nil, 0},
		{"sprint sin historias", []HistoriaSprint{}, 0},
		{"suma todas las historias asignadas", []HistoriaSprint{
			{StoryPoints: 3, Completada: true},
			{StoryPoints: 5, Completada: false},
			{StoryPoints: 2, Completada: true},
		}, 10},
		{"una historia sin estimar", []HistoriaSprint{
			{StoryPoints: 0, Completada: false},
		}, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := SPPlanificados(c.historias); got != c.quiere {
				t.Errorf("SPPlanificados(%v) = %d, se esperaba %d", c.historias, got, c.quiere)
			}
		})
	}
}

func TestSPCompletados(t *testing.T) {
	casos := []struct {
		nombre    string
		historias []HistoriaSprint
		quiere    int
	}{
		{"lista nula", nil, 0},
		{"ninguna completada", []HistoriaSprint{
			{StoryPoints: 3, Completada: false},
			{StoryPoints: 8, Completada: false},
		}, 0},
		{"solo las completadas", []HistoriaSprint{
			{StoryPoints: 3, Completada: true},
			{StoryPoints: 5, Completada: false},
			{StoryPoints: 2, Completada: true},
		}, 5},
		{"todas completadas", []HistoriaSprint{
			{StoryPoints: 1, Completada: true},
			{StoryPoints: 2, Completada: true},
		}, 3},
		{"completada sin estimar", []HistoriaSprint{
			{StoryPoints: 0, Completada: true},
		}, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := SPCompletados(c.historias); got != c.quiere {
				t.Errorf("SPCompletados(%v) = %d, se esperaba %d", c.historias, got, c.quiere)
			}
		})
	}
}

func TestVelocidad(t *testing.T) {
	casos := []struct {
		nombre      string
		completados []int
		quiere      float64
	}{
		{"sin sprints cerrados", nil, 0},
		{"lista vacía", []int{}, 0},
		{"un sprint cerrado con 13 SP", []int{13}, 13},
		{"tres sprints cerrados", []int{10, 13, 8}, 31.0 / 3.0},
		{"un sprint cerrado con 0 SP cuenta en el promedio", []int{10, 0}, 5},
		{"todos los sprints con 0 SP", []int{0, 0}, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := Velocidad(c.completados)
			if math.Abs(got-c.quiere) > 1e-9 {
				t.Errorf("Velocidad(%v) = %v, se esperaba %v", c.completados, got, c.quiere)
			}
		})
	}
}
