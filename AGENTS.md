# Instrucciones para agentes de IA

Este archivo es para cualquier asistente de IA (Claude, Codex, Cursor, Copilot, Gemini, ChatGPT…) que trabaje en el
repositorio. Las personas también lo pueden leer: resume las reglas de [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Contexto

Trabajo Práctico Integrador **Software Metrics & Estimation** (Ingeniería y Calidad de Software, UTN FRSR, 2026):
una aplicación web en Go para estimar, planificar y medir proyectos de software. Cinco personas trabajan en paralelo
en sprints de 2–3 días; la presentación final es el **lunes 2/11/2026**. Se aplica **SDD** (specs antes del código),
**BDD** (escenarios Dado/Cuando/Entonces con godog) y **TDD** (RED → GREEN → REFACTOR). Los profesores evalúan el
proceso mirando el historial de Git, tanto como el producto.

**La persona con la que trabajás es responsable de todo lo que subas.** Antes de cada commit explicale qué hiciste
y por qué, mostrale el diff y esperá su confirmación. Si no lo entiende, no se commitea.

## Qué leer antes de escribir código

1. [`docs/adr/0002-contratos-de-dominio.md`](docs/adr/0002-contratos-de-dominio.md): tipos de cada épica, dependencias
   permitidas entre paquetes, convención de errores y orden de merge. **Obligatorio.**
2. [`CONTRIBUTING.md`](CONTRIBUTING.md): ramas, commits, Pull Requests, estructura del código, cómo agregar una pantalla.
3. [`features/README.md`](features/README.md): cómo escribir los pasos BDD y **qué pasos compartidos implementa cada uno**.
4. [`docs/scrum/definition-of-ready-done.md`](docs/scrum/definition-of-ready-done.md): Definition of Done.
5. **El issue de la historia** (responsable, sprint, criterios) y su **spec** y **.feature**:
   `specs/HU-XX-*.md`, `features/hu-xx-*.feature`.
6. Código de referencia: `internal/store/memory/store.go` (persistencia), `internal/server/` (handlers),
   `features/mundo_test.go` (mundo BDD).

## Reglas

- **Go 1.22**, solo biblioteca estándar más lo que ya está en `go.mod`. No agregues dependencias.
- **Identificadores del dominio en español**, como en el ADR 0002 (`Proyecto`, `Nuevo`, `ErrNombreInvalido`).
- **Trabajá solo en los paquetes de la épica de tu persona** (ver la tabla del ADR 0002). Lo de otra épica se
  referencia por ID (`internal/domain/ids.go`) o se importa según las dependencias permitidas. No modifiques
  paquetes de otra persona ni los tipos del ADR.
- **La spec es la fuente de verdad.** Implementá exactamente sus reglas (`RN-XX.n`), casos límite y errores.
  - Los **mensajes de error son idénticos** a los de la spec: los escenarios BDD los comparan textualmente.
  - **No agregues funcionalidades** que la spec no pide.
  - **Si algo es ambiguo, no lo decidas:** frená y preguntá. Las specs las mantiene el Agile Enabler.
  - No edites specs ni `.feature`, salvo para quitar la etiqueta `@pendiente`.
- **Errores:** variables exportadas por paquete (`var ErrX = errors.New("mensaje de la spec")`), comparadas en los
  tests con `errors.Is`. **Constructores que validan:** `Nuevo(...) (Tipo, error)`.
- **Fechas:** `time.Time` sin hora. Si una regla depende de "hoy", la fecha entra como parámetro: no uses
  `time.Now()` en el dominio.
- **Persistencia:** el servicio declara la interfaz del repositorio que necesita; la implementación va en
  `internal/store/memory` con `memory.Store[T]`. Los IDs los genera el repositorio (`P1`, `I1`, `HU-1`, `S1`, `E1`, `D1`).
- **Operaciones que tocan dos épicas:** cada entidad valida sus reglas en su paquete y el servicio coordina y guarda
  ambas (ADR 0002, sección 4).
- **Antes de cada commit:** `gofmt -w . && go vet ./... && go test ./...`. Única excepción: el commit RED, donde los
  tests nuevos fallan a propósito.

## Flujo de cada historia (un commit por fase)

1. **Rama** desde `main` actualizado: `git switch main && git pull && git switch -c feature/HU-XX-titulo-corto`.
   La persona mueve el issue a *In progress*.
2. **Dominio con TDD** (`internal/domain/<paquete>/`):
   1. Tests de tabla que cubran cada regla, caso límite y error de la spec, con los stubs mínimos para compilar.
      Verificá que fallen. → `test(HU-XX): RED - <qué se prueba>`
   2. El código mínimo para que pasen. → `feat(HU-XX): GREEN - <qué se implementó>`
   3. Mejoras sin cambiar el comportamiento, si hacen falta. → `refactor(HU-XX): <qué mejoró>`
3. **Servicio y repositorio** (`internal/service/<epica>.go`, `internal/store/memory/<epica>.go`): el mismo ciclo.
4. **Pasos BDD** (`features/huXX_steps_test.go`), siguiendo [`features/README.md`](features/README.md):
   - registrar `registrarPasosHUXX(sc, m)` en `inicializarEscenario` y agregar el servicio al `mundo`;
   - "Dado" prepara; "Cuando" guarda el error con `m.resultado(err)`; "Entonces" verifica;
   - nombres del escenario ↔ IDs reales con `m.recordar(...)` y `m.id(...)`;
   - **pasos compartidos:** solo los que la tabla asigna a tu persona, en `features/comun_<epica>_steps_test.go`.
     Nunca dupliques uno que ya existe: godog falla por paso ambiguo. Si falta uno de otra persona, frená y avisá;
   - quitar `@pendiente` del `.feature`: `test(HU-XX): RED - pasos BDD` → `feat(HU-XX): GREEN - escenarios BDD`.
5. **Pantallas web:** desde el Sprint 2, salvo que la persona lo pida (`CONTRIBUTING.md`, sección 6).
6. **Cierre:**
   - `go test -race ./...` en verde;
   - push y Pull Request con la plantilla completa: `Closes #N`, tabla de trazabilidad y checklist de la DoD;
   - una fila en [`docs/ia/registro.md`](docs/ia/registro.md): fecha, persona, herramienta, qué se pidió, qué se usó,
     cómo se validó.

## Nunca

- Push directo a `main`, force push o reescribir commits ya subidos.
- Commits que mezclen fases (tests y código juntos) o que no compilen, salvo el RED a propósito.
- Decidir lo que la spec no dice, cambiar mensajes de error o tocar paquetes de otra épica.
- Subir código que la persona no haya leído y entendido.

## Quién hace qué en el Sprint 1

| Persona (GitHub) | Historias (issue) | Paquetes | Notas |
|---|---|---|---|
| Federico Sosa (`sosa-federico`) | HU-01 (#8), HU-03 (#10), HU-10 (#17), HU-11 (#18) | `domain/project`, `domain/sprint` | **Prioridad:** tipos `Proyecto`, `Integrante` y `Sprint` mergeados el sábado 24 a las 12:00 |
| Ignacio Bianchi (`Igna-Bianchi`) | HU-05 (#12), HU-07 (#14), HU-15 (#22), HU-08 (#15) | `domain/backlog`, `domain/estimation` | **Prioridad:** tipo `Historia` mergeado el sábado 24 a las 12:00 |
| Agustín Lara (`LaraAgustin`) | HU-19 (#26), HU-20 (#27), HU-22 (#29) | `domain/effort`, `domain/defect` | Dominio ya; servicio y pasos BDD cuando estén los tipos de Sosa |
| Agustín Giorlando (`AgustinGiorlando`) | HU-24 (#31), HU-25 (#32) | `internal/metrics` | Funciones puras, sin dependencias: arranca ya |
| Branko Almeira (`branko-utn`) | HU-12 (#19) | `domain/sprint` (con Sosa) | Cuando estén `Historia` y `Sprint` |
