# TEC-20 · Base compartida del Sprint 3

- **Issue:** #99 · **Responsable:** Branko Almeira · **Story Points:** 1
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de la planning del Sprint 3
- **Depende de:** nada: funciones puras del dominio y de las métricas
- **Estado:** implementada (Sprint 3)

## Objetivo

Tres funciones que necesitan varias historias del Sprint 3, de distintos responsables. Si cada uno las escribiera por su
cuenta, quedarían duplicadas o alguien tendría que esperar a otro. Como TEC-18 en el Sprint 2, entran primero y nadie
depende de una historia de otro.

| Función | La usan |
|---|---|
| `project.Proyecto.EstadoAl(hoy time.Time) project.Estado` | HU-04 (Sosa), HU-28 (Giorlando), HU-30 (Branko) |
| `project.Proyecto.DiasRestantes(hoy time.Time) int` | HU-04 (Sosa), HU-28 (Giorlando), HU-30 (Branko) |
| `metrics.HorasEnPeriodo(registros []metrics.RegistroHoras, inicio, fin time.Time) float64` | HU-28 (Giorlando), HU-29 (Branko) |

## Reglas

- **EstadoAl:** antes del inicio, `Planificado`; entre el inicio y el fin, ambos incluidos, `EnCurso`; después del fin,
  `Finalizado` (RN-04.1). Las fechas se comparan por día, sin hora.
- **DiasRestantes:** los días desde hoy hasta la fecha de fin, sin contar hoy; 0 si hoy es el fin o ya pasó (RN-04.2).
- **HorasEnPeriodo:** suma las horas de los registros con fecha entre `inicio` y `fin`, ambos incluidos (RN-28.3 y
  RN-29.3). `RegistroHoras` tiene solo `Fecha` y `Horas`: el servicio ya filtró las historias del sprint. Sin
  registros en el período devuelve 0.

## Criterio de terminado

- [x] Las tres funciones con tests de tabla (incluidos los bordes: día de inicio, día de fin, registros fuera del período)
- [x] Cobertura del 100 % en las funciones nuevas
- [ ] Mergeado antes de que se abra el primer PR de una historia del Sprint 3
