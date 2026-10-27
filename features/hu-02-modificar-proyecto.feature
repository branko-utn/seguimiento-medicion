# language: es
@HU-02
Característica: Modificar proyecto
  Como líder del proyecto
  Quiero modificar los datos de un proyecto
  Para corregirlos o actualizarlos

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"

  @CA-02.1 @normal
  Escenario: modificar un proyecto
    Cuando modifico el proyecto "Demo" con el nombre "Software Metrics" del "2026-10-22" al "2026-11-03"
    Entonces el proyecto "Software Metrics" existe
    Y el proyecto "Software Metrics" está en estado "Planificado"

  @CA-02.1 @limite
  Escenario: cambiar solo las mayúsculas del nombre
    Cuando modifico el proyecto "Demo" con el nombre "DEMO" del "2026-10-21" al "2026-11-02"
    Entonces el proyecto "DEMO" existe

  @CA-02.1 @error
  Escenario: nombre inválido
    Cuando modifico el proyecto "Demo" con el nombre "AB" del "2026-10-21" al "2026-11-02"
    Entonces veo el error "el nombre debe tener entre 3 y 100 caracteres"

  @CA-02.1 @error
  Escenario: fin anterior al inicio
    Cuando modifico el proyecto "Demo" con el nombre "Demo" del "2026-11-02" al "2026-10-21"
    Entonces veo el error "la fecha de fin no puede ser anterior a la de inicio"

  @CA-02.1 @error
  Escenario: nombre de otro proyecto
    Dado que existe un proyecto "Otro"
    Cuando modifico el proyecto "Demo" con el nombre "otro" del "2026-10-21" al "2026-11-02"
    Entonces veo el error "ya existe un proyecto con ese nombre"

  @CA-02.2 @error
  Escenario: proyecto inexistente
    Cuando modifico el proyecto "No existe" con el nombre "Nuevo" del "2026-10-21" al "2026-11-02"
    Entonces veo el error "el proyecto no existe"

  @CA-02.3 @error
  Escenario: proyecto finalizado
    Dado que el proyecto "Demo" está finalizado
    Cuando modifico el proyecto "Demo" con el nombre "Demo 2" del "2026-10-21" al "2026-11-02"
    Entonces veo el error "un proyecto finalizado no se puede modificar"
