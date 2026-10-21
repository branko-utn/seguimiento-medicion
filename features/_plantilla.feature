# language: es
# Plantilla: copiar como hu-xx-<titulo>.feature. Un escenario por criterio de aceptación como mínimo,
# cubriendo caso normal, alternativo, límite y error. Etiquetar con el ID de la historia y del criterio.
@HU-XX
Característica: <Título de la historia>
  Como <rol>
  Quiero <funcionalidad>
  Para <beneficio>

  @CA-XX.1 @normal
  Escenario: <caso normal>
    Dado <contexto>
    Cuando <acción>
    Entonces <resultado esperado>

  @CA-XX.2 @error
  Escenario: <caso de error>
    Dado <contexto>
    Cuando <acción inválida>
    Entonces <error esperado>
