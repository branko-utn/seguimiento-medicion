# HU-11 · Asignar historias a un sprint

- **Issue:** #18 · **Épica:** E3 Sprints · **Responsable:** Federico Sosa · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-11-asignar-historias-a-un-sprint.feature`](../features/hu-11-asignar-historias-a-un-sprint.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Federico Sosa (al aprobar el PR)
- **Depende de:** HU-05 y HU-07 (`Historia`, Bianchi), HU-10 (`Sprint`), HU-15 (Story Points, Bianchi)
- **Estado:** revisada

## Historia

**Como** líder del proyecto **quiero** asignar historias del backlog a un sprint **para** armar el Sprint Backlog.

## Objetivo

Armar el Sprint Backlog solo con historias listas para trabajar (estimadas y con criterios de aceptación), sin que una
historia quede en dos sprints a la vez.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Sprint | `domain.SprintID` | Sí | Debe existir y no estar cerrado |
| Historia | `domain.HistoriaID` | Sí | Debe existir, ser del mismo proyecto, estar estimada y tener criterios |

## Salidas esperadas

- La historia queda en `Sprint.Historias`.
- La historia queda con `Historia.SprintID` igual al sprint.
- Los Story Points planificados del sprint (calculados por HU-24) suman los de la historia.

## Reglas de negocio

- **RN-11.1:** solo se asignan historias estimadas (`StoryPoints > 0`).
- **RN-11.2:** solo se asignan historias con al menos un criterio de aceptación.
- **RN-11.3:** una historia no puede estar en dos sprints abiertos: si ya tiene `SprintID`, no se asigna a otro.
- **RN-11.4:** no se asignan historias a un sprint cerrado.
- **RN-11.5:** no se asignan historias completadas.
- **RN-11.6:** la historia y el sprint deben ser del mismo proyecto.

## Restricciones

- `Sprint.Asignar(h backlog.Historia) error` valida RN-11.1, RN-11.2, RN-11.4, RN-11.5 y RN-11.6 (`sprint` importa `backlog`, ADR 0002).
- `Historia.AsignarASprint(id domain.SprintID) error` valida RN-11.3 y guarda el `SprintID`.
- El servicio coordina las dos **en este orden**: primero `Sprint.Asignar(h)` y después `h.AsignarASprint(id)`. Así una historia completada (que siempre tiene `SprintID`) informa "la historia ya está completada" y no "ya está asignada". Guarda ambas solo si las dos validaciones pasan; si una falla, no se modifica ninguna.
- Al cerrar un sprint (HU-13), las historias no completadas vuelven al backlog con `SprintID` vacío, y se pueden asignar a otro.

## Casos límite

- Asignar la misma historia dos veces al mismo sprint: error de RN-11.3 (ya asignada), sin duplicarla en la lista.
- Asignar a un sprint planificado o activo: válido. A uno cerrado: error.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Historia sin estimar | `sprint.ErrHistoriaSinEstimar`: "la historia no está estimada" |
| Historia sin criterios | `sprint.ErrHistoriaSinCriterios`: "la historia no tiene criterios de aceptación" |
| Historia ya asignada a un sprint | `backlog.ErrYaAsignada`: "la historia ya está asignada a un sprint" |
| Sprint cerrado | `sprint.ErrSprintCerrado`: "el sprint está cerrado y no se puede modificar" |
| Historia completada | `sprint.ErrHistoriaCompletada`: "la historia ya está completada" |
| Historia de otro proyecto | `sprint.ErrOtroProyecto`: "la historia pertenece a otro proyecto" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-11.1 | Solo se asignan historias estimadas y con criterios de aceptación | `Escenario: asignar una historia lista`, `Escenario: historia sin estimar`, `Escenario: historia sin criterios` |
| CA-11.2 | Una historia no puede estar en dos sprints abiertos | `Escenario: historia ya asignada a otro sprint` |
| CA-11.3 | No se asignan historias a un sprint cerrado | `Escenario: sprint cerrado` |
