# HU-10 · Crear sprint con Sprint Goal

- **Issue:** #17 · **Épica:** E3 Sprints · **Responsable:** Federico Sosa · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-10-crear-sprint.feature`](../features/hu-10-crear-sprint.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Federico Sosa (al aprobar el PR)
- **Depende de:** HU-01 (`Proyecto`)
- **Estado:** revisada

## Historia

**Como** líder del proyecto **quiero** crear un sprint con fechas y un Sprint Goal **para** planificar el trabajo.

## Objetivo

Crear sprints dentro del período del proyecto, sin superponerse, y poder iniciar uno a la vez. El sprint es el
contenedor donde se asignan (HU-11), completan (HU-12) y cierran (HU-13) las historias.

## Entradas

**Crear sprint**

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Sprint Goal | `string` | Sí | No vacío (sin contar espacios al inicio y al final) |
| Fecha de inicio | `time.Time` (día) | Sí | Dentro de las fechas del proyecto |
| Fecha de fin | `time.Time` (día) | Sí | No anterior a la de inicio; dentro de las fechas del proyecto |

**Iniciar sprint**: el `domain.SprintID` del sprint a iniciar.

## Salidas esperadas

- Crear: un `sprint.Sprint` guardado, con `ID` generado (`S1`, `S2`, …), estado `Planificado` y sin historias.
- Iniciar: el sprint pasa a estado `Activo`.

## Reglas de negocio

- **RN-10.1:** el Sprint Goal es obligatorio.
- **RN-10.2:** la fecha de fin no puede ser anterior a la de inicio.
- **RN-10.3:** las fechas del sprint están dentro de las del proyecto (inicio y fin incluidos).
- **RN-10.4:** un sprint no se superpone con otro sprint del mismo proyecto. Dos sprints se superponen si comparten al
  menos un día: un sprint puede empezar el día siguiente al fin del anterior, pero no el mismo día.
- **RN-10.5:** un sprint nuevo queda `Planificado`. Solo un sprint del proyecto puede estar `Activo` a la vez.
- **RN-10.6:** solo se inicia un sprint `Planificado`.

## Restricciones

- RN-10.1, RN-10.2 y RN-10.5 se validan en `sprint.Nuevo(proyectoID, goal, inicio, fin) (Sprint, error)`.
- RN-10.3 recibe el `project.Proyecto` y RN-10.4 los otros sprints del proyecto: los coordina el servicio
  (`sprint` puede importar `project` según el ADR 0002).
- `Sprint.Iniciar()` cambia el estado; el servicio verifica antes que no haya otro sprint activo.

## Casos límite

- Sprint que empieza el día de inicio del proyecto y termina el día de fin: válido.
- Sprint de un solo día (inicio = fin): válido.
- S1 del 23 al 25 y S2 del 26 al 27: válido. S2 del 25 al 27: superpuesto.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| Sprint Goal vacío | `sprint.ErrGoalVacio`: "el Sprint Goal es obligatorio" |
| Fin anterior al inicio | `sprint.ErrFechasInvalidas`: "la fecha de fin no puede ser anterior a la de inicio" |
| Fechas fuera del proyecto | `sprint.ErrFueraDelProyecto`: "las fechas del sprint deben estar dentro de las del proyecto" |
| Se superpone con otro sprint | `sprint.ErrSuperpuesto`: "las fechas se superponen con otro sprint del proyecto" |
| Iniciar con otro sprint activo | `sprint.ErrOtroActivo`: "ya hay un sprint activo en el proyecto" |
| Iniciar un sprint que no está planificado | `sprint.ErrNoPlanificado`: "solo se puede iniciar un sprint planificado" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-10.1 | El Sprint Goal es obligatorio | `Escenario: crear un sprint`, `Escenario: Sprint Goal vacío` |
| CA-10.2 | Las fechas están dentro de las del proyecto y no se superponen con otro sprint | `Escenario: fechas fuera del proyecto`, `Escenario: sprints superpuestos`, `Escenario: sprints consecutivos` |
| CA-10.3 | Solo puede haber un sprint activo a la vez | `Escenario: iniciar un sprint`, `Escenario: ya hay un sprint activo` |
