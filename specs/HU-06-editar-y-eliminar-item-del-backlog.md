# HU-06 · Editar y eliminar ítem del backlog

- **Issue:** #13 · **Épica:** E2 Backlog · **Responsable:** Agustín Lara · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-06-editar-y-eliminar-item-del-backlog.feature`](../features/hu-06-editar-y-eliminar-item-del-backlog.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-05 (`Historia`, validaciones de título y prioridad), HU-11 (asignación a sprint), HU-08 (estados)
- **Estado:** borrador para la planning del Sprint 3

## Historia

**Como** líder del proyecto **quiero** editar o eliminar historias **para** mantener el backlog actualizado.

## Objetivo

Corregir el título, la descripción o la prioridad de una historia, y sacar del backlog las que ya no se van a hacer,
sin romper lo que ya se planificó o se terminó.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Editar | `domain.HistoriaID`, título, descripción, prioridad | Las mismas que al crear (HU-05); la historia existe y no está completada |
| Eliminar | `domain.HistoriaID` | La historia existe, no está asignada a un sprint y no está completada |

## Salidas esperadas

- Editar: la historia queda guardada con el título (sin espacios sobrantes), la descripción y la prioridad nuevos. El
  ID, el estado, los Story Points, los criterios y el sprint no cambian.
- Eliminar: la historia deja de estar en el backlog del proyecto. Su ID no se reutiliza.

## Reglas de negocio

- **RN-06.1:** se editan solo el título, la descripción y la prioridad, con las validaciones de HU-05 (título
  obligatorio; prioridad Alta, Media o Baja).
- **RN-06.2:** una historia completada no se edita: es parte del resultado de un sprint.
- **RN-06.3:** no se elimina una historia asignada a un sprint ni una completada.
- **RN-06.4:** los IDs son únicos en todo el sistema (HU-05): una historia eliminada no libera su ID.

## Restricciones

- `Historia.Editar(titulo, descripcion string, prioridad Prioridad) error` en `internal/domain/backlog` valida RN-06.1
  y RN-06.2 sin cambiar la historia si falla.
- `Historia.PuedeEliminarse() error` valida RN-06.3; el servicio la llama antes de `Eliminar` en el repositorio.
- El repositorio de historias suma `Eliminar(id domain.HistoriaID) error` (`memory.Historias`).

## Casos límite

- Editar sin cambiar nada: se guarda igual, sin error.
- Editar una historia asignada a un sprint (no completada): se permite.
- Eliminar la única historia del backlog: el backlog queda vacío.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Título vacío | `backlog.ErrTituloVacio`: "el título es obligatorio" |
| Prioridad inválida | `backlog.ErrPrioridadInvalida`: "la prioridad debe ser Alta, Media o Baja" |
| Editar una historia completada | `backlog.ErrCompletadaNoSeEdita`: "una historia completada no se puede editar" |
| Eliminar una historia asignada a un sprint | `backlog.ErrAsignadaNoSeElimina`: "no se puede eliminar una historia asignada a un sprint" |
| Eliminar una historia completada | `backlog.ErrCompletadaNoSeElimina`: "no se puede eliminar una historia completada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-06.1 | Se pueden editar título, descripción y prioridad | `Escenario: editar una historia`, `Escenario: editar con título vacío` |
| CA-06.2 | No se puede eliminar una historia asignada a un sprint ni una completada | `Escenario: eliminar una historia del backlog`, `Escenario: eliminar una historia asignada a un sprint`, `Escenario: eliminar una historia completada` |

## Preguntas para la planning

- **Historias con esfuerzo o defectos.** Una historia sin sprint puede tener esfuerzo registrado (HU-20) o defectos
  (HU-22). Propuesta: tampoco se elimina (`"no se puede eliminar una historia con esfuerzo o defectos registrados"`),
  para no dejar registros huérfanos. ¿Lo sumamos o alcanza con RN-06.3?
