# HU-14 · Consultar sprints anteriores

- **Issue:** #21 · **Épica:** E3 Sprints · **Responsable:** Agustín Giorlando · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-14-consultar-sprints-anteriores.feature`](../features/hu-14-consultar-sprints-anteriores.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el responsable al aprobar el PR de la planning del Sprint 3
- **Depende de:** HU-10 (`Sprint`), HU-12 (`Sprint.Completadas`), HU-13 (`Sprint.Cerrar`, SP fijos al cerrar)
- **Estado:** aprobada en la planning del Sprint 3

## Historia

**Como** integrante del equipo **quiero** ver los sprints cerrados con sus resultados **para** comparar el desempeño.

## Objetivo

Listar los sprints que ya terminaron con lo que se planificó y lo que se completó en cada uno, y poder abrir uno para ver
sus historias. El dashboard (HU-28) y el reporte de proyecto (HU-30) leen los mismos datos del sprint cerrado
(`SPPlanificados`, `SPCompletados`, `Historias`, `Completadas`), sin esperar a esta historia.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Listar sprints cerrados | `domain.ProyectoID` | El proyecto existe |
| Ver el detalle de un sprint cerrado | `domain.SprintID` | El sprint existe y está cerrado |

## Salidas esperadas

- Listado: por cada sprint cerrado del proyecto, sus fechas, el Sprint Goal y los SP planificados y completados que
  quedaron fijos al cerrarlo (HU-13).
- Detalle: las historias que tuvo el sprint (`Sprint.Historias`), con título, Story Points y si se completaron en ese
  sprint (están en `Sprint.Completadas`).

## Reglas de negocio

- **RN-14.1:** solo se listan los sprints en estado `Cerrado`; los planificados y el activo no aparecen.
- **RN-14.2:** el listado se ordena por fecha de inicio, del más antiguo al más reciente.
- **RN-14.3:** los SP que se muestran son los fijados al cerrar (`SPPlanificados` y `SPCompletados` del sprint), no se
  recalculan: si una historia cambió de Story Points después, el resultado del sprint no cambia.
- **RN-14.4:** el detalle muestra todas las historias que se planificaron en el sprint, también las que volvieron al
  backlog al cerrarlo (HU-13), marcadas como "no completada".

## Restricciones

- El servicio de sprints suma `Cerrados(proyectoID) ([]sprint.Sprint, error)` y `Detalle(sprintID) (DetalleSprint, error)`.
- `DetalleSprint` tiene el sprint y una lista de historias con `ID`, `Titulo`, `StoryPoints` y `Completada bool`.
- **Pantalla:** la lista de sprints (`/sprints`, TEC-14) suma la sección **Sprints cerrados** con una tabla por proyecto
  (fechas, Sprint Goal, SP planificados y completados, enlace al detalle). El detalle (`/sprints/{sprint}`) ya muestra
  las historias, los SP y las métricas (TEC-14 y TEC-17); HU-14 agrega la marca "No completada" a las historias que
  volvieron al backlog (RN-14.4).
- Al cerrar, HU-13 deja `Sprint.Historias` como estaba: el detalle se arma con esa lista, no con el `SprintID` de las
  historias (que se vacía en las que vuelven al backlog).

## Casos límite

- Proyecto sin sprints cerrados: listado vacío; la pantalla muestra "Todavía no hay sprints cerrados".
- Sprint cerrado sin historias: detalle con la lista vacía y 0 SP.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| El sprint no está cerrado | `sprint.ErrNoCerrado`: "el sprint todavía no está cerrado" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-14.1 | Se listan con fechas, Sprint Goal, SP planificados y completados | `Escenario: listar sprints cerrados`, `Escenario: sin sprints cerrados` |
| CA-14.2 | Se puede ver el detalle de las historias de cada sprint | `Escenario: detalle de un sprint cerrado`, `Escenario: detalle de un sprint que no está cerrado` |
