# language: es
@HU-32 @HU-32.3 @pendiente
Característica: Interfaz de esfuerzo y defectos
  Como integrante del equipo
  Quiero ver la severidad de los defectos y los totales de esfuerzo de un vistazo
  Para saber dónde está el problema

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
    Y existe un sprint activo "S1" con las historias "HU-1"

  @CA-32.1 @normal
  Escenario: la severidad y el estado de los defectos se ven con insignias
    Dado registro el defecto "No guarda" de severidad "Crítica" en "HU-1" detectado en "S1"
    Cuando hago un pedido "GET" a "/defectos"
    Entonces la página contiene "insignia"

  @CA-32.2 @normal
  Escenario: los totales de esfuerzo se ven en tarjetas
    Dado que "Ana Pérez" es integrante del proyecto "Demo"
    Y hoy es "2026-10-31"
    Y que "Ana Pérez" registró 2.5 horas de "Desarrollo" en "HU-1" el "2026-10-22"
    Cuando hago un pedido "GET" a "/esfuerzo"
    Entonces la página contiene "tarjeta-metrica"
    Y la página contiene "2,5"
