# HU-22 · Registrar defecto

- **Issue:** #29 · **Épica:** E6 Defectos · **Responsable:** Agustín Lara · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-22-registrar-defecto.feature`](../features/hu-22-registrar-defecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Agustín Lara (al aprobar el PR)
- **Depende de:** solo `HistoriaID` y `SprintID` (ADR 0002)
- **Estado:** revisada

## Historia

**Como** integrante del equipo **quiero** registrar un defecto **para** hacer el seguimiento de la calidad.

## Objetivo

Registrar los defectos encontrados, con su severidad y en qué sprint se detectaron. Su resolución es HU-23 y su
conteo HU-27.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Descripción | `string` | Sí | No vacía |
| Severidad | `defect.Severidad` | Sí | `Critica`, `Alta`, `Media` o `Baja` |
| Historia relacionada | `domain.HistoriaID` | Sí | Debe existir |
| Sprint de detección | `domain.SprintID` | Sí | Debe existir y ser del proyecto de la historia |

## Salidas esperadas

Un `defect.Defecto` guardado con `ID` generado (`D1`, `D2`, …), estado `Abierto` y sin sprint de resolución.

## Reglas de negocio

- **RN-22.1:** descripción, severidad, historia relacionada y sprint de detección son obligatorios.
- **RN-22.2:** la severidad es Crítica, Alta, Media o Baja.
- **RN-22.3:** un defecto nuevo queda en estado `Abierto`, sin sprint de resolución.
- **RN-22.4:** el sprint de detección es del mismo proyecto que la historia.

## Restricciones

- `defect.Nuevo(descripcion, severidad, historiaID, sprintDeteccion) (Defecto, error)` valida RN-22.1 a RN-22.3.
- RN-22.4 y la existencia de la historia y del sprint las controla el servicio.

## Casos límite

- Defecto en una historia ya completada: válido (es justamente cuando suelen aparecer).
- Defecto detectado en un sprint cerrado: válido (se registra tarde, pero se detectó en ese sprint).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Descripción vacía | `defect.ErrDescripcionVacia`: "la descripción es obligatoria" |
| Severidad fuera de la lista | `defect.ErrSeveridadInvalida`: "la severidad debe ser Crítica, Alta, Media o Baja" |
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| El sprint no existe | `sprint.ErrNoExiste`: "el sprint no existe" |
| Sprint de otro proyecto | `defect.ErrSprintDeOtroProyecto`: "el sprint de detección no pertenece al proyecto de la historia" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-22.1 | Descripción, severidad, historia relacionada y sprint de detección son obligatorios | `Escenario: descripción vacía`, `Escenario: severidad inválida`, `Escenario: historia inexistente` |
| CA-22.2 | Un defecto nuevo queda en estado "Abierto" | `Escenario: registrar un defecto` |
