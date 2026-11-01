# HU-27 · Defectos detectados y resueltos

- **Issue:** #34 · **Épica:** E7 Métricas · **Responsable:** Agustín Lara (reparto del Sprint 2) · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-27-defectos-detectados-y-resueltos.feature`](../features/hu-27-defectos-detectados-y-resueltos.feature)
- **Depende de:** nada: función pura con tipos de entrada propios. No espera a HU-23: el servicio conecta los datos reales cuando HU-23 esté en `main`.
- **Redactada por:** Branko Almeira (Agile Enabler) después de la planning del Sprint 2 · **Validada por:** Agustín Lara (al aprobar el PR)
- **Estado:** implementada (Sprint 2)

## Historia

**Como** líder del proyecto **quiero** ver la cantidad de defectos detectados y resueltos **para** medir la calidad.

## Objetivo

Contar los defectos detectados y resueltos por sprint y en total. Lo usan la vista de métricas (TEC-17), el dashboard
(HU-28) y los reportes (E9).

## Entradas

```go
// Un defecto, visto desde las métricas.
type DefectoSprint struct {
	SprintDeteccion  string
	SprintResolucion string // "" mientras no esté resuelto
}
```

| Función | Salida |
|---|---|
| `DefectosDetectados(ds []DefectoSprint, sprint string) int` | Defectos cuyo sprint de detección es `sprint` |
| `DefectosResueltos(ds []DefectoSprint, sprint string) int` | Defectos cuyo sprint de resolución es `sprint` |
| `TotalDetectados(ds []DefectoSprint) int` · `TotalResueltos(ds []DefectoSprint) int` | Totales del proyecto |

## Reglas de negocio

- **RN-27.1:** un defecto cuenta como detectado en el sprint de detección.
- **RN-27.2:** un defecto cuenta como resuelto en el sprint de resolución; si no tiene sprint de resolución, no está resuelto.
- **RN-27.3:** los totales son la cantidad de defectos y la cantidad de resueltos de todo el proyecto.

## Restricciones

- Funciones puras en `internal/metrics`; los sprints se identifican con su ID como texto (no se importa `domain`).

## Casos límite

- Sin defectos: todo es 0.
- Un defecto detectado en S1 y resuelto en S2 cuenta como detectado en S1 y resuelto en S2.

## Condiciones de error

Las funciones no devuelven error: reciben datos ya validados por el dominio (HU-22 y HU-23).

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-27.1 | Se cuentan por sprint (según el sprint de detección y el de resolución) y en total | `Escenario: defectos por sprint`, `Escenario: totales del proyecto`, `Escenario: sin defectos` |
