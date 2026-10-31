# Sprint 3 · Review

- **Fecha:** sábado 31/10/2026 · **Duración:** 20'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando. No participó
  ningún Product Architect.
- **Sprint:** jue 29/10 – vie 30/10 (Funcionalidad y calidad). HU-09, HU-14, HU-18, HU-28 y TEC-06 se mergearon el
  sábado 31 a la mañana (ver la retro).

## Sprint Goal

> Ver el estado del proyecto de un vistazo en el dashboard y generar los reportes de sprint y de proyecto, con el
> Planning Poker completo.

**¿Se cumplió?** Sí. En la demo se cargaron los datos de nuestro propio proyecto (`make demo`) y se recorrió:
1. el **dashboard**, con el estado, el sprint activo, las métricas y los gráficos de los Sprints 1 a 3;
2. el **reporte de proyecto** y el de cada sprint, guardados como PDF desde el navegador;
3. el **Planning Poker** completo, con nueva ronda y estimación acordada.

## Incremento demostrado

| Issue | Ítem | SP | Estado | Quién lo mostró |
|---|---|---|---|---|
| #99 | TEC-20 · Base compartida del Sprint 3 (PR #101) | 1 | Done | Branko Almeira |
| #36 | HU-29 · Reporte de sprint (PR #102) | 5 | Done | Branko Almeira |
| #37 | HU-30 · Reporte de proyecto (PR #108) | 3 | Done | Branko Almeira |
| #39 | TEC-05 · Umbral de cobertura en el CI (PR #103) | 2 | Done | Ignacio Bianchi |
| #25 | HU-18 · Nueva ronda y estimación acordada (PR #105) | 5 | Done | Ignacio Bianchi |
| #11 | HU-04 · Consultar estado del proyecto (PR #104) | 3 | Done | Federico Sosa |
| #13 | HU-06 · Editar y eliminar ítem del backlog (PR #106) | 3 | Done | Agustín Lara |
| #28 | HU-21 · Consultar el esfuerzo registrado (PR #107) | 3 | Done | Agustín Lara |
| #16 | HU-09 · Ver el backlog ordenado por prioridad (PR #110) | 2 | Done | Agustín Lara |
| #21 | HU-14 · Consultar sprints anteriores (PR #109) | 3 | Done | Agustín Giorlando |
| #35 | HU-28 · Dashboard del proyecto (PR #111) | 8 | Done | Agustín Giorlando |
| #40 | TEC-06 · Datos de demo: cargar nuestro propio proyecto (PR #112) | 3 | Done | Branko Almeira |

Los dos ítems recortables (HU-09 y HU-30) entraron.

## Métricas del sprint

| Métrica | Valor |
|---|---|
| SP planificados | 41 (36 firmes + 5 recortables) |
| SP completados (velocidad) | 41 |
| % de historias completadas | 100 % (12 de 12) |
| Velocidad promedio (Sprints 1 a 3) | 42 SP |
| Defectos detectados / resueltos | 7 / 6, todos en la revisión de PRs |
| PRs mergeados | 14 (de #100 a #112, más las actas #98) |
| Escenarios BDD en verde | 203 |

Defectos detectados en la revisión:

- **#103 (TEC-05):**
  - el umbral de cobertura dejaba afuera todo `internal/metrics`;
  - contaba bloques en lugar de sentencias.
- **#104 (HU-04):** los escenarios de HU-02 iban a fallar desde el 03/11 por usar la fecha real.
- **#105 (HU-18):**
  - acordar la estimación no cerraba la votación, así que la historia no se podía volver a votar;
  - faltaban los tests de la pantalla.
- **#111 (HU-28):** el dashboard seguía figurando como "en construcción" y había archivos sin `gofmt`.
- **#109 y #111:** cerraban issues equivocados (HU-07 y HU-21). Por eso HU-21 volvió a *In progress* en el tablero.
- **Abierto:** el dashboard muestra los números con punto decimal en lugar de coma. Pasa al Sprint 4.

## Feedback recibido y cambios al Product Backlog

- **El dashboard con nuestros datos es lo que se muestra en la presentación.** Hay que pulir el formato de los
  números (Sprint 4).
- **HU-31 (exportar a PDF):** el botón "Imprimir o guardar como PDF" de los reportes ya cubre gran parte. Se define
  su alcance en la planning del Sprint 4.
- **Sin cambios de alcance:** quedan HU-31, TEC-07 (manual y documentación técnica) y TEC-08 (informe de métricas y
  cobertura).
