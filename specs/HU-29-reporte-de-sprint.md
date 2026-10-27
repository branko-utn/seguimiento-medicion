# HU-29 · Reporte de sprint

- **Issue:** #36 · **Épica:** E9 Reportes · **Responsable:** Branko Almeira · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-29-reporte-de-sprint.feature`](../features/hu-29-reporte-de-sprint.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-12 y HU-13 (sprint y SP), HU-19 y HU-20 (horas), HU-22 y HU-23 (defectos), HU-24 a HU-27 (métricas)
- **Estado:** borrador para la planning del Sprint 3

## Historia

**Como** líder del proyecto **quiero** generar el reporte de un sprint **para** presentarlo en la Review.

## Objetivo

Armar en una página imprimible todo lo que se muestra en la Sprint Review: qué se planificó, qué se completó, cuánto se
estimó, cuánto se trabajó, las métricas y los defectos del sprint. Se guarda como PDF desde el navegador.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Sprint | `domain.SprintID` | Sí | Debe existir y estar activo o cerrado |

## Salidas esperadas

Un reporte con:

1. **Datos del sprint:** Sprint Goal, fechas y estado. Si está activo, se marca como **parcial**.
2. **Historias planificadas:** por cada una, título, Story Points, horas estimadas, horas reales en el sprint y si se
   completó en el sprint.
3. **Métricas:** SP planificados y completados, porcentaje de historias completadas (HU-25 por sprint), horas estimadas,
   reales y desviación (HU-26).
4. **Esfuerzo registrado:** los registros del sprint y el total por integrante.
5. **Defectos:** los detectados en el sprint (descripción, severidad, estado) y la cantidad de detectados y de resueltos
   en el sprint (HU-27).

## Reglas de negocio

- **RN-29.1:** solo hay reporte de sprints activos o cerrados; uno planificado todavía no tiene resultado.
- **RN-29.2:** en un sprint cerrado, los SP son los fijados al cerrar (RN-14.3). En uno activo se calculan en el momento
  con `metrics.SPPlanificados` y `metrics.SPCompletados`.
- **RN-29.3:** horas del sprint: las estimadas son la suma de las horas estimadas de las historias planificadas; las
  reales, la suma de los registros de esas historias con fecha entre el inicio y el fin del sprint (la misma regla que
  RN-28.3). Así, una historia que pasó por dos sprints no suma sus horas dos veces.
- **RN-29.4:** la desviación sin horas estimadas se muestra como "no calculable" (RN-26.3), sin error en el reporte.
- **RN-29.5:** las historias se listan en el orden en que se asignaron al sprint (`Sprint.Historias`).

## Restricciones

- `report.ReporteSprint` vive en `internal/report` (ADR 0002): el paquete arma el reporte a partir de los datos que
  recibe y usa `internal/metrics`; no accede a repositorios.
- Un servicio nuevo, `service.Reportes`, con `Sprint(sprintID) (report.ReporteSprint, error)`, busca el sprint, sus
  historias, horas, registros y defectos, y se los pasa a `internal/report`.
- Pantalla: `GET /reportes/sprints/{sprint}`, desde la sección **Reportes**. Una hoja de estilos `@media print` oculta la
  navegación y el botón "Imprimir o guardar como PDF" (`window.print()`). Sin bibliotecas de PDF: solo biblioteca
  estándar (ADR 0001). El PDF lo genera el navegador al imprimir.

## Casos límite

- Sprint sin historias: listas vacías, 0 SP, 0 % y desviación "no calculable".
- Sprint sin defectos: "No se registraron defectos en este sprint".
- Registros de una historia fuera de las fechas del sprint: no se suman (RN-29.3).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" (404) |
| El sprint está planificado | `service.ErrSprintSinEmpezar`: "el sprint todavía no empezó" (400) |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-29.1 | Incluye historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas y defectos | `Escenario: reporte de un sprint cerrado`, `Escenario: las horas fuera del sprint no se suman`, `Escenario: reporte parcial de un sprint activo` |
