# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-25 @pendiente
Característica: Porcentaje de historias completadas
  Como líder del proyecto
  Quiero ver el porcentaje de historias completadas
  Para conocer el avance

  @CA-25.1 @normal @limite
  Esquema del escenario: porcentaje de completadas
    Cuando calculo el porcentaje con <completadas> historias completadas de <total>
    Entonces el porcentaje de historias completadas es <porcentaje>

    Ejemplos:
      | completadas | total | porcentaje |
      | 3           | 4     | 75.0       |
      | 1           | 3     | 33.3       |
      | 5           | 5     | 100.0      |
      | 0           | 5     | 0.0        |

  @CA-25.2 @limite
  Escenario: sin historias
    Cuando calculo el porcentaje con 0 historias completadas de 0
    Entonces el porcentaje de historias completadas es 0.0

  @error
  Escenario: más completadas que el total
    Cuando calculo el porcentaje con 6 historias completadas de 5
    Entonces veo el error "las historias completadas no pueden superar al total"
