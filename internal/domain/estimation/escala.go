package estimation

import "errors"

// ErrStoryPointsInvalidos: RN-15.1.
var ErrStoryPointsInvalidos = errors.New("los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21")

// Escala es la escala de Planning Poker: los únicos valores de Story Points válidos (spec HU-15, RN-15.1).
// El Planning Poker (HU-16 a HU-18) reutiliza esta misma lista.
var Escala = []int{1, 2, 3, 5, 8, 13, 21}

// ValidarStoryPoints informa si sp pertenece a la escala de Planning Poker (RN-15.1).
// 0 significa "sin estimar" y no es un valor asignable (RN-15.3).
func ValidarStoryPoints(sp int) error {
	return nil // stub RED
}
