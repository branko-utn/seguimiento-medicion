# language: es
@HU-29 @pendiente
Característica: Reporte de sprint
  Como líder del proyecto
  Quiero generar el reporte de un sprint
  Para presentarlo en la Review

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título         | story points |
      | HU-1     | Crear proyecto | 3            |
      | HU-2     | Crear sprint   | 5            |
    Y que la historia "HU-1" tiene 4 horas estimadas
    Y que la historia "HU-2" tiene 6 horas estimadas
    Y hoy es "2026-10-27"
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2"

  @CA-29.1 @normal
  Escenario: reporte de un sprint cerrado
    Dado que "Ana Pérez" registró 3 horas de "Desarrollo" en "HU-1" el "2026-10-21"
    Y que "Ana Pérez" registró 2 horas de "Pruebas" en "HU-1" el "2026-10-22"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y registro el defecto "Falla al guardar" de severidad "Alta" en "HU-1" detectado en "S1"
    Y cierro el sprint "S1"
    Cuando genero el reporte del sprint "S1"
    Entonces el reporte tiene las historias:
      | historia | story points | horas estimadas | horas reales | completada |
      | HU-1     | 3            | 4.0             | 5.0          | sí         |
      | HU-2     | 5            | 6.0             | 0.0          | no         |
    Y el reporte muestra 8 Story Points planificados y 3 completados
    Y el reporte muestra 50.0 % de historias completadas
    Y el reporte muestra 10.0 horas estimadas, 5.0 reales y una desviación de -50.0 %
    Y el reporte muestra 5.0 horas de "Ana Pérez"
    Y el reporte muestra 1 defecto detectado y 0 resueltos

  @CA-29.1 @limite
  Escenario: las horas fuera del sprint no se suman
    Dado que "Ana Pérez" registró 3 horas de "Desarrollo" en "HU-1" el "2026-10-21"
    Y que "Ana Pérez" registró 2 horas de "Desarrollo" en "HU-1" el "2026-10-26"
    Cuando genero el reporte del sprint "S1"
    Entonces el reporte muestra 3.0 horas reales
    Pero la historia "HU-1" tiene 5 horas reales

  @CA-29.1 @alternativo
  Escenario: reporte parcial de un sprint activo
    Dado que la historia "HU-2" ya está completada en el sprint "S1"
    Cuando genero el reporte del sprint "S1"
    Entonces el reporte es parcial
    Y el reporte muestra 8 Story Points planificados y 5 completados

  @error
  Escenario: reporte de un sprint que no empezó
    Dado que existe un sprint "S2" en "Demo"
    Cuando genero el reporte del sprint "S2"
    Entonces veo el error "el sprint todavía no empezó"

  @error
  Escenario: reporte de un sprint inexistente
    Cuando genero el reporte del sprint "S9"
    Entonces veo el error "el sprint no existe"
