# Sprint 4 · Planning

- **Fecha:** sábado 31/10/2026 · **Duración:** 25'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Sprint:** El Cierre y la Entrega Final · sábado 31/10 (corrido un día en la planning del Sprint 3)

## Sprint Goal

> Producto pulido, documentado y listo para presentar.

## Fechas clave

| Cuándo | Qué |
|---|---|
| Sábado 31/10, 23:59 | **Code freeze.** Después solo se corrigen defectos críticos |
| Domingo 1/11 | Ensayo de la demo (90') con los datos de `make demo` |
| Lunes 2/11 | Presentación final |

## Sprint Backlog

| Responsable | Ítem | SP |
|---|---|---|
| Branko Almeira | HU-31 · Exportar reporte a PDF (#38) | 5 |
| Agustín Giorlando | TEC-08 · Informe de métricas y cobertura de pruebas (#42) — **pasa de Lara a Giorlando** | 3 |
| | Defecto abierto: números del dashboard con punto decimal | — |
| Todo el equipo | TEC-07 · Manual de usuario y documentación técnica (#41): cada responsable escribe la sección de su épica y Branko la parte técnica | 5 |
| | **Total** | **13** |

Velocidad: 44, 41 y 41 SP en los Sprints 1 a 3 (promedio 42). El Sprint 4 dura un día y es de cierre, así que se
comprometen 13 SP.

## Carga acumulada al cierre del proyecto

| Integrante | Sprints 1 a 3 | Sprint 4 | **Total** |
|---|---|---|---|
| Federico Sosa | 27 | — | **27** |
| Ignacio Bianchi | 26 | — | **26** |
| Agustín Lara | 25 | — | **25** |
| Agustín Giorlando | 24 | TEC-08 (3) | **27** |
| Branko Almeira | 30 | HU-31 (5) | 35 (+ Sprint 0, facilitación y auditorías) |

Sin contar los ítems de todo el equipo (TEC-06 y TEC-07). **Los cuatro Product Builders terminan entre 25 y 27 SP.**
Para lograrlo, TEC-08 pasa de Lara a Giorlando, como estaba previsto en la planning del Sprint 3. Además, el
cronograma ya ponía el informe de métricas en la épica de Giorlando (E7).

## Decisiones

- **HU-31 se cumple con la impresión del navegador.** El botón "Imprimir o guardar como PDF" de los reportes, junto
  con la hoja `@media print`, ya produce un PDF con el mismo contenido que la pantalla (CA-31.1). Así lo pide el ADR
  0001: sin bibliotecas de PDF. Lo que falta:
  - la spec;
  - el escenario que verifica que cada reporte ofrece la exportación;
  - el mismo botón en el dashboard;
  - revisar que los gráficos y las tablas se impriman bien.
- **El manual es uno solo** (`docs/manual-de-usuario.md`), con una sección por épica escrita por su responsable.
- **Code freeze a las 23:59.** Después del freeze, cualquier corrección necesita un PR revisado y la aprobación del
  Agile Enabler.

## Riesgos

- **Un solo día de sprint.** Mitigación: lo de código es chico (HU-31 y el defecto) y entra antes del mediodía; la
  tarde queda para la documentación.
- **Sin tiempo para corregir después de la demo.** Mitigación: el ensayo del domingo usa `make demo` sobre una
  carpeta `data/` limpia, igual que en la presentación.

## Cambio de alcance: mejora de la interfaz (31/10, a pedido del equipo)

Después de ver la demo con nuestros datos, el equipo coincidió en que la interfaz se ve básica para la presentación.
Se suman dos ítems, armados como en el Sprint 3 para que nadie dependa de otro:

| Responsable | Ítem | SP | Depende de |
|---|---|---|---|
| Branko Almeira | TEC-21 · Base visual: paleta, tarjetas, insignias, estados vacíos, portada (**entra primero**) | 2 | — |
| Federico Sosa | HU-32.1 · Interfaz de proyectos y sprints — **recortable** | 2 | TEC-21 |
| Ignacio Bianchi | HU-32.2 · Interfaz del backlog y del Planning Poker — **recortable** | 2 | TEC-21 |
| Agustín Lara | HU-32.3 · Interfaz de esfuerzo y defectos — **recortable** | 2 | TEC-21 |
| Agustín Giorlando | HU-32.4 · Interfaz del dashboard (resuelve el defecto del punto decimal) — **recortable** | 2 | TEC-21 |
| Todo el equipo | TEC-22 · Presentación final y demo (entregable 13) | 3 | — |
| | **Total del Sprint 4** | **26** | |

- **Cada parte toca solo las plantillas y los handlers de su épica** (tabla en `specs/HU-32-*.md`): no hay
  conflictos entre personas.
- **Las partes son recortables una por una:** no cambian comportamiento, así que si una no llega antes del freeze,
  esa pantalla queda como está.
- **Orden de Branko:** TEC-21 a primera hora, después HU-31, que usa la misma base.

**Carga final con el cambio:** Sosa 29, Bianchi 28, Lara 27, Giorlando 29. Los cuatro Product Builders siguen
parejos, entre 27 y 29 SP.

### Specs y plantillas del Sprint 4

- `specs/HU-31-exportar-reporte-a-pdf.md` y su feature.
- `specs/TEC-21-base-visual.md`.
- `specs/HU-32-interfaz-clara-y-consistente.md`, con cuatro features, una por parte.
- TEC-07: `docs/manual-de-usuario.md`, con una sección por épica y su responsable.
- TEC-08: `docs/informe-de-metricas.md`, con las tablas a completar y de dónde sale cada dato.

## Corrección: HU-31 genera el PDF en el sistema (31/10)

Al revisar la consigna punto por punto aparecieron tres cambios:

1. **HU-31 genera el archivo PDF en Go.** La consigna dice "el sistema permitirá exportar un informe del proyecto
   **a un archivo en formato PDF**". La impresión del navegador no alcanza, porque el PDF lo hace el navegador, no
   nuestro sistema. Además, el ADR 0001 ya preveía `github.com/go-pdf/fpdf` para el Sprint 4: la decisión de arriba
   ("sin bibliotecas de PDF") era un error de lectura. Los reportes suman "Descargar PDF"
   (`/reportes/proyecto.pdf` y `/reportes/sprints/{sprint}/pdf`), y la impresión queda como alternativa. Siguen siendo
   5 SP.
2. **TEC-22 · Presentación final (entregable 13).** No había un ítem para prepararla. Incluye:
   - las diapositivas;
   - el guion de la demo con `make demo`;
   - el recorrido completo de una funcionalidad que pide la consigna (historia → spec → criterios → BDD → tests →
     código Go). Se propone HU-13 (cerrar sprint), que tiene los commits RED → GREEN más claros;
   - cada integrante presenta su épica.

   Se ensaya el domingo 1/11.
3. **Dailies:** se reconstruyeron las de los Sprints 1 a 3 (`docs/scrum/sprint-N/dailies.md`) a partir de los
   commits y PRs de cada día. Desde el Sprint 4, el Agile Enabler las vuelca el mismo día.
