# language: es
@HU-31 @pendiente
Característica: Exportar reporte a PDF
  Como líder del proyecto
  Quiero exportar un reporte a PDF
  Para compartirlo fuera de la aplicación

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
    Y existe un sprint activo "S1" con las historias "HU-1"

  @CA-31.1 @normal
  Escenario: descargar el informe del proyecto en PDF
    Cuando hago un pedido "GET" a "/reportes/proyecto.pdf?proyecto=P1"
    Entonces la respuesta tiene código 200
    Y la respuesta es un PDF llamado "reporte-proyecto-P1.pdf"

  @CA-31.1 @alternativo
  Escenario: descargar el informe de un sprint en PDF
    Cuando hago un pedido "GET" a "/reportes/sprints/S1/pdf"
    Entonces la respuesta tiene código 200
    Y la respuesta es un PDF llamado "reporte-sprint-S1.pdf"

  @CA-31.1 @normal
  Escenario: los reportes ofrecen la descarga
    Cuando hago un pedido "GET" a "/reportes/proyecto?proyecto=P1"
    Entonces la página contiene "Descargar PDF"

  @error
  Escenario: sin reporte no hay PDF
    Cuando hago un pedido "GET" a "/reportes/proyecto.pdf?proyecto=P9"
    Entonces la respuesta tiene código 404
