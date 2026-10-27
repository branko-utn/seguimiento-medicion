# language: es
@HU-16 @pendiente
Característica: Votar en Planning Poker con votos ocultos
  Como integrante del equipo
  Quiero votar la estimación de una historia sin ver los votos de los demás
  Para no influirnos

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y que "Beto Gómez" es integrante del proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y se abrió la votación de la historia "HU-1"

  @CA-16.1 @normal
  Escenario: votar
    Cuando "Ana Pérez" vota 5 en la historia "HU-1"
    Entonces en la votación de "HU-1" votaron "Ana Pérez"

  @CA-16.2 @normal
  Escenario: los votos quedan ocultos
    Dado que "Ana Pérez" votó 5 en la historia "HU-1"
    Y que "Beto Gómez" votó 8 en la historia "HU-1"
    Cuando consulto los votos de la historia "HU-1"
    Entonces veo el error "los votos se muestran recién al revelar las cartas"
    Y en la votación de "HU-1" votaron "Ana Pérez, Beto Gómez"

  @limite
  Escenario: votación sin votos
    Entonces en la votación de "HU-1" no votó nadie

  @CA-16.1 @error
  Escenario: votar dos veces en la misma ronda
    Dado que "Ana Pérez" votó 5 en la historia "HU-1"
    Cuando "Ana Pérez" vota 8 en la historia "HU-1"
    Entonces veo el error "el integrante ya votó en esta ronda"

  @CA-16.1 @error
  Escenario: votar sin ser del proyecto
    Dado que existe un proyecto "Otro"
    Y que "Carla Ruiz" es integrante del proyecto "Otro"
    Cuando "Carla Ruiz" vota 5 en la historia "HU-1"
    Entonces veo el error "el integrante no pertenece al proyecto de la historia"

  @error
  Escenario: votar una carta fuera de la escala
    Cuando "Ana Pérez" vota 4 en la historia "HU-1"
    Entonces veo el error "los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21"

  @error
  Escenario: abrir otra votación para la misma historia
    Cuando abro la votación de la historia "HU-1"
    Entonces veo el error "la historia ya tiene una votación abierta"
