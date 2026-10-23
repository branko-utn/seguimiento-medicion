# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-19 @pendiente
Característica: Registrar horas estimadas de una historia
  Como integrante del equipo
  Quiero registrar las horas estimadas de una historia
  Para compararlas después con las reales

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"

  @CA-19.1 @normal
  Escenario: registrar horas estimadas
    Cuando estimo la historia "HU-1" en 4 horas
    Entonces la historia "HU-1" tiene 4 horas estimadas

  @CA-19.1 @limite
  Escenario: fracciones de hora
    Cuando estimo la historia "HU-1" en 0.5 horas
    Entonces la historia "HU-1" tiene 0.5 horas estimadas

  @alternativo
  Escenario: re-estimar reemplaza la estimación anterior
    Dado que la historia "HU-1" tiene 4 horas estimadas
    Cuando estimo la historia "HU-1" en 6 horas
    Entonces la historia "HU-1" tiene 6 horas estimadas

  @CA-19.1 @error
  Esquema del escenario: horas inválidas
    Cuando estimo la historia "HU-1" en <horas> horas
    Entonces veo el error "las horas deben ser mayores que cero"

    Ejemplos:
      | horas |
      | 0     |
      | -2    |

  @error
  Escenario: historia inexistente
    Cuando estimo la historia "HU-99" en 4 horas
    Entonces veo el error "la historia no existe"
