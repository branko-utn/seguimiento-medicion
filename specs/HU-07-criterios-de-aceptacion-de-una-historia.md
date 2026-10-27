# HU-07 · Criterios de aceptación de una historia

- **Issue:** #14 · **Épica:** E2 Backlog · **Responsable:** Ignacio Bianchi · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-07-criterios-de-aceptacion-de-una-historia.feature`](../features/hu-07-criterios-de-aceptacion-de-una-historia.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Ignacio Bianchi (al aprobar el PR)
- **Depende de:** HU-05 (`Historia`)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** integrante del equipo **quiero** agregar y quitar criterios de aceptación a una historia **para** saber cuándo está terminada.

## Objetivo

Mantener la lista de criterios de aceptación de cada historia. Una historia sin criterios no está lista para un sprint.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Agregar | `domain.HistoriaID`, texto del criterio | La historia existe; texto no vacío |
| Quitar | `domain.HistoriaID`, número del criterio (1, 2, …) | La historia existe; el número corresponde a un criterio |

## Salidas esperadas

- Agregar: el criterio queda al final de `Historia.Criterios`, con el texto sin espacios sobrantes.
- Quitar: el criterio sale de la lista y los siguientes se renumeran (el 3 pasa a ser el 2).
- `Historia.TieneCriterios() bool` indica si tiene al menos uno (lo usa HU-11).

## Reglas de negocio

- **RN-07.1:** un criterio no puede estar vacío.
- **RN-07.2:** una historia sin criterios no puede asignarse a un sprint (lo controla HU-11 con `TieneCriterios()`).
- **RN-07.3:** los criterios se identifican por su posición, empezando en 1.

## Restricciones

- `Historia.AgregarCriterio(texto string) error` y `Historia.QuitarCriterio(numero int) error` en `internal/domain/backlog`.
- `backlog.Criterio` tiene al menos el campo `Texto string`.

## Casos límite

- Quitar el único criterio: la historia queda sin criterios (válido) y `TieneCriterios()` devuelve `false`.
- Quitar el criterio 0 o uno mayor que la cantidad: error.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Criterio vacío | `backlog.ErrCriterioVacio`: "el criterio de aceptación no puede estar vacío" |
| Número de criterio inexistente | `backlog.ErrCriterioNoExiste`: "el criterio no existe" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-07.1 | Un criterio no puede estar vacío | `Escenario: agregar un criterio`, `Escenario: criterio vacío` |
| CA-07.2 | Una historia sin criterios no puede asignarse a un sprint | `Escenario: quitar el único criterio` (el rechazo al asignar se prueba en HU-11) |
