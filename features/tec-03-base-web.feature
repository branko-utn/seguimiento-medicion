# language: es
@TEC-03
Característica: Base web
  Como integrante del equipo
  Quiero una aplicación web con navegación común
  Para llegar a cada sección del proyecto desde cualquier pantalla

  @CA-TEC03.1 @normal
  Esquema del escenario: la navegación lleva a cada sección
    Cuando hago un pedido "GET" a "/"
    Entonces la respuesta tiene código 200
    Y la página contiene un enlace a "<ruta>" con el texto "<sección>"

    Ejemplos:
      | sección        | ruta            |
      | Proyectos      | /proyectos      |
      | Backlog        | /backlog        |
      | Sprints        | /sprints        |
      | Planning Poker | /planning-poker |
      | Esfuerzo       | /esfuerzo       |
      | Defectos       | /defectos       |
      | Dashboard      | /dashboard      |
      | Reportes       | /reportes       |

  # Con HU-28 ya no queda ninguna sección en construcción: el aviso lo cubre TestSeccionEnConstruccion_avisaElSprint.
  @CA-TEC03.2 @alternativo
  Escenario: el dashboard ya no está en construcción
    Cuando hago un pedido "GET" a "/dashboard"
    Entonces la respuesta tiene código 200
    Y la página contiene "Dashboard"

  @CA-TEC03.3 @normal
  Escenario: los recursos estáticos se sirven sin internet
    Cuando hago un pedido "GET" a "/static/vendor/htmx.min.js"
    Entonces la respuesta tiene código 200

  @CA-TEC03.4 @error
  Escenario: ruta inexistente
    Cuando hago un pedido "GET" a "/no-existe"
    Entonces la respuesta tiene código 404
