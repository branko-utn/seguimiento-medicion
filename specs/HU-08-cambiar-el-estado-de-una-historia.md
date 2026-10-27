# HU-08 · Cambiar el estado de una historia

- **Issue:** #15 · **Épica:** E2 Backlog · **Responsable:** Ignacio Bianchi · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-08-cambiar-el-estado-de-una-historia.feature`](../features/hu-08-cambiar-el-estado-de-una-historia.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Ignacio Bianchi (al aprobar el PR)
- **Relacionada con:** HU-11 (asignar, define `SprintID`) y HU-12 (completar en el sprint)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** integrante del equipo **quiero** cambiar el estado de una historia **para** reflejar su avance.

## Objetivo

Reflejar el avance de cada historia con tres estados y transiciones controladas, coherentes con el sprint al que
pertenece.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Historia | `domain.HistoriaID` | Sí | Debe existir |
| Estado nuevo | `backlog.Estado` | Sí | `Pendiente`, `EnProgreso` o `Completada` |

## Salidas esperadas

La historia queda en el estado nuevo.

## Reglas de negocio

- **RN-08.1:** los estados son Pendiente → En progreso → Completada.
- **RN-08.2:** transiciones válidas:

  | Desde | Hacia | Condición |
  |---|---|---|
  | Pendiente | En progreso | La historia está asignada a un sprint |
  | En progreso | Pendiente | — (se puede volver atrás) |
  | Pendiente o En progreso | Completada | La historia está asignada a un sprint |

- **RN-08.3:** una historia completada no cambia de estado.
- **RN-08.4:** cambiar al mismo estado en el que ya está es un error (no hay cambio que registrar).

## Restricciones

- `Historia.CambiarEstado(nuevo Estado) error` valida RN-08.2 a RN-08.4. `Historia.Completar() error` aplica la misma
  regla para pasar a Completada y es la que usa HU-12.
- **El paso a Completada desde la aplicación se hace con HU-12** (completar en el sprint), para que el sprint también lo
  registre. El servicio de HU-08 delega en el de HU-12 cuando el estado nuevo es Completada.

## Casos límite

- En progreso → Pendiente: válido.
- Pendiente → Completada estando en un sprint: válido (alguien la terminó sin marcarla "en progreso").

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Estado fuera de la lista | `backlog.ErrEstadoInvalido`: "el estado debe ser Pendiente, En progreso o Completada" |
| Transición no permitida (desde Completada, o al mismo estado) | `backlog.ErrTransicionInvalida`: "no se puede pasar de <estado actual> a <estado nuevo>" |
| En progreso o Completada sin estar en un sprint | `backlog.ErrSinSprint`: "la historia debe estar asignada a un sprint" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-08.1 | Los estados son Pendiente → En progreso → Completada | `Escenario: empezar una historia del sprint`, `Escenario: volver a pendiente` |
| CA-08.2 | Las transiciones inválidas informan un error | `Escenario: completar sin estar en un sprint`, `Escenario: empezar sin estar en un sprint`, `Escenario: cambiar una historia completada` |
