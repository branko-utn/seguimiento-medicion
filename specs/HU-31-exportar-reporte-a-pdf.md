# HU-31 · Exportar reporte a PDF

- **Issue:** #38 · **Épica:** E9 Reportes · **Responsable:** Agustín Giorlando · **Story Points:** 5
- **Escenarios BDD:** [`features/hu-31-exportar-reporte-a-pdf.feature`](../features/hu-31-exportar-reporte-a-pdf.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs del Sprint 4
- **Depende de:** solo código en `main`: HU-29 y HU-30 (`report.ReporteSprint` y `report.ReporteProyecto`)
- **Estado:** aprobada para el Sprint 4

## Historia

**Como** líder del proyecto **quiero** exportar un reporte a PDF **para** compartirlo fuera de la aplicación.

## Objetivo

Que **el sistema** genere un archivo PDF con el informe del proyecto y con el de cada sprint, como pide la consigna
para el Sprint 4 ("el sistema permitirá exportar un informe del proyecto a un archivo en formato PDF"). El PDF se arma
en Go con los mismos datos del reporte en pantalla (HU-29 y HU-30) y se descarga desde el botón **"Descargar PDF"**.

## Entradas

| Archivo | Ruta | Datos |
|---|---|---|
| Informe del proyecto | `GET /reportes/proyecto.pdf?proyecto={proyecto}` | `Reportes.Proyecto(id, hoy)` |
| Informe de un sprint | `GET /reportes/sprints/{sprint}/pdf` | `Reportes.Sprint(id)` |

## Salidas esperadas

Una respuesta `200` con:
- `Content-Type: application/pdf`;
- `Content-Disposition: attachment; filename="reporte-proyecto-P1.pdf"` (o `reporte-sprint-S3.pdf`);
- un PDF A4 con el título, el nombre del proyecto o el Sprint Goal, "Generado el dd/mm/aaaa", y las mismas secciones
  que la pantalla:
  - **proyecto:** estado y días restantes, resumen de sprints (tabla), métricas acumuladas y defectos abiertos;
  - **sprint:** historias planificadas (tabla), métricas, esfuerzo por integrante y defectos, con la marca
    "Reporte parcial" si el sprint está activo.

## Reglas de negocio

- **RN-31.1:** el PDF tiene el mismo contenido que el reporte en pantalla: se arma con el mismo `ReporteProyecto` o
  `ReporteSprint`, no recalcula nada.
- **RN-31.2:** los números van con coma decimal, igual que en pantalla (`unDecimal`).
- **RN-31.3:** los textos con acentos y eñes se ven bien (traducción a la codificación de las fuentes del PDF).
- **RN-31.4:** las tablas que no entran en una página siguen en la siguiente, repitiendo el encabezado.

## Restricciones

- **Biblioteca:** `github.com/go-pdf/fpdf`, la prevista en el ADR 0001 para el Sprint 4: puro Go, sin dependencias del
  sistema. Se agrega a `go.mod`.
- **Paquete:** `internal/report/pdf`, con `Proyecto(io.Writer, report.ReporteProyecto, hoy) error` y
  `Sprint(io.Writer, report.ReporteSprint, hoy) error`. No accede a servicios (ADR 0002).
- **Pantallas:** los reportes suman el botón "Descargar PDF" junto a "Imprimir", que queda para el dashboard y para
  quien prefiera imprimir.
- **Tests:** el contenido del PDF se verifica en los tests de `internal/report/pdf` generándolo sin compresión
  (`SetCompression(false)`) y buscando los textos. Los tests del handler verifican los encabezados y que el cuerpo
  empiece con `%PDF-`.

## Casos límite

- Proyecto sin sprints: el PDF dice "Todavía no hay sprints con resultado".
- Sprint activo: el PDF lleva "Reporte parcial".
- Historias con títulos largos: se cortan en varias líneas dentro de su celda.

## Condiciones de error

| Situación | Respuesta |
|---|---|
| El proyecto no existe | 404 con la página de error: "el proyecto no existe" |
| El sprint no existe | 404: "el sprint no existe" |
| El sprint está planificado | 400: "el sprint todavía no empezó" (HU-29) |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-31.1 | El PDF tiene el mismo contenido que el reporte en pantalla | `Escenario: descargar el informe del proyecto en PDF`, `Escenario: descargar el informe de un sprint en PDF`, `Escenario: sin reporte no hay PDF` |
