# Sprint 4 · Retrospectiva

> **Propuesta para validar.** La preparó el Agile Enabler el 01/11 a partir de la evidencia del sprint (planning,
> dailies, PRs, revisiones y la auditoría de cumplimiento). Cada punto cita de dónde sale. El equipo la revisa en una
> retro corta antes de la presentación: lo que no se confirme se borra, y lo que falte se agrega. Al validarla se
> completan la fecha, los asistentes y se quita este aviso.

- **Fecha:** a completar al validarla · **Duración:** 15'
- **Asistentes:** a completar
- **Formato:** Seguir / Dejar de hacer / Probar

## Acciones de la retro anterior

| Acción (retro del Sprint 3) | ¿Se hizo? | Evidencia |
|---|---|---|
| Code freeze el sábado 31/10 a las 23:59 | Parcialmente: 11 PRs se abrieron en las últimas horas antes del freeze y se mergearon el 1/11, después de la auditoría | PRs #127 a #136 |
| Escribir la resolución de conflictos de los archivos compartidos | Sí | [`CONTRIBUTING.md`](../../../CONTRIBUTING.md) y anexo A.7 del manual |
| Resolver el defecto del punto decimal del dashboard | Sí, con HU-32.4 | PR #134 |
| Ensayar la demo el domingo 1/11 | A confirmar | — |

## Seguir haciendo

- **Una base compartida que entra primero.** TEC-21 (base visual) entró antes que las cuatro partes de HU-32, y cada
  uno mejoró sus pantallas sin esperar a los demás (PRs #121, #123, #124, #129, #134).
- **Cortes con hora y responsable de reemplazo.** HU-31 tenía corte a las 18:00: pasó a Lara y llegó al freeze
  (planning, "Riesgos y cortes"; PR #132).
- **Decisiones escritas antes de ejecutarlas.** El despliegue tuvo specs, ADR y comentario de decisión en el issue
  cuando apareció el problema del entorno express (#140, #141, comentario en #138).
- **La seguridad como condición para publicar,** no como algo de después: TEC-24 entró en el mismo PR que TEC-23
  (#142).

## Dejar de hacer

- **Abrir los PRs al límite del freeze.** Los 11 PRs de las últimas horas del 31/10 obligaron a auditar y mergear todo el
  1/11, el día del ensayo.
- **Planificar dos veces en el mismo día.** La planning del Sprint 4 tuvo cuatro actualizaciones el 31/10 (interfaz,
  PDF en Go, reparto final, presentación): el plan cambió más rápido de lo que el equipo podía leerlo.

## Probar

- **Revisar la consigna punto por punto al planificar, no al cerrar.** Que el PDF lo tiene que generar el sistema se
  vio el último día.
- **Volcar la daily el mismo día, sin excepción.** La del 1/11 se volvió a reconstruir.

## Acciones

Es el último sprint, así que las acciones son para el próximo proyecto del equipo:

| Acción | Responsable |
|---|---|
| Auditoría de la consigna en la planning de cada sprint | Agile Enabler |
| Corte para abrir PRs 24 horas antes del freeze | Todo el equipo |
| Daily volcada en el día por quien la facilita | Agile Enabler |
