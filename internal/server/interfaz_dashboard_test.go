package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// HU-32.4: el dashboard usa los componentes de TEC-21: resumen en tarjetas, estado con insignia, los cuatro gráficos
// en una grilla de dos columnas y coma decimal.

type entornoDashboard struct {
	router        http.Handler
	sprints       *service.Sprints
	repoProyectos *memory.Proyectos
	repoHistorias *memory.Historias
	repoEstim     *memory.Estimaciones
	repoRegistros *memory.Registros
}

// nuevoEntornoDashboard arma el servidor sin datos. Hoy es el 27/10/2026, dentro de los proyectos de prueba, para que
// el estado no dependa de la fecha real.
func nuevoEntornoDashboard(t *testing.T) entornoDashboard {
	t.Helper()
	original := hoy
	hoy = func() time.Time { return time.Date(2026, 10, 27, 0, 0, 0, 0, time.Local) }
	t.Cleanup(func() { hoy = original })

	repoProyectos, _ := memory.NuevoProyectos("")
	repoSprints, _ := memory.NuevoSprints("")
	repoHistorias, _ := memory.NuevoHistorias("")
	repoEstimaciones, _ := memory.NuevoEstimaciones("")
	repoRegistros, _ := memory.NuevoRegistros("")
	repoDefectos, _ := memory.NuevoDefectos("")
	sprintsSvc := service.NuevoSprints(repoSprints, repoProyectos, repoHistorias)
	svc := Servicios{
		Proyectos: service.NuevoProyectos(repoProyectos),
		Sprints:   sprintsSvc,
		Dashboard: service.NuevoDashboard(repoProyectos, repoSprints, repoHistorias, repoEstimaciones, repoRegistros, repoDefectos),
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}
	return entornoDashboard{router: router, sprints: sprintsSvc, repoProyectos: repoProyectos, repoHistorias: repoHistorias, repoEstim: repoEstimaciones, repoRegistros: repoRegistros}
}

// crearDemo agrega el proyecto Demo (21/10 al 02/11/2026) y devuelve su ID.
func (e entornoDashboard) crearDemo(t *testing.T) domain.ProyectoID {
	t.Helper()
	p, err := e.repoProyectos.Agregar(project.Proyecto{Nombre: "Demo", Inicio: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), Fin: time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local)})
	if err != nil {
		t.Fatal(err)
	}
	return p.ID
}

// cargarSprintCerrado deja a Demo con un sprint cerrado: HU-1 (3 SP, 4 h estimadas, 5 h reales) completada, HU-2
// (5 SP, 6 h estimadas) sin completar y HU-3 fuera del sprint. Velocidad 3, 33,3 % de historias completadas (1 de 3),
// 10 h estimadas, 5 h reales, desviación -50 %.
func (e entornoDashboard) cargarSprintCerrado(t *testing.T, pid domain.ProyectoID) {
	t.Helper()
	h1, _ := e.repoHistorias.Agregar(backlog.Historia{ProyectoID: pid, Titulo: "Crear proyecto", StoryPoints: 3, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	h2, _ := e.repoHistorias.Agregar(backlog.Historia{ProyectoID: pid, Titulo: "Crear sprint", StoryPoints: 5, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}})
	if _, err := e.repoHistorias.Agregar(backlog.Historia{ProyectoID: pid, Titulo: "Registrar defecto", StoryPoints: 2, Estado: backlog.Pendiente, Criterios: []backlog.Criterio{{Texto: "ok"}}}); err != nil {
		t.Fatal(err)
	}
	for _, est := range []effort.EstimacionHoras{{HistoriaID: h1.ID, Horas: 4}, {HistoriaID: h2.ID, Horas: 6}} {
		if err := e.repoEstim.Guardar(est); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := e.sprints.Crear(pid, "S1", time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), time.Date(2026, 10, 23, 0, 0, 0, 0, time.Local)); err != nil {
		t.Fatal(err)
	}
	pasos := []func() error{
		func() error { return e.sprints.Iniciar("S1") },
		func() error { return e.sprints.Asignar("S1", h1.ID) },
		func() error { return e.sprints.Asignar("S1", h2.ID) },
		func() error {
			_, err := e.repoRegistros.Agregar(effort.Registro{HistoriaID: h1.ID, Fecha: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), Horas: 5})
			return err
		},
		func() error { return e.sprints.CompletarHistoria("S1", h1.ID) },
		func() error { return e.sprints.CerrarSprint("S1") },
	}
	for _, paso := range pasos {
		if err := paso(); err != nil {
			t.Fatal(err)
		}
	}
}

