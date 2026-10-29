# language: es
@HU-21
Característica: Consultar el esfuerzo registrado
  Como integrante del equipo
  Quiero ver el esfuerzo registrado por historia y por integrante
  Para saber en qué se fue el tiempo

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y que "Beto Gómez" es integrante del proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y que el backlog de "Demo" tiene la historia "HU-2" "Crear sprint"
    Y hoy es "2026-10-27"

  @CA-21.1 @normal
  Escenario: totales por historia y por integrante
    Dado que "Ana Pérez" registró 3 horas de "Desarrollo" en "HU-1" el "2026-10-22"
    Y que "Beto Gómez" registró 2.5 horas de "Pruebas" en "HU-1" el "2026-10-23"
    Y que "Ana Pérez" registró 4 horas de "Desarrollo" en "HU-2" el "2026-10-24"
    Cuando consulto el esfuerzo registrado en "Demo"
    Entonces los totales por historia son:
      | historia | horas |
      | HU-1     | 5.5   |
      | HU-2     | 4.0   |
    Y los totales por integrante son:
      | integrante | horas |
      | Ana Pérez  | 7.0   |
      | Beto Gómez | 2.5   |
    Y el esfuerzo total del proyecto es 9.5 horas

  @CA-21.1 @normal
  Escenario: los registros se listan del más reciente al más antiguo
    Dado que "Ana Pérez" registró 3 horas de "Desarrollo" en "HU-1" el "2026-10-22"
    Y que "Beto Gómez" registró 2 horas de "Pruebas" en "HU-2" el "2026-10-25"
    Cuando consulto el esfuerzo registrado en "Demo"
    Entonces los registros son:
      | fecha      | integrante | historia | actividad  | horas |
      | 2026-10-25 | Beto Gómez | HU-2     | Pruebas    | 2.0   |
      | 2026-10-22 | Ana Pérez  | HU-1     | Desarrollo | 3.0   |

  @CA-21.1 @limite
  Escenario: proyecto sin esfuerzo registrado
    Cuando consulto el esfuerzo registrado en "Demo"
    Entonces no hay esfuerzo registrado
    Y el esfuerzo total del proyecto es 0.0 horas

  @error
  Escenario: proyecto inexistente
    Cuando consulto el esfuerzo registrado en "Otro"
    Entonces veo el error "el proyecto no existe"
