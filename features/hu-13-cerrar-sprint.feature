# language: es
@HU-13
Característica: Cerrar sprint
  Como líder del proyecto
  Quiero cerrar un sprint
  Para registrar su resultado

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2, HU-3"

  @CA-13.1 @CA-13.2 @normal
  Escenario: cerrar un sprint con historias sin terminar
    Dado que la historia "HU-1" ya está completada en el sprint "S1"
    Y que la historia "HU-2" está en estado "En progreso"
    Cuando cierro el sprint "S1"
    Entonces el sprint "S1" está en estado "Cerrado"
    Y el sprint "S1" cerró con 10 Story Points planificados y 3 completados
    Y la historia "HU-2" está en estado "Pendiente"
    Y la historia "HU-2" no está asignada a ningún sprint
    Y la historia "HU-3" no está asignada a ningún sprint
    Y la historia "HU-1" está en estado "Completada"

  @CA-13.2 @limite
  Escenario: cerrar un sprint con todas las historias completadas
    Dado que la historia "HU-1" ya está completada en el sprint "S1"
    Y que la historia "HU-2" ya está completada en el sprint "S1"
    Y que la historia "HU-3" ya está completada en el sprint "S1"
    Cuando cierro el sprint "S1"
    Entonces el sprint "S1" cerró con 10 Story Points planificados y 10 completados

  @CA-13.2 @limite
  Escenario: cerrar un sprint sin historias
    Dado que el sprint "S1" está cerrado
    Y que existe un sprint activo "S2" en "Demo"
    Cuando cierro el sprint "S2"
    Entonces el sprint "S2" cerró con 0 Story Points planificados y 0 completados

  @CA-13.3 @error
  Escenario: cerrar un sprint ya cerrado
    Dado que el sprint "S1" está cerrado
    Cuando cierro el sprint "S1"
    Entonces veo el error "el sprint está cerrado y no se puede modificar"

  @error
  Escenario: cerrar un sprint que no empezó
    Dado que existe un sprint "S2" en "Demo"
    Cuando cierro el sprint "S2"
    Entonces veo el error "solo se puede cerrar un sprint activo"
