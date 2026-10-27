# HU-17 · Revelar votos y detectar diferencias

- **Issue:** #24 · **Épica:** E4 Estimación · **Responsable:** Ignacio Bianchi · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-17-revelar-votos-y-detectar-diferencias.feature`](../features/hu-17-revelar-votos-y-detectar-diferencias.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 2 (aprobada)
- **Depende de:** HU-16 (`estimation.Sesion`)
- **Estado:** aprobada en la planning del Sprint 2

## Historia

**Como** líder del proyecto **quiero** revelar los votos y ver si hay diferencias **para** que se discutan.

## Objetivo

Mostrar todos los votos de la ronda a la vez y señalar si coinciden. Si no coinciden, indicar quiénes votaron el valor
más alto y el más bajo, que son los que explican su voto antes de una nueva ronda (HU-18).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Votación | `estimation.Sesion` de la historia | Sí | Abierta y con al menos un voto |

## Salidas esperadas

Un resultado de la ronda con:
- todos los votos (integrante y carta);
- si hay **consenso** (todos votaron lo mismo);
- si no hay consenso: el valor más alto y quiénes lo votaron, y el valor más bajo y quiénes lo votaron.

Después de revelar, la ronda queda cerrada: no se aceptan más votos en ella.

## Reglas de negocio

- **RN-17.1:** al revelar se muestran todos los votos de la ronda.
- **RN-17.2:** hay consenso cuando todos los votos son iguales. Con un solo voto también hay consenso.
- **RN-17.3:** si no hay consenso, se informan el valor más alto y el más bajo, con todos los integrantes que votaron cada uno.
- **RN-17.4:** solo se revela una ronda con al menos un voto; una ronda revelada no acepta votos nuevos.

## Restricciones

- `Sesion.Revelar() (Resultado, error)` en `internal/domain/estimation`. `Resultado` tiene al menos `Votos`, `Consenso`,
  `Maximo`, `Minimo`, `VotaronMaximo` y `VotaronMinimo`.
- Los valores se comparan como números de la escala (HU-16 solo acepta cartas de `estimation.Escala`).

## Casos límite

- Un solo voto: consenso.
- Dos integrantes votan el valor más alto: los dos aparecen en `VotaronMaximo`.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Revelar una ronda sin votos | `estimation.ErrSinVotos`: "no hay votos para revelar" |
| Revelar una ronda ya revelada | `estimation.ErrYaRevelada`: "la ronda ya fue revelada" |
| Votar en una ronda ya revelada | `estimation.ErrYaRevelada`: "la ronda ya fue revelada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-17.1 | Al revelar, se muestran todos los votos de la ronda | `Escenario: revelar con consenso`, `Escenario: revelar con diferencias` |
| CA-17.2 | Si los votos no coinciden, se señala la diferencia y quiénes votaron el valor más alto y el más bajo | `Escenario: revelar con diferencias`, `Escenario: varios votaron el valor más alto` |
