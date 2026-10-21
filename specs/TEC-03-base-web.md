# TEC-03 · Base web: layout, navegación y estilos

- **Issue:** #4 · **Responsable:** Branko Almeira
- **Escenarios BDD:** [`features/tec-03-base-web.feature`](../features/tec-03-base-web.feature)
- **Estado:** implementada

## Objetivo

Que cada Product Builder pueda agregar sus pantallas desde el primer día del Sprint 1 con un layout, una navegación
y un estilo comunes, sin pisarse con los demás y sin depender de internet durante la demo.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Pedido HTTP `GET` a una ruta de la aplicación | — | Sí | — |
| Servicios de cada épica (`server.Servicios`) | struct | No | Cada épica agrega su campo |

## Salidas esperadas

- `GET /` muestra la página de inicio con la navegación a todas las secciones.
- Cada sección que todavía no está construida muestra una página que avisa en qué sprint llega.
- `GET /static/...` sirve los estilos y scripts incluidos en el repositorio (Pico CSS, HTMX, Chart.js).

## Reglas de negocio

- **RN-TEC03.1:** todas las páginas comparten el layout base (`web/templates/base.html`): encabezado, navegación y pie.
- **RN-TEC03.2:** la navegación lista las secciones definidas en `internal/server/secciones.go`; una sección marcada
  como en construcción responde con la página de aviso hasta que su responsable registre la ruta real.

## Restricciones

- Las plantillas y los archivos estáticos se embeben en el binario (`embed`): la aplicación funciona sin internet.
- El contenido dinámico se escapa siempre (`html/template`).

## Casos límite

- Ruta inexistente: `404 Not Found`.
- Una página que no existe en las plantillas: `500` con una página de error, sin mostrar detalles internos.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Las plantillas no se pueden cargar al iniciar | `NewRouter` devuelve el error y la aplicación no arranca |
| Error al renderizar una página | `500 Internal Server Error` con mensaje genérico |
| Mensaje de error con HTML (`<script>`) | Se muestra escapado |

## Criterios de aceptación

| ID | Criterio | Escenario BDD / test |
|---|---|---|
| CA-TEC03.1 | La página de inicio muestra la navegación a todas las secciones | `Esquema del escenario: la navegación lleva a cada sección` |
| CA-TEC03.2 | Una sección en construcción avisa en qué sprint llega | `Escenario: una sección en construcción avisa cuándo llega` |
| CA-TEC03.3 | Los estilos y scripts se sirven desde la aplicación | `Escenario: los recursos estáticos se sirven sin internet` |
| CA-TEC03.4 | Una ruta inexistente responde 404 | `Escenario: ruta inexistente` |
| CA-TEC03.5 | Los errores se muestran en una página propia, escapados | `TestRenderError` |
