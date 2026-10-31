# Informe de métricas y cobertura de pruebas

> Plantilla de TEC-08 · Responsable: Agustín Lara. Los números salen de la app con `make demo` (dashboard y
> reporte de proyecto), del tablero y de `make cover`. Cada tabla dice de dónde sale cada dato.

## 1. Resumen

<!-- 4 o 5 líneas: velocidad, SP comprometidos vs. completados, desviación de horas, defectos y cobertura. -->

## 2. Métricas por sprint

| Sprint | SP planificados | SP completados | % historias | Horas estimadas | Horas reales | Defectos detectados / resueltos |
|---|---|---|---|---|---|---|
| 1 · MVP | | | | | | |
| 2 · Interfaz | | | | | | |
| 3 · Funcionalidad y calidad | | | | | | |
| 4 · Cierre | | | | | | |

Fuente: reporte de proyecto (`/reportes/proyecto?proyecto=P1`). Aclarar que la app cuenta solo historias (el
tablero también cuenta los ítems técnicos) y que las horas de la demo son estimadas (ver README).

## 3. Velocidad y estimación

<!-- Velocidad por sprint y promedio; comparar con lo comprometido en cada planning (actas en docs/scrum/). -->

## 4. Calidad: defectos

<!-- Defectos por sprint y por severidad; de dónde salieron (auditorías de PRs, actas de review). -->

## 5. Cobertura de pruebas

| Paquete | Cobertura |
|---|---|
| | |

Fuente: `make cover` y el resumen del CI (TEC-05: umbral de 80 % en `internal/domain` + `internal/metrics`).
Agregar la cantidad de escenarios BDD (`make bdd`).

## 6. Conclusiones

<!-- Qué aprendimos de las métricas y qué haríamos distinto. -->
