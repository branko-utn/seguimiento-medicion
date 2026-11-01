# Presentación final · Guion (TEC-22)

- **Cuándo:** lunes 2/11/2026, ante los Product Architects. **Ensayo:** domingo 1/11, con la app en una carpeta de
  datos limpia.
- **Preparan:** Branko Almeira (Agile Enabler): diapositivas, guion de la demo y ensayo. Cada integrante presenta
  su épica. Ignacio Bianchi hace el recorrido completo de HU-16
  ([`recorrido-hu-16.md`](recorrido-hu-16.md)).
- **Duración objetivo:** unos 25 minutos más preguntas. Los tiempos de abajo son orientativos: se ajustan en el ensayo
  si la cátedra fija otra duración.

## Reglas para el día

- **Cada uno defiende lo que hizo.** Si preguntan por una parte, responde su responsable (tabla de épicas abajo).
- **La demo no se improvisa:** se sigue el guion de la sección "Demo" sobre una carpeta recién cargada con `make demo`.
- **Plan B:** si la demo falla, se muestran las capturas del manual (`docs/manual/img/`) y el reporte del proyecto
  en PDF ya descargado.
- **Cifras:** se toman del dashboard con los datos de demo, el mismo día. Las de esta guía son las del 31/10.

## Diapositivas

| # | Diapositiva | Contenido | Presenta | Min |
|---|---|---|---|---|
| 1 | Portada | Software Metrics & Estimation · UTN · Ingeniería y Calidad de Software · el equipo con sus roles | Branko | 0,5 |
| 2 | El problema | Un equipo ágil necesita estimar, planificar, seguir el avance y medir la calidad en un solo lugar | Branko | 1 |
| 3 | Cómo trabajamos | Scrum con un Sprint 0 y 4 mini-sprints; tablero en GitHub Projects; SDD → BDD → TDD; revisión cruzada; IA registrada en `docs/ia/registro.md` | Branko | 2 |
| 4 | Arquitectura | Go 1.22 + `net/http` + `html/template` + HTMX + Chart.js; capas `domain → service → server`; ADRs 0001 a 0004 (diagrama del anexo A.2 del manual) | Branko | 2 |
| 5 | Proyectos e integrantes (E1) | Estado calculado por fechas, resumen del proyecto, roles | Federico Sosa | 1,5 |
| 6 | Backlog (E2) | Historias, criterios de aceptación, transiciones de estado / edición, eliminación y orden por prioridad | Ignacio Bianchi / Agustín Lara | 2 |
| 7 | Sprints (E3) | Sprint Goal, asignación, inicio, completar historias, cierre e historial | Federico Sosa (Branko: completar y cerrar; Giorlando: historial) | 1,5 |
| 8 | Planning Poker (E4) | Votos ocultos, revelar, diferencias, nueva ronda y estimación acordada | Ignacio Bianchi | 1,5 |
| 9 | Recorrido completo de HU-16 | Historia → spec → criterios → escenarios BDD → tests → código Go, con los commits RED → GREEN | Ignacio Bianchi | 3 |
| 10 | Esfuerzo y defectos (E5, E6) | Horas por integrante y actividad, defectos con severidad y su ciclo de estados | Agustín Lara | 2 |
| 11 | Métricas y dashboard (E7, E8) | SP planificados y completados, velocidad, % completado, desviación de horas, defectos; gráficos | Agustín Giorlando | 2 |
| 12 | Reportes (E9) | Reporte de sprint y de proyecto; exportación a PDF generada en Go | Agustín Giorlando (Branko: reportes) | 1,5 |
| — | **Demo en vivo** | Ver la sección "Demo" | Branko conduce; cada uno muestra su pantalla | 5 |
| 13 | Calidad | Tests, cobertura, BDD y CI (cifras abajo) | Agustín Lara | 1 |
| 14 | Métricas de nuestro propio proyecto | Velocidad, SP por sprint y defectos del equipo, sacados de la app | Agustín Giorlando | 1 |
| 15 | Retrospectiva | Qué funcionó, qué mejoraríamos (de las retros de los Sprints 1 a 3) | Branko | 1 |
| 16 | Cierre | Gracias · preguntas | Todos | — |

### Cifras para las diapositivas (al 31/10, a actualizar el 1/11)

