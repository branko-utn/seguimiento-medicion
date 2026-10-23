# language: es
# @pendiente: se quita al implementar los pasos (commit RED).
@HU-22 @pendiente
Característica: Registrar defecto
  Como integrante del equipo
  Quiero registrar un defecto
  Para hacer el seguimiento de la calidad

  Antecedentes:
    Dado que existe un proyecto "Demo"
    Y que el backlog de "Demo" tiene la historia "HU-1" "Crear proyecto"
    Y existe un sprint "S1" en "Demo"

  @CA-22.2 @normal
  Escenario: registrar un defecto
    Cuando registro el defecto "Acepta fechas invertidas" de severidad "Alta" en "HU-1" detectado en "S1"
    Entonces existe el defecto "D1" en estado "Abierto"
    Y el defecto "D1" no tiene sprint de resolución

  @CA-22.1 @error
  Escenario: descripción vacía
    Cuando registro el defecto "  " de severidad "Alta" en "HU-1" detectado en "S1"
    Entonces veo el error "la descripción es obligatoria"

  @CA-22.1 @error
  Escenario: severidad inválida
    Cuando registro el defecto "Acepta fechas invertidas" de severidad "Urgente" en "HU-1" detectado en "S1"
    Entonces veo el error "la severidad debe ser Crítica, Alta, Media o Baja"

  @CA-22.1 @error
  Escenario: historia inexistente
    Cuando registro el defecto "Acepta fechas invertidas" de severidad "Alta" en "HU-99" detectado en "S1"
    Entonces veo el error "la historia no existe"

  @limite
  Escenario: defecto en una historia completada
    Dado que la historia "HU-1" está completada
    Cuando registro el defecto "Acepta fechas invertidas" de severidad "Crítica" en "HU-1" detectado en "S1"
    Entonces existe el defecto "D1" en estado "Abierto"

  @error
  Escenario: sprint de otro proyecto
    Dado que existe un proyecto "Otro"
    Y existe un sprint "S9" en "Otro"
    Cuando registro el defecto "Acepta fechas invertidas" de severidad "Alta" en "HU-1" detectado en "S9"
    Entonces veo el error "el sprint de detección no pertenece al proyecto de la historia"
