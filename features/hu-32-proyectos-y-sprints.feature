# language: es
@HU-32 @HU-32.1 @pendiente
Característica: Interfaz de proyectos y sprints
  Como integrante del equipo
  Quiero ver proyectos y sprints con estados y métricas claros
  Para entender cómo vamos de un vistazo

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
    Y existe un sprint activo "S1" con las historias "HU-1"

  @CA-32.1 @normal
  Escenario: el estado del proyecto y del sprint se ven con insignias
    Cuando hago un pedido "GET" a "/proyectos/P1"
    Entonces la página contiene "insignia"
    Cuando hago un pedido "GET" a "/sprints/S1"
    Entonces la página contiene "insignia"

  @CA-32.2 @normal
  Escenario: las métricas del sprint se ven en tarjetas
    Cuando hago un pedido "GET" a "/sprints/S1"
    Entonces la página contiene "tarjeta-metrica"
