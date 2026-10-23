# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-20 @pendiente
Característica: Registrar esfuerzo trabajado
  Como integrante del equipo
  Quiero registrar las horas que trabajé en una historia
  Para conocer el esfuerzo real

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y hoy es "2026-10-24"

  @CA-20.1 @normal
  Escenario: registrar esfuerzo
    Cuando "Ana Pérez" registra 3 horas de "Desarrollo" en "HU-1" el "2026-10-24"
    Entonces la historia "HU-1" tiene 3 horas reales
    Y "Ana Pérez" tiene 3 horas registradas el "2026-10-24"

  @CA-20.1 @error
  Escenario: actividad vacía
    Cuando "Ana Pérez" registra 3 horas de "  " en "HU-1" el "2026-10-24"
    Entonces veo el error "la actividad es obligatoria"

  @CA-20.2 @error
  Escenario: horas en cero
    Cuando "Ana Pérez" registra 0 horas de "Desarrollo" en "HU-1" el "2026-10-24"
    Entonces veo el error "las horas deben ser mayores que cero"

  @CA-20.2 @limite
  Escenario: 24 horas justas en un día
    Dado que "Ana Pérez" registró 10 horas de "Desarrollo" en "HU-1" el "2026-10-23"
    Cuando "Ana Pérez" registra 14 horas de "Tests" en "HU-1" el "2026-10-23"
    Entonces "Ana Pérez" tiene 24 horas registradas el "2026-10-23"

  @CA-20.2 @error
  Escenario: superar 24 horas en un día
    Dado que "Ana Pérez" registró 24 horas de "Desarrollo" en "HU-1" el "2026-10-23"
    Cuando "Ana Pérez" registra 0.5 horas de "Tests" en "HU-1" el "2026-10-23"
    Entonces veo el error "el integrante superaría las 24 horas en ese día"

  @CA-20.3 @error
  Escenario: fecha futura
    Cuando "Ana Pérez" registra 2 horas de "Desarrollo" en "HU-1" el "2026-10-25"
    Entonces veo el error "la fecha no puede ser futura"

  @error
  Escenario: integrante de otro proyecto
    Dado que existe un proyecto "Otro"
    Y que "Beto Gómez" es integrante del proyecto "Otro"
    Cuando "Beto Gómez" registra 2 horas de "Desarrollo" en "HU-1" el "2026-10-24"
    Entonces veo el error "el integrante no pertenece al proyecto de la historia"
