package estimation

import (
	"errors"
	"reflect"
	"testing"
)

func TestEscala(t *testing.T) {
	// La escala de Planning Poker de la spec HU-15 (RN-15.1).
	quiere := []int{1, 2, 3, 5, 8, 13, 21}
	if !reflect.DeepEqual(Escala, quiere) {
		t.Errorf("Escala = %v, se esperaba %v", Escala, quiere)
	}
}

func TestValidarStoryPoints(t *testing.T) {
	// Todos los valores de la escala son válidos; los casos límite son 1 y 21.
	for _, sp := range Escala {
		if err := ValidarStoryPoints(sp); err != nil {
			t.Errorf("ValidarStoryPoints(%d) = %v, se esperaba nil", sp, err)
		}
	}

	casos := []struct {
		nombre string
		sp     int
	}{
		{"cero es sin estimar (RN-15.3)", 0},
		{"fuera de la escala", 4},
		{"negativo", -1},
		{"muy grande", 34},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if err := ValidarStoryPoints(caso.sp); !errors.Is(err, ErrStoryPointsInvalidos) {
				t.Fatalf("se esperaba ErrStoryPointsInvalidos, se obtuvo %v", err)
			}
		})
	}
}
