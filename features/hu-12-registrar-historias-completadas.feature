# language: es
@HU-12
Característica: Registrar historias completadas
  Como integrante del equipo
  Quiero marcar historias del sprint como completadas
  Para medir el avance

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2"

  @CA-12.1 @normal
  Escenario: completar una historia del sprint
    Cuando completo la historia "HU-1" en el sprint "S1"
    Entonces el sprint "S1" tiene completadas las historias "HU-1"
    Y la historia "HU-1" está en estado "Completada"

  @CA-12.2 @normal
  Escenario: los Story Points de la historia cuentan como completados
    Cuando completo la historia "HU-2" en el sprint "S1"
    Entonces el sprint "S1" tiene 5 Story Points completados de 8 planificados

  @CA-12.2 @limite
  Escenario: completar la última historia del sprint
    Dado que la historia "HU-1" ya está completada en el sprint "S1"
    Cuando completo la historia "HU-2" en el sprint "S1"
    Entonces el sprint "S1" tiene 8 Story Points completados de 8 planificados

  @CA-12.1 @error
  Escenario: completar una historia que no está en el sprint
    Cuando completo la historia "HU-3" en el sprint "S1"
    Entonces veo el error "la historia no está asignada a este sprint"
    Y la historia "HU-3" está en estado "Pendiente"

  @limite @error
  Escenario: completar una historia que ya está completada
    Dado que la historia "HU-1" ya está completada en el sprint "S1"
    Cuando completo la historia "HU-1" en el sprint "S1"
    Entonces veo el error "la historia ya está completada"
    Y el sprint "S1" tiene 3 Story Points completados de 8 planificados

  @error
  Escenario: completar una historia en un sprint cerrado
    Dado que el sprint "S1" está cerrado
    Cuando completo la historia "HU-1" en el sprint "S1"
    Entonces veo el error "el sprint está cerrado y no se puede modificar"
