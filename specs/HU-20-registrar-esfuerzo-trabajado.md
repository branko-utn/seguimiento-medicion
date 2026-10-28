# HU-20 · Registrar esfuerzo trabajado

- **Issue:** #27 · **Épica:** E5 Esfuerzo · **Responsable:** Agustín Lara · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-20-registrar-esfuerzo-trabajado.feature`](../features/hu-20-registrar-esfuerzo-trabajado.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Agustín Lara (al aprobar el PR)
- **Depende de:** solo `HistoriaID` e `IntegranteID` (ADR 0002)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** integrante del equipo **quiero** registrar las horas que trabajé en una historia **para** conocer el esfuerzo real.

## Objetivo

Registrar el esfuerzo real por integrante, fecha y actividad. La suma de estos registros son las horas reales de cada
historia (HU-26) y alimenta los reportes (E9).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Integrante | `domain.IntegranteID` | Sí | Debe existir y ser del proyecto de la historia |
| Historia | `domain.HistoriaID` | Sí | Debe existir |
| Fecha | `time.Time` (día) | Sí | No posterior a hoy |
| Actividad | `string` | Sí | No vacía (por ejemplo "Desarrollo", "Tests", "Revisión de PR") |
| Horas | `float64` | Sí | Mayor que cero |

## Salidas esperadas

Un `effort.Registro` guardado con `ID` generado (`E1`, `E2`, …).

## Reglas de negocio

- **RN-20.1:** se registra integrante, fecha, actividad, horas e historia; todos son obligatorios.
- **RN-20.2:** las horas son mayores que cero.
- **RN-20.3:** un integrante no puede superar 24 horas en un mismo día, sumando todos sus registros de ese día y el nuevo.
- **RN-20.4:** la fecha no puede ser futura; hoy es válido.
- **RN-20.5:** el integrante debe pertenecer al proyecto de la historia.

## Restricciones

- `effort.NuevoRegistro(integranteID, historiaID, fecha, actividad, horas, hoy) (Registro, error)` valida RN-20.1,
  RN-20.2 y RN-20.4. **La fecha de hoy se recibe como parámetro** para que los tests no dependan del reloj.
- RN-20.3 y RN-20.5 las controla el servicio con los registros existentes y los integrantes del proyecto.

## Casos límite

- Registros de 10 + 14 horas el mismo día: válido (24 justas). Un registro más de 0,5 horas ese día: error.
- 24 horas en un día y 1 hora al día siguiente: válido.
- Fecha de hoy: válida. Mañana: inválida.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El integrante no existe | `project.ErrIntegranteNoExiste`: "el integrante no existe" |
| El integrante es de otro proyecto | `effort.ErrIntegranteDeOtroProyecto`: "el integrante no pertenece al proyecto de la historia" |
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Actividad vacía | `effort.ErrActividadVacia`: "la actividad es obligatoria" |
| Fecha vacía | `effort.ErrFechaVacia`: "la fecha es obligatoria" (lo agregó TEC-15, #90) |
| Horas iguales a cero o negativas | `effort.ErrHorasInvalidas`: "las horas deben ser mayores que cero" |
| Superaría 24 horas en el día | `effort.ErrSuperaHorasDiarias`: "el integrante superaría las 24 horas en ese día" |
| Fecha futura | `effort.ErrFechaFutura`: "la fecha no puede ser futura" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-20.1 | Se registra integrante, fecha, actividad, horas e historia | `Escenario: registrar esfuerzo`, `Escenario: actividad vacía` |
| CA-20.2 | Las horas son mayores que cero y no superan 24 por integrante y por día | `Escenario: horas en cero`, `Escenario: 24 horas justas en un día`, `Escenario: superar 24 horas en un día` |
| CA-20.3 | La fecha no puede ser futura | `Escenario: fecha futura` |
