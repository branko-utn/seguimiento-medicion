package main

import (
	"math"
	"strings"
	"testing"
)

// perfil es un coverage.out chico: cada línea es "archivo:inicio,fin sentencias veces".
const perfil = `mode: atomic
github.com/branko-utn/seguimiento-medicion/internal/domain/backlog/historia.go:10.2,12.3 3 1
github.com/branko-utn/seguimiento-medicion/internal/domain/backlog/historia.go:14.2,15.3 1 0
github.com/branko-utn/seguimiento-medicion/internal/metrics/horas.go:14.2,18.3 4 2
github.com/branko-utn/seguimiento-medicion/internal/metrics/horas.go:20.2,21.3 2 0
github.com/branko-utn/seguimiento-medicion/internal/server/sprints.go:30.2,31.3 5 0
`

func TestCalcular_incluyeMetricsYPesaPorSentencias(t *testing.T) {
	paquetes, combinada, err := calcular(strings.NewReader(perfil))
	if err != nil {
		t.Fatal(err)
	}
	// internal/metrics tiene que contar (su directorio no termina en "/"), e internal/server no.
	if len(paquetes) != 2 || paquetes[0].nombre != "internal/domain/backlog" || paquetes[1].nombre != "internal/metrics" {
		t.Fatalf("paquetes = %+v, se esperaban internal/domain/backlog e internal/metrics", paquetes)
	}
	// Por sentencias, como go tool cover: backlog 3 de 4, metrics 4 de 6, combinada 7 de 10.
	if paquetes[0].cubiertas != 3 || paquetes[0].totales != 4 || paquetes[1].cubiertas != 4 || paquetes[1].totales != 6 {
		t.Errorf("sentencias = %+v", paquetes)
	}
	if math.Abs(combinada-70) > 1e-9 {
		t.Errorf("combinada = %v, se esperaba 70", combinada)
	}
}

func TestCalcular_sinPaquetesMedidos(t *testing.T) {
	if _, _, err := calcular(strings.NewReader("mode: atomic\n")); err == nil {
		t.Error("un perfil sin internal/domain ni internal/metrics tiene que dar error")
	}
}
