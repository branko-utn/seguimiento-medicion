# language: es
@HU-09
Característica: Ver el backlog ordenado por prioridad
  Como integrante del equipo
  Quiero ver el backlog ordenado por prioridad
  Para saber qué sigue

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título            | story points | prioridad |
      | HU-1     | Registrar defecto | 2            | Baja      |
      | HU-2     | Crear sprint      | 5            | Alta      |
      | HU-3     | Ver dashboard     | 8            | Media     |
      | HU-4     | Crear proyecto    | 3            | Alta      |

  @CA-09.1 @normal
  Escenario: ordenar por prioridad e identificador
    Cuando veo el backlog de "Demo"
    Entonces el backlog se muestra en el orden "HU-2, HU-4, HU-3, HU-1"

  @CA-09.1 @limite
  Escenario: a igual prioridad, el identificador se compara como número
    Cuando ordeno las historias:
      | historia | prioridad |
      | HU-10    | Alta      |
      | HU-1     | Media     |
      | HU-2     | Alta      |
    Entonces quedan en el orden "HU-2, HU-10, HU-1"

  @CA-09.2 @normal
  Escenario: filtrar por estado
    Dado que la historia "HU-3" está en estado "En progreso"
    Cuando veo el backlog de "Demo" con estado "Pendiente"
    Entonces el backlog se muestra en el orden "HU-2, HU-4, HU-1"

  @CA-09.2 @limite
  Escenario: ninguna historia en el estado pedido
    Cuando veo el backlog de "Demo" con estado "Completada"
    Entonces el backlog se muestra vacío

  @error
  Escenario: proyecto inexistente
    Cuando veo el backlog de "Otro"
    Entonces veo el error "el proyecto no existe"
