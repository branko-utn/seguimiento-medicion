# HU-05 · Crear ítem del Product Backlog

- **Issue:** #12 · **Épica:** E2 Backlog · **Responsable:** Ignacio Bianchi · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-05-crear-item-del-product-backlog.feature`](../features/hu-05-crear-item-del-product-backlog.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Ignacio Bianchi (al aprobar el PR)
- **Depende de:** HU-01 (`Proyecto`)
- **Estado:** implementada (Sprint 1) · **Cambio 24/10/2026:** IDs únicos en el sistema (antes, por proyecto), para que `Buscar(HistoriaID)` no sea ambiguo

## Historia

**Como** líder del proyecto **quiero** agregar historias al Product Backlog **para** registrar lo que hay que construir.

## Objetivo

Crear las historias del Product Backlog de un proyecto. `backlog.Historia` es el tipo central que usan sprints,
estimación, esfuerzo, defectos, métricas y reportes (ADR 0002).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Título | `string` | Sí | No vacío (sin contar espacios al inicio y al final) |
| Descripción | `string` | No | — |
| Prioridad | `backlog.Prioridad` | Sí | `Alta`, `Media` o `Baja` |

## Salidas esperadas

Una `backlog.Historia` guardada con:
- `ID` generado y único en todo el sistema: `HU-1`, `HU-2`, … en orden de creación (la numeración no se reinicia por proyecto).
- Título sin espacios sobrantes, estado `Pendiente`, `StoryPoints` en 0 (sin estimar), sin criterios y sin sprint.

## Reglas de negocio

- **RN-05.1:** el identificador se genera automáticamente y es único en todo el sistema: `HU-1`, `HU-2`, … sin reiniciar por proyecto. Así las demás historias pueden buscar una historia solo por su ID (ADR 0002).
- **RN-05.2:** el título es obligatorio; la descripción es opcional.
- **RN-05.3:** la prioridad es Alta, Media o Baja.
- **RN-05.4:** una historia nueva queda en estado `Pendiente`, sin Story Points, sin criterios y sin sprint.

## Restricciones

- RN-05.2, RN-05.3 y RN-05.4 se validan en `backlog.Nueva(proyectoID, titulo, descripcion, prioridad) (Historia, error)`.
- RN-05.1: el ID lo asigna el repositorio al guardar, según la cantidad total de historias guardadas.

## Casos límite

- Primera historia del sistema: `HU-1`. Si ya hay 3 historias (de cualquier proyecto), la siguiente es `HU-4`, aunque sea la primera de otro proyecto.
- Título `"  Login  "`: se guarda como `"Login"`.
- Descripción vacía: válido.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| Título vacío | `backlog.ErrTituloVacio`: "el título es obligatorio" |
| Prioridad fuera de la lista | `backlog.ErrPrioridadInvalida`: "la prioridad debe ser Alta, Media o Baja" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-05.1 | El identificador se genera automáticamente y es único | `Escenario: los identificadores son correlativos`, `Escenario: los identificadores no se repiten entre proyectos` |
| CA-05.2 | El título es obligatorio; la descripción es opcional | `Escenario: crear una historia sin descripción`, `Escenario: título vacío` |
| CA-05.3 | La prioridad es Alta, Media o Baja | `Escenario: prioridad inválida` |
| CA-05.4 | Una historia nueva queda en estado "Pendiente" y sin Story Points | `Escenario: crear una historia` |
