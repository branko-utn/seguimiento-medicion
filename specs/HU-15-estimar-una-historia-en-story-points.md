# HU-15 · Estimar una historia en Story Points

- **Issue:** #22 · **Épica:** E4 Estimación · **Responsable:** Ignacio Bianchi · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-15-estimar-una-historia-en-story-points.feature`](../features/hu-15-estimar-una-historia-en-story-points.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Ignacio Bianchi (al aprobar el PR)
- **Depende de:** HU-05 (`Historia`)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** integrante del equipo **quiero** asignar Story Points a una historia **para** conocer su tamaño.

## Objetivo

Guardar la estimación de cada historia con la escala de Planning Poker. Es la base de los Story Points planificados,
completados y de la velocidad (E7). El Planning Poker (HU-16 a HU-18) termina usando esta misma operación.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Historia | `domain.HistoriaID` | Sí | Debe existir y no estar completada |
| Story Points | `int` | Sí | 1, 2, 3, 5, 8, 13 o 21 |

## Salidas esperadas

`Historia.StoryPoints` queda con el valor indicado. Si ya tenía una estimación, se reemplaza.

## Reglas de negocio

- **RN-15.1:** solo se aceptan valores de la escala 1, 2, 3, 5, 8, 13, 21.
- **RN-15.2:** una historia se puede re-estimar mientras no esté completada.
- **RN-15.3:** no se puede volver a "sin estimar" (0) asignando Story Points.

## Restricciones

- La escala está en `internal/domain/estimation` (`estimation.Escala` y `estimation.ValidarStoryPoints(sp int) error`),
  para que el Planning Poker la reutilice. `backlog` importa `estimation` (ADR 0002).
- `Historia.Estimar(sp int) error` valida RN-15.1 a RN-15.3.

## Casos límite

- 1 y 21 (extremos de la escala): válidos.
- 4, 0, -1 y 34: inválidos.
- Re-estimar una historia asignada a un sprint activo: válido (cambian los SP planificados del sprint).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Valor fuera de la escala | `estimation.ErrStoryPointsInvalidos`: "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21" |
| Historia completada | `backlog.ErrHistoriaCompletada`: "no se puede re-estimar una historia completada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-15.1 | Solo se aceptan valores de la escala 1, 2, 3, 5, 8, 13, 21 | `Esquema del escenario: valores de la escala` |
| CA-15.2 | No se puede re-estimar una historia completada | `Escenario: re-estimar una historia`, `Escenario: re-estimar una historia completada` |
