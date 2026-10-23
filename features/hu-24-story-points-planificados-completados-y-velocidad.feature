# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-24
Característica: Story Points planificados, completados y velocidad
  Como líder del proyecto
  Quiero ver los SP planificados y completados de cada sprint y la velocidad del equipo
  Para planificar mejor

  @CA-24.1 @CA-24.2 @normal
  Escenario: SP planificados y completados de un sprint
    Dado un sprint con las historias:
      | story points | completada |
      | 3            | sí         |
      | 5            | no         |
      | 2            | sí         |
    Entonces los Story Points planificados son 10
    Y los Story Points completados son 5

  @CA-24.2 @limite
  Escenario: sprint sin historias completadas
    Dado un sprint con las historias:
      | story points | completada |
      | 3            | no         |
      | 8            | no         |
    Entonces los Story Points completados son 0

  @limite
  Escenario: sprint sin historias
    Dado un sprint sin historias
    Entonces los Story Points planificados son 0
    Y los Story Points completados son 0

  @CA-24.3 @normal
  Escenario: velocidad del equipo
    Dado sprints cerrados con 10, 13 y 8 Story Points completados
    Entonces la velocidad del equipo es 10.3

  @CA-24.3 @limite
  Escenario: un sprint cerrado sin historias completadas cuenta en el promedio
    Dado sprints cerrados con 10 y 0 Story Points completados
    Entonces la velocidad del equipo es 5.0

  @CA-24.3 @alternativo
  Escenario: sin sprints cerrados
    Dado que no hay sprints cerrados
    Entonces la velocidad del equipo es 0.0
