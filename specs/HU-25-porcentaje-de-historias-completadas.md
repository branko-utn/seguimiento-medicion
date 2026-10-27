# HU-25 · Porcentaje de historias completadas

- **Issue:** #32 · **Épica:** E7 Métricas · **Responsable:** Agustín Giorlando · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-25-porcentaje-de-historias-completadas.feature`](../features/hu-25-porcentaje-de-historias-completadas.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Agustín Giorlando (al aprobar el PR)
- **Depende de:** nada: función pura (ADR 0002, sección 3)
- **Estado:** implementada (Sprint 1)

## Historia

**Como** líder del proyecto **quiero** ver el porcentaje de historias completadas **para** conocer el avance.

## Objetivo

Calcular el avance en porcentaje, por sprint y por proyecto. Lo usan el estado del proyecto (HU-04), el dashboard
(HU-28) y los reportes (E9).

## Entradas

| Función | Entrada | Salida |
|---|---|---|
| `PorcentajeCompletadas(completadas, total int) (float64, error)` | Cantidad de historias completadas y total de referencia | Porcentaje entre 0 y 100 |

El servicio elige el total según lo que se mide:
- **Por sprint:** total = historias asignadas al sprint.
- **Por proyecto:** total = historias del Product Backlog del proyecto.

## Salidas esperadas

Un porcentaje entre 0 y 100. Se muestra con un decimal; la función devuelve el valor sin redondear.

## Reglas de negocio

- **RN-25.1:** porcentaje = completadas / total × 100.
- **RN-25.2:** si el total es 0, el porcentaje es 0 (no se divide por cero).
- **RN-25.3:** las completadas no pueden ser más que el total, ni ninguno de los dos negativo.

## Restricciones

- Función pura en `internal/metrics`, sin dependencias de otros paquetes.

## Casos límite

- 0 de 0: 0 %. 0 de 5: 0 %. 5 de 5: 100 %.
- 1 de 3: 33,333… (se muestra 33,3 %).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Completadas mayores que el total | `metrics.ErrDatosInconsistentes`: "las historias completadas no pueden superar al total" |
| Valores negativos | `metrics.ErrDatosInconsistentes` |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-25.1 | Porcentaje = completadas / planificadas × 100, por sprint y por proyecto | `Esquema del escenario: porcentaje de completadas` |
| CA-25.2 | Si no hay historias planificadas, el porcentaje es 0 | `Escenario: sin historias` |
