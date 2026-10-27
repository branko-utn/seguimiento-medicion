# HU-28 · Dashboard del proyecto

- **Issue:** #35 · **Épica:** E8 Dashboard · **Responsable:** Agustín Giorlando · **Story Points:** 8
- **Escenarios BDD:** [`features/hu-28-dashboard-del-proyecto.feature`](../features/hu-28-dashboard-del-proyecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-04 (estado del proyecto), HU-14 (sprints cerrados), HU-24 a HU-27 (métricas), HU-19 y HU-20 (horas),
  HU-22 y HU-23 (defectos), TEC-17 (vista de métricas del sprint: se reutilizan sus gráficos)
- **Estado:** borrador para la planning del Sprint 3

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

1. **Resumen** (HU-04): estado, fechas, días restantes, sprint activo y porcentaje de historias completadas.
2. **Métricas principales:** velocidad (HU-24), horas estimadas, reales y desviación del proyecto (HU-26), defectos
   detectados, resueltos y abiertos (HU-27).
3. **Gráficos**, con un punto por sprint cerrado, en orden de inicio:
   - velocidad por sprint: SP completados de cada sprint y la velocidad promedio como línea de referencia;
   - SP planificados vs. completados;
   - horas estimadas vs. reales;
   - defectos detectados vs. resueltos.

## Reglas de negocio

- **RN-28.1:** los gráficos usan solo los sprints **cerrados** (HU-14): sus valores ya no cambian. El sprint activo se
  ve en el resumen.
- **RN-28.2:** los SP de cada sprint son los fijados al cerrarlo (RN-14.3).
- **RN-28.3:** horas de un sprint: las **estimadas** son la suma de las horas estimadas (HU-19) de las historias
  planificadas en el sprint; las **reales**, la suma de los registros de esas historias con fecha entre el inicio y el
  fin del sprint (HU-20). Es la misma regla del reporte de sprint (RN-29.3).
- **RN-28.4:** defectos de un sprint: detectados y resueltos según HU-27 (sprint de detección y de resolución).
- **RN-28.5:** la desviación del proyecto se informa como "no calculable" si no hay horas estimadas (RN-26.3), sin error
  en la pantalla.
- **RN-28.6:** un gráfico sin datos no se dibuja vacío: en su lugar se muestra un mensaje (CA-28.3):
  - sin sprints cerrados: "Todavía no hay sprints cerrados para graficar" (en los cuatro gráficos);
  - sin horas estimadas ni registradas en los sprints cerrados: "Todavía no hay horas registradas";
  - sin defectos en los sprints cerrados: "No hay defectos registrados".

## Restricciones

- Un servicio nuevo, `service.Dashboard`, con `Del(proyectoID, hoy) (Dashboard, error)`, arma los datos reutilizando los
  servicios y las funciones de `internal/metrics`; no calcula métricas por su cuenta.
- `Dashboard` trae las series ya armadas (`[]PuntoSprint` con sprint, SP planificados, SP completados, horas estimadas,
  horas reales, defectos detectados y resueltos) y un mensaje por gráfico cuando no hay datos.
- La pantalla usa Chart.js 4.4.7, embebido en `web/` (TEC-03), y los mismos gráficos de la vista de métricas (TEC-17).
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

## Preguntas para la planning

- **Sprint activo en los gráficos.** Propuesta: no (RN-28.1). ¿Lo queremos como último punto, marcado "en curso"?
- **Tamaño.** Con 8 SP es la historia más grande del sprint. Si TEC-17 ya deja los gráficos hechos, alcanza; si no,
  conviene partirla en "resumen y métricas" y "gráficos".