| Dato | Valor | De dónde sale |
|---|---|---|
| Historias de usuario | 32 (HU-01 a HU-32) en 9 épicas | `specs/` |
| Escenarios BDD | 183 en 39 archivos `.feature` | `features/` |
| Funciones de test | 389 (más de 1.100 casos contando subtests) | `go test ./... -v` |
| Cobertura de `internal/domain` + `internal/metrics` | 98,6 % (umbral del CI: 80 %) | `go run ./tools/checkcoverage` |
| Pull Requests mergeados | 67 | GitHub |
| Commits RED de TDD | 120 | `git log --grep RED` |
| Velocidad del equipo (Sprints 1 a 3) | la que muestre el dashboard con `make demo` | Dashboard |

## Demo

La guía detallada de la demo (entorno, verificación de pantallas, qué decir en cada una con las cifras y cómo
dejar los datos como al principio) está en [`docs/demo.md`](../../demo.md), de Agustín Giorlando. Acá va el resumen.

**Preparación (antes de entrar):**

```bash
git switch main && git pull
DATA_DIR=demo-presentacion make demo   # carpeta limpia, solo para la presentación
DATA_DIR=demo-presentacion make run    # http://localhost:8080
```

Abrir el navegador en la portada, con el zoom al 125 % para que se lea desde lejos.

**Recorrido (5 minutos, una sola pestaña):**

| Paso | Pantalla | Qué mostrar | Quién |
|---|---|---|---|
| 1 | Portada | Las secciones como tarjetas; la navegación | Branko |
| 2 | Proyectos → el proyecto | Estado calculado (En curso), días restantes, sprint activo, % completadas | Sosa |
| 3 | Backlog | Insignias de prioridad y estado; abrir una historia y mostrar sus criterios | Bianchi |
| 4 | Planning Poker | Abrir una votación de una historia sin estimar, votar como dos integrantes, revelar, mostrar la diferencia y acordar | Bianchi |
| 5 | Sprints → Sprint 4 | Historias asignadas, métricas en tarjetas; completar una historia | Sosa / Branko |
| 6 | Esfuerzo | Registrar horas en esa historia; ver los totales por integrante | Lara |
| 7 | Defectos | Registrar un defecto y pasarlo a En progreso | Lara |
| 8 | Dashboard | Gráfico de velocidad y de SP por sprint, desviación de horas, defectos | Giorlando |
| 9 | Reportes | Reporte del proyecto y **Descargar PDF** | Giorlando |

**Si algo sale mal:** no se debuggea en vivo. Se dice qué se esperaba ver y se pasa a las capturas del plan B.

## Ensayo del domingo 1/11

- [ ] Diapositivas armadas y compartidas con el equipo (Branko).
- [ ] Cada integrante revisó la diapositiva de su épica y sabe qué dice.
- [ ] Recorrido de HU-16 practicado con los tiempos (Bianchi).
- [ ] Demo completa dos veces sobre una carpeta limpia, con los tiempos.
- [ ] Cifras actualizadas con el `main` final (después del code freeze del 31/10 a las 23:59).
- [ ] PDF del reporte del proyecto descargado y capturas listas como plan B.
- [ ] Preguntas probables repartidas (abajo).

## Preguntas probables

| Pregunta | Responde |
|---|---|
| ¿Cómo aseguraron que el TDD fue real? | Historial de commits RED → GREEN y la regla de no hacer *squash* (Bianchi, con HU-16) |
| ¿Por qué Go y una app web? ¿Por qué sin base de datos? | ADR 0001 (Branko) |
| ¿Cómo se calcula la velocidad / la desviación? | Giorlando |
| ¿Qué pasa con una historia que no se completó al cerrar el sprint? | Vuelve al backlog y no suma SP completados (Branko, HU-13) |
| ¿Cómo evitaron depender unos de otros? | Contratos por ID entre épicas (ADR 0002) y specs armadas para avanzar sin bloqueos (Branko) |
| ¿Cómo usaron la IA? | Registro en `docs/ia/registro.md`: qué se pidió, qué se usó y cómo se validó; cada uno explica su código (todos) |
| ¿Cómo se identifica quién vota sin login? | ADR 0004 (Bianchi) |
