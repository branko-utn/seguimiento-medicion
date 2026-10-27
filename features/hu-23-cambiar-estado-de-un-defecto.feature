# language: es
@HU-23
Característica: Cambiar estado de un defecto y registrar su resolución
  Como integrante del equipo
  Quiero cambiar el estado de un defecto
  Para saber cuáles están resueltos

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y que en "Demo" existe un sprint "S1" del "2026-10-23" al "2026-10-25"
    Y que en "Demo" existe un sprint "S2" del "2026-10-26" al "2026-10-27"
    Y que existe el defecto "D1" en "HU-1" detectado en "S2"

  @CA-23.1 @normal
  Escenario: empezar a trabajar un defecto
    Cuando cambio el estado del defecto "D1" a "En progreso"
    Entonces existe el defecto "D1" en estado "En progreso"

  @CA-23.1 @normal
  Escenario: resolver un defecto
    Dado que el defecto "D1" está en estado "En progreso"
    Cuando resuelvo el defecto "D1" en el sprint "S2"
    Entonces existe el defecto "D1" en estado "Resuelto"
    Y el defecto "D1" se resolvió en el sprint "S2"

  @CA-23.2 @limite
  Escenario: resolver en el mismo sprint directamente desde Abierto
    Cuando resuelvo el defecto "D1" en el sprint "S2"
    Entonces existe el defecto "D1" en estado "Resuelto"

  @CA-23.2 @error
  Escenario: resolver en un sprint anterior al de detección
    Cuando resuelvo el defecto "D1" en el sprint "S1"
    Entonces veo el error "el sprint de resolución no puede ser anterior al de detección"
    Y existe el defecto "D1" en estado "Abierto"

  @CA-23.2 @error
  Escenario: resolver sin sprint
    Cuando cambio el estado del defecto "D1" a "Resuelto"
    Entonces veo el error "para resolver un defecto hay que indicar el sprint de resolución"

  @CA-23.1 @error
  Escenario: cambiar un defecto resuelto
    Dado que el defecto "D1" se resolvió en el sprint "S2"
    Cuando cambio el estado del defecto "D1" a "En progreso"
    Entonces veo el error "no se puede pasar de Resuelto a En progreso"

  @error
  Escenario: defecto inexistente
    Cuando cambio el estado del defecto "D9" a "En progreso"
    Entonces veo el error "el defecto no existe"
