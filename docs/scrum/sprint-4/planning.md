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
