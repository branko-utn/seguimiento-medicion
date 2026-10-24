# language: es
@HU-10
Característica: Crear sprint con Sprint Goal
  Como líder del proyecto
  Quiero crear un sprint con fechas y un Sprint Goal
  Para planificar el trabajo

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"

  @CA-10.1 @normal
  Escenario: crear un sprint
    Cuando creo en "Demo" un sprint del "2026-10-23" al "2026-10-25" con el objetivo "Primer MVP"
    Entonces el proyecto "Demo" tiene un sprint "S1" en estado "Planificado" con el objetivo "Primer MVP"

  @CA-10.1 @error
  Escenario: Sprint Goal vacío
    Cuando creo en "Demo" un sprint del "2026-10-23" al "2026-10-25" con el objetivo "  "
    Entonces veo el error "el Sprint Goal es obligatorio"

  @error
  Escenario: fin anterior al inicio
    Cuando creo en "Demo" un sprint del "2026-10-25" al "2026-10-23" con el objetivo "Primer MVP"
    Entonces veo el error "la fecha de fin no puede ser anterior a la de inicio"

  @CA-10.2 @error
  Escenario: fechas fuera del proyecto
    Cuando creo en "Demo" un sprint del "2026-10-30" al "2026-11-03" con el objetivo "Cierre"
    Entonces veo el error "las fechas del sprint deben estar dentro de las del proyecto"

  @CA-10.2 @limite
  Escenario: sprint que ocupa todo el proyecto
    Cuando creo en "Demo" un sprint del "2026-10-21" al "2026-11-02" con el objetivo "Todo"
    Entonces el proyecto "Demo" tiene 1 sprints

  @CA-10.2 @error
  Escenario: sprints superpuestos
    Dado que en "Demo" existe un sprint del "2026-10-23" al "2026-10-25"
    Cuando creo en "Demo" un sprint del "2026-10-25" al "2026-10-27" con el objetivo "Interfaz"
    Entonces veo el error "las fechas se superponen con otro sprint del proyecto"

  @CA-10.2 @limite
  Escenario: sprints consecutivos
    Dado que en "Demo" existe un sprint del "2026-10-23" al "2026-10-25"
    Cuando creo en "Demo" un sprint del "2026-10-26" al "2026-10-27" con el objetivo "Interfaz"
    Entonces el proyecto "Demo" tiene 2 sprints

  @CA-10.3 @normal
  Escenario: iniciar un sprint
    Dado que en "Demo" existe un sprint "S1" del "2026-10-23" al "2026-10-25"
    Cuando inicio el sprint "S1"
    Entonces el sprint "S1" está en estado "Activo"

  @CA-10.3 @error
  Escenario: ya hay un sprint activo
    Dado que en "Demo" existe un sprint activo "S1" del "2026-10-23" al "2026-10-25"
    Y que en "Demo" existe un sprint "S2" del "2026-10-26" al "2026-10-27"
    Cuando inicio el sprint "S2"
    Entonces veo el error "ya hay un sprint activo en el proyecto"
