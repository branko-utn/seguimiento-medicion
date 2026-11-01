# HU-26 · Horas estimadas, reales y desviación

- **Issue:** #33 · **Épica:** E7 Métricas · **Responsable:** Agustín Giorlando · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-26-horas-estimadas-reales-y-desviacion.feature`](../features/hu-26-horas-estimadas-reales-y-desviacion.feature)
- **Depende de:** nada: funciones puras con tipos de entrada propios (ADR 0002, sección 3)
- **Redactada por:** Branko Almeira (Agile Enabler) después de la planning del Sprint 2 · **Validada por:** Agustín Giorlando (al aprobar el PR)
- **Estado:** implementada (Sprint 2)

## Historia

**Como** líder del proyecto **quiero** comparar las horas estimadas con las reales **para** mejorar nuestras estimaciones.

## Objetivo

Calcular las horas reales (suma del esfuerzo registrado, HU-20) y la desviación porcentual respecto de las horas
estimadas (HU-19). Lo usan la vista de métricas (TEC-17), el dashboard (HU-28) y los reportes (E9).

## Entradas

| Función | Entrada | Salida |
|---|---|---|
| `HorasReales(horas []float64) float64` | Horas de cada registro de esfuerzo | Suma |
| `Desviacion(estimadas, reales float64) (float64, error)` | Horas estimadas y reales | Desviación en % |

El servicio arma las entradas por historia o por sprint (sumando las historias del sprint).

## Salidas esperadas

La desviación en porcentaje, sin redondear (la pantalla la muestra con un decimal y con signo).

## Reglas de negocio

- **RN-26.1:** horas reales = suma de las horas registradas.
- **RN-26.2:** desviación (%) = (reales − estimadas) / estimadas × 100. Positiva si se trabajó de más, negativa si de menos.
- **RN-26.3:** sin horas estimadas (0), la desviación **no es calculable**: se informa con un error y la pantalla muestra "no calculable".
- **RN-26.4:** las horas no pueden ser negativas. `HorasReales` recibe horas ya validadas por el dominio (HU-20 solo
  acepta horas mayores que cero) y las suma todas, sin filtrar.

## Restricciones

- Funciones puras en `internal/metrics`, sin importar paquetes de dominio.

## Casos límite

- Sin registros: 0 horas reales; con estimadas, la desviación es −100 %.
- Reales iguales a las estimadas: 0 %.
- Estimadas 0 y reales 0: no calculable.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Horas estimadas en 0 | `metrics.ErrDesviacionNoCalculable`: "la desviación no es calculable sin horas estimadas" |
| Horas estimadas o reales negativas | `metrics.ErrHorasNegativas`: "las horas no pueden ser negativas" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-26.1 | Horas reales = suma del esfuerzo registrado | `Escenario: horas reales` |
| CA-26.2 | Desviación (%) = (reales − estimadas) / estimadas × 100; positiva si se trabajó de más | `Esquema del escenario: desviación` |
| CA-26.3 | Sin horas estimadas, la desviación se informa como "no calculable" | `Escenario: sin horas estimadas` |
