# Sprint 1 · Planning

- **Fecha:** viernes 23/10/2026
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Sprint:** vie 23/10 – dom 25/10 (3 días)

## Sprint Goal

> Poder crear un proyecto, armar su backlog, planificar un sprint y ver sus primeras métricas.

## Sprint Backlog

| Issue | Historia | SP | Responsable |
|---|---|---|---|
| #8 | HU-01 · Crear proyecto | 3 | Federico Sosa |
| #10 | HU-03 · Registrar integrantes | 3 | Federico Sosa |
| #17 | HU-10 · Crear sprint con Sprint Goal | 5 | Federico Sosa |
| #18 | HU-11 · Asignar historias a un sprint | 5 | Federico Sosa |
| #12 | HU-05 · Crear ítem del Product Backlog | 3 | Ignacio Bianchi |
| #14 | HU-07 · Criterios de aceptación de una historia | 2 | Ignacio Bianchi |
| #15 | HU-08 · Cambiar el estado de una historia | 3 | Ignacio Bianchi |
| #22 | HU-15 · Estimar una historia en Story Points | 2 | Ignacio Bianchi |
| #26 | HU-19 · Registrar horas estimadas de una historia | 1 | Agustín Lara |
| #27 | HU-20 · Registrar esfuerzo trabajado | 5 | Agustín Lara |
| #29 | HU-22 · Registrar defecto | 3 | Agustín Lara |
| #31 | HU-24 · Story Points planificados, completados y velocidad | 5 | Agustín Giorlando |
| #32 | HU-25 · Porcentaje de historias completadas | 2 | Agustín Giorlando |
| #19 | HU-12 · Registrar historias completadas | 2 | Branko Almeira |
| | **Total** | **44** | |

Ítems técnicos del sprint, a cargo de Branko: TEC-09 · contratos de dominio y guía de specs (#43) y TEC-10 · infraestructura
BDD común (#46). Primer sprint: todavía no hay velocidad de referencia.

| Responsable | SP en el Sprint 1 |
|---|---|
| Federico Sosa | 16 |
| Ignacio Bianchi | 10 |
| Agustín Lara | 9 |
| Agustín Giorlando | 7 |
| Branko Almeira | 2 (+ TEC-09, TEC-10 y facilitación) |

## Decisiones

- **Specs y escenarios BDD:** los redacta el Agile Enabler para unificar criterios entre épicas; cada responsable los
  valida (PR #45). Los Builders escriben los tests, los pasos BDD y el código con TDD.
- **Contratos de dominio** (ADR 0002): tipos de cada épica, referencias por ID entre épicas y dependencias permitidas.
- **Pasos BDD compartidos:** cada uno lo implementa una sola persona, según la tabla de `features/README.md` (TEC-10).
- **Uso de IA:** cada integrante usa su asistente con las reglas de `CONTRIBUTING.md`, el ADR 0002 y la guía
  `docs/guias/de-issue-a-spec.md`, y registra el uso en `docs/ia/registro.md`.

## Riesgos y dependencias

- **Sosa tiene la mayor carga (16 SP) y los tipos que necesitan los demás.** Orden de merge acordado (ADR 0002, sección 6):
  `Proyecto` e `Integrante` (HU-01 y HU-03), `Historia` (HU-05) y `Sprint` (HU-10), **mergeados el sábado 24 a las 12:00**.
- HU-12 (Branko) depende de HU-05 y HU-10. HU-20 y HU-22 (Lara) dependen de los tipos de Sosa para el servicio y los pasos BDD.
- Lara y Giorlando no dependen de nadie para el dominio: arrancan el viernes.
- Si el sábado al mediodía los tipos no están mergeados, se replanifica en la Daily del domingo. **Los sprints no se estiran.**

## Ceremonias del sprint

- **Daily:** por el grupo, todos los días antes de las 10:00 (`dailies.md`).
- **Cierre del Sprint 1** (review + retro + planning del Sprint 2): domingo 25/10 por la noche, 60'.
