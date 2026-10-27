# HU-18 · Nueva ronda y estimación acordada

- **Issue:** #25 · **Épica:** E4 Estimación · **Responsable:** Ignacio Bianchi · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-18-nueva-ronda-y-estimacion-acordada.feature`](../features/hu-18-nueva-ronda-y-estimacion-acordada.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-16 (`estimation.Sesion`, votar), HU-17 (`Sesion.Revelar`), HU-15 (`Historia.Estimar`)
- **Estado:** borrador para la planning del Sprint 3

## Historia

**Como** líder del proyecto **quiero** abrir nuevas rondas y registrar la estimación acordada **para** cerrar la
estimación de la historia.

## Objetivo

Completar el Planning Poker: si al revelar no hubo acuerdo, se discute y se vota de nuevo en otra ronda; cuando el
equipo acuerda un valor, queda como Story Points de la historia y la votación termina.

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Nueva ronda | votación | Abierta y con la ronda actual revelada |
| Registrar la estimación acordada | votación, Story Points (`int`) | Abierta, con la ronda actual revelada; el valor es de `estimation.Escala` |
| Consultar el historial | votación | — |

## Salidas esperadas

- Nueva ronda: la votación pasa a la ronda siguiente (2, 3, …) sin votos. Las rondas anteriores quedan en el historial
  con sus votos revelados.
- Estimación acordada: la historia queda con esos Story Points (como en HU-15) y la votación queda **cerrada**.
- Historial: las rondas, cada una con su número y sus votos (los de la ronda actual, solo si ya se reveló).

## Reglas de negocio

- **RN-18.1:** solo se abre una ronda nueva después de revelar la actual (HU-17). Una ronda nueva empieza sin votos.
- **RN-18.2:** las rondas anteriores no se modifican: quedan en el historial con sus votos.
- **RN-18.3:** la estimación acordada se registra después de revelar la ronda actual. Puede ser cualquier valor de la
  escala, aunque nadie lo haya votado (es lo que el equipo acordó al discutir).
- **RN-18.4:** al registrar la estimación, la historia toma esos Story Points con `Historia.Estimar` (HU-15) y la
  votación se cierra: no acepta votos, rondas ni otra estimación.
- **RN-18.5:** con la votación cerrada, la historia puede abrir otra votación nueva (RN-16.5 solo impide dos abiertas).

## Restricciones

- En `internal/domain/estimation`: `Sesion.NuevaRonda() error`, `Sesion.Acordar(sp int) error`, `Sesion.Rondas()` (el
  historial) y `Sesion.Cerrada() bool`. `estimation` no importa `backlog` (ADR 0002).
- El servicio de votaciones coordina RN-18.4: valida y cierra la sesión, llama a `Historia.Estimar` y guarda las dos. Si
  la historia no se puede estimar (por ejemplo, porque se completó), no guarda nada.

## Casos límite

- Acordar en la ronda 1 con consenso: la votación se cierra con una sola ronda en el historial.
- Varias rondas sin acuerdo: el historial las conserva todas, en orden.
- Acordar un valor que nadie votó: válido (RN-18.3).

## Condiciones de error

| Situación | Error esperado |
|---|---|
| La votación no existe | `estimation.ErrVotacionNoExiste`: "la votación no existe" |
| Nueva ronda o acuerdo sin revelar la ronda actual | `estimation.ErrRondaSinRevelar`: "primero hay que revelar los votos de la ronda" |
| Operar sobre una votación cerrada | `estimation.ErrVotacionCerrada`: "la votación ya está cerrada" |
| El valor acordado no es de la escala | `estimation.ErrStoryPointsInvalidos`: "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21" |
| La historia se completó mientras se votaba | `backlog.ErrHistoriaCompletada` (HU-15): "no se puede re-estimar una historia completada" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-18.1 | Una ronda nueva empieza sin votos; las anteriores quedan en el historial | `Escenario: nueva ronda después de revelar`, `Escenario: nueva ronda sin revelar` |
| CA-18.2 | La estimación acordada se guarda como Story Points de la historia | `Escenario: registrar la estimación acordada`, `Escenario: acordar un valor que nadie votó`, `Escenario: acordar sin revelar`, `Escenario: acordar un valor fuera de la escala`, `Escenario: operar sobre una votación cerrada` |

## Preguntas para la planning

- **¿Quién acuerda?** Sin login (ADR 0004), cualquiera puede registrar la estimación. Propuesta: aceptarlo, como en HU-16.
