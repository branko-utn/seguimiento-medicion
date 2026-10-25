# Escenarios BDD

Los `.feature` describen el comportamiento de cada historia (Dado / Cuando / Entonces) y se ejecutan con
[godog](https://github.com/cucumber/godog) en `go test ./features/`. Las specs y los escenarios ya están escritos;
**cada Builder escribe los pasos de sus historias** junto con su código (commit RED).

## Cómo implementar los pasos de una historia

1. Crear `features/huXX_steps_test.go` con una función `registrarPasosHUXX(sc *godog.ScenarioContext, m *mundo)` y
   llamarla desde `inicializarEscenario` (`features_test.go`).
2. Agregar el servicio de la épica a `mundo` (`mundo_test.go`) y crearlo en `nuevoMundo()` con repositorios en memoria
   (`memory.Open[...]("")`), si todavía no está.
3. Quitar `@pendiente` del `.feature` en el commit RED: los escenarios empiezan a fallar y pasan a verde con el código.
4. Seguir la convención del mundo:
   - **Dado**: prepara datos; si falla, devuelve el error y el escenario falla.
   - **Cuando**: ejecuta la acción y guarda su error con `m.resultado(err)`, sin fallar (puede ser el error esperado).
   - **Entonces**: verifica. Los pasos de resultado exitoso llaman primero a `m.sinError()`.
   - Si el escenario **ya verificó el error** con `veo el error`, los "Entonces" siguientes **no llaman a
     `sinError()`** (el error esperado ya fue verificado): por ejemplo `la historia … está en estado …`,
     `la historia … tiene N Story Points` o `el sprint … no tiene historias`.
5. **Nombres del escenario e IDs reales:** los escenarios usan nombres ("Demo", "S9", "HU-1"), pero los IDs los genera el
   repositorio. Al crear algo, guardar la relación con `m.recordar("sprint", "S9", string(sp.ID))`; al usarlo,
   traducir con `m.id("sprint", "S9")`. Un nombre sin alias se usa tal cual (sirve para "HU-99", que no existe).

## Pasos compartidos: quién los implementa

Godog falla si dos funciones definen el mismo paso (*paso ambiguo*). Los pasos que usan historias de **distintos
responsables** los implementa **una sola persona**, en `features/comun_<epica>_steps_test.go`, y los demás los
reutilizan. Si tu historia necesita uno que todavía no existe, avisá en la Daily: lo implementa su responsable o lo
implementás vos en ese archivo, pero no en el tuyo.

| Paso (texto en el `.feature`) | Responsable | Archivo | Lo usan |
|---|---|---|---|
| `veo el error "…"` | Branko (TEC-10) | `mundo_test.go` ✅ | Todas |
| `hoy es "AAAA-MM-DD"` | Branko (TEC-10) | `mundo_test.go` ✅ | HU-20 |
| `que existe un proyecto "…"` | Sosa | `comun_proyectos_steps_test.go` ✅ | Todas |
| `que existe un proyecto "…" del "…" al "…"` | Sosa | `comun_proyectos_steps_test.go` ✅ | HU-10 |
| `que "…" es integrante del proyecto "…"` | Sosa | `comun_proyectos_steps_test.go` ✅ | HU-03, HU-20 |
| `existe un sprint "…" en "…"` · `existe un sprint activo "…" en "…"` | Sosa | `comun_sprints_steps_test.go` ✅ | HU-08, HU-11, HU-22 |
| `que la historia "…" está asignada al sprint "…"` | Sosa | `comun_sprints_steps_test.go` ✅ | HU-08, HU-11 |
| `que el backlog de "…" tiene la historia "…" "…"` | Bianchi | `comun_backlog_steps_test.go` ✅ | HU-07, HU-08, HU-15, HU-19, HU-20, HU-22 |
| `el proyecto tiene las historias:` (tabla; `criterios` opcional, vale 1) | Bianchi | `comun_backlog_steps_test.go` ✅ | HU-11, HU-12 |
| `la historia "…" está en estado "…"` (Dado y Entonces) | Bianchi | `comun_backlog_steps_test.go` ✅ | HU-05, HU-08, HU-12 |
| `la historia "…" tiene N Story Points` (Dado y Entonces) | Bianchi | `comun_backlog_steps_test.go` ✅ | HU-15 |
| `que la historia "…" ya está completada en el sprint "…"` | Branko | `comun_sprints_cierre_steps_test.go` ✅ | HU-08, HU-12 |
| `que la historia "…" está completada` | Branko | `comun_sprints_cierre_steps_test.go` ✅ | HU-15, HU-22 |
| `que el sprint "…" está cerrado` | Branko | `comun_sprints_cierre_steps_test.go` ✅ | HU-11, HU-12 |

Valores por defecto de los pasos de proyectos y sprints (cuando el escenario no los indica):

- `que existe un proyecto "…"`: del 2026-10-21 al 2026-11-02 (el período del TP).
- `que "…" es integrante del proyecto "…"`: rol Product Builder.
- `existe un sprint "…" en "…"` (y `activo`): 3 días a partir del día siguiente al último sprint del proyecto (el
  primero empieza el día de inicio del proyecto), con el Sprint Goal "Sprint de prueba". El nombre ("S1", "S9") queda
  como alias: se usa con `m.id("sprint", "S9")`.

Los demás pasos aparecen en una sola historia: los implementa su responsable en su propio archivo.

Los pasos "Dado" de backlog fijan Story Points, criterios y estado directamente en el repositorio (`m.historias`): son preparación del escenario, no los casos de uso de HU-07, HU-08 ni HU-15.

**Orden:** los pasos de Sosa y Bianchi desbloquean a los demás, igual que sus tipos (ADR 0002, sección 6).
