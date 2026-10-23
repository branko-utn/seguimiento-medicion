# HU-24 · Story Points planificados, completados y velocidad

- **Issue:** #31 · **Épica:** E7 Métricas · **Responsable:** Agustín Giorlando · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-24-story-points-planificados-completados-y-velocidad.feature`](../features/hu-24-story-points-planificados-completados-y-velocidad.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Agustín Giorlando (al aprobar el PR)
- **Depende de:** nada: funciones puras con tipos de entrada propios (ADR 0002, sección 3)
- **Estado:** revisada

## Historia

**Como** líder del proyecto **quiero** ver los SP planificados y completados de cada sprint y la velocidad del equipo **para** planificar mejor.

## Objetivo

Calcular las tres métricas básicas de Scrum. Son funciones puras en `internal/metrics`, sin dependencias de otros
paquetes: el servicio traduce sprints e historias a sus tipos de entrada.

## Entradas

```go
// Una historia del sprint, vista desde las métricas.
type HistoriaSprint struct {
	StoryPoints int
	Completada  bool
}
```

| Función | Entrada | Salida |
|---|---|---|
| `SPPlanificados(hs []HistoriaSprint) int` | Historias asignadas al sprint | Suma de sus Story Points |
| `SPCompletados(hs []HistoriaSprint) int` | Historias asignadas al sprint | Suma de los Story Points de las completadas |
| `Velocidad(completados []int) float64` | SP completados de cada sprint **cerrado** | Promedio |

## Salidas esperadas

Números listos para mostrar en el estado del proyecto, el dashboard (HU-28) y los reportes (E9). La velocidad se
muestra con un decimal; la función devuelve el valor sin redondear.

## Reglas de negocio

- **RN-24.1:** SP planificados = suma de los SP de las historias asignadas al sprint.
- **RN-24.2:** SP completados = suma de los SP de las historias completadas del sprint.
- **RN-24.3:** velocidad = promedio de SP completados de los sprints cerrados.
- **RN-24.4:** sin sprints cerrados, la velocidad es 0. Un sprint sin historias suma 0.

## Restricciones

- Las funciones no acceden a repositorios ni a otros paquetes de dominio.
- El servicio decide qué sprints están cerrados y arma las entradas; las métricas no conocen estados de sprint.

## Casos límite

- Lista vacía de historias: planificados 0, completados 0.
- Ninguna historia completada: completados 0.
- Un solo sprint cerrado con 13 SP: velocidad 13.
- Sprints cerrados con 10, 13 y 8 SP: velocidad 10,333… (se muestra 10,3).
- Un sprint cerrado con 0 SP completados cuenta en el promedio: 10 y 0 dan 5.

## Condiciones de error

Las funciones no devuelven error: reciben datos ya validados por el dominio (los SP siempre son 0 o de la escala, HU-15).

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-24.1 | SP planificados = suma de los SP de las historias asignadas al sprint | `Escenario: SP planificados y completados de un sprint` |
| CA-24.2 | SP completados = suma de los SP de las historias completadas del sprint | `Escenario: SP planificados y completados de un sprint`, `Escenario: sprint sin historias completadas` |
| CA-24.3 | Velocidad = promedio de SP completados de los sprints cerrados; sin sprints cerrados es 0 | `Escenario: velocidad del equipo`, `Escenario: sin sprints cerrados` |
