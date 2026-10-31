// Comando checkcoverage hace cumplir el umbral de cobertura de la Definition of Done (TEC-05): la cobertura
// COMBINADA de internal/domain e internal/metrics tiene que ser de al menos 80 %.
// Uso: go run ./tools/checkcoverage -profile=coverage.out
// Si $GITHUB_STEP_SUMMARY está definida (CI), escribe ahí la tabla por paquete y la combinada.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

const umbral = 80.0

// modulo es el prefijo que se saca del nombre de los paquetes en la tabla.
const modulo = "github.com/branko-utn/seguimiento-medicion/"

// medidos son los paquetes que cuentan para el umbral (con sus subpaquetes).
var medidos = []string{"internal/domain", "internal/metrics"}

// paquete es la cobertura de un paquete, en sentencias, como la calcula go tool cover.
type paquete struct {
	nombre             string
	cubiertas, totales int64
}

func main() {
	profile := flag.String("profile", "coverage.out", "archivo de perfil de cobertura")
	flag.Parse()

	f, err := os.Open(*profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error abriendo %s: %v\n", *profile, err)
		os.Exit(1)
	}
	defer f.Close()
	paquetes, combinada, err := calcular(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	tabla := armarTabla(paquetes, combinada)
	fmt.Println(tabla)
	if resumen := os.Getenv("GITHUB_STEP_SUMMARY"); resumen != "" {
		if s, err := os.OpenFile(resumen, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644); err == nil {
			fmt.Fprintln(s, tabla)
			s.Close()
		}
	}
	if combinada < umbral {
		fmt.Fprintf(os.Stderr, "❌ Cobertura combinada de internal/domain e internal/metrics = %.1f %% (< %.0f %%)\n", combinada, umbral)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "✅ Cobertura combinada de internal/domain e internal/metrics = %.1f %% (>= %.0f %%)\n", combinada, umbral)
}

// calcular lee un perfil de cobertura ("archivo:inicio,fin sentencias veces") y devuelve la cobertura de cada
// paquete medido, ordenados por nombre, y la combinada en porcentaje. Pesa por sentencias, como go tool cover.
func calcular(r io.Reader) ([]paquete, float64, error) {
	porNombre := make(map[string]*paquete)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		campos := strings.Fields(sc.Text())
		if len(campos) != 3 || strings.HasPrefix(campos[0], "mode:") {
			continue
		}
		sentencias, err1 := strconv.ParseInt(campos[1], 10, 64)
		veces, err2 := strconv.ParseInt(campos[2], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		archivo, _, _ := strings.Cut(campos[0], ":")
		nombre := strings.TrimPrefix(path.Dir(archivo), modulo)
		if !esMedido(nombre) {
			continue
		}
		p := porNombre[nombre]
		if p == nil {
			p = &paquete{nombre: nombre}
			porNombre[nombre] = p
		}
		p.totales += sentencias
		if veces > 0 {
			p.cubiertas += sentencias
		}
	}
	if err := sc.Err(); err != nil {
		return nil, 0, err
	}
	if len(porNombre) == 0 {
		return nil, 0, errors.New("el perfil no tiene paquetes de internal/domain ni de internal/metrics")
	}

	paquetes := make([]paquete, 0, len(porNombre))
	var cubiertas, totales int64
	for _, p := range porNombre {
		paquetes = append(paquetes, *p)
		cubiertas += p.cubiertas
		totales += p.totales
	}
	sort.Slice(paquetes, func(i, j int) bool { return paquetes[i].nombre < paquetes[j].nombre })
	return paquetes, porcentaje(cubiertas, totales), nil
}

// esMedido indica si el paquete es uno de los medidos o un subpaquete suyo.
func esMedido(nombre string) bool {
	for _, m := range medidos {
		if nombre == m || strings.HasPrefix(nombre, m+"/") {
			return true
		}
	}
	return false
}

func porcentaje(cubiertas, totales int64) float64 {
	if totales == 0 {
		return 0
	}
	return float64(cubiertas) * 100 / float64(totales)
}

// armarTabla arma la tabla markdown del resumen del CI.
func armarTabla(paquetes []paquete, combinada float64) string {
	var b strings.Builder
	b.WriteString("### Cobertura de internal/domain e internal/metrics (umbral 80 %)\n\n")
	b.WriteString("| Paquete | Sentencias cubiertas | Sentencias | Cobertura |\n|---|---|---|---|\n")
	var cubiertas, totales int64
	for _, p := range paquetes {
		fmt.Fprintf(&b, "| %s | %d | %d | %.1f %% |\n", p.nombre, p.cubiertas, p.totales, porcentaje(p.cubiertas, p.totales))
		cubiertas += p.cubiertas
		totales += p.totales
	}
	fmt.Fprintf(&b, "| **Combinada** | **%d** | **%d** | **%.1f %%** |\n", cubiertas, totales, combinada)
	return b.String()
}
