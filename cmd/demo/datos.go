package main

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// Datos de nuestro propio proyecto, tomados del tablero (GitHub Projects #1), las specs, las actas y las revisiones
// de PRs. Las horas no se registraron durante el proyecto: se estiman a 2 horas por Story Point y las reales varían
// con un factor por historia (ver horasReales). Todo lo demás es tal cual pasó.

func dia(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }

var proyectoDemo = struct {
	nombre, descripcion string
	inicio, fin         time.Time
}{
	nombre:      "Software Metrics & Estimation",
	descripcion: "TPI de Ingeniería y Calidad de Software (UTN FRSR): seguimiento de proyectos Scrum con métricas y estimación.",
	inicio:      dia(21),
	fin:         time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local),
}

type integranteDemo struct {
	nombre string
	rol    project.Rol
}

// integrantesDemo: el índice es la clave que usan las historias (responsable).
var integrantesDemo = map[string]integranteDemo{
	"branko":    {"Branko Almeira", project.AgileEnabler},
	"sosa":      {"Federico Sosa", project.ProductBuilder},
	"bianchi":   {"Ignacio Bianchi", project.ProductBuilder},
	"lara":      {"Agustín Lara", project.ProductBuilder},
	"giorlando": {"Agustín Giorlando", project.ProductBuilder},
}

// ordenIntegrantes fija el orden de alta (los mapas de Go no tienen orden).
var ordenIntegrantes = []string{"branko", "sosa", "bianchi", "lara", "giorlando"}

type sprintDemo struct {
	goal        string
	inicio, fin time.Time
	cerrar      bool
}

// sprintsDemo son los cuatro sprints, con las fechas de las actas (el Sprint 3 se corrió un día en su planning).
var sprintsDemo = []sprintDemo{
	{"Poder crear un proyecto, armar su backlog, planificar un sprint y ver sus primeras métricas.", dia(23), dia(25), true},
	{"Usar la aplicación desde el navegador, de proyecto a métricas, y estimar con Planning Poker.", dia(26), dia(27), true},
	{"Ver el estado del proyecto de un vistazo en el dashboard y generar los reportes de sprint y de proyecto.", dia(29), dia(30), true},
	{"Producto pulido, documentado y listo para presentar.", dia(31), dia(31), false},
}

type historiaDemo struct {
	titulo      string
	sp          int
	responsable string
	sprint      int // 1 a 4; las del sprint 4 quedan sin completar
	criterio    string
	factor      float64 // horas reales / horas estimadas
}

