# language: es
@HU-26 @pendiente
Característica: Horas estimadas, reales y desviación
  Como líder del proyecto
  Quiero comparar las horas estimadas con las reales
  Para mejorar nuestras estimaciones

  @CA-26.1 @normal
  Escenario: horas reales
    Dado registros de esfuerzo de 3, 2.5 y 4 horas
    Entonces las horas reales son 9.5

  @CA-26.1 @limite
  Escenario: sin registros de esfuerzo
    Dado que no hay registros de esfuerzo
    Entonces las horas reales son 0.0

  @CA-26.2 @normal @limite
  Esquema del escenario: desviación
    Cuando calculo la desviación con <estimadas> horas estimadas y <reales> reales
    Entonces la desviación es <desviacion> %

    Ejemplos:
      | estimadas | reales | desviacion |
      | 10        | 12     | 20.0       |
      | 10        | 8      | -20.0      |
      | 10        | 10     | 0.0        |
      | 8         | 0      | -100.0     |

  @CA-26.3 @error
  Escenario: sin horas estimadas
    Cuando calculo la desviación con 0 horas estimadas y 5 reales
    Entonces veo el error "la desviación no es calculable sin horas estimadas"
