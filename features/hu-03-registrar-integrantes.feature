# language: es
@HU-03
Característica: Registrar integrantes
  Como líder del proyecto
  Quiero registrar los integrantes del proyecto con su rol
  Para asignarles trabajo, votos y horas

  Antecedentes:
    Dado que existe un proyecto "Demo"

  @CA-03.1 @normal
  Escenario: registrar un integrante
    Cuando registro a "Ana Pérez" como "Product Builder" en el proyecto "Demo"
    Entonces el proyecto "Demo" tiene los integrantes:
      | nombre    | rol             |
      | Ana Pérez | Product Builder |

  @CA-03.3 @normal
  Escenario: consultar los integrantes
    Dado que el proyecto "Demo" tiene los integrantes:
      | nombre      | rol               |
      | Ana Pérez   | Agile Enabler     |
      | Beto Gómez  | Product Builder   |
      | Carla Ruiz  | Product Architect |
    Entonces el proyecto "Demo" tiene 3 integrantes

  @CA-03.3 @limite
  Escenario: proyecto sin integrantes
    Entonces el proyecto "Demo" tiene 0 integrantes

  @CA-03.1 @error
  Escenario: nombre vacío
    Cuando registro a "   " como "Product Builder" en el proyecto "Demo"
    Entonces veo el error "el nombre del integrante es obligatorio"

  @CA-03.1 @error
  Escenario: rol inválido
    Cuando registro a "Ana Pérez" como "Tester" en el proyecto "Demo"
    Entonces veo el error "el rol debe ser Agile Enabler, Product Builder o Product Architect"

  @CA-03.2 @error
  Escenario: integrante duplicado
    Dado que "Ana Pérez" es integrante del proyecto "Demo"
    Cuando registro a "ana pérez" como "Product Builder" en el proyecto "Demo"
    Entonces veo el error "el integrante ya está registrado en el proyecto"

  @CA-03.2 @alternativo
  Escenario: mismo integrante en otro proyecto
    Dado que existe un proyecto "Otro"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Cuando registro a "Ana Pérez" como "Product Builder" en el proyecto "Otro"
    Entonces el proyecto "Otro" tiene 1 integrantes

  @error
  Escenario: proyecto inexistente
    Cuando registro a "Ana Pérez" como "Product Builder" en el proyecto "No existe"
    Entonces veo el error "el proyecto no existe"
