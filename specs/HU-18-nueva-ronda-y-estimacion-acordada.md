# HU-18 · Nueva ronda y estimación acordada

- **Issue:** #25 · **Épica:** E4 Estimación · **Responsable:** Ignacio Bianchi · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-18-nueva-ronda-y-estimacion-acordada.feature`](../features/hu-18-nueva-ronda-y-estimacion-acordada.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el responsable al aprobar el PR de la planning del Sprint 3
- **Depende de:** solo código en `main`: HU-16 (`estimation.Sesion`, `service.PlanningPoker`), HU-17 (`Sesion.Revelar`),
  HU-15 (`Historia.Estimar`)
- **Estado:** aprobada en la planning del Sprint 3

## Historia

**Como** líder del proyecto **quiero** abrir nuevas rondas y registrar la estimación acordada **para** cerrar la
estimación de la historia.

## Objetivo

Completar el Planning Poker: si al revelar no hubo acuerdo, se discute y se vota de nuevo en otra ronda; cuando el
equipo acuerda un valor, queda como Story Points de la historia y la votación termina.

## Entradas

Como en HU-16 y HU-17, la votación se identifica por su historia: cada historia tiene como máximo una votación abierta.

| Operación | Datos | Validación |
|---|---|---|
| Nueva ronda | `domain.HistoriaID` | La historia tiene una votación abierta con la ronda actual revelada |
| Registrar la estimación acordada | `domain.HistoriaID`, Story Points (`int`) | Votación abierta con la ronda actual revelada; el valor es de `estimation.Escala` |
| Consultar el historial | `domain.HistoriaID` | La historia tiene una votación abierta |

## Salidas esperadas

- Nueva ronda: la votación pasa a la ronda siguiente (2, 3, …) sin votos. Las rondas anteriores quedan en el historial
  con sus votos revelados.
- Estimación acordada: la historia queda con esos Story Points (como en HU-15) y la votación se **cierra**: deja de
  haber una votación abierta para la historia.
- Historial: las rondas ya reveladas, cada una con su número y sus votos.

## Reglas de negocio

- **RN-18.1:** solo se abre una ronda nueva después de revelar la actual (HU-17). Una ronda nueva empieza sin votos.
- **RN-18.2:** las rondas anteriores no se modifican: quedan en el historial con sus votos.
- **RN-18.3:** la estimación acordada se registra después de revelar la ronda actual. Puede ser cualquier valor de la
  escala, aunque nadie lo haya votado (es lo que el equipo acordó al discutir).
- **RN-18.4:** al registrar la estimación, la historia toma esos Story Points con `Historia.Estimar` (HU-15) y la
  votación se cierra: no acepta votos, rondas ni otra estimación.
- **RN-18.5:** con la votación cerrada, la historia puede abrir otra votación nueva (RN-16.5 solo impide dos abiertas).

## Restricciones

- En `internal/domain/estimation`:
  - `Sesion` guarda las rondas reveladas anteriores (`rondas []Ronda`, con número y votos), y las persiste en el JSON
    igual que los votos (`sesionJSON`).
  - `Sesion.NuevaRonda() error`: archiva la ronda actual en el historial y abre la siguiente sin votos.
  - `Sesion.Acordar(sp int) error`: valida que la ronda actual esté revelada y que el valor sea de la escala.
  - `Sesion.Rondas() []Ronda`: el historial, con la ronda actual al final si ya se reveló.
  - `estimation` no importa `backlog` (ADR 0002).
- En `service.PlanningPoker`: `NuevaRonda(historiaID)`, `Acordar(historiaID, sp)` y `Rondas(historiaID)`.
  - `Acordar` coordina RN-18.4: valida con `Sesion.Acordar`, estima la historia con `Historia.Estimar` y la guarda, y
    **elimina la votación abierta** (`RepoSesiones` suma `Eliminar(historiaID) error`). Si la historia no se puede
    estimar, no cambia nada.
  - "Votación cerrada" es no tener votación abierta: no hace falta un estado nuevo y la historia puede abrir otra
    votación (RN-18.5) con `IniciarVotacion`.
- **Pantalla:** en `/planning-poker`, una vez revelada la ronda, se ofrecen el botón "Nueva ronda" y un desplegable con
  la escala para "Registrar estimación acordada"; debajo, el historial de rondas anteriores.

## Casos límite

- Acordar en la ronda 1 con consenso: la votación se cierra con una sola ronda en el historial.
- Varias rondas sin acuerdo: el historial las conserva todas, en orden.
- Acordar un valor que nadie votó: válido (RN-18.3).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La historia no tiene una votación abierta (nunca se abrió o ya se acordó) | `estimation.ErrSinVotacionAbierta` (HU-16): "la historia no tiene una votación abierta" |
| Nueva ronda o acuerdo sin revelar la ronda actual | `estimation.ErrRondaSinRevelar`: "primero hay que revelar los votos de la ronda" |
| El valor acordado no es de la escala | `estimation.ErrStoryPointsInvalidos`: "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21" |
| La historia se completó mientras se votaba | `backlog.ErrHistoriaCompletada` (HU-15): "no se puede re-estimar una historia completada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-18.1 | Una ronda nueva empieza sin votos; las anteriores quedan en el historial | `Escenario: nueva ronda después de revelar`, `Escenario: nueva ronda sin revelar` |
| CA-18.2 | La estimación acordada se guarda como Story Points de la historia | `Escenario: registrar la estimación acordada`, `Escenario: acordar un valor que nadie votó`, `Escenario: acordar sin revelar`, `Escenario: acordar un valor fuera de la escala`, `Escenario: operar sobre una votación cerrada` |

## Preguntas para la planning

- **¿Quién acuerda?** Sin login (ADR 0004), cualquiera puede registrar la estimación. Propuesta: aceptarlo, como en HU-16.
