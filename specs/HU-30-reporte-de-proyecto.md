# HU-30 · Reporte de proyecto

- **Issue:** #37 · **Épica:** E9 Reportes · **Responsable:** Branko Almeira · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-30-reporte-de-proyecto.feature`](../features/hu-30-reporte-de-proyecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el responsable al aprobar el PR de la planning del Sprint 3
- **Depende de:** HU-29 (`service.Reportes` y las filas por sprint; mismo responsable, va después) y código en `main`:
  HU-24 a HU-27 (métricas) y TEC-20 (`Proyecto.EstadoAl` y `Proyecto.DiasRestantes`). **No espera a HU-04 ni a HU-21.**
- **Estado:** aprobada en la planning del Sprint 3

## Historia

**Como** líder del proyecto **quiero** generar el reporte del proyecto completo **para** informar su estado.

## Objetivo

Armar en una página imprimible el estado del proyecto, un resumen de todos sus sprints, las métricas acumuladas y los
defectos. Es el informe que acompaña la presentación final.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Hoy | `time.Time` | Sí | La fecha del sistema (para el estado y los días restantes) |

## Salidas esperadas

1. **Proyecto:** nombre, fechas, estado y días restantes (TEC-20).
2. **Resumen de sprints:** una fila por sprint activo o cerrado, en orden de inicio, con Sprint Goal, estado, SP
   planificados y completados, porcentaje de historias completadas, horas estimadas y reales, y defectos detectados y
   resueltos. Son los mismos valores que el reporte de cada sprint (HU-29).
3. **Métricas acumuladas:** velocidad (HU-24), SP completados en total, porcentaje de historias completadas del proyecto
   (HU-25), horas estimadas, reales y desviación del proyecto (HU-26).
4. **Defectos:** totales detectados, resueltos y abiertos (HU-27), y la lista de los abiertos con su severidad.

## Reglas de negocio

- **RN-30.1:** el resumen incluye los sprints activos y cerrados; los planificados no (no tienen resultado, RN-29.1).
- **RN-30.2:** cada fila del resumen se calcula con las reglas de HU-29 (RN-29.2 y RN-29.3).
- **RN-30.3:** la velocidad usa solo los sprints cerrados (RN-24.3).
- **RN-30.4:** las horas acumuladas son de todo el proyecto: estimadas, la suma de las horas estimadas de todas sus
  historias; reales, la suma de todos los registros de esas historias (`Registros.DeHistoria` y `metrics.HorasReales`),
  también lo registrado fuera de un sprint. Por eso pueden ser más que la suma de las filas.
- **RN-30.5:** los defectos abiertos son los que no están resueltos (Abierto o En progreso), ordenados por severidad
  (Crítica, Alta, Media, Baja).

## Restricciones

- `report.ReporteProyecto` vive en `internal/report` (ADR 0002), junto a `ReporteSprint`; las filas del resumen se arman
  con la misma función que el reporte de sprint.
- `Reportes.Proyecto(proyectoID, hoy) (report.ReporteProyecto, error)` en el mismo servicio que HU-29.
- Pantalla: `GET /reportes/proyecto`, con el proyecto elegido como en las demás secciones, y la misma hoja de impresión que
  HU-29.

## Casos límite

- Proyecto sin sprints: resumen vacío ("Todavía no hay sprints con resultado"), velocidad 0 y desviación "no calculable"
  si no hay horas estimadas.
- Proyecto sin defectos: totales en 0 y "No hay defectos abiertos".

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" (404) |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-30.1 | Incluye el resumen de todos los sprints, métricas acumuladas y defectos | `Escenario: reporte con sprints cerrados y activo`, `Escenario: defectos abiertos por severidad`, `Escenario: proyecto sin sprints` |
