# language: es
@HU-11
Característica: Asignar historias a un sprint
  Como líder del proyecto
  Quiero asignar historias del backlog a un sprint
  Para armar el Sprint Backlog

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título           | story points | criterios |
      | HU-1     | Crear proyecto   | 3            | 2         |
      | HU-2     | Crear sprint     | 0            | 1         |
      | HU-3     | Cerrar sprint    | 5            | 0         |
      | HU-4     | Ver métricas     | 2            | 1         |
    Y existe un sprint "S1" en "Demo"

  @CA-11.1 @normal
  Escenario: asignar una historia lista
    Cuando asigno la historia "HU-1" al sprint "S1"
    Entonces el sprint "S1" tiene las historias "HU-1"
    Y el sprint "S1" tiene 3 Story Points planificados

  @CA-11.1 @error
  Escenario: historia sin estimar
    Cuando asigno la historia "HU-2" al sprint "S1"
    Entonces veo el error "la historia no está estimada"

  @CA-11.1 @error
  Escenario: historia sin criterios
    Cuando asigno la historia "HU-3" al sprint "S1"
    Entonces veo el error "la historia no tiene criterios de aceptación"

  @CA-11.2 @error
  Escenario: historia ya asignada a otro sprint
    Dado que existe un sprint "S2" en "Demo"
    Y que la historia "HU-1" está asignada al sprint "S1"
    Cuando asigno la historia "HU-1" al sprint "S2"
    Entonces veo el error "la historia ya está asignada a un sprint"
    Y el sprint "S2" no tiene historias

  @CA-11.2 @limite
  Escenario: asignar dos veces la misma historia al mismo sprint
    Dado que la historia "HU-1" está asignada al sprint "S1"
    Cuando asigno la historia "HU-1" al sprint "S1"
    Entonces veo el error "la historia ya está asignada a un sprint"
    Y el sprint "S1" tiene las historias "HU-1"

  @CA-11.3 @error
  Escenario: sprint cerrado
    Dado que el sprint "S1" está cerrado
    Cuando asigno la historia "HU-4" al sprint "S1"
    Entonces veo el error "el sprint está cerrado y no se puede modificar"
