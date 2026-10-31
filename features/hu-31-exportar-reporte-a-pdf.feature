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
    Y hoy es "2026-10-31"
    Y existe un sprint activo "S1" con las historias "HU-1"

  @CA-31.1 @normal
  Esquema del escenario: las páginas de reporte se pueden exportar
    Cuando hago un pedido "GET" a "<ruta>"
    Entonces la respuesta tiene código 200
    Y la página contiene "Imprimir o guardar como PDF"
    Y la página contiene "no-imprimir"

    Ejemplos:
      | ruta                           |
      | /reportes/sprints/S1           |
      | /reportes/proyecto?proyecto=P1 |
      | /dashboard?proyecto=P1         |

  @CA-31.1 @normal
  Escenario: la exportación muestra la fecha de generación
    Cuando hago un pedido "GET" a "/reportes/proyecto?proyecto=P1"
    Entonces la página contiene "Generado el"

  @error
  Escenario: sin reporte no hay exportación
    Cuando hago un pedido "GET" a "/reportes/sprints/S9"
    Entonces la respuesta tiene código 404
