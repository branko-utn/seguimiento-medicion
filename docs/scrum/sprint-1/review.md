# Sprint 1 · Review

- **Fecha:** domingo 25/10/2026 por la noche (cierre del sprint) · **Duración:** 20'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando. No participó
  ningún Product Architect.
- **Sprint:** vie 23/10 – dom 25/10

## Sprint Goal

> Poder crear un proyecto, armar su backlog, planificar un sprint y ver sus primeras métricas.

**¿Se cumplió?** Sí, como núcleo funcional: todo el recorrido funciona en el dominio y los servicios y está
verificado con los escenarios BDD, pero todavía no tiene pantallas. El equipo acordó que el MVP usable desde el
navegador sea el objetivo del Sprint 2.

## Incremento demostrado

Cada responsable mostró sus escenarios BDD corriendo (`go test ./features`) y el código en `main`.

| Issue | Historia | SP | Estado | Quién lo mostró |
|---|---|---|---|---|
| #8 | HU-01 · Crear proyecto (PR #51) | 3 | Done | Federico Sosa |
| #10 | HU-03 · Registrar integrantes (PR #52) | 3 | Done | Federico Sosa |
| #17 | HU-10 · Crear sprint con Sprint Goal (PR #55) | 5 | Done | Federico Sosa |
| #18 | HU-11 · Asignar historias a un sprint (PR #61) | 5 | Done | Federico Sosa |
| #12 | HU-05 · Crear ítem del Product Backlog (PRs #54 y #56) | 3 | Done | Ignacio Bianchi |
| #14 | HU-07 · Criterios de aceptación de una historia (PRs #57 y #58) | 2 | Done | Ignacio Bianchi |
| #15 | HU-08 · Cambiar el estado de una historia (PRs #68 y #70) | 3 | Done | Ignacio Bianchi |
| #22 | HU-15 · Estimar una historia en Story Points (PRs #63 y #67) | 2 | Done | Ignacio Bianchi |
| #26 | HU-19 · Registrar horas estimadas de una historia (PR #60) | 1 | Done | Agustín Lara |
| #27 | HU-20 · Registrar esfuerzo trabajado (PR #62) | 5 | Done | Agustín Lara |
| #29 | HU-22 · Registrar defecto (PR #65) | 3 | Done | Agustín Lara |
| #31 | HU-24 · Story Points planificados, completados y velocidad (PR #49) | 5 | Done | Agustín Giorlando |
| #32 | HU-25 · Porcentaje de historias completadas (PR #53) | 2 | Done | Agustín Giorlando |
| #19 | HU-12 · Registrar historias completadas (PRs #59, #64 y #69) | 2 | Done | Branko Almeira |

Ítems técnicos: TEC-09 · contratos de dominio y guía de specs (PR #44), TEC-10 · infraestructura BDD común (PR #47) y
ADR 0003 · enumeradores del dominio (PR #50, de Ignacio Bianchi).

## Métricas del sprint

| Métrica | Valor |
|---|---|
| SP planificados | 44 |
| SP completados (velocidad) | 44 |
| % de historias completadas | 100 % (14 de 14) |
| Defectos detectados / resueltos | 5 / 5, todos en la revisión de PRs y antes del merge |
| PRs mergeados | 25 (de #44 a #71; #48 se cerró sin mergear y #66 entró en el Sprint 2) |

Defectos detectados en la revisión:

- **#51:** el estado del proyecto era un `string`; se pasó a `int` con `iota + 1` para seguir el ADR 0003.
- **#53 y #54:** conflictos con `main` en `features/features_test.go` y `docs/ia/registro.md` que borraban
  registros de otros.
- **#54:** `Closes #12` en un PR parcial habría cerrado una historia no terminada; se cambió a `Parte de #12`.
- **#60:** faltaban los escenarios BDD de HU-19, que la Definition of Done exige.

## Feedback recibido y cambios al Product Backlog

- **El equipo pidió ver la aplicación en el navegador:** el Sprint 2 se arma alrededor de las pantallas (TEC-13 a
  TEC-18), que se apoyan en los servicios ya probados.
- **HU-27 pasa de Agustín Giorlando a Agustín Lara:** usa los defectos de HU-22 y así nadie espera código de otro.
- **Pasan al Sprint 3:** HU-29 (reporte de sprint), HU-04 (estado del proyecto), HU-06 (editar y eliminar del
  backlog) y HU-09 (ordenar el backlog, recortable).
- **Desde el Sprint 2, "probado en la aplicación" deja de ser N/A en la Definition of Done.**
