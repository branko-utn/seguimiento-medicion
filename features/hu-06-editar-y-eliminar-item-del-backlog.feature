# language: es
@HU-06 @pendiente
Característica: Editar y eliminar ítem del backlog
  Como líder del proyecto
  Quiero editar o eliminar historias
  Para mantener el backlog actualizado

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
      | HU-2     | Crear sprint   | 5            |

  @CA-06.1 @normal
  Escenario: editar una historia
    Cuando edito la historia "HU-1" con título "Crear y listar proyectos", descripción "Alta y listado" y prioridad "Alta"
    Entonces la historia "HU-1" tiene el título "Crear y listar proyectos" y prioridad "Alta"
    Y la historia "HU-1" tiene 3 Story Points
    Y la historia "HU-1" está en estado "Pendiente"

  @CA-06.1 @error
  Escenario: editar con título vacío
    Cuando edito la historia "HU-1" con título "   ", descripción "" y prioridad "Alta"
    Entonces veo el error "el título es obligatorio"
    Y la historia "HU-1" tiene el título "Crear proyecto" y prioridad "Media"

  @CA-06.1 @error
  Escenario: editar una historia completada
    Dado que la historia "HU-1" está completada
    Cuando edito la historia "HU-1" con título "Otro título", descripción "" y prioridad "Alta"
    Entonces veo el error "una historia completada no se puede editar"

  @CA-06.2 @normal
  Escenario: eliminar una historia del backlog
    Cuando elimino la historia "HU-2"
    Entonces el backlog de "Demo" no tiene la historia "HU-2"
    Y el backlog de "Demo" tiene 1 historias

  @CA-06.2 @error
  Escenario: eliminar una historia asignada a un sprint
    Dado que existe un sprint "S1" en "Demo"
    Y que la historia "HU-2" está asignada al sprint "S1"
    Cuando elimino la historia "HU-2"
    Entonces veo el error "no se puede eliminar una historia asignada a un sprint"
    Y el backlog de "Demo" tiene 2 historias

  @CA-06.2 @error
  Escenario: eliminar una historia completada
    Dado que la historia "HU-1" está completada
    Cuando elimino la historia "HU-1"
    Entonces veo el error "no se puede eliminar una historia completada"

  @CA-06.2 @error
  Escenario: eliminar una historia con esfuerzo registrado
    Dado que "Ana Pérez" es integrante del proyecto "Demo"
    Y que "Ana Pérez" registró 2 horas de "Desarrollo" en "HU-2" el "2026-10-27"
    Cuando elimino la historia "HU-2"
    Entonces veo el error "no se puede eliminar una historia con esfuerzo o defectos registrados"
    Y el backlog de "Demo" tiene 2 historias

  @error
  Escenario: eliminar una historia inexistente
    Cuando elimino la historia "HU-99"
    Entonces veo el error "la historia no existe"
