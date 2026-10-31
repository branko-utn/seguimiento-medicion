# HU-31 · Exportar reporte a PDF

- **Issue:** #38 · **Épica:** E9 Reportes · **Responsable:** Branko Almeira · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-31-exportar-reporte-a-pdf.feature`](../features/hu-31-exportar-reporte-a-pdf.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs del Sprint 4
- **Depende de:** solo código en `main`: HU-29 y HU-30 (reportes), HU-28 (dashboard)
- **Estado:** aprobada para el Sprint 4

## Historia

**Como** líder del proyecto **quiero** exportar un reporte a PDF **para** compartirlo fuera de la aplicación.

## Objetivo

Poder guardar como PDF el reporte de un sprint, el del proyecto y el dashboard, con el mismo contenido que se ve en
pantalla, sin instalar nada y sin bibliotecas de PDF (ADR 0001): el PDF lo genera el navegador al imprimir.

## Entradas

| Página | Ruta |
|---|---|
| Reporte de sprint | `GET /reportes/sprints/{sprint}` |
| Reporte de proyecto | `GET /reportes/proyecto?proyecto={proyecto}` |
| Dashboard | `GET /dashboard?proyecto={proyecto}` |

## Salidas esperadas

En cada página, un botón **"Imprimir o guardar como PDF"** (`window.print()`). Al imprimir:
- no salen la navegación, el pie ni los botones (`.no-imprimir`);
- sí salen el título, la fecha de generación y todo el contenido del reporte;
- las tablas y los gráficos no se cortan a la mitad entre páginas.

## Reglas de negocio

- **RN-31.1:** el PDF tiene el mismo contenido que la página: no hay una versión aparte para exportar.
- **RN-31.2:** cada página exportable muestra al pie del encabezado "Generado el dd/mm/aaaa", con la fecha del día
  (`hoy()` del servidor), para saber de cuándo es el PDF. La fecha exacta la verifican los tests del handler, que
  fijan `hoy`; el escenario BDD usa el servidor real y verifica solo el texto.
- **RN-31.3:** los gráficos del dashboard se imprimen como imagen: Chart.js dibuja en `<canvas>`, que el navegador
  imprime tal cual. Antes de imprimir, cada gráfico se redimensiona al ancho de la hoja (`beforeprint`).

## Restricciones

- Sin dependencias nuevas (ADR 0001). Solo `web/static/app.css` (`@media print`, ya existe desde HU-29) y un script
  chico en el dashboard para `beforeprint`.
- El botón y la fecha salen de un bloque compartido `{{template "exportar" .}}` definido en `web/templates/base.html`,
  para que las tres páginas lo muestren igual.

## Casos límite

- Dashboard sin sprints cerrados: se imprime con los mensajes en lugar de los gráficos (HU-28).
- Reporte de sprint activo: el PDF lleva la marca "Reporte parcial" (HU-29).

## Condiciones de error

Las mismas de cada página: un sprint o proyecto inexistente da 404, y un sprint planificado da 400 (HU-29). En esos
casos no hay botón de exportar.

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-31.1 | El PDF tiene el mismo contenido que el reporte en pantalla | `Esquema del escenario: las páginas de reporte se pueden exportar`, `Escenario: la exportación muestra la fecha de generación` |
