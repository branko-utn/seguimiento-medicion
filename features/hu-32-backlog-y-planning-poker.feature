# language: es
@HU-32 @HU-32.2
Característica: Interfaz del backlog y del Planning Poker
  Como integrante del equipo
  Quiero ver prioridades y estados con colores y votar con cartas
  Para trabajar el backlog sin esfuerzo

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |

  @CA-32.1 @normal
  Escenario: la prioridad y el estado de las historias se ven con insignias
    Cuando hago un pedido "GET" a "/backlog"
    Entonces la página contiene "insignia"

  @CA-32.1 @normal
  Escenario: el Planning Poker no usa estilos en línea
    Cuando hago un pedido "GET" a "/planning-poker"
    Entonces la respuesta tiene código 200
    Y la página no contiene "style="
