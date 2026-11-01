# HU-28 · Dashboard del proyecto

- **Issue:** #35 · **Épica:** E8 Dashboard · **Responsable:** Agustín Giorlando · **Story Points:** 8
- **Escenarios BDD:** [`features/hu-28-dashboard-del-proyecto.feature`](../features/hu-28-dashboard-del-proyecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el responsable al aprobar el PR de la planning del Sprint 3
- **Depende de:** solo código en `main`: HU-24 a HU-27 (métricas), HU-19 y HU-20 (horas), HU-22 y HU-23 (defectos),
  TEC-17 (velocidad del proyecto en la pantalla del sprint) y TEC-20 (`Proyecto.EstadoAl`, `Proyecto.DiasRestantes` y
  `metrics.HorasEnPeriodo`). **No espera a HU-04 ni a HU-14**: arma su resumen con el dominio y los servicios que ya
  existen.
- **Estado:** implementada (Sprint 3)

## Historia

**Como** integrante del equipo **quiero** un dashboard con el estado del proyecto y gráficos **para** ver todo de un
vistazo.

## Objetivo

Reunir en la sección **Dashboard** (`/dashboard`) el estado del proyecto, el sprint activo, las métricas principales y la
evolución sprint a sprint en cuatro gráficos. Es la pantalla que se muestra en la presentación final.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir. En la pantalla se elige igual que en las demás secciones (TEC-13) |
| Hoy | `time.Time` | Sí | La fecha del sistema |

## Salidas esperadas

1. **Resumen:** estado y días restantes (TEC-20), fechas, sprint activo y porcentaje de historias completadas del
   proyecto (HU-25). Es el mismo dato que muestra HU-04 en el detalle del proyecto, calculado con las mismas funciones.
2. **Métricas principales:** velocidad (HU-24), horas estimadas, reales y desviación del proyecto (HU-26), defectos
   detectados, resueltos y abiertos (HU-27).
3. **Gráficos**, con un punto por sprint cerrado, en orden de inicio:
   - velocidad por sprint: SP completados de cada sprint y la velocidad promedio como línea de referencia;
   - SP planificados vs. completados;
   - horas estimadas vs. reales;
   - defectos detectados vs. resueltos.

## Reglas de negocio

- **RN-28.1:** los gráficos usan solo los sprints **cerrados** (estado `Cerrado`): sus valores ya no cambian. El sprint activo se
  ve en el resumen.
- **RN-28.2:** los SP de cada sprint son los fijados al cerrarlo (`Sprint.SPPlanificados` y `Sprint.SPCompletados`, HU-13).
- **RN-28.3:** horas de un sprint: las **estimadas** son la suma de las horas estimadas (HU-19) de las historias
  planificadas en el sprint; las **reales**, la suma de los registros de esas historias con fecha entre el inicio y el
  fin del sprint (HU-20), calculadas con `metrics.HorasEnPeriodo` (TEC-20). Es la misma función que usa el reporte de
  sprint (RN-29.3): los dos dan el mismo número sin depender uno del otro.
- **RN-28.4:** defectos de un sprint: detectados y resueltos según HU-27 (sprint de detección y de resolución).
- **RN-28.5:** la desviación del proyecto se informa como "no calculable" si no hay horas estimadas (RN-26.3), sin error
  en la pantalla.
- **RN-28.6:** un gráfico sin datos no se dibuja vacío: en su lugar se muestra un mensaje (CA-28.3):
  - sin sprints cerrados: "Todavía no hay sprints cerrados para graficar" (en los cuatro gráficos);
  - sin horas estimadas ni registradas en los sprints cerrados: "Todavía no hay horas registradas";
  - sin defectos en los sprints cerrados: "No hay defectos registrados".

## Restricciones

- Un servicio nuevo, `service.Dashboard`, con `Del(proyectoID, hoy) (Dashboard, error)`, arma los datos con los
  repositorios y servicios que ya están en `main` (proyectos, sprints, historias, estimaciones de horas, registros y
  defectos) y las funciones de `internal/metrics`; no calcula métricas por su cuenta.
- `Dashboard` trae las series ya armadas (`[]PuntoSprint` con sprint, SP planificados, SP completados, horas estimadas,
  horas reales, defectos detectados y resueltos) y un mensaje por gráfico cuando no hay datos.
- La pantalla usa Chart.js 4.4.7, que ya está embebido en `web/static/vendor/chart.umd.min.js` y cargado en
  `base.html` (TEC-03), con un `<canvas>` por gráfico y los datos de las series en un `<script type="application/json">`
  generado por la plantilla. Los tests del handler verifican esos datos, no el dibujo.
- **Dos PRs:** primero el servicio con el resumen y las métricas (CA-28.1), después los gráficos (CA-28.2 y CA-28.3).
  Así el primero ya suma valor aunque el segundo se atrase.
- La sección deja de estar "en construcción" (`EnConstruccion: false` solo en Dashboard).

## Casos límite

- Proyecto recién creado: resumen completo, métricas en 0, desviación "no calculable" y los cuatro mensajes.
- Un solo sprint cerrado: cada gráfico tiene un punto; la velocidad es la de ese sprint.
- Sprints cerrados sin defectos pero con horas: se dibujan tres gráficos y el de defectos muestra su mensaje.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" (la pantalla responde 404) |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-28.1 | Muestra estado del proyecto, sprint activo y métricas principales | `Escenario: resumen y métricas principales` |
| CA-28.2 | Incluye gráficos de velocidad por sprint, SP planificados vs. completados, horas estimadas vs. reales y defectos | `Escenario: series por sprint cerrado` |
| CA-28.3 | Si no hay datos, muestra un mensaje en lugar de gráficos vacíos | `Escenario: proyecto sin sprints cerrados`, `Escenario: sprints cerrados sin defectos` |

## Decisiones de la planning

- **El sprint activo no entra en los gráficos** (RN-28.1): se ve en el resumen.
- **Se mantiene en 8 SP y en dos PRs** (ver Restricciones). Si en la daily del segundo día el PR de gráficos no está
  abierto, el equipo decide si se reasigna una parte.
