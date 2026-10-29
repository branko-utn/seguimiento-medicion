# Sprint 3 · Planning

- **Fecha:** jueves 29/10/2026 · **Duración:** 25'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Sprint:** Funcionalidad y calidad · jue 29/10 – vie 30/10

## Sprint Goal

> Ver el estado del proyecto de un vistazo en el dashboard y generar los reportes de sprint y de proyecto, con el
> Planning Poker completo.

## Fechas

El Sprint 2 cerró el 28/10, un día tarde (ver su retro). Para no perder un día del Sprint 3, el equipo decide:

| Sprint | Antes | Ahora |
|---|---|---|
| 3 · Funcionalidad y calidad | mié 28 – jue 29/10 | **jue 29 – vie 30/10** |
| 4 · Cierre | vie 30 – sáb 31/10 | **sáb 31/10** (code freeze sin cambios: sáb 31 a las 23:59) |

El Sprint 4 queda para la exportación a PDF, la documentación y el informe; el ensayo de la demo sigue el domingo 1/11.

## Criterios del reparto

1. **Esfuerzo parejo entre los Product Builders al terminar el proyecto.** Giorlando recupera en este sprint los SP
   que tuvo de menos en los Sprints 1 y 2.
2. **Cada historia depende solo de código que ya está en `main`.** Lo que necesitan varias historias de distintas
   personas entra antes, en TEC-20 (Branko). Si una historia depende de otra, es del mismo responsable.
3. **Cada uno hace lo suyo de punta a punta:** dominio, servicio, escenarios BDD y pantalla.

## Sprint Backlog

| Responsable | Ítem | SP | Depende de |
|---|---|---|---|
| Agustín Giorlando | HU-14 · Consultar sprints anteriores (#21) | 3 | `main` |
| | HU-28 · Dashboard del proyecto (#35), en dos PRs: resumen y métricas, después gráficos | 8 | `main` + TEC-20 |
| Agustín Lara | HU-06 · Editar y eliminar ítem del backlog (#13) | 3 | `main` |
| | HU-21 · Consultar el esfuerzo registrado (#28) | 3 | `main` |
| | HU-09 · Ver el backlog ordenado por prioridad (#16) — **recortable** | 2 | `main` |
| Ignacio Bianchi | HU-18 · Nueva ronda y estimación acordada (#25) | 5 | `main` |
| | TEC-05 · Umbral de cobertura en el CI (#39) | 2 | `main` |
| Federico Sosa | HU-04 · Consultar estado del proyecto (#11) | 3 | `main` + TEC-20 |
| Branko Almeira | TEC-20 · Base compartida del Sprint 3 (#99, **entra primero**) | 1 | — |
| | HU-29 · Reporte de sprint (#36) | 5 | `main` + TEC-20 |
| | HU-30 · Reporte de proyecto (#37) — **recortable** (pasa al Sprint 4) | 3 | HU-29 (propia) |
| Todo el equipo | TEC-06 · Datos de demo: cargar nuestro propio proyecto (#40) | 3 | todo lo anterior |
| | **Total** | **41** | |

Velocidad: 44 SP en el Sprint 1 y 41 en el Sprint 2 (promedio 42,5). Se comprometen **36 SP** firmes más los dos
recortables (HU-09 y HU-30). Si en la daily del viernes un ítem no tiene PR abierto, se recorta en este orden: HU-09,
HU-30 y los gráficos de HU-28.

## Carga acumulada

| Integrante | Sprints 1 y 2 | Sprint 3 | Sprint 4 (previsto) | **Total** |
|---|---|---|---|---|
| Federico Sosa | 24 | 3 | — | **27** |
| Ignacio Bianchi | 19 | 7 | — | **26** |
| Agustín Lara | 17 | 8 | TEC-08 · Informe de métricas (3) | **28** |
| Agustín Giorlando | 13 | 11 | — | **24** |
| Branko Almeira | 20 | 9 | HU-31 · Exportar a PDF (5) | 34 (+ Sprint 0 y facilitación) |

Sin contar los ítems de todo el equipo (TEC-06 y TEC-07). Los cuatro Builders quedan entre 24 y 28 SP. En la planning
del Sprint 4 se puede pasar TEC-08 de Lara a Giorlando si hace falta emparejar más.

**Cambios respecto del tablero:** TEC-05 pasa de Branko a Bianchi (Bianchi tenía 5 SP en el sprint y el umbral de
cobertura no depende de nadie). HU-09 sigue con Lara, que también tiene HU-06: las dos tocan el backlog y así no hay
conflictos entre personas.

## Decisiones

- **Specs aprobadas** con este PR, una por historia, cada una sin dependencias de historias de otros. Las preguntas de
  los borradores quedaron resueltas en la sección "Decisiones de la planning" de cada spec.
- **TEC-20 entra primero:** `Proyecto.EstadoAl`, `Proyecto.DiasRestantes` y `metrics.HorasEnPeriodo`.
- **Estado del proyecto calculado con las fechas** (HU-04); HU-02 pasa a usar lo mismo.
- **HU-06:** tampoco se elimina una historia con esfuerzo o defectos registrados.
- **HU-18:** la votación se cierra eliminando la votación abierta; no hay un estado nuevo.
- **HU-28:** el sprint activo no entra en los gráficos; Chart.js por CDN (ADR 0001).
- **Acciones de la retro del Sprint 2:** las dailies se vuelcan en `dailies.md` y el corte se hace en la daily del
  viernes.

## Riesgos

- **Giorlando tiene la carga más alta del sprint (11 SP)** y el ítem más visible de la demo. Mitigación: HU-28 en dos
  PRs, HU-14 primero (es la más chica) y revisión del avance en la daily del viernes; si los gráficos no están, se
  reasigna una parte en ese momento.
- **Un día menos para el Sprint 4.** Mitigación: HU-30 queda recortable al Sprint 4 y la documentación se escribe en
  paralelo desde el viernes.
- **TEC-06 depende de todo.** Se carga el sábado a la mañana con lo que esté en `main`.
