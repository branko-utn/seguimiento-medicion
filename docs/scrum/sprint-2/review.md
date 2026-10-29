# Sprint 2 · Review

- **Fecha:** miércoles 28/10/2026 por la noche · **Duración:** 20'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando. No participó
  ningún Product Architect.
- **Sprint:** lun 26/10 – mar 27/10 (La Interfaz). La review se corrió un día porque TEC-16 y TEC-17 entraron el 28/10
  (ver la retro).

## Sprint Goal

> Usar la aplicación desde el navegador, recorriendo proyecto → backlog → estimación → sprint → métricas, y estimar
> con Planning Poker.

**¿Se cumplió?** Sí. En la demo se recorrió todo el flujo en el navegador, con datos que quedan guardados en disco:

1. crear el proyecto y sus integrantes;
2. cargar las historias con sus criterios de aceptación;
3. estimarlas con Planning Poker;
4. armar e iniciar el sprint;
5. completar historias y cerrar el sprint;
6. ver las métricas del sprint.

## Incremento demostrado

| Issue | Ítem | SP | Estado | Quién lo mostró |
|---|---|---|---|---|
| #74 | TEC-13 · Pantallas de proyectos e integrantes (PR #87) | 3 | Done | Federico Sosa |
| #75 | TEC-14 · Pantalla de sprints: crear, iniciar y asignar (PR #89) | 3 | Done | Federico Sosa |
| #9 | HU-02 · Modificar proyecto (PR #84) | 2 | Done | Federico Sosa |
| #23 | HU-16 · Votar en Planning Poker con votos ocultos (PR #85) | 5 | Done | Ignacio Bianchi |
| #24 | HU-17 · Revelar votos y detectar diferencias (PR #92) | 3 | Done | Ignacio Bianchi |
| #76 | TEC-15 · Pantallas de esfuerzo y defectos (PR #90) | 3 | Done | Agustín Lara |
| #30 | HU-23 · Cambiar estado de un defecto (PR #86) | 3 | Done | Agustín Lara |
| #34 | HU-27 · Defectos detectados y resueltos (PR #88) | 2 | Done | Agustín Lara |
| #33 | HU-26 · Horas estimadas, reales y desviación (PR #91) | 3 | Done | Agustín Giorlando |
| #77 | TEC-16 · Pantalla del backlog (PR #96) | 3 | Done | Agustín Giorlando y Branko Almeira |
| #78 | TEC-17 · Vista de métricas del sprint (PR #97) | 2 | Done | Branko Almeira |
| #79 | TEC-18 · Conectar los servicios a la app web (PR #80) | 1 | Done | Branko Almeira |
| #20 | HU-13 · Cerrar sprint (PRs #81, #82 y #93) | 5 | Done | Branko Almeira |
| #95 | TEC-19 · Specs SDD y escenarios BDD del Sprint 2 (PRs #66 y #94) | 3 | Done | Branko Almeira |

- **TEC-19:** el trabajo se hizo antes de la planning (PR #66). Se cargó en el tablero el 27/10 para que la
  facilitación quedara visible, y por eso el sprint pasa de 38 a 41 SP.
- **TEC-16 (backlog):** Agustín Giorlando hizo la lista de historias, y Branko completó el alta y el detalle de cada
  historia (criterios, estimación y estado).
- **TEC-17 (métricas):** se reasignó a Branko para cerrar el sprint.

## Métricas del sprint

| Métrica | Valor |
|---|---|
| SP planificados | 41 (38 comprometidos en la planning + 3 de TEC-19) |
| SP completados (velocidad) | 41 |
| % de historias completadas | 100 % (14 de 14) |
| Velocidad promedio (Sprints 1 y 2) | 42,5 SP |
| Defectos detectados / resueltos | 4 / 4, todos en la revisión de PRs y antes del merge |
| PRs mergeados | 17 (#66, #80 a #82 y #84 a #97; #83 se cerró sin mergear) |
| Escenarios BDD en verde | 151 |

Defectos detectados en la revisión:

- **#91 (HU-26):**
  - las horas negativas se descartaban en silencio en lugar de devolver un error;
  - el mensaje de error era equivocado;
  - faltaban los tests unitarios.
- **#94 (specs):** las specs de HU-26 y HU-20 no definían el error de horas negativas ni el de fecha vacía. Salió de
  las revisiones de #90 y #91, y se corrigieron las specs.
- **#96 (TEC-16):**
  - cubría solo la lista, sin alta ni detalle;
  - cerraba el issue de HU-09 (Sprint 3) en lugar del de TEC-16.
- **#82 / #93 (HU-13):** la pantalla del sprint no ofrecía los botones para completar historias y cerrar el sprint,
  aunque las rutas existían.

## Feedback recibido y cambios al Product Backlog

- **El MVP ya se usa de punta a punta.** El equipo propone que el Sprint 3 apunte al dashboard (HU-28), que puede
  partir de la vista de métricas del sprint (TEC-17), y a los reportes (E9).
- **Siguen para el Sprint 3, como se acordó en la planning del Sprint 2:** HU-29, HU-04, HU-06 y HU-09 (recortable).
- **HU-09 (#16) sigue en el backlog:** se había movido sola a *In progress* por el `Closes` equivocado de #96.
- **Sin cambios de alcance en el Product Backlog** por esta review.
