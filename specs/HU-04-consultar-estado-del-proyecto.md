# HU-04 · Consultar estado del proyecto

- **Issue:** #11 · **Épica:** E1 Proyectos · **Responsable:** Federico Sosa · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-04-consultar-estado-del-proyecto.feature`](../features/hu-04-consultar-estado-del-proyecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el responsable al aprobar el PR de la planning del Sprint 3
- **Depende de:** solo código en `main`: HU-01, HU-10, HU-05, HU-08, HU-25 y TEC-20 (`Proyecto.EstadoAl` y
  `Proyecto.DiasRestantes`, que entran antes de que arranque el sprint)
- **Estado:** aprobada en la planning del Sprint 3

## Historia

**Como** integrante del equipo **quiero** ver el estado del proyecto **para** saber cómo vamos.

## Objetivo

Mostrar en un solo lugar en qué etapa está el proyecto, cuánto tiempo le queda, qué sprint está en curso y cuánto del
backlog ya se completó. El dashboard (HU-28) y el reporte de proyecto (HU-30) muestran los mismos datos con las mismas
funciones de TEC-20 y HU-25, sin esperar a esta historia.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Hoy | `time.Time` | Sí | La fecha del sistema (en los escenarios, el paso `hoy es "…"`) |

## Salidas esperadas

Un resumen del proyecto con:

- el estado: "Planificado", "En curso" o "Finalizado";
- las fechas de inicio y fin;
- los días restantes;
- el sprint activo (ID y Sprint Goal), si hay uno;
- el porcentaje de historias completadas del proyecto.

La consulta no modifica nada.

## Reglas de negocio

- **RN-04.1:** el estado sale de las fechas del proyecto respecto de hoy: antes del inicio, "Planificado"; entre el
  inicio y el fin (ambos incluidos), "En curso"; después del fin, "Finalizado".
- **RN-04.2:** los días restantes son los días que faltan desde hoy hasta la fecha de fin, sin contar hoy. En un proyecto
  finalizado son 0.
- **RN-04.3:** el sprint activo es el sprint del proyecto en estado `Activo` (HU-10 permite como máximo uno). Si no hay,
  se informa "sin sprint activo".
- **RN-04.4:** el porcentaje de completadas es el de HU-25 **por proyecto**: historias del backlog en estado `Completada`
  sobre el total de historias del backlog.

## Restricciones

- RN-04.1 y RN-04.2 ya están en el dominio (TEC-20): `Proyecto.EstadoAl(hoy)` y `Proyecto.DiasRestantes(hoy)`. HU-04
  no los reimplementa.
- `Proyectos.Resumen(id domain.ProyectoID, hoy time.Time) (ResumenProyecto, error)` en `internal/service` arma el
  resumen: busca el proyecto, sus sprints (sprint activo) y sus historias, y usa `metrics.PorcentajeCompletadas`.
- **Pantalla:** el detalle del proyecto (`/proyectos/{proyecto}`, TEC-13) muestra el resumen arriba de los integrantes.
- **HU-02 usa el estado calculado:** "un proyecto finalizado no se puede modificar" se valida con `EstadoAl(hoy)`, no con
  el campo guardado. El campo `Estado` queda solo por compatibilidad con los datos ya guardados.

## Casos límite

- Hoy es el día de inicio: "En curso".
- Hoy es el día de fin: "En curso", con 0 días restantes.
- Proyecto sin historias: 0 % (RN-25.2).
- Proyecto sin sprints o con todos cerrados: "sin sprint activo".

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-04.1 | El estado es "Planificado", "En curso" o "Finalizado" | `Esquema del escenario: estado según la fecha` |
| CA-04.2 | Se muestran las fechas, los días restantes y el sprint activo, si hay uno | `Escenario: proyecto en curso con sprint activo`, `Escenario: proyecto sin sprint activo` |
| CA-04.3 | Se muestra el porcentaje de historias completadas del proyecto | `Escenario: proyecto en curso con sprint activo`, `Escenario: proyecto sin historias` |

## Decisiones de la planning

- **El estado se calcula siempre con las fechas** (`EstadoAl(hoy)`); no hay botón "Finalizar proyecto".
