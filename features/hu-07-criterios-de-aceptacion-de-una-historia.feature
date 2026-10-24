# language: es
@HU-07
Característica: Criterios de aceptación de una historia
  Como integrante del equipo
  Quiero agregar y quitar criterios de aceptación a una historia
  Para saber cuándo está terminada

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"

  @CA-07.1 @normal
  Escenario: agregar un criterio
    Cuando agrego a "HU-1" el criterio "El nombre es obligatorio"
    Y agrego a "HU-1" el criterio "  La fecha de fin no es anterior a la de inicio  "
    Entonces la historia "HU-1" tiene los criterios:
      | número | criterio                                      |
      | 1      | El nombre es obligatorio                      |
      | 2      | La fecha de fin no es anterior a la de inicio |

  @CA-07.1 @error
  Escenario: criterio vacío
    Cuando agrego a "HU-1" el criterio "   "
    Entonces veo el error "el criterio de aceptación no puede estar vacío"

  @alternativo
  Escenario: quitar un criterio renumera los siguientes
    Dado que la historia "HU-1" tiene los criterios "A, B, C"
    Cuando quito de "HU-1" el criterio 2
    Entonces la historia "HU-1" tiene los criterios:
      | número | criterio |
      | 1      | A        |
      | 2      | C        |

  @CA-07.2 @limite
  Escenario: quitar el único criterio
    Dado que la historia "HU-1" tiene los criterios "A"
    Cuando quito de "HU-1" el criterio 1
    Entonces la historia "HU-1" no tiene criterios

  @error
  Escenario: criterio inexistente
    Dado que la historia "HU-1" tiene los criterios "A"
    Cuando quito de "HU-1" el criterio 2
    Entonces veo el error "el criterio no existe"
