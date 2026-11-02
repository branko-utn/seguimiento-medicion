# Sprint 4 · Retrospectiva

- **Fecha:** domingo 01/11 a lunes 02/11/2026 · **Formato:** asíncrona, por el chat del grupo (como la daily)
- **Participaron:** Federico Sosa, Ignacio Bianchi, Agustín Lara y Agustín Giorlando. Branko Almeira (Agile Enabler)
  la facilitó.
- **Dinámica:** el Agile Enabler armó una propuesta con la evidencia del sprint (planning, dailies, PRs, revisiones y
  la auditoría de la consigna) y la compartió. Cada integrante respondió tres preguntas: si estaba de acuerdo con los
  puntos, un punto propio y si se había hecho el ensayo de la demo. Abajo, lo que aportó cada uno va con su nombre.
- **Formato del contenido:** Seguir / Dejar de hacer / Probar

## Acciones de la retro anterior

| Acción (retro del Sprint 3) | ¿Se hizo? | Evidencia |
|---|---|---|
| Code freeze el sábado 31/10 a las 23:59 | Parcialmente: 11 PRs se abrieron en las últimas horas antes del freeze y se mergearon el 1/11, después de la auditoría | PRs #127 a #136 |
| Escribir la resolución de conflictos de los archivos compartidos | Sí | [`CONTRIBUTING.md`](../../../CONTRIBUTING.md) y anexo A.7 del manual |
| Resolver el defecto del punto decimal del dashboard | Sí, con HU-32.4 | PR #134 |
| Ensayar la demo el domingo 1/11 | **No.** Los cuatro Product Builders respondieron que todavía no se hizo | Respuestas de la retro |

## Acuerdo con la propuesta

| Integrante | Respuesta |
|---|---|
| Federico Sosa | De acuerdo |
| Ignacio Bianchi | De acuerdo con "Dejar de hacer" y "Probar". Con un matiz en "Seguir" (abajo) |
| Agustín Lara | De acuerdo: "están bien reflejados" |
| Agustín Giorlando | De acuerdo |

## Seguir haciendo

- **Una base compartida que entra primero, pero sin depender de una sola persona.** TEC-21 (base visual) entró antes
  que las cuatro partes de HU-32, y cada uno mejoró sus pantallas sin esperar a los demás (PRs #121, #123, #124,
  #129, #134). **Matiz de Bianchi:** la base es útil, pero mientras no entra se vuelve a caer en el bloqueo de esperar
  a alguien (la daily del 31/10 registra a Bianchi, Lara y Giorlando esperando TEC-21).
- **Cortes con hora y responsable de reemplazo.** HU-31 tenía corte a las 18:00: pasó a Lara y llegó al freeze
  (planning, "Riesgos y cortes"; PR #132).
- **Decisiones escritas antes de ejecutarlas.** El despliegue tuvo specs, ADR y comentario de decisión en el issue
  cuando apareció el problema del entorno express (#140, #141, comentario en #138).
- **La seguridad como condición para publicar,** no como algo de después: TEC-24 entró en el mismo PR que TEC-23
  (#142).
- **Destrabar entre todos los bloqueos de infraestructura** *(Sosa)*: la buena comunicación permitió resolver temas
  que nadie conocía, como Azure.
- **Mantener la comunicación fluida en el equipo** *(Lara)*.

## Dejar de hacer

- **Abrir los PRs al límite del freeze.** Los 11 PRs de las últimas horas del 31/10 obligaron a auditar y mergear todo
  el 1/11, el día previsto para el ensayo.
- **Planificar dos veces en el mismo día.** La planning del Sprint 4 tuvo cuatro actualizaciones el 31/10 (interfaz,
  PDF en Go, reparto final, presentación).

## Probar

- **Revisar la consigna punto por punto al planificar, no al cerrar.** Que el PDF lo tiene que generar el sistema se
  vio el último día.
- **Volcar la daily el mismo día, sin excepción.** La del 1/11 se volvió a reconstruir.
- **Actualizar el tablero en el momento de cada reasignación** *(Giorlando)*, para que no quede desfasado de Git. Pasó
  con HU-31 (#38) y HU-27 (#34): los assignees cambiaron, pero el cuerpo de los issues siguió con el responsable
  anterior.

## Acciones

| Acción | Responsable | Cuándo |
|---|---|---|
| **Ensayar la demo completa antes de presentar**, sobre una carpeta limpia y con la URL de Azure como respaldo | Todo el equipo | Lunes 02/11, antes de la presentación |
| Auditoría de la consigna en la planning de cada sprint | Agile Enabler | Próximo proyecto |
| Corte para abrir PRs 24 horas antes del freeze | Todo el equipo | Próximo proyecto |
| Daily volcada en el día por quien la facilita | Agile Enabler | Próximo proyecto |
| Al reasignar un ítem, actualizar en el momento el assignee, el cuerpo del issue y la spec | Quien reasigna | Próximo proyecto |
