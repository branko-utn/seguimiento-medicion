# Informe de métricas y cobertura de pruebas

> TEC-08 · Responsable: Agustín Lara. Los números salen de la app con `make demo` (reporte de proyecto y pantalla de
> defectos), de las actas de cada sprint (`docs/scrum/`) y de `make cover`. Cada tabla dice de dónde sale cada dato.
> Medido el 01/11/2026 sobre `main` en el commit `fdd25fa`, con el Sprint 4 cerrando: todo mergeado salvo el manual
> (TEC-07) y la presentación (TEC-22), que terminan con la presentación del 2/11.

## 1. Resumen

- **Velocidad estable:** 44, 41 y 41 Story Points completados en los Sprints 1 a 3 según el tablero, con un promedio
  de 42. En cada uno se completó todo lo comprometido. El Sprint 4 tiene 26 de 34 SP mergeados al 01/11, y lo que
  falta (manual y presentación) se cierra el 2/11.
- **Historias:** la app cuenta 102 Story Points de historias completados en los Sprints 1 a 3 (velocidad de 34,0) y
  30 de 31 historias terminadas (96,8 %). La diferencia con el tablero son los ítems técnicos (TEC), que la app no
  cuenta.
- **Horas:** 214,0 estimadas y 227,7 reales, con una desviación de 6,4 %. Las horas de la demo son **simuladas**
  (sección 2), así que muestran que el cálculo funciona, no cómo estimamos.
