# Registro de uso de IA

La consigna permite usar IA, pero **todo resultado debe ser comprendido, revisado y validado por el equipo**.
Registramos cada uso relevante: qué se pidió, qué se usó del resultado y cómo se validó.

| Fecha | Integrante | Herramienta | Para qué | Qué se usó | Cómo se validó |
|---|---|---|---|---|---|
| 21/10/2026 | Branko Almeira | Claude Code | Análisis de la consigna y de buenas prácticas de otros repositorios del curso; cronograma de sprints | Cronograma, reparto de épicas y ceremonias (`docs/plan/`) | Revisado por Branko y aprobado por el equipo |
| 21/10/2026 | Branko Almeira | Claude Code | Esqueleto del repositorio: estructura, CI, plantillas, guías de Scrum y ejemplo TEC-01 | Archivos del PR `chore/TEC-01-esqueleto` | Revisión del PR por otro integrante; CI en verde |
| 21/10/2026 | Branko Almeira | Claude Code | TEC-03 (base web) y TEC-04 (persistencia): specs, tests, código y guía para agregar pantallas | Archivos de los PR #6 y #7 | Branko revisó cada fase (spec → RED → GREEN); revisión del PR por otro integrante; CI en verde |
| 23/10/2026 | Branko Almeira | Claude Code | TEC-09: ADR 0002 de contratos de dominio, guía "De issue a spec" e IDs compartidos | Archivos del PR de TEC-09 | Revisión de Branko; revisión del PR por otro integrante |
| 23/10/2026 | Branko Almeira | Claude Code | Specs SDD y escenarios BDD de las 13 historias P0 del Sprint 1 (por decisión del Agile Enabler, para unificar criterios) | `specs/HU-*.md` y `features/hu-*.feature` del PR de specs del Sprint 1 | Cada responsable revisa y aprueba las specs de sus historias en el PR |
| 23/10/2026 | Branko Almeira | Claude Code | TEC-10: mundo compartido de los escenarios BDD, pasos comunes y tabla de pasos compartidos | Archivos del PR de TEC-10 | Revisión de Branko; revisión del PR por otro integrante; CI en verde |
| 23/10/2026 | Ignacio Bianchi | OpenCode | HU-05 (dominio): tipo Historia, Prioridad y Estado con TDD (RED→GREEN) | `internal/domain/backlog/historia.go` e `historia_test.go`; ADR 0003 | Tests de tabla con 100 % de cobertura; Ignacio revisó y entendió el diff antes de cada commit; verificado con `go vet` y `go test ./...` en local |
