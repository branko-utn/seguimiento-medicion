# language: es
@HU-28 @pendiente
Característica: Dashboard del proyecto
  Como integrante del equipo
  Quiero un dashboard con el estado del proyecto y gráficos
  Para ver todo de un vistazo

  Antecedentes:
    Dado que existe un proyecto "Demo" del "2026-10-21" al "2026-11-02"
    Y que "Ana Pérez" es integrante del proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |
    Y hoy es "2026-10-27"

  @CA-28.1 @normal
  Escenario: resumen y métricas principales
    Dado existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"
    Y existe un sprint activo "S2" con las historias "HU-2"
    Cuando veo el dashboard de "Demo"
    Entonces el dashboard muestra el estado "En curso" con 6 días restantes
    Y el dashboard muestra el sprint activo "S2"
    Y el dashboard muestra una velocidad de 3.0
    Y el dashboard muestra 33.3 % de historias completadas

  @CA-28.2 @normal
  Escenario: series por sprint cerrado
    Dado que la historia "HU-1" tiene 4 horas estimadas
    Y que la historia "HU-2" tiene 6 horas estimadas
    Y existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que "Ana Pérez" registró 5 horas de "Desarrollo" en "HU-1" el "2026-10-21"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y registro el defecto "Falla al guardar" de severidad "Alta" en "HU-1" detectado en "S1"
    Y cierro el sprint "S1"
    Cuando veo el dashboard de "Demo"
    Entonces el dashboard tiene los puntos por sprint:
      | sprint | sp planificados | sp completados | horas estimadas | horas reales | defectos detectados | defectos resueltos |
      | S1     | 8               | 3              | 10.0            | 5.0          | 1                   | 0                  |

  @CA-28.3 @limite
  Escenario: proyecto sin sprints cerrados
    Dado existe un sprint activo "S1" con las historias "HU-1"
    Cuando veo el dashboard de "Demo"
    Entonces los gráficos muestran el mensaje "Todavía no hay sprints cerrados para graficar"
    Y el dashboard muestra la desviación como "no calculable"

  @CA-28.3 @alternativo
  Escenario: sprints cerrados sin defectos
    Dado que la historia "HU-1" tiene 4 horas estimadas
    Y existe un sprint activo "S1" con las historias "HU-1"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"
    Cuando veo el dashboard de "Demo"
    Entonces el gráfico de defectos muestra el mensaje "No hay defectos registrados"
    Y el gráfico de horas se dibuja

  @error
  Escenario: proyecto inexistente
    Cuando veo el dashboard de "Otro"
    Entonces veo el error "el proyecto no existe"
