# Sprint 4 · Dailies

Formato por integrante: **hice · voy a hacer · me bloquea**. El Agile Enabler la vuelca el mismo día.

## Sábado 31/10

| Integrante | Hice | Voy a hacer | Me bloquea |
|---|---|---|---|
| Sosa | Correcciones de HU-04; revisión de la planning | TEC-21 (base visual) primero, después HU-32.1 y su sección del manual | — |
| Bianchi | Correcciones de HU-18 | HU-32.2 (backlog y Planning Poker), su sección del manual y el recorrido de HU-16 para la presentación | Esperar TEC-21 |
| Lara | HU-09 | TEC-08 (informe de métricas), HU-32.3 y sus secciones del manual | Esperar TEC-21 para HU-32.3 |
| Giorlando | HU-14 y HU-28 | HU-31 (PDF), después HU-32.4 y su sección del manual | Esperar TEC-21 para HU-32.4 |
| Almeira | Auditorías, TEC-06 y actas del Sprint 3 | Facilitar, auditar y mergear; anexo técnico del manual; diapositivas, guion de la demo y ensayo de la presentación. No ejecuta ítems de producto | — |

**Acciones sobre los bloqueos:** Sosa sube TEC-21 apenas se aprueben las specs. Hasta entonces, cada uno avanza con
su sección del manual (TEC-07). HU-31 tiene corte a las 18:00: sin PR en borrador, pasa a Lara.

## Domingo 1/11

> Reconstruida por el Agile Enabler a partir de los commits, los PRs y las revisiones del día: la daily no se volcó
> en el momento. Como no hay registro de lo que cada uno dijo, la tabla no usa "voy a hacer": muestra lo que hizo
> cada integrante el día anterior y el día de hoy, según el repositorio.

| Integrante | Sábado 31/10 | Domingo 1/11 (commits y PRs) | Bloqueos |
|---|---|---|---|
| Sosa | HU-32.1 (interfaz de proyectos y sprints) y su sección del manual | TEC-23 · Despliegue en Azure (#138): Dockerfile, `containerapp.yaml`, workflow `cd-azure.yml` y `docs/despliegue.md` (PR #142) | 19:52: el entorno `utnapp-env` salió en modo express (abajo) |
| Bianchi | HU-32.2 (backlog y Planning Poker) y el recorrido de HU-16 (#127) | Mergeado #127; revisión y aprobación de #142 | — |
| Lara | HU-31 (PDF, #132, tomada en el corte de las 18:00), TEC-08 (#133), HU-32.3 (#129) y sus secciones del manual (#131) | TEC-24 · Seguridad del despliegue (#139): proxy Caddy, `govulncheck` y pruebas de seguridad (PR #142) | — |
| Giorlando | HU-32.4 (#134) | Sección del manual del dashboard (#135) y guía de la demo (#136), abiertas la noche del 31/10 | — |
| Almeira | Guion de la presentación (#130), §1 y anexo del manual (#128), auditorías | Auditoría y merge de los 11 PRs del freeze (#127 a #136); estado de las specs (#137); specs del despliegue y ADR 0005 (#140, #141); secretos de Azure; auditoría de #142 | — |

**Durante el día:**

- Se mergearon los 11 PRs que quedaron del code freeze, después de la auditoría. HU-31 sumó una corrección del
  Agile Enabler (caracteres fuera de cp1252 en el PDF, con su test RED → GREEN).
- La cátedra pidió desplegar la app y entregó una cuenta de Azure compartida. Se sumaron TEC-23 (Sosa) y TEC-24
  (Lara): ver la planning, "Cambio de alcance: despliegue en Azure".
- **Bloqueo de Sosa:** el entorno `utnapp-env` salió en modo express, que no admite init containers ni sidecars. Se
  decidió usar un entorno propio, `grupo02-env` (comentario en #138). Esa noche se mergeó #142 y la app quedó
  desplegada.
