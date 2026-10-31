# language: es
@HU-30
Característica: Reporte de proyecto
  Como líder del proyecto
  Quiero generar el reporte del proyecto completo
  Para informar su estado

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |
    Y hoy es "2026-10-27"

  @CA-30.1 @normal
  Escenario: reporte con sprints cerrados y activo
    Dado existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"
    Y existe un sprint activo "S2" con las historias "HU-2, HU-3"
    Y que la historia "HU-2" ya está completada en el sprint "S2"
    Y existe un sprint "S3" en "Demo"
    Cuando genero el reporte del proyecto "Demo"
    Entonces el reporte del proyecto muestra el estado "En curso"
    Y el resumen de sprints es:
      | sprint | estado  | sp planificados | sp completados |
      | S1     | Cerrado | 8               | 3              |
      | S2     | Activo  | 7               | 5              |
    Y el reporte del proyecto muestra una velocidad de 3.0
    Y el reporte del proyecto muestra 66.7 % de historias completadas

  @CA-30.1 @normal
  Escenario: defectos abiertos por severidad
    Dado existe un sprint activo "S1" con las historias "HU-1"
    Y registro el defecto "Texto cortado" de severidad "Baja" en "HU-1" detectado en "S1"
    Y registro el defecto "No guarda" de severidad "Crítica" en "HU-1" detectado en "S1"
    Cuando genero el reporte del proyecto "Demo"
    Entonces el reporte del proyecto muestra 2 defectos detectados, 0 resueltos y 2 abiertos
    Y los defectos abiertos son "No guarda, Texto cortado"

  @CA-30.1 @limite
  Escenario: proyecto sin sprints
    Cuando genero el reporte del proyecto "Demo"
    Entonces el resumen de sprints está vacío
    Y el reporte del proyecto muestra una velocidad de 0.0
    Y el reporte del proyecto muestra la desviación como "no calculable"

  @error
  Escenario: proyecto inexistente
    Cuando genero el reporte del proyecto "Otro"
    Entonces veo el error "el proyecto no existe"