func (e entornoDashboard) pedir(t *testing.T, ruta string) string {
	t.Helper()
	rec := servir(e.router, httptest.NewRequest(http.MethodGet, ruta, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: código = %d, se esperaba 200", ruta, rec.Code)
	}
	return rec.Body.String()
}

// iconoDashboard devuelve el SVG de la sección Dashboard tal como lo define secciones.go (TEC-03).
func iconoDashboard(t *testing.T) string {
	t.Helper()
	for _, s := range secciones {
		if s.Ruta == "/dashboard" {
			return `<svg class="icono" viewBox="0 0 24 24" aria-hidden="true"><path d="` + s.Icono + `"/></svg>`
		}
	}
	t.Fatal("no está la sección /dashboard")
	return ""
}

func TestHU32_4_resumenEnTarjetasConComaDecimal(t *testing.T) {
	e := nuevoEntornoDashboard(t)
	pid := e.crearDemo(t)
	e.cargarSprintCerrado(t, pid)

	cuerpo := e.pedir(t, "/dashboard?proyecto="+string(pid))

	contieneTodo(t, cuerpo, `<div class="tarjetas">`,
		`<span class="valor">3,0</span><span class="etiqueta">velocidad del proyecto</span>`,
		`<span class="valor">33,3 %</span><span class="etiqueta">historias completadas</span>`,
		`<span class="valor">10,0</span><span class="etiqueta">horas estimadas</span>`,
		`<span class="valor">5,0</span><span class="etiqueta">horas reales</span>`,
		`<span class="valor">-50,0 %</span><span class="etiqueta">desviación</span>`,
		`<span class="valor">0</span><span class="etiqueta">defectos abiertos</span>`,
	)
	if n := strings.Count(cuerpo, `class="tarjeta-metrica"`); n < 8 {
		t.Errorf("hay %d tarjetas, se esperaban al menos 8 (una por métrica del resumen)", n)
	}
	// RN-32.3: ningún número con punto decimal.
	for _, conPunto := range []string{"3.0", "33.3", "10.0", "5.0", "-50.0"} {
		if strings.Contains(cuerpo, conPunto) {
			t.Errorf("la página muestra %q con punto decimal (RN-32.3)", conPunto)
		}
	}
}

func TestHU32_4_sinSprintsCerradosMuestraCerosYNoCalculable(t *testing.T) {
	e := nuevoEntornoDashboard(t)
	pid := e.crearDemo(t)

	cuerpo := e.pedir(t, "/dashboard?proyecto="+string(pid))

	// Caso límite: sin datos, los valores son 0,0 con coma, sin sprint activo y la desviación no se puede calcular.
	contieneTodo(t, cuerpo,
		`<span class="valor">0,0</span><span class="etiqueta">velocidad del proyecto</span>`,
		`<span class="valor">0,0 %</span><span class="etiqueta">historias completadas</span>`,
		`<span class="valor">Ninguno</span><span class="etiqueta">sprint activo</span>`,
		`<span class="valor">no calculable</span><span class="etiqueta">desviación</span>`,
		"Todavía no hay sprints cerrados para graficar",
	)
}

func TestHU32_4_estadoDelProyectoConInsignia(t *testing.T) {
	e := nuevoEntornoDashboard(t)
	pid := e.crearDemo(t)

	cuerpo := e.pedir(t, "/dashboard?proyecto="+string(pid))

	// El tono lo decide el handler con el enum (RN-32.1): el 27/10 Demo está en curso.
	contieneTodo(t, cuerpo, `<span class="insignia tono-info">En curso</span>`)
}

func TestHU32_4_graficosEnGrillaDeDosColumnas(t *testing.T) {
	e := nuevoEntornoDashboard(t)
	pid := e.crearDemo(t)
	e.cargarSprintCerrado(t, pid)

	cuerpo := e.pedir(t, "/dashboard?proyecto="+string(pid))

	inicio := strings.Index(cuerpo, `<div class="grilla-graficos">`)
	if inicio < 0 {
		t.Fatal("falta la grilla de gráficos")
	}
	grilla := cuerpo[inicio:]
	if fin := strings.Index(grilla, `<script`); fin >= 0 {
		grilla = grilla[:fin]
	}
	for _, id := range []string{"grafico-velocidad", "grafico-sp", "grafico-horas", "grafico-defectos"} {
		if !strings.Contains(grilla, `id="`+id+`"`) {
			t.Errorf("el gráfico %s no está dentro de la grilla", id)
		}
	}
	if n := strings.Count(grilla, `<article class="grafico">`); n != 4 {
		t.Errorf("la grilla tiene %d gráficos, se esperaban 4", n)
	}
}

func TestHU32_4_pantallasConEncabezadoIconoYSinEstilosEnLinea(t *testing.T) {
	e := nuevoEntornoDashboard(t)
	pid := e.crearDemo(t)
	e.cargarSprintCerrado(t, pid)

	for _, ruta := range []string{"/dashboard", "/dashboard?proyecto=" + string(pid)} {
		t.Run(ruta, func(t *testing.T) {
			cuerpo := e.pedir(t, ruta)
			contieneTodo(t, cuerpo, `class="encabezado"`, iconoDashboard(t))
			if strings.Contains(cuerpo, `style="`) {
				t.Error("la página tiene estilos en línea (RN-32.5)")
			}
			for _, emoji := range emojis {
				if strings.Contains(cuerpo, emoji) {
					t.Errorf("la página usa %q como ícono (RN-32.5)", emoji)
				}
			}
		})
	}
}

func TestHU32_4_selectorSinProyectosMuestraVacioConAccion(t *testing.T) {
	e := nuevoEntornoDashboard(t)

	cuerpo := e.pedir(t, "/dashboard")

	// RN-32.4: sin proyectos no hay un selector vacío sino el estado vacío con el enlace para crear uno.
	contieneTodo(t, cuerpo, `class="vacio"`, "Todavía no hay proyectos", `href="/proyectos/nuevo"`)
	if strings.Contains(cuerpo, "<select") {
		t.Error("sin proyectos no debe mostrarse el selector")
	}
}

func TestHU32_4_proyectoInexistenteSigueDando404(t *testing.T) {
	e := nuevoEntornoDashboard(t)

	rec := servir(e.router, httptest.NewRequest(http.MethodGet, "/dashboard?proyecto=P9", nil))

	// RN-32.6: el comportamiento ante errores no cambia.
	if rec.Code != http.StatusNotFound {
		t.Fatalf("código = %d, se esperaba 404", rec.Code)
	}
	contieneTodo(t, rec.Body.String(), "el proyecto no existe")
}
