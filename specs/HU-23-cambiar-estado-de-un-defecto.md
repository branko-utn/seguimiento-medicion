# HU-23 · Cambiar estado de un defecto y registrar su resolución

- **Issue:** #30 · **Épica:** E6 Defectos · **Responsable:** Agustín Lara · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-23-cambiar-estado-de-un-defecto.feature`](../features/hu-23-cambiar-estado-de-un-defecto.feature)
- **Depende de:** HU-22 (`Defecto`, en `main`), HU-10 (`Sprint`, en `main`)
- **Redactada por:** Branko Almeira (Agile Enabler) después de la planning del Sprint 2 · **Validada por:** Agustín Lara (al aprobar el PR)
- **Estado:** implementada (Sprint 2)

## Historia

**Como** integrante del equipo **quiero** cambiar el estado de un defecto **para** saber cuáles están resueltos.

## Objetivo

Seguir el avance de cada defecto y, al resolverlo, registrar en qué sprint se resolvió. Con eso HU-27 cuenta los
defectos resueltos por sprint. Incluye la pantalla de cambio de estado.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Empezar a trabajarlo / volver a abrirlo | `domain.DefectoID`, estado nuevo (`EnProgreso` o `Abierto`) | El defecto existe; transición válida |
| Resolverlo | `domain.DefectoID`, sprint de resolución (`domain.SprintID`) | El defecto existe y no está resuelto; el sprint existe, es del mismo proyecto y no es anterior al de detección |

## Salidas esperadas

- El defecto queda en el estado nuevo.
- Al resolverlo, queda `Resuelto` con `SprintResolucion` igual al sprint indicado.

## Reglas de negocio

- **RN-23.1:** los estados son Abierto → En progreso → Resuelto. Transiciones válidas: Abierto → En progreso,
  En progreso → Abierto, y Abierto o En progreso → Resuelto.
- **RN-23.2:** un defecto resuelto no cambia de estado; cambiar al mismo estado es un error.
- **RN-23.3:** para resolver, el sprint de resolución es obligatorio, debe existir y ser del mismo proyecto que el de detección.
- **RN-23.4:** el sprint de resolución no puede ser anterior al de detección: su fecha de inicio no puede ser anterior a
  la del sprint de detección. Puede ser el mismo sprint.

## Restricciones

- `Defecto.CambiarEstado(nuevo Estado) error` y `Defecto.Resolver(sprintID domain.SprintID) error` en `internal/domain/defect`.
  `CambiarEstado` no acepta `Resuelto`: para resolver se usa `Resolver`, que exige el sprint.
- RN-23.3 y RN-23.4 las controla el servicio con los dos sprints (`defect` no importa `sprint`, ADR 0002).
- `Defectos.Buscar` hoy devuelve `(defecto, encontrado)`; para esta historia el servicio devuelve `defect.ErrNoExiste`.

## Casos límite

- Resolver en el mismo sprint de detección: válido.
- Resolver directamente desde Abierto: válido.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El defecto no existe | `defect.ErrNoExiste`: "el defecto no existe" |
| Transición no permitida (desde Resuelto o al mismo estado) | `defect.ErrTransicionInvalida`: "no se puede pasar de <estado actual> a <estado nuevo>" |
| Pasar a Resuelto con CambiarEstado | `defect.ErrResolverSinSprint`: "para resolver un defecto hay que indicar el sprint de resolución" |
| El sprint de resolución no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| Sprint de resolución de otro proyecto | `defect.ErrResolucionDeOtroProyecto`: "el sprint de resolución no pertenece al proyecto del defecto" |
| Sprint de resolución anterior al de detección | `defect.ErrResolucionAnterior`: "el sprint de resolución no puede ser anterior al de detección" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-23.1 | Los estados son Abierto → En progreso → Resuelto | `Escenario: empezar a trabajar un defecto`, `Escenario: resolver un defecto`, `Escenario: cambiar un defecto resuelto` |
| CA-23.2 | Al resolverlo, el sprint de resolución es obligatorio y no puede ser anterior al de detección | `Escenario: resolver en un sprint anterior al de detección`, `Escenario: resolver en el mismo sprint directamente desde Abierto`, `Escenario: resolver sin sprint` |
