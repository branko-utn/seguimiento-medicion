# HU-16 · Votar en Planning Poker con votos ocultos

- **Issue:** #23 · **Épica:** E4 Estimación · **Responsable:** Ignacio Bianchi · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-16-votar-en-planning-poker.feature`](../features/hu-16-votar-en-planning-poker.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 2 (aprobada)
- **Depende de:** HU-03 (`Integrante`), HU-05 (`Historia`), HU-15 (`estimation.Escala`)
- **Estado:** implementada (Sprint 2)

## Historia

**Como** integrante del equipo **quiero** votar la estimación de una historia sin ver los votos de los demás **para** no influirnos.

## Objetivo

Abrir una votación de Planning Poker para una historia y registrar el voto de cada integrante sin mostrarlo hasta que
se revelen las cartas (HU-17). Es la primera parte del Planning Poker; HU-17 revela y HU-18 abre rondas y registra la
estimación acordada.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Iniciar votación | `domain.HistoriaID` | La historia existe, no está completada y no tiene otra votación abierta |
| Votar | votación, `domain.IntegranteID`, carta (`int`) | El integrante es del proyecto de la historia; la carta es de `estimation.Escala`; no votó en esta ronda |
| Consultar quién votó | votación | — |

## Salidas esperadas

- Iniciar: una `estimation.Sesion` abierta en la ronda 1, sin votos.
- Votar: el voto queda registrado en la ronda actual.
- Consultar: la lista de integrantes que ya votaron, **sin el valor de sus votos**.

## Reglas de negocio

- **RN-16.1:** solo votan los integrantes del proyecto de la historia. En la pantalla, quien vota se elige de esa lista (ADR 0004).
- **RN-16.2:** cada integrante vota una vez por ronda.
- **RN-16.3:** la carta tiene que ser un valor de la escala de Story Points (1, 2, 3, 5, 8, 13, 21; `estimation.Escala`).
- **RN-16.4:** mientras la votación está abierta, los valores de los votos no se pueden consultar: solo quién votó.
- **RN-16.5:** una historia tiene como máximo una votación abierta, y no se vota una historia completada.

## Restricciones

- `estimation.Sesion` (ADR 0002) referencia la historia y los integrantes **por ID**; `estimation` no importa `backlog`
  ni `project`. Que la historia exista, no esté completada y que el integrante sea de su proyecto lo controla el servicio.
- `Sesion.Votar(integranteID, carta) error`, `Sesion.QuienesVotaron() []domain.IntegranteID`.
- La sesión no expone los valores de los votos mientras está abierta (los expone HU-17 al revelar).

## Casos límite

- Votación sin votos: `QuienesVotaron()` devuelve una lista vacía.
- Los extremos de la escala (1 y 21): válidos.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no existe | `backlog.ErrNoExiste`: "la historia no existe" |
| La historia está completada | `estimation.ErrHistoriaCompletada`: "no se puede votar una historia completada" |
| Ya hay una votación abierta para la historia | `estimation.ErrVotacionAbierta`: "la historia ya tiene una votación abierta" |
| El integrante no existe | `project.ErrIntegranteNoExiste`: "el integrante no existe" |
| El integrante no es del proyecto | `estimation.ErrNoEsDelProyecto`: "el integrante no pertenece al proyecto de la historia" |
| La carta no es de la escala | `estimation.ErrStoryPointsInvalidos`: "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21" |
| El integrante ya votó en esta ronda | `estimation.ErrYaVoto`: "el integrante ya votó en esta ronda" |
| Consultar los valores de los votos antes de revelar | `estimation.ErrVotosOcultos`: "los votos se muestran recién al revelar las cartas" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-16.1 | Solo votan los integrantes del proyecto, una vez por ronda | `Escenario: votar`, `Escenario: votar dos veces en la misma ronda`, `Escenario: votar sin ser del proyecto` |
| CA-16.2 | Mientras la votación está abierta, los votos no se muestran (solo quién ya votó) | `Escenario: los votos quedan ocultos` |

## Decisiones de la planning del Sprint 2

- **Identidad sin login (ADR 0004):** quien vota se elige de la lista de integrantes del proyecto ("Votás como: …").
  No hay autenticación: cualquiera puede elegir otro nombre, y es una limitación aceptada para una herramienta interna.
- **Sin cartas "?" ni "☕":** solo se vota con la escala 1, 2, 3, 5, 8, 13, 21 (RN-16.3).
- **No se cambia el voto antes de revelar:** cada integrante vota una vez por ronda (RN-16.2).
