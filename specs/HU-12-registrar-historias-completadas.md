# HU-12 · Registrar historias completadas

- **Issue:** #19 · **Épica:** E3 Sprints · **Responsable:** Branko Almeira · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-12-registrar-historias-completadas.feature`](../features/hu-12-registrar-historias-completadas.feature)
- **Depende de:** HU-05 (`Historia`, Bianchi), HU-10 (`Sprint`, Sosa), HU-11 (asignar historias, Sosa)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** integrante del equipo **quiero** marcar historias del sprint como completadas **para** medir el avance.

## Objetivo

Registrar que una historia asignada a un sprint se terminó, para que cuente en el avance del sprint: sus Story Points
pasan a sumar como completados y la historia queda en estado "Completada" en el backlog.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Sprint | `domain.SprintID` | Sí | Debe existir |
| Historia | `domain.HistoriaID` | Sí | Debe existir |

## Salidas esperadas

- La historia queda en la lista `Completadas` del sprint.
- La historia queda en estado `Completada` en el backlog.
- Los Story Points completados del sprint (calculados por HU-24) suman los de la historia.

## Reglas de negocio

- **RN-12.1:** solo se completan historias **asignadas a ese sprint** (que están en `Sprint.Historias`).
- **RN-12.2:** no se completan historias en un sprint **cerrado** (HU-13: un sprint cerrado no se modifica).
- **RN-12.3:** una historia ya completada no se vuelve a completar.
- **RN-12.4:** una historia del sprint se puede completar desde "Pendiente" o "En progreso".
- **RN-12.5:** al completarla, sus Story Points cuentan como completados en ese sprint.

## Restricciones

- Las reglas del sprint (RN-12.1, RN-12.2, RN-12.3) viven en `internal/domain/sprint` (`Sprint.Completar`) y las de la
  historia (RN-12.4) en `internal/domain/backlog` (`Historia.Completar`). El servicio coordina las dos y guarda ambas
  (ADR 0002, sección 4).
- Si alguna regla falla, no se modifica ninguna de las dos entidades.

## Casos límite

- Completar la última historia pendiente del sprint: los SP completados igualan a los planificados.
- Una historia de 0 SP (sin estimar) no puede estar en un sprint (HU-11), así que no se contempla aquí.
- Completar una historia que ya estaba completada: error, y los SP completados no cambian (no se cuentan dos veces).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| La historia no está asignada a ese sprint | `sprint.ErrHistoriaNoAsignada`: "la historia no está asignada a este sprint" |
| El sprint está cerrado | `sprint.ErrSprintCerrado`: "el sprint está cerrado y no se puede modificar" |
| La historia ya está completada | `sprint.ErrHistoriaCompletada`: "la historia ya está completada" (el mismo error que usa HU-11) |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-12.1 | Solo se completan historias asignadas al sprint | `Escenario: completar una historia del sprint`, `Escenario: completar una historia que no está en el sprint` |
| CA-12.2 | Al completarla, sus Story Points cuentan como completados en el sprint | `Escenario: los Story Points de la historia cuentan como completados` |

Escenarios adicionales por las reglas de negocio: `sprint cerrado` (RN-12.2), `historia ya completada` (RN-12.3).
