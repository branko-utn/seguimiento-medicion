# language: es
@HU-01
Característica: Crear proyecto
  Como líder del proyecto
  Quiero crear un proyecto con nombre, descripción y fechas de inicio y fin
  Para empezar a administrarlo

  @CA-01.4 @normal
  Escenario: crear un proyecto válido
    Cuando creo el proyecto "Software Metrics" del "2026-10-21" al "2026-11-02"
    Entonces el proyecto "Software Metrics" existe
    Y el proyecto "Software Metrics" está en estado "Planificado"

  @CA-01.1 @limite
  Esquema del escenario: largo del nombre
    Cuando creo un proyecto con un nombre de <largo> caracteres del "2026-10-21" al "2026-11-02"
    Entonces <resultado>

    Ejemplos:
      | largo | resultado                                                      |
      | 3     | el proyecto se crea                                            |
      | 100   | el proyecto se crea                                            |
      | 2     | veo el error "el nombre debe tener entre 3 y 100 caracteres"   |
      | 101   | veo el error "el nombre debe tener entre 3 y 100 caracteres"   |

  @CA-01.1 @alternativo
  Escenario: el nombre se guarda sin espacios sobrantes
    Cuando creo el proyecto "  Demo  " del "2026-10-21" al "2026-11-02"
    Entonces el proyecto "Demo" existe

  @CA-01.2 @error
  Escenario: fin anterior al inicio
    Cuando creo el proyecto "Demo" del "2026-11-02" al "2026-10-21"
    Entonces veo el error "la fecha de fin no puede ser anterior a la de inicio"

  @CA-01.2 @limite
  Escenario: inicio y fin el mismo día
    Cuando creo el proyecto "Demo" del "2026-10-21" al "2026-10-21"
    Entonces el proyecto "Demo" existe

  @CA-01.3 @error
  Escenario: nombre duplicado
    Dado que existe un proyecto "Demo"
    Cuando creo el proyecto "demo" del "2026-10-21" al "2026-11-02"
    Entonces veo el error "ya existe un proyecto con ese nombre"
