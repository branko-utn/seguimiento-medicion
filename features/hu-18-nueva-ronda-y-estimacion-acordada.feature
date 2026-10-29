# language: es
@HU-18 @pendiente
Característica: Nueva ronda y estimación acordada
  Como líder del proyecto
  Quiero abrir nuevas rondas y registrar la estimación acordada
  Para cerrar la estimación de la historia

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y que "Beto Gómez" es integrante del proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y se abrió la votación de la historia "HU-1"
    Y que "Ana Pérez" votó 3 en la historia "HU-1"
    Y que "Beto Gómez" votó 8 en la historia "HU-1"

  @CA-18.1 @normal
  Escenario: nueva ronda después de revelar
    Dado que se revelaron los votos de la historia "HU-1"
    Cuando abro una nueva ronda en la historia "HU-1"
    Entonces la votación de la historia "HU-1" está en la ronda 2 sin votos
    Y el historial de la historia "HU-1" tiene la ronda 1 con los votos:
      | integrante | voto |
      | Ana Pérez  | 3    |
      | Beto Gómez | 8    |

  @CA-18.1 @error
  Escenario: nueva ronda sin revelar
    Cuando abro una nueva ronda en la historia "HU-1"
    Entonces veo el error "primero hay que revelar los votos de la ronda"
    Y la votación de la historia "HU-1" está en la ronda 1

  @CA-18.2 @normal
  Escenario: registrar la estimación acordada
    Dado que se revelaron los votos de la historia "HU-1"
    Y que se abrió una nueva ronda en la historia "HU-1"
    Y que "Ana Pérez" votó 5 en la historia "HU-1"
    Y que "Beto Gómez" votó 5 en la historia "HU-1"
    Y que se revelaron los votos de la historia "HU-1"
    Cuando registro 5 como estimación acordada de la historia "HU-1"
    Entonces la historia "HU-1" tiene 5 Story Points
    Y la historia "HU-1" no tiene una votación abierta

  @CA-18.2 @alternativo
  Escenario: acordar un valor que nadie votó
    Dado que se revelaron los votos de la historia "HU-1"
    Cuando registro 5 como estimación acordada de la historia "HU-1"
    Entonces la historia "HU-1" tiene 5 Story Points

  @CA-18.2 @error
  Escenario: acordar sin revelar
    Cuando registro 5 como estimación acordada de la historia "HU-1"
    Entonces veo el error "primero hay que revelar los votos de la ronda"
    Y la historia "HU-1" no tiene Story Points

  @CA-18.2 @error
  Escenario: acordar un valor fuera de la escala
    Dado que se revelaron los votos de la historia "HU-1"
    Cuando registro 4 como estimación acordada de la historia "HU-1"
    Entonces veo el error "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"
    Y la historia "HU-1" no tiene Story Points

  @CA-18.2 @error
  Escenario: operar sobre una votación cerrada
    Dado que se revelaron los votos de la historia "HU-1"
    Y que se registró 8 como estimación acordada de la historia "HU-1"
    Cuando abro una nueva ronda en la historia "HU-1"
    Entonces veo el error "la historia no tiene una votación abierta"
    Y la historia "HU-1" tiene 8 Story Points