- **Defectos:** las reviews registran 16 defectos, todos encontrados en la revisión de PRs, y **los 16 están
  resueltos**. El último, D9 (punto decimal en el dashboard), se resolvió en el Sprint 4 con HU-32.4 (#134). La demo
  carga 9 de ellos (1 crítico).
- **Cobertura:** 88,7 % de las sentencias de todo el código. En `internal/domain` + `internal/metrics`, la que exige
  el CI, es de 98,6 % contra un umbral de 80 %. 218 escenarios BDD en verde y ninguno pendiente.
- **Despliegue:** la app está publicada en Azure (TEC-23), compilada con Go 1.26.10, sin vulnerabilidades conocidas
  según `govulncheck`, y con los controles de seguridad de TEC-24 (sección 6).

## 2. Métricas por sprint

| Sprint | SP planificados | SP completados | % historias | Horas estimadas | Horas reales | Defectos detectados / resueltos |
|---|---|---|---|---|---|---|
| 1 · MVP | 44 | 44 | 100,0 % | 88,0 | 92,1 | 3 / 3 |
| 2 · Interfaz | 23 | 23 | 100,0 % | 46,0 | 50,8 | 2 / 2 |
| 3 · Funcionalidad y calidad | 35 | 35 | 100,0 % | 70,0 | 82,8 | 4 / 3 |
| 4 · Cierre (en curso) | 5 | 0 | 0,0 % | 10,0 | 2,0 | 0 / 0 |
| **Proyecto** | | **102** | **96,8 %** (30 de 31) | **214,0** | **227,7** | **9 / 8** |

Fuente: reporte de proyecto (`/reportes/proyecto?proyecto=P1`) con los datos de `make demo`. Hay que tener en cuenta:

- **La app cuenta solo historias (HU).** El tablero también cuenta los ítems técnicos (TEC), por eso la velocidad del
  tablero (sección 3) es mayor en los Sprints 2 y 3.
- **Las horas son simuladas.** No se registraron durante el proyecto: `make demo` carga las estimadas a 2 horas por
  Story Point y las reales con una variación por historia (ver el README y `cmd/demo/datos.go`).
- **Los totales del proyecto suman todo el esfuerzo**, también el registrado fuera de las fechas de un sprint. En la
  demo no hay horas fuera de los sprints, así que coinciden con la suma de las filas: 214,0 y 227,7.
- **El Sprint 4 sale parcial en la app:** la demo carga los datos como estaban al empezar el sprint, con una sola
  historia, HU-31 (5 SP), sin completar. Por eso también figura abierto D9, que se resolvió en el Sprint 4. Lo que
  pasó en el Sprint 4 está en las secciones 3 y 4.

Desviación de horas por sprint, calculada como (reales − estimadas) / estimadas:

| Sprint | Desviación |
|---|---|
| 1 | 4,7 % |
| 2 | 10,4 % |
| 3 | 18,3 % |
| Proyecto | 6,4 % |

## 3. Velocidad y estimación

| Sprint | Duración | SP comprometidos en la planning | SP completados (tablero) | SP de historias (app) |
|---|---|---|---|---|
| 1 · MVP | 3 días | 44 | 44 | 44 |
| 2 · Interfaz | 2 días | 38 | 41 | 23 |
| 3 · Funcionalidad y calidad | 2 días | 41 (36 firmes + 5 recortables) | 41 | 35 |
| 4 · Cierre | 1 día | 34 (13 iniciales + 13 de la interfaz + 8 del despliegue) | 26 al 01/11 (faltan TEC-07 y TEC-22) | sin datos en la demo |

Fuente: `planning.md` y `review.md` de cada sprint en `docs/scrum/`. La velocidad promedio del tablero es **42 SP**
en los Sprints 1 a 3. La de la app es **34,0 SP** porque solo suma historias.

- **Sprint 1:** sin velocidad de referencia. Se comprometieron 44 SP y se completaron todos.
- **Sprint 2:** con la velocidad del Sprint 1 se estimaron unos 30 SP para un sprint de 2 días, pero se comprometieron
  38. Se completaron 41: TEC-19 (las specs del Sprint 2, 3 SP) se hizo antes de la planning y se cargó en el tablero
  durante el sprint.
- **Sprint 3:** se comprometieron 36 SP firmes más 5 recortables (HU-09 y HU-30). Entraron los dos recortables y se
  completaron los 41.
- **Sprint 4:** al ser de un día y de cierre, se comprometieron 13 SP. Dos cambios de alcance lo llevaron a 34: la
  mejora de la interfaz el 31/10 (13 SP: TEC-21, las cuatro partes de HU-32 y TEC-22) y el despliegue en Azure que
  pidió la cátedra el 01/11 (8 SP: TEC-23 y TEC-24). Al 01/11 están mergeados 26 SP. Los otros 8 son el manual
  (TEC-07, 5 SP) y la presentación (TEC-22, 3 SP), que se cierran con la presentación del 2/11.

En los Sprints 1 a 3, lo comprometido siempre se cumplió, y en el Sprint 2 se superó. La planning de cada sprint usó
la velocidad anterior como referencia, y el reparto buscó que los cuatro Product Builders terminaran el proyecto con
una carga pareja.

## 4. Calidad: defectos

**Defectos registrados en las reviews**, todos encontrados en la revisión de PRs:

| Sprint | Detectados | Resueltos | Fuente |
|---|---|---|---|
| 1 | 5 | 5 | `docs/scrum/sprint-1/review.md` |
| 2 | 4 | 4 | `docs/scrum/sprint-2/review.md` |
| 3 | 7 | 6 | `docs/scrum/sprint-3/review.md` |
| 4 | — | 1 (D9, del Sprint 3) | HU-32.4 (#134) |
| **Total** | **16** | **16** | |

**Defectos cargados en la app** (`make demo`, pantalla `/defectos`). Son 9 de los 16 de las reviews:

| Sprint de detección | Defectos | Severidad |
|---|---|---|
| 1 | D1, D2, D3 | Media, Alta, Media |
| 2 | D4, D5 | Alta, Media |
| 3 | D6, D7, D8, D9 | Alta, Crítica, Media, Baja |

| Severidad | Cantidad |
|---|---|
| Crítica | 1 (D7: acordar la estimación no cerraba la votación del Planning Poker) |
| Alta | 3 |
| Media | 4 |
| Baja | 1 (D9: el dashboard mostraba punto decimal; se resolvió en el Sprint 4 con HU-32.4, pero la demo lo carga abierto) |

**De dónde salieron:** todos aparecieron en la revisión de PRs y en las auditorías del Agile Enabler. En los
Sprints 1 y 2, todos se corrigieron antes del merge. En el Sprint 3 quedó uno abierto, D9, que se resolvió en el
Sprint 4 con HU-32.4 (#134).
Algunos casos que se repitieron:
- merges que pisaban registros de otras historias (D2);
- tests que dependían de la fecha real (D6);
- reglas de la spec que no estaban cubiertas por un escenario (D3).

## 5. Cobertura de pruebas

Cobertura de sentencias por paquete (`go test -cover ./...`):

| Paquete | Cobertura |
|---|---|
| `internal/domain/defect` | 100,0 % |
| `internal/domain/effort` | 100,0 % |
| `internal/domain/project` | 100,0 % |
| `internal/domain/sprint` | 100,0 % |
| `internal/metrics` | 100,0 % |
| `internal/domain/estimation` | 97,2 % |
| `internal/domain/backlog` | 96,4 % |
| `internal/report/pdf` | 97,5 % |
| `internal/report` | 93,7 % |
| `internal/service` | 93,2 % |
| `internal/server` | 87,2 % |
| `internal/store/memory` | 87,2 % |
| `internal/app` | 72,7 % |
| `cmd/demo` | 65,8 % |
| `tools/checkcoverage` | 52,8 % |
| `cmd/server` | 0,0 % (solo arranca el servidor) |
| `internal/dashboard` | sin tests propios (lo prueban los tests de `internal/server` y los escenarios de HU-28) |
| **Total** | **88,7 %** |

**Umbral del CI (TEC-05):** `internal/domain` + `internal/metrics` tienen 346 de 351 sentencias cubiertas, es decir,
**98,6 %**, contra un mínimo de 80 %. Si baja de ese umbral, el CI falla.

**Cantidad de pruebas:**
- **429 funciones de test:** 412 unitarias en `internal/`, `cmd/` y `tools/`, y 17 en `features/`. Muchas son
  tablas con varios casos cada una.
- **218 escenarios BDD** (1214 pasos) en 38 archivos `.feature`, todos en verde (`make bdd`). Ninguno sigue con
  `@pendiente`: la etiqueta solo queda en la plantilla (`features/_plantilla.feature`).

Fuente: `make cover` (`go tool cover -func`), `go run ./tools/checkcoverage` y `go test ./features -v`.

## 6. Despliegue y seguridad

La app de `main` está publicada en Azure Container Apps (TEC-23, ADR 0005) y se vuelve a desplegar en cada merge que
cambia código. Leer es público; escribir pide una clave.

| Control (TEC-24) | Resultado |
|---|---|
| Vulnerabilidades conocidas | La imagen se compila con **Go 1.26.10**. `govulncheck` no encuentra ninguna, y si encuentra alguna el workflow no publica. Con Go 1.22.2 eran 45 de la biblioteca estándar |
| Escrituras | Un proxy Caddy pide usuario y clave para todo lo que no sea lectura: sin clave, 401 |
| Transporte | Solo HTTPS: `http://` redirige a `https://` (301), con `Strict-Transport-Security` |
| Cabeceras | CSP sin `'unsafe-inline'` (el único script y el único estilo en línea se autorizan por hash), `nosniff`, `DENY`, `Referrer-Policy` y sin `Server` |
| Pedidos | Hasta 1 MB y con tiempo límite de lectura |
| Imagen | Distroless, sin shell y sin root. Las imágenes base y las acciones del workflow se fijan por hash, y se despliega por digest |
| Secretos | Solo en el Environment `azure` de GitHub. `git grep` no encuentra credenciales en el repo |
| Costo | Una réplica como máximo, y escala a cero sin tráfico |

Fuente: corrida 41349200582 del workflow `cd-azure.yml` (prueba de humo y pruebas de seguridad en verde) y
verificación desde afuera con `curl` y Chrome sobre la URL pública. El detalle está en
`specs/TEC-24-seguridad-del-despliegue.md` y en `docs/despliegue.md`.

## 7. Conclusiones

- **La velocidad fue estable y predecible.** Con 41 a 44 SP por sprint y todo lo comprometido completado, la
  planning pudo usar la velocidad anterior para comprometer cada sprint. Los ítems recortables (HU-09 y HU-30 en el
  Sprint 3) dieron margen sin poner en riesgo el Sprint Goal, y al final entraron.
- **La revisión de PRs fue el principal filtro de calidad.** Los 16 defectos se encontraron ahí, y se corrigieron
  todos. Sostenerla costó tiempo de auditoría, pero los errores se encontraron antes de la demo.
- **El trabajo con SDD → BDD → TDD se nota en la cobertura.** El dominio y las métricas, escritos test primero, están
  entre 96 % y 100 %. Lo menos cubierto es lo que no tiene reglas de negocio: el arranque del servidor, el armado de
  la aplicación y las herramientas.
- **El despliegue se aseguró sin tocar el código congelado.** El proxy, la imagen y el workflow suman los controles
  que la app no tenía (clave para escribir, cabeceras, límites), y `govulncheck` mostró que actualizar Go era
  condición para publicar.
- **Qué haríamos distinto:**
  - **Registrar las horas reales durante el proyecto**, no reconstruirlas al final. Con horas simuladas, la
    desviación muestra que el cálculo funciona, pero no nos dice nada de cómo estimamos.
  - **Cargar en la app todos los defectos al momento de encontrarlos.** Así la app y las actas tendrían los mismos
    números (hoy 9 contra 16).
  - **Resolver antes los conflictos en los archivos compartidos** (el registro de IA, el wiring de los servicios).
    Fue la acción pendiente de tres retros seguidas.
