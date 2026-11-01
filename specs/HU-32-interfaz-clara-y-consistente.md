# HU-32 · Interfaz clara y consistente

- **Issues:** uno por parte (se crean con la aprobación) · **Épica:** transversal · **Story Points:** 8 (2 por parte)
- **Escenarios BDD:** uno por parte, en `features/hu-32-<parte>.feature`
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs del Sprint 4
- **Depende de:** TEC-21 (base visual), que entra antes. Cada parte depende solo de TEC-21 y de sus propias pantallas.
- **Estado:** implementada (Sprint 4)

## Historia

**Como** integrante del equipo **quiero** una interfaz clara y consistente **para** entender el estado del proyecto
de un vistazo y presentarlo sin vergüenza.

## Objetivo

Que todas las pantallas usen los mismos componentes de TEC-21: estados con insignias de color, métricas en tarjetas,
tablas prolijas, estados vacíos con una acción y números con coma decimal. Es la cara del producto en la
presentación del lunes.

## Partes y responsables

Cada uno toca **solo las plantillas y los handlers de su épica**. Así no hay conflictos entre personas.

| Parte | Responsable | Plantillas | Handlers |
|---|---|---|---|
| HU-32.1 · Proyectos y sprints | Federico Sosa | `proyectos`, `proyecto`, `proyecto_form`, `sprints`, `sprint`, `sprint_form` | `proyectos.go`, `sprints.go`, `acciones_sprint.go` |
| HU-32.2 · Backlog y Planning Poker | Ignacio Bianchi | `backlog`, `historia`, `historia_form`, `planning_poker` | `backlog.go`, `planning_poker.go` |
| HU-32.3 · Esfuerzo y defectos | Agustín Lara | `esfuerzo`, `defectos` | `esfuerzo.go`, `defectos.go` |
| HU-32.4 · Dashboard | Agustín Giorlando | `dashboard`, `dashboard_seleccionar` | `dashboard.go` |

La portada, la cabecera, la página de error y la de "en construcción" los hace Sosa en TEC-21. Los reportes suman
el botón "Descargar PDF" en HU-31 (Giorlando).

## Entradas

HU-32 no agrega datos nuevos: cada pantalla recibe lo mismo que ya recibía (los modelos de sus handlers). Lo único
nuevo que entra a las plantillas es lo que el handler resuelve para mostrar:

| Dato | Tipo | De dónde sale |
|---|---|---|
| Tono de una insignia | `string` (`tono-neutro`, `tono-info`, `tono-exito`, `tono-alerta`, `tono-peligro`) | `tonoEstado`, `tonoSeveridad` y `tonoPrioridad` a partir del enumerador (ADR 0003) |
| Números con un decimal | `string` con coma decimal | `unDecimal` |

## Salidas esperadas

- Las mismas pantallas, con los componentes de TEC-21: insignias de color para estados, severidades y prioridades;
  métricas en tarjetas; tablas `table.tabla` con números alineados; estados vacíos con un enlace a la acción.
- Ningún cambio de comportamiento: las mismas rutas, formularios, validaciones y respuestas (RN-32.6).

## Reglas de negocio

- **RN-32.1:** todo estado se muestra con `.insignia` y el tono que decide el handler con `tonoEstado` /
  `tonoSeveridad` (TEC-21, ADR 0003). La plantilla no compara textos de estado.
- **RN-32.2:** las métricas de una pantalla (SP, velocidad, %, horas, defectos) se muestran en `.tarjetas` con
  `.tarjeta-metrica`.
- **RN-32.3:** los números con decimales usan coma (`unDecimal`); no queda `printf "%.1f"` en las plantillas.
- **RN-32.4:** toda lista vacía usa `.vacio`, con un enlace o botón a la acción que la llena.
- **RN-32.5:** no quedan estilos en línea (`style="…"`) ni emojis como íconos.
- **RN-32.6:** no cambia ningún comportamiento: mismas rutas, mismos formularios, mismos mensajes. Los tests y
  escenarios existentes siguen en verde sin tocarlos, salvo los que comparan HTML literal que cambió de forma.

## Restricciones

- No se agregan clases a `app.css` fuera de TEC-21. Si una parte necesita una, la agrega **al final** del archivo,
  bajo un comentario con su parte (`/* HU-32.3 */`); en un conflicto se quedan los dos lados.
- `tonoEstado` y `tonoSeveridad` ya están en `main` (TEC-21): cada parte agrega un campo `Tono` a sus filas.
- Pantalla por pantalla, lo mínimo esperado:
  - **HU-32.1:** detalle del proyecto con el resumen (HU-04) en tarjetas; estado del proyecto y de cada sprint con
    insignia; métricas del sprint (TEC-17) en tarjetas.
  - **HU-32.2:** prioridad y estado de cada historia con insignia; el Planning Poker sin estilos en línea ni emojis,
    con las cartas como botones de la escala.
  - **HU-32.3:** severidad y estado de cada defecto con insignia; totales de esfuerzo (HU-21) en tarjetas.
  - **HU-32.4:** resumen en tarjetas, los cuatro gráficos en una grilla de dos columnas y coma decimal (resuelve el
    defecto abierto del Sprint 3).

## Casos límite

- Pantallas sin datos: muestran `.vacio`, no una tabla vacía.
- Pantallas angostas (celular): las tarjetas pasan a una columna y las tablas se desplazan en horizontal
  (`.overflow-auto`).

## Condiciones de error

No agrega errores: los mensajes y códigos de cada pantalla no cambian (RN-32.6).

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-32.1 | Los estados se ven con insignias de color | uno por parte en `features/hu-32-*.feature` |
| CA-32.2 | Las métricas se ven en tarjetas y con coma decimal | uno por parte en `features/hu-32-*.feature` |
| CA-32.3 | El comportamiento no cambia: los tests y escenarios existentes siguen en verde | toda la suite |
