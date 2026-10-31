# language: es
@HU-32 @HU-32.4
Característica: Interfaz del dashboard
  Como integrante del equipo
  Quiero el dashboard con las métricas en tarjetas y los gráficos ordenados
  Para presentar el estado del proyecto de un vistazo

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
      | HU-2     | Crear sprint   | 5            |
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"

  @CA-32.2 @normal
  Escenario: el resumen se ve en tarjetas y con coma decimal
    Cuando hago un pedido "GET" a "/dashboard?proyecto=P1"
    Entonces la página contiene "tarjeta-metrica"
    Y la página contiene "3,0"
    Y la página no contiene "3.0"

  @CA-32.2 @limite
  Escenario: un porcentaje con decimales se muestra con coma
    Cuando hago un pedido "GET" a "/dashboard?proyecto=P1"
    Entonces la página contiene "37,5 %"
    Y la página no contiene "37.5"

  @CA-32.1 @normal
  Escenario: el estado del proyecto se ve con una insignia
    Cuando hago un pedido "GET" a "/dashboard?proyecto=P1"
    Entonces la página contiene "insignia"

  @CA-32.2 @normal
  Escenario: los cuatro gráficos se ven en una grilla de dos columnas
    Cuando hago un pedido "GET" a "/dashboard?proyecto=P1"
    Entonces la página contiene "grilla-graficos"
    Y la página no contiene "style="

  @CA-32.3 @error
  Escenario: un proyecto que no existe sigue dando el mismo error
    Cuando hago un pedido "GET" a "/dashboard?proyecto=P9"
    Entonces la respuesta tiene código 404
    Y la página contiene "el proyecto no existe"
