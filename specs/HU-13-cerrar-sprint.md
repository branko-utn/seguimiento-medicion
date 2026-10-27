# HU-13 · Cerrar sprint

- **Issue:** #20 · **Épica:** E3 Sprints · **Responsable:** Branko Almeira · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-13-cerrar-sprint.feature`](../features/hu-13-cerrar-sprint.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 2 (aprobada)
- **Depende de:** HU-10 y HU-11 (`Sprint`), HU-12 (`Sprint.Completadas`), HU-08 (estados de la historia), HU-24 (SP)
- **Estado:** aprobada en la planning del Sprint 2

## Historia

**Como** líder del proyecto **quiero** cerrar un sprint **para** registrar su resultado.

## Objetivo

Terminar un sprint dejando registrado cuánto se planificó y cuánto se completó, y devolver al Product Backlog lo que no
se terminó para que se pueda planificar en otro sprint. Es la base de la velocidad (HU-24), de los sprints anteriores
(HU-14) y del reporte de sprint (HU-29).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Sprint | `domain.SprintID` | Sí | Debe existir y estar activo |

## Salidas esperadas

- El sprint queda en estado `Cerrado`, con `SPPlanificados` y `SPCompletados` fijos.
- Las historias del sprint que no se completaron vuelven al backlog: estado `Pendiente` y sin `SprintID`.
- Las historias completadas quedan como están (estado `Completada`, asignadas al sprint).

## Reglas de negocio

- **RN-13.1:** solo se cierra un sprint **activo**.
- **RN-13.2:** al cerrar, quedan fijos los Story Points planificados (suma de los SP de `Sprint.Historias`) y los
  completados (suma de los SP de `Sprint.Completadas`), calculados con `metrics.SPPlanificados` y `metrics.SPCompletados`.
- **RN-13.3:** las historias no completadas vuelven al backlog como `Pendiente`, sin sprint asignado (también las que
  estaban "En progreso").
- **RN-13.4:** un sprint cerrado no se puede modificar: no se le asignan ni completan historias (ya lo controlan
  HU-11 y HU-12 con `ErrSprintCerrado`) ni se vuelve a cerrar.

## Restricciones

- `Sprint.Cerrar(spPlanificados, spCompletados int) error` valida RN-13.1 y fija los valores; `Historia.VolverAlBacklog()`
  (paquete `backlog`) deja la historia en `Pendiente` y sin sprint.
- El servicio coordina: busca el sprint y sus historias, calcula los SP, cierra el sprint, devuelve las no completadas
  y guarda todo. Si alguna validación falla, no guarda nada.

## Casos límite

- Sprint sin historias: se cierra con 0 SP planificados y 0 completados.
- Todas las historias completadas: ninguna vuelve al backlog; planificados = completados.
- Ninguna completada: todas vuelven al backlog; 0 SP completados.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| El sprint está cerrado | `sprint.ErrSprintCerrado`: "el sprint está cerrado y no se puede modificar" |
| El sprint está planificado (todavía no empezó) | `sprint.ErrNoActivo`: "solo se puede cerrar un sprint activo" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-13.1 | Las historias no completadas vuelven al backlog como Pendientes | `Escenario: cerrar un sprint con historias sin terminar` |
| CA-13.2 | Quedan registrados los SP planificados y completados del sprint | `Escenario: cerrar un sprint con historias sin terminar`, `Escenario: cerrar un sprint sin historias` |
| CA-13.3 | Un sprint cerrado no se puede modificar | `Escenario: cerrar un sprint ya cerrado` |

## Decisiones de la planning del Sprint 2

- Un sprint **planificado** (que nunca se inició) no se puede cerrar: RN-13.1 queda como está.
