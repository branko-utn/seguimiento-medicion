package backlog

import (
	"sort"
	"strconv"
	"strings"
)

// Ordenar devuelve una copia de las historias en el orden en que conviene trabajarlas (spec HU-09, RN-09.1): por
// prioridad (Alta, Media, Baja) y, a igual prioridad, por el número del identificador (HU-2 antes que HU-10). No
// modifica la entrada.
func Ordenar(hs []Historia) []Historia {
	copia := append([]Historia{}, hs...)
	sort.SliceStable(copia, func(i, j int) bool {
		if copia[i].Prioridad != copia[j].Prioridad {
			return copia[i].Prioridad < copia[j].Prioridad // PrioridadAlta es la de menor valor.
		}
		return numeroDeHistoria(copia[i]) < numeroDeHistoria(copia[j])
	})
	return copia
}

// numeroDeHistoria devuelve el número del identificador ("HU-10" → 10), para compararlo como número y no como texto.
func numeroDeHistoria(h Historia) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(string(h.ID), "HU-"))
	return n
}