// historiasDemo es el Product Backlog: HU-01 a HU-31, en orden, con los SP y responsables del tablero y el primer
// criterio de aceptación de su spec.
var historiasDemo = []historiaDemo{
	{"Crear proyecto", 3, "sosa", 1, "El nombre es obligatorio y tiene entre 3 y 100 caracteres", 1.2},
	{"Modificar proyecto", 2, "sosa", 2, "Se pueden editar nombre, descripción y fechas, con las mismas validaciones del alta", 1.0},
	{"Registrar integrantes", 3, "sosa", 1, "El nombre es obligatorio y el rol es uno de los tres de la consigna", 0.9},
	{"Consultar estado del proyecto", 3, "sosa", 3, `El estado es "Planificado", "En curso" o "Finalizado"`, 1.3},
	{"Crear ítem del Product Backlog", 3, "bianchi", 1, "El identificador se genera automáticamente y es único", 1.1},
	{"Editar y eliminar ítem del backlog", 3, "lara", 3, "Se pueden editar título, descripción y prioridad", 1.2},
	{"Criterios de aceptación de una historia", 2, "bianchi", 1, "Un criterio no puede estar vacío", 0.8},
	{"Cambiar el estado de una historia", 3, "bianchi", 1, "Los estados son Pendiente → En progreso → Completada", 1.0},
	{"Ver el backlog ordenado por prioridad", 2, "lara", 3, "Se ordena por prioridad (Alta, Media, Baja) y, a igual prioridad, por identificador", 0.9},
	{"Crear sprint con Sprint Goal", 5, "sosa", 1, "El Sprint Goal es obligatorio", 1.1},
	{"Asignar historias a un sprint", 5, "sosa", 1, "Solo se asignan historias estimadas y con criterios de aceptación", 1.2},
	{"Registrar historias completadas", 2, "branko", 1, "Solo se completan historias asignadas al sprint", 1.0},
	{"Cerrar sprint", 5, "branko", 2, "Las historias no completadas vuelven al backlog como Pendientes", 1.1},
	{"Consultar sprints anteriores", 3, "giorlando", 3, "Se listan con fechas, Sprint Goal, SP planificados y completados", 1.4},
	{"Estimar una historia en Story Points", 2, "bianchi", 1, "Solo se aceptan valores de la escala 1, 2, 3, 5, 8, 13, 21", 0.9},
	{"Votar en Planning Poker con votos ocultos", 5, "bianchi", 2, "Solo votan los integrantes del proyecto, una vez por ronda", 1.2},
	{"Revelar votos y detectar diferencias", 3, "bianchi", 2, "Al revelar, se muestran todos los votos de la ronda", 1.0},
	{"Nueva ronda y estimación acordada", 5, "bianchi", 3, "Una ronda nueva empieza sin votos; las anteriores quedan en el historial", 1.5},
	{"Registrar horas estimadas de una historia", 1, "lara", 1, "Las horas estimadas son mayores que cero", 1.0},
	{"Registrar esfuerzo trabajado", 5, "lara", 1, "Se registra integrante, fecha, actividad, horas e historia", 1.1},
	{"Consultar el esfuerzo registrado", 3, "lara", 3, "Se listan los registros con sus totales por historia y por integrante", 1.0},
	{"Registrar defecto", 3, "lara", 1, "Descripción, severidad, historia relacionada y sprint de detección son obligatorios", 0.9},
	{"Cambiar estado de un defecto y registrar su resolución", 3, "lara", 2, "Los estados son Abierto → En progreso → Resuelto", 1.0},
	{"Story Points planificados, completados y velocidad", 5, "giorlando", 1, "SP planificados = suma de los SP de las historias asignadas al sprint", 1.0},
	{"Porcentaje de historias completadas", 2, "giorlando", 1, "Porcentaje = completadas / planificadas × 100, por sprint y por proyecto", 1.2},
	{"Horas estimadas, reales y desviación", 3, "giorlando", 2, "Horas reales = suma del esfuerzo registrado", 1.4},
	{"Defectos detectados y resueltos", 2, "lara", 2, "Se cuentan por sprint (según el sprint de detección y el de resolución) y en total", 0.8},
	{"Dashboard del proyecto", 8, "giorlando", 3, "Muestra estado del proyecto, sprint activo y métricas principales", 1.3},
	{"Reporte de sprint", 5, "branko", 3, "Incluye historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas y defectos", 0.9},
	{"Reporte de proyecto", 3, "branko", 3, "Incluye el resumen de todos los sprints, métricas acumuladas y defectos", 0.8},
	{"Exportar reporte a PDF", 5, "branko", 4, "El reporte de sprint o de proyecto se puede guardar como PDF", 0.2},
}

// horasPorSP es el criterio con el que se cargan las horas estimadas de la demo.
const horasPorSP = 2.0

// prioridadDelSprint: lo de los primeros sprints era lo más prioritario.
func prioridadDelSprint(sprint int) backlog.Prioridad {
	switch sprint {
	case 1, 2:
		return backlog.PrioridadAlta
	case 3:
		return backlog.PrioridadMedia
	}
	return backlog.PrioridadBaja
}

type defectoDemo struct {
	descripcion string
	severidad   defect.Severidad
	historia    int // número de HU
	deteccion   int // sprint 1 a 4
	resolucion  int // 0 = abierto
}

// defectosDemo son defectos reales encontrados en las revisiones de PRs (ver las actas de review).
var defectosDemo = []defectoDemo{
	{"El estado del proyecto era un texto y no un enumerador (ADR 0003)", defect.Media, 1, 1, 1},
	{"Un merge borraba los registros de pasos BDD de otras historias", defect.Alta, 25, 1, 1},
	{"Faltaban los escenarios BDD de las horas estimadas", defect.Media, 19, 1, 1},
	{"Las horas negativas se descartaban en silencio en lugar de dar error", defect.Alta, 26, 2, 2},
	{"La pantalla del sprint no ofrecía completar historias ni cerrar el sprint", defect.Media, 13, 2, 2},
	{"Los escenarios de modificar proyecto iban a fallar desde el 03/11 por usar la fecha real", defect.Alta, 4, 3, 3},
	{"Acordar la estimación no cerraba la votación: la historia no se podía volver a votar", defect.Critica, 18, 3, 3},
	{"El dashboard seguía figurando como \"en construcción\" en la portada", defect.Media, 28, 3, 3},
	{"El dashboard muestra los números con punto decimal en lugar de coma", defect.Baja, 28, 3, 0},
}
