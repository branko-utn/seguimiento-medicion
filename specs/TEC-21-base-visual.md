# TEC-21 · Base visual

- **Issue:** #115 · **Responsable:** Federico Sosa · **Story Points:** 2
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs del Sprint 4
- **Depende de:** nada
- **Estado:** implementada (Sprint 4)

## Objetivo

Hoy la interfaz es Pico CSS casi sin ajustes: se ve básica y cada pantalla resuelve lo visual a su manera (Planning
Poker con estilos en línea y emojis, el dashboard con punto decimal, estados como texto plano). TEC-21 deja en `main`
los componentes visuales comunes, para que en HU-32 cada responsable mejore **sus** pantallas usando las mismas
piezas, sin esperar a nadie y sin tocar las plantillas de otros. Es lo mismo que hizo TEC-20 con las funciones.

## Componentes (clases en `web/static/app.css`)

| Componente | Clases | Uso |
|---|---|---|
| Paleta | variables `--sme-*` sobre las de Pico (`--pico-primary` etc.), en modo claro | Un color de marca, uno de éxito, uno de alerta y uno de peligro, con contraste AA |
| Cabecera | `header.container` con la marca y la navegación; la sección activa resaltada (ya existe `aria-current`) | En todas las páginas (base.html) |
| Encabezado de página | `.encabezado` (título, subtítulo y acciones a la derecha) | Arriba de cada pantalla |
| Tarjetas de métricas | `.tarjetas` (grilla) y `.tarjeta-metrica` con `.valor` y `.etiqueta` | Números grandes: SP, velocidad, %, horas, defectos |
| Insignias | `.insignia` + tono: `.tono-neutro`, `.tono-info`, `.tono-exito`, `.tono-alerta`, `.tono-peligro` | Estados (historias, sprints, defectos, proyecto) y severidades |
| Estado vacío | `.vacio` | "Todavía no hay …", con un enlace a la acción que lo resuelve |
| Tablas | `table.tabla` (rayada, encabezado fijo, números alineados a la derecha con `.num`) | Todas las listas |
| Acciones en línea | `form.en-linea` (ya existe) y `.acciones` | Botones dentro de tablas |
| Impresión | `@media print` (ya existe, HU-29/HU-31) | Reportes y dashboard |

## Reglas

- **El tono de una insignia lo decide el handler** con el enum (ADR 0003): la plantilla recibe la clase ya resuelta
  (`Tono string`), no compara textos. TEC-21 agrega en `internal/server` la función `tonoEstado` para los estados de
  historia, sprint, defecto y proyecto, y `tonoSeveridad` para defectos, con tests.
- **Formato de números:** `unDecimal` (ya existe) se usa en todas las pantallas; nada de `printf "%.1f"` en
  plantillas (así se resuelve el defecto del punto decimal).
- **Sin estilos en línea ni dependencias nuevas.** Todo en `app.css`, sin frameworks extra (ADR 0001).
- **Paso BDD compartido:** TEC-21 agrega `la página no contiene "…"` en `features/http_steps_test.go`, que usan los
  escenarios de HU-32. Así ninguna parte tiene que agregarlo ni esperar a otra.
- **La portada** (`inicio.html`) pasa a mostrar las secciones como tarjetas con un ícono SVG simple por sección.

## Criterio de terminado

- [ ] Componentes y paleta en `app.css`, con una página de muestra en el test de `inicio.html`
- [ ] `tonoEstado` y `tonoSeveridad` con tests de tabla
- [ ] Paso `la página no contiene "…"` en verde
- [ ] Portada y cabecera con el estilo nuevo, probadas en la app
- [ ] Mergeado antes del primer PR de HU-32
