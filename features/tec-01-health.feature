# language: es
@TEC-01
Característica: Estado del servicio
  Como integrante del equipo
  Quiero consultar si la aplicación está levantada
  Para verificarla en el CI y antes de la demo

  @CA-TEC01.1 @normal
  Escenario: el servicio está levantado
    Cuando hago un pedido "GET" a "/health"
    Entonces la respuesta tiene código 200
    Y el cuerpo es el JSON {"estado":"ok"}

  @CA-TEC01.2 @error
  Escenario: método no permitido
    Cuando hago un pedido "POST" a "/health"
    Entonces la respuesta tiene código 405
