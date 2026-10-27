# language: es
@HU-15
Característica: Estimar una historia en Story Points
  Como integrante del equipo
  Quiero asignar Story Points a una historia
  Para conocer su tamaño

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"

  @CA-15.1 @normal @limite @error
  Esquema del escenario: valores de la escala
    Cuando estimo la historia "HU-1" en <sp> Story Points
    Entonces <resultado>

    Ejemplos:
      | sp | resultado                                                                 |
      | 1  | la historia "HU-1" tiene 1 Story Points                                   |
      | 5  | la historia "HU-1" tiene 5 Story Points                                   |
      | 21 | la historia "HU-1" tiene 21 Story Points                                  |
      | 4  | veo el error "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"          |
      | 0  | veo el error "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"          |
      | 34 | veo el error "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"          |

  @CA-15.2 @alternativo
  Escenario: re-estimar una historia
    Dado que la historia "HU-1" tiene 3 Story Points
    Cuando estimo la historia "HU-1" en 5 Story Points
    Entonces la historia "HU-1" tiene 5 Story Points

  @CA-15.2 @error
  Escenario: re-estimar una historia completada
    Dado que la historia "HU-1" tiene 3 Story Points
    Y que la historia "HU-1" está completada
    Cuando estimo la historia "HU-1" en 5 Story Points
    Entonces veo el error "no se puede re-estimar una historia completada"
    Y la historia "HU-1" tiene 3 Story Points
