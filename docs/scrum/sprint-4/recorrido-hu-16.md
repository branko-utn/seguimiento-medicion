# TEC-22 · Recorrido completo con HU-16: Votar en Planning Poker con votos ocultos

- **Issue:** #120 (TEC-22 · Presentación final y demo) · **Responsable:** Ignacio Bianchi
- **Funcionalidad elegida:** HU-16 · Votar en Planning Poker con votos ocultos (issue #23), historia propia del
  Sprint 2, implementada con el flujo completo SDD → BDD → TDD ([PR #85](https://github.com/branko-utn/seguimiento-medicion/pull/85)).
- **Meta del recorrido:** contar la funcionalidad de punta a punta en la presentación: historia → spec →
  criterios → escenarios BDD → tests → código Go, y demostrarla en vivo con los datos de demo.

## 1. La historia

> **Como** integrante del equipo **quiero** votar la estimación de una historia sin ver los votos de los demás
> **para** no influirnos.

La idea del Planning Poker es que cada uno estime de forma independiente: si se ven las cartas, los primeros
votos arrastran a los últimos. HU-16 es la primera parte de la épica E4 (Estimación): HU-17 revela los votos y
HU-18 abre nuevas rondas y registra la estimación acordada, las tres construyendo sobre el mismo tipo
`estimation.Sesion`.

## 2. La spec SDD (fuente de verdad)

Especificación completa en [`specs/HU-16-votar-en-planning-poker.md`](../../../specs/HU-16-votar-en-planning-poker.md),
redactada como borrador antes del Sprint 2 (`3e89e8f`), validada por el equipo en la planning y aprobada en su PR.

**Entradas y operaciones:**

| Operación | Datos | Validación |
|---|---|---|
| Iniciar votación | historia | Existe, no está completada, no tiene otra votación abierta |
| Votar | votación, integrante, carta | Integrante del proyecto de la historia; carta de la escala; no votó en esta ronda |
| Consultar quién votó | votación | — |

**Reglas que cuenta la presentación:**

| Regla | Cómo se dice en la demo |
|---|---|
| RN-16.1 | "Solo votan los integrantes del proyecto de la historia" |
| RN-16.2 | "Cada uno vota una vez por ronda" |
| RN-16.3 | "Las cartas son siempre 1, 2, 3, 5, 8, 13 o 21" |
| RN-16.4 | "Los valores quedan ocultos: se ve quién votó, no qué votó" |
| RN-16.5 | "Una sola votación abierta por historia, y una completada no se vota" |

**Identidad sin login (ADR 0004):** en la pantalla cada uno elige su nombre de la lista del proyecto
("Votás como…"). No hay autenticación; es una limitación aceptada para una herramienta interna y es
decisión de la planning, junto con dos más: sin cartas "?" ni café, y no se cambia el voto antes de revelar.

## 3. Criterios de aceptación y escenarios BDD

Los escenarios viven en [`features/hu-16-votar-en-planning-poker.feature`](../../../features/hu-16-votar-en-planning-poker.feature)
(base compartida: proyecto con Ana y Beto, la historia HU-1 y la votación ya abierta):

| Criterio | Escenarios |
|---|---|
| **CA-16.1** — solo votan los integrantes del proyecto, una vez por ronda | `votar`, `votar dos veces en la misma ronda` (error), `votar sin ser del proyecto` (error) |
| **CA-16.2** — con la votación abierta, los valores quedan ocultos | `los votos quedan ocultos` |
| Casos extra | `votación sin votos` (límite), `votar una carta fuera de la escala` y `abrir otra votación para la misma historia` (error) |

El escenario que mejor resume la historia es `los votos quedan ocultos`: dos integrantes votan, se consultan
**los valores** y responde con el error "los votos se muestran recién al revelar las cartas"; en cambio,
consultar **quién votó** sí muestra "Ana Pérez, Beto Gómez". Ver valores y ver nombres son dos consultas distintas
por diseño.

## 4. El ciclo TDD: RED → GREEN por capa

Cada capa se escribió test-first: primero el test que falla (RED), después el código mínimo para pasarlo (GREEN).
Historial de commits de HU-16:

| Fase | Commit | Qué fija |
|---|---|---|
| Spec (antes del sprint) | `3e89e8f` docs | Borrador de la spec y los escenarios, para la planning |
| Repositorio (RED) | `6de4499` test | Las sesiones se guardan y buscan una por historia |
| Repositorio (GREEN) | `3a6022f` feat | `memory.Sesiones` sobre el store genérico de JSON |
| Servicio (RED) | `802c151` test | El servicio valida historia, integrante y pertenencia al proyecto |
| Servicio (GREEN) | `35678d6` feat | Casos de uso `IniciarVotacion`, `Votar`, `QuienesVotaron`, `ConsultarVotos` |
| Pantalla | `62823dc` feat | Handlers y plantilla con los votos ocultos ("Votaron: …") |
| BDD (RED) | `3ea9997` test | Pasos `.feature` conectados al mundo con las APIs reales |
| BDD (GREEN) | `947fd5d` feat | Los 7 escenarios del `.feature` en verde |
| Cierre | `0cc56d6` docs | Registro de uso de IA; PR #85 revisado y mergeado |

## 5. Los tests, capas

| Capa | Archivo | Tests de HU-16 |
|---|---|---|
| Dominio | `internal/domain/estimation/sesion_test.go` | `NuevaSesion` (abierta en ronda 1, sin votos), `Votar_QuienesVotaron`, `Votar_DosVecesMismaRonda`, `Votar_CartaFueraDeEscala`, `Votar_ExtremosEscala` (1 y 21), `VotosOcultosMientrasAbierta`, `JSONRoundTrip` (persistencia) |
| Repositorio | `internal/store/memory/sesiones_test.go` | `TestSesionesGuardarBuscar`: reemplaza la sesión de la misma historia, no duplica |
| Servicio | `internal/service/planning_poker_test.go` | `IniciarVotacion` (y `_errores`, `_historiaCompletada`), `Votar` (y `_errores`, `_historiaInexistente`, `_sinVotacionAbierta`), `QuienesVotaron_sinVotacionAbierta`, `ConsultarVotos_ocultos`, `ConsultarVotos_sinVotacionAbierta` |
| BDD | `features/hu16_steps_test.go` | Los 8 pasos que conectan el `.feature` al mundo con las APIs reales (`seAbrioLaVotacion`, `votaEnLaHistoria`, `enLaVotacionVotaron`, `consultoLosVotos`, …) |

Los nombres describen **comportamiento**, no numeración de tickets (convenio del repo): en el archivo el test se
llama `TestSesion_Votar_CartaFueraDeEscala`, nunca "testHU16".

## 6. El código Go, capa por capa

### Dominio — `internal/domain/estimation` (ADR 0002)

La sesión `estimation.Sesion` es por historia, arranca abierta en la ronda 1 y sin votos. Los votos van a un campo
**no exportado**: el ocultamiento no es una regla de la pantalla, es imposibilidad de consultarlos desde fuera
del tipo.

```go
// NuevaSesion crea una votación abierta en la ronda 1 y sin votos.
func NuevaSesion(historiaID domain.HistoriaID) (Sesion, error) {
    return Sesion{HistoriaID: historiaID, Ronda: 1, rondas: make([]Ronda, 0)}, nil
}

// Votar registra el voto: una ronda revelada no acepta votos (RN-17.4), la carta es de la
// escala (RN-16.3) y cada integrante vota una sola vez por ronda (RN-16.2).
func (s *Sesion) Votar(integranteID domain.IntegranteID, carta int) error {
    if s.revelada { return ErrYaRevelada }
    if s.yaVoto(integranteID) { return ErrYaVoto }
    if err := ValidarStoryPoints(carta); err != nil { return err }
    s.votos = append(s.votos, Voto{IntegranteID: integranteID, Carta: carta})
    return nil
}

// QuienesVotaron: los integrantes en orden, SIN los valores (RN-16.4).
func (s *Sesion) QuienesVotaron() []domain.IntegranteID { ... }

// Votos: los valores, solo si ya se reveló (RN-16.4); si no, ErrVotosOcultos.
func (s *Sesion) Votos() ([]Voto, error) {
    if !s.revelada { return nil, ErrVotosOcultos }
    return append([]Voto{}, s.votos...), nil
}
```

Los errores son valores del paquete (`var ErrYaVoto = errors.New("el integrante ya votó en esta ronda")`, etc.),
que el servicio y la pantalla muestran tal cual.

### Repositorio — `internal/store/memory/sesiones.go`

Guarda **una sesión por historia** (RN-16.5) sobre el store genérico de JSON (`data/sesiones.json` en la app):

```go
// Guardar crea o reemplaza la votación de la historia de la sesión (RN-16.5).
func (r *Sesiones) Guardar(s estimation.Sesion) error { ... }

// Buscar devuelve la votación y si existe: que no exista no es un error, es un estado válido.
func (r *Sesiones) Buscar(id domain.HistoriaID) (estimation.Sesion, bool) { ... }
```

### Servicio — `internal/service/planning_poker.go`

Coordina el dominio con los repositorios de otras épicas: la historia viene de E2 (HU-05) y los integrantes de
E1 (HU-03). `Votar` valida en orden: la historia existe, hay votación abierta, el integrante existe y **es del
proyecto de la historia** (RN-16.1), y recién entonces delega en `Sesion.Votar` (RN-16.2 y 16.3). Sin login, la
"identidad" se elige en la pantalla, pero la pertenencia al proyecto la valida el servicio, no la buena fe.

### Interfaz — `internal/server/planning_poker.go` + plantilla

`POST /planning-poker/votar` registra el voto con el formulario (integrante del desplegable + carta de la escala);
la fila de la tabla muestra, en vivo, "Votaron: …" con los nombres — nunca con los valores — hasta que alguien
revela (HU-17). En el Sprint 4 la pantalla recibió el pulido de HU-32.2: cartas como botones de la escala e insignias.

## 7. Para la demo en vivo (guion corto)

Con `make demo` sobre una carpeta `data/` limpia y `make run`, en `/planning-poker`:

1. **Abrir votación** en una historia sin estimar → fila en ronda 1, "Votaron: todavía nadie". *(RN-16.5, salidas)*
2. **Votar** como Ana la carta 5 y como Beto la carta 8, cada uno sin ver el otro → la fila muestra "Votaron:
   Ana Pérez, Beto Gómez" y nada más. *(RN-16.1, 16.2, 16.4)*
3. **Revelar** como Agile Enabler → tabla con los dos valores y, sin consenso, el más alto y el más bajo.
   *(Puente a HU-17)*
4. Si sobra tiempo: **la escala sin fisuras** — las cartas solo son 1, 2, 3, 5, 8, 13, 21: votar 4 no es posible
   ni en la pantalla ni en la API. *(RN-16.3)*

Cada paso de la demo mapea a una regla de la spec: eso es lo que se quiere mostrar en el recorrido — la
funcionalidad se comportó exactamente como dice la spec porque los tests la fijaron antes de existir.

## 8. Cómo llegó a Done

- Escenarios BDD en verde desde el PR, y siguen en verde en `main` (los pasos son del mundo compartido).
- Revisado y mergeado en el [PR #85](https://github.com/branko-utn/seguimiento-medicion/pull/85); el resto de la
  épica (HU-17, HU-18) y el pulido de HU-32.2 crecieron sobre el mismo diseño sin rehacer nada.
- Registrado en `docs/ia/registro.md` (fila del 27/10/2026).
