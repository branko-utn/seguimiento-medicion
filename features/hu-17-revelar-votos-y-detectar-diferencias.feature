# language: es
@HU-17 @pendiente
Característica: Revelar votos y detectar diferencias
  Como líder del proyecto
  Quiero revelar los votos y ver si hay diferencias
  Para que se discutan

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y que "Beto Gómez" es integrante del proyecto "Demo"
    Y que "Carla Ruiz" es integrante del proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y se abrió la votación de la historia "HU-1"

  @CA-17.1 @normal
  Escenario: revelar con consenso
    Dado que "Ana Pérez" votó 5 en la historia "HU-1"
    Y que "Beto Gómez" votó 5 en la historia "HU-1"
    Cuando revelo los votos de la historia "HU-1"
    Entonces los votos de la historia "HU-1" son:
      | integrante | voto |
      | Ana Pérez  | 5    |
      | Beto Gómez | 5    |
    Y hay consenso en 5

  @CA-17.1 @CA-17.2 @normal
  Escenario: revelar con diferencias
    Dado que "Ana Pérez" votó 3 en la historia "HU-1"
    Y que "Beto Gómez" votó 8 en la historia "HU-1"
    Y que "Carla Ruiz" votó 5 en la historia "HU-1"
    Cuando revelo los votos de la historia "HU-1"
    Entonces no hay consenso
    Y el voto más alto es 8, de "Beto Gómez"
    Y el voto más bajo es 3, de "Ana Pérez"

  @CA-17.2 @limite
  Escenario: varios votaron el valor más alto
    Dado que "Ana Pérez" votó 8 en la historia "HU-1"
    Y que "Beto Gómez" votó 8 en la historia "HU-1"
    Y que "Carla Ruiz" votó 2 en la historia "HU-1"
    Cuando revelo los votos de la historia "HU-1"
    Entonces el voto más alto es 8, de "Ana Pérez, Beto Gómez"
    Y el voto más bajo es 2, de "Carla Ruiz"

  @limite
  Escenario: un solo voto
    Dado que "Ana Pérez" votó 3 en la historia "HU-1"
    Cuando revelo los votos de la historia "HU-1"
    Entonces hay consenso en 3

  @error
  Escenario: revelar sin votos
    Cuando revelo los votos de la historia "HU-1"
    Entonces veo el error "no hay votos para revelar"

  @error
  Escenario: votar después de revelar
    Dado que "Ana Pérez" votó 3 en la historia "HU-1"
    Y que se revelaron los votos de la historia "HU-1"
    Cuando "Beto Gómez" vota 5 en la historia "HU-1"
    Entonces veo el error "la ronda ya fue revelada"
