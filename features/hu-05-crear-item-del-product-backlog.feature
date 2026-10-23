# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-05 @pendiente
Característica: Crear ítem del Product Backlog
  Como líder del proyecto
  Quiero agregar historias al Product Backlog
  Para registrar lo que hay que construir

  Antecedentes:
    Dado que existe un proyecto "Demo"

  @CA-05.4 @normal
  Escenario: crear una historia
    Cuando agrego al proyecto "Demo" la historia "Crear proyecto" con prioridad "Alta"
    Entonces el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y la historia "HU-1" está en estado "Pendiente"
    Y la historia "HU-1" no tiene Story Points

  @CA-05.1 @normal
  Escenario: los identificadores son correlativos
    Cuando agrego al proyecto "Demo" la historia "Crear proyecto" con prioridad "Alta"
    Y agrego al proyecto "Demo" la historia "Crear sprint" con prioridad "Media"
    Entonces el backlog de "Demo" tiene la historia "HU-2" "Crear sprint"

  @CA-05.1 @limite
  Escenario: cada proyecto numera desde HU-1
    Dado que existe un proyecto "Otro"
    Y que el backlog de "Demo" tiene 3 historias
    Cuando agrego al proyecto "Otro" la historia "Login" con prioridad "Baja"
    Entonces el backlog de "Otro" tiene la historia "HU-1" "Login"

  @CA-05.2 @alternativo
  Escenario: crear una historia sin descripción
    Cuando agrego al proyecto "Demo" la historia "  Login  " con prioridad "Baja"
    Entonces el backlog de "Demo" tiene la historia "HU-1" "Login"

  @CA-05.2 @error
  Escenario: título vacío
    Cuando agrego al proyecto "Demo" la historia "   " con prioridad "Alta"
    Entonces veo el error "el título es obligatorio"

  @CA-05.3 @error
  Escenario: prioridad inválida
    Cuando agrego al proyecto "Demo" la historia "Login" con prioridad "Urgente"
    Entonces veo el error "la prioridad debe ser Alta, Media o Baja"
