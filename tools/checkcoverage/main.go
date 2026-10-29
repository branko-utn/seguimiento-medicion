// Check que la cobertura COMBINADA (aggregate) de internal/domain + internal/metrics sea >= 80%.
// Uso: go run ./tools/checkcoverage -profile=coverage.out
// Si $GITHUB_STEP_SUMMARY está seteado (CI), escribe tabla por paquete + combinada ahí.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	profile := flag.String("profile", "coverage.out", "archivo de perfil de cobertura")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Verificando umbral de cobertura combinada (>= 80%%) en internal/domain + internal/metrics...\n")

	checkCoverage(*profile)
}

func checkCoverage(profilePath string) {
	file, err := os.Open(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error abriendo %s: %v\n", profilePath, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Perfil vacío\n")
		os.Exit(1)
	}

	type pkgStats struct{ covered, total int64 }
	stats := make(map[string]*pkgStats)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		loc := parts[0]
		countStr := parts[2]

		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			continue
		}

		filePath := strings.Split(loc, ":")[0]
		// Normalizar separadores de ruta (Windows usa \, Unix usa /)
		filePath = filepath.ToSlash(filePath)
		dir := filepath.ToSlash(filepath.Dir(filePath))

		if !strings.Contains(dir, "/internal/domain/") && !strings.Contains(dir, "/internal/metrics/") {
			continue
		}

		pkgPath := dir
		if stats[pkgPath] == nil {
			stats[pkgPath] = &pkgStats{}
		}
		stats[pkgPath].total++
		if count > 0 {
			stats[pkgPath].covered++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo perfil: %v\n", err)
		os.Exit(1)
	}

	if len(stats) == 0 {
		fmt.Fprintf(os.Stderr, "No se encontraron paquetes internal/domain ni internal/metrics en el perfil\n")
		os.Exit(1)
	}

	// Calcular por paquete y combinada
	var totalCovered, totalTotal int64
	rows := make([]string, 0, len(stats))

	// Ordenar por nombre de paquete para salida determinística
	pkgs := make([]string, 0, len(stats))
	for p := range stats {
		pkgs = append(pkgs, p)
	}
	// Sort simple
	for i := 0; i < len(pkgs)-1; i++ {
		for j := i + 1; j < len(pkgs); j++ {
			if pkgs[i] > pkgs[j] {
				pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
			}
		}
	}

	for _, pkg := range pkgs {
		s := stats[pkg]
		if s.total == 0 {
			continue
		}
		pct := float64(s.covered) * 100.0 / float64(s.total)
		totalCovered += s.covered
		totalTotal += s.total
		short := strings.TrimPrefix(pkg, "github.com/branko-utn/seguimiento-medicion/")
		rows = append(rows, fmt.Sprintf("| %s | %d | %d | %.1f%% |", short, s.covered, s.total, pct))
	}

	combinedPct := float64(totalCovered) * 100.0 / float64(totalTotal)
	combinedRow := fmt.Sprintf("| **Combinada (domain+metrics)** | **%d** | **%d** | **%.1f%%** |", totalCovered, totalTotal, combinedPct)

	// Output tabla markdown
	table := "### Cobertura por paquete (internal/domain + internal/metrics)\n\n"
	table += "| Paquete | Cubiertas | Totales | Cobertura |\n"
	table += "|---------|-----------|---------|-----------|\n"
	table += strings.Join(rows, "\n") + "\n"
	table += combinedRow + "\n"

	// Imprimir a stdout (logs CI) y a $GITHUB_STEP_SUMMARY si existe
	fmt.Println(table)
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		f, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(table + "\n")
		}
	}

	// Verificar umbral COMBINADO
	if combinedPct < 80.0 {
		fmt.Fprintf(os.Stderr, "\n❌ Cobertura combinada internal/domain + internal/metrics = %.1f%% (< 80%%)\n", combinedPct)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\n✅ Cobertura combinada internal/domain + internal/metrics = %.1f%% (>= 80%%)\n", combinedPct)
}
