# HU-19 · Registrar horas estimadas de una historia

- **Issue:** #26 · **Épica:** E5 Esfuerzo · **Responsable:** Agustín Lara · **Story Points:** 1
- **Escenarios BDD:** [`features/hu-19-registrar-horas-estimadas-de-una-historia.feature`](../features/hu-19-registrar-horas-estimadas-de-una-historia.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Agustín Lara (al aprobar el PR)
- **Depende de:** solo el `HistoriaID` (ADR 0002): se puede desarrollar sin esperar a HU-05
- **Estado:** revisada

## Historia

**Como** integrante del equipo **quiero** registrar las horas estimadas de una historia **para** compararlas después con las reales.

## Objetivo

Guardar cuántas horas se espera que lleve cada historia. Con las horas reales (HU-20) se calcula la desviación
entre esfuerzo estimado y real (HU-26).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Historia | `domain.HistoriaID` | Sí | Debe existir |
| Horas estimadas | `float64` | Sí | Mayor que cero |

## Salidas esperadas

Una `effort.EstimacionHoras` guardada para la historia. Cada historia tiene una sola estimación: registrar otra la reemplaza.

## Reglas de negocio

- **RN-19.1:** las horas estimadas son mayores que cero.
- **RN-19.2:** una historia tiene una sola estimación de horas; la nueva reemplaza a la anterior.
- **RN-19.3:** se aceptan fracciones de hora (por ejemplo 1,5).

## Restricciones

- `effort.NuevaEstimacion(historiaID, horas) (EstimacionHoras, error)` valida RN-19.1.
- El servicio verifica que la historia exista y guarda (o reemplaza) la estimación.

## Casos límite

- 0,5 horas: válido. 0: inválido. -2: inválido.
- Re-estimar de 4 a 6 horas: queda 6.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| Horas iguales a cero o negativas | `effort.ErrHorasInvalidas`: "las horas deben ser mayores que cero" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-19.1 | Las horas estimadas son mayores que cero | `Escenario: registrar horas estimadas`, `Esquema del escenario: horas inválidas` |
