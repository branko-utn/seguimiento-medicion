# language: es
@HU-04
Característica: Consultar estado del proyecto
  Como integrante del equipo
  Quiero ver el estado del proyecto
  Para saber cómo vamos

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"

  @CA-04.1 @normal @limite
  Esquema del escenario: estado según la fecha
    Dado hoy es "<hoy>"
    Cuando consulto el estado del proyecto "Demo"
    Entonces el estado del proyecto es "<estado>"
    Y quedan <dias> días para el fin del proyecto

    Ejemplos:
      | hoy        | estado      | dias |
      | 2026-10-20 | Planificado | 13   |
      | 2026-10-21 | En curso    | 12   |
      | 2026-10-27 | En curso    | 6    |
      | 2026-11-02 | En curso    | 0    |
      | 2026-11-03 | Finalizado  | 0    |

  @CA-04.2 @CA-04.3 @normal
  Escenario: proyecto en curso con sprint activo
    Dado hoy es "2026-10-27"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |
      | HU-4     | Ver dashboard     | 8            |
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Cuando consulto el estado del proyecto "Demo"
    Entonces el proyecto va del "2026-10-21" al "2026-11-02"
    Y el sprint activo del proyecto es "S1"
    Y el porcentaje de historias completadas del proyecto es 25.0

  @CA-04.2 @alternativo
  Escenario: proyecto sin sprint activo
    Dado hoy es "2026-10-27"
    Y existe un sprint "S1" en "Demo"
    Cuando consulto el estado del proyecto "Demo"
    Entonces el proyecto no tiene sprint activo

  @CA-04.3 @limite
  Escenario: proyecto sin historias
    Dado hoy es "2026-10-27"
    Cuando consulto el estado del proyecto "Demo"
    Entonces el porcentaje de historias completadas del proyecto es 0.0

  @error
  Escenario: proyecto inexistente
    Dado hoy es "2026-10-27"
    Cuando consulto el estado del proyecto "Otro"
    Entonces veo el error "el proyecto no existe"
