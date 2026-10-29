# language: es
@HU-14
Característica: Consultar sprints anteriores
  Como integrante del equipo
  Quiero ver los sprints cerrados con sus resultados
  Para comparar el desempeño

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y el proyecto tiene las historias:
      | historia | título            | story points |
      | HU-1     | Crear proyecto    | 3            |
      | HU-2     | Crear sprint      | 5            |
      | HU-3     | Registrar defecto | 2            |

  @CA-14.1 @normal
  Escenario: listar sprints cerrados
    Dado existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"
    Y existe un sprint activo "S2" con las historias "HU-2, HU-3"
    Cuando consulto los sprints cerrados de "Demo"
    Entonces veo los sprints cerrados:
      | sprint | sprint goal      | story points planificados | story points completados |
      | S1     | Sprint de prueba | 8                         | 3                        |

  @CA-14.1 @limite
  Escenario: sin sprints cerrados
    Dado existe un sprint activo "S1" con las historias "HU-1"
    Cuando consulto los sprints cerrados de "Demo"
    Entonces no hay sprints cerrados

  @CA-14.2 @normal
  Escenario: detalle de un sprint cerrado
    Dado existe un sprint activo "S1" con las historias "HU-1, HU-2"
    Y que la historia "HU-1" ya está completada en el sprint "S1"
    Y cierro el sprint "S1"
    Cuando veo el detalle del sprint "S1"
    Entonces el detalle tiene las historias:
      | historia | story points | completada |
      | HU-1     | 3            | sí         |
      | HU-2     | 5            | no         |

  @CA-14.2 @error
  Escenario: detalle de un sprint que no está cerrado
    Dado existe un sprint activo "S1" con las historias "HU-1"
    Cuando veo el detalle del sprint "S1"
    Entonces veo el error "el sprint todavía no está cerrado"

  @error
  Escenario: detalle de un sprint inexistente
    Cuando veo el detalle del sprint "S9"
    Entonces veo el error "el sprint no existe"
