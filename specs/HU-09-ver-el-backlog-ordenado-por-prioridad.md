# HU-09 · Ver el backlog ordenado por prioridad

- **Issue:** #16 · **Épica:** E2 Backlog · **Responsable:** Agustín Lara · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-09-ver-el-backlog-ordenado-por-prioridad.feature`](../features/hu-09-ver-el-backlog-ordenado-por-prioridad.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-05 (`Historia`, `Prioridad`), HU-08 (estados)
- **Estado:** borrador para la planning del Sprint 3

## Historia

**Como** integrante del equipo **quiero** ver el backlog ordenado por prioridad **para** saber qué sigue.

## Objetivo

Listar el Product Backlog de un proyecto en el orden en que conviene trabajarlo y poder quedarse solo con las historias
de un estado (por ejemplo, las pendientes). La pantalla del backlog (TEC-16) usa este orden.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Estado | `backlog.Estado` | No | Si se indica, Pendiente, En progreso o Completada |

## Salidas esperadas

Las historias del proyecto ordenadas según RN-09.1, filtradas por estado si se indicó. La consulta no modifica nada.

## Reglas de negocio

- **RN-09.1:** se ordena por prioridad (Alta, Media, Baja) y, a igual prioridad, por identificador **numérico**
  ascendente: HU-2 va antes que HU-10.
- **RN-09.2:** si se indica un estado, solo se listan las historias en ese estado, en el mismo orden.
- **RN-09.3:** solo se listan las historias del proyecto indicado.

## Restricciones

- `backlog.Ordenar(hs []Historia) []Historia` en `internal/domain/backlog`: función pura que devuelve una copia ordenada
  (no modifica la entrada).
- El servicio filtra por proyecto y estado y después ordena: `Backlog.DelProyecto(proyectoID, estado *backlog.Estado)`.
  Con `estado == nil` no filtra.

## Casos límite

- Backlog vacío: lista vacía, sin error.
- Ninguna historia en el estado pedido: lista vacía, sin error.
- Todas con la misma prioridad: quedan ordenadas por identificador.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| Estado inválido | `backlog.ErrEstadoInvalido`: "el estado debe ser Pendiente, En progreso o Completada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-09.1 | Se ordena por prioridad (Alta, Media, Baja) y, a igual prioridad, por identificador | `Escenario: ordenar por prioridad e identificador` |
| CA-09.2 | Se puede filtrar por estado | `Escenario: filtrar por estado`, `Escenario: ninguna historia en el estado pedido` |

## Nota para los pasos

La tabla compartida `el proyecto tiene las historias:` suma la columna opcional `prioridad` (si falta, Media). La agrega
su responsable en `comun_backlog_steps_test.go`, no en el archivo de HU-09.
