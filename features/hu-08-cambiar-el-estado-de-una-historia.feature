# language: es
@HU-08
Característica: Cambiar el estado de una historia
  Como integrante del equipo
  Quiero cambiar el estado de una historia
  Para reflejar su avance

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y existe un sprint activo "S1" en "Demo"

  @CA-08.1 @normal
  Escenario: empezar una historia del sprint
    Dado que la historia "HU-1" está asignada al sprint "S1"
    Cuando cambio el estado de "HU-1" a "En progreso"
    Entonces la historia "HU-1" está en estado "En progreso"

  @CA-08.1 @alternativo
  Escenario: volver a pendiente
    Dado que la historia "HU-1" está asignada al sprint "S1"
    Y que la historia "HU-1" está en estado "En progreso"
    Cuando cambio el estado de "HU-1" a "Pendiente"
    Entonces la historia "HU-1" está en estado "Pendiente"

  @CA-08.2 @error
  Escenario: completar sin estar en un sprint
    Cuando cambio el estado de "HU-1" a "Completada"
    Entonces veo el error "la historia debe estar asignada a un sprint"
    Y la historia "HU-1" está en estado "Pendiente"

  @CA-08.2 @error
  Escenario: empezar sin estar en un sprint
    Cuando cambio el estado de "HU-1" a "En progreso"
    Entonces veo el error "la historia debe estar asignada a un sprint"

  @CA-08.2 @error
  Escenario: cambiar una historia completada
    Dado que la historia "HU-1" está asignada al sprint "S1"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Cuando cambio el estado de "HU-1" a "En progreso"
    Entonces veo el error "no se puede pasar de Completada a En progreso"

  @limite @error
  Escenario: cambiar al mismo estado
    Cuando cambio el estado de "HU-1" a "Pendiente"
    Entonces veo el error "no se puede pasar de Pendiente a Pendiente"
