# language: es
@TEC-16
Característica: Pantalla del backlog
  Como líder del proyecto
  Quiero ver el Product Backlog en la aplicación web
  Para revisar las historias, su estado, prioridad, Story Points y criterios de aceptación

  Antecedentes:
    Dado que existe un proyecto "Demo"

  @CA-TEC16.1 @normal
  Escenario: ver el backlog vacío
    Cuando hago un pedido "GET" a "/backlog"
    Entonces la respuesta tiene código 200
    Y la página contiene "Backlog"

  @CA-TEC16.2 @normal
  Escenario: ver las historias del backlog
    Dado que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Cuando hago un pedido "GET" a "/backlog"
    Entonces la respuesta tiene código 200
    Y la página contiene "HU-1"
    Y la página contiene "Crear proyecto"
