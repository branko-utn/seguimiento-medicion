# language: es
@HU-27 @pendiente
Característica: Defectos detectados y resueltos
  Como líder del proyecto
  Quiero ver la cantidad de defectos detectados y resueltos
  Para medir la calidad

  Antecedentes:
    Dado los defectos:
      | defecto | detectado en | resuelto en |
      | D1      | S1           | S1          |
      | D2      | S1           | S2          |
      | D3      | S2           |             |

  @CA-27.1 @normal
  Escenario: defectos por sprint
    Entonces en el sprint "S1" se detectaron 2 defectos y se resolvieron 1
    Y en el sprint "S2" se detectaron 1 defectos y se resolvieron 1

  @CA-27.1 @normal
  Escenario: totales del proyecto
    Entonces en total se detectaron 3 defectos y se resolvieron 2

  @CA-27.1 @limite
  Escenario: sprint sin defectos
    Entonces en el sprint "S3" se detectaron 0 defectos y se resolvieron 0
