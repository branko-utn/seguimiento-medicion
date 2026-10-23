# HU-01 · Crear proyecto

- **Issue:** #8 · **Épica:** E1 Proyectos · **Responsable:** Federico Sosa · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-01-crear-proyecto.feature`](../features/hu-01-crear-proyecto.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Federico Sosa (al aprobar el PR)
- **Estado:** revisada

## Historia

**Como** líder del proyecto **quiero** crear un proyecto con nombre, descripción y fechas de inicio y fin **para** empezar a administrarlo.

## Objetivo

Dar de alta un proyecto válido, que es el punto de partida de todo lo demás: integrantes, backlog, sprints y métricas.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Nombre | `string` | Sí | Entre 3 y 100 caracteres, sin contar espacios al inicio y al final |
| Descripción | `string` | No | — |
| Fecha de inicio | `time.Time` (día) | Sí | — |
| Fecha de fin | `time.Time` (día) | Sí | No anterior a la de inicio |

## Salidas esperadas

Un `project.Proyecto` guardado, con `ID` generado (`P1`, `P2`, …), el nombre sin espacios sobrantes y estado `Planificado`.

## Reglas de negocio

- **RN-01.1:** el nombre tiene entre 3 y 100 caracteres después de quitar los espacios al inicio y al final.
- **RN-01.2:** la fecha de fin no puede ser anterior a la de inicio; pueden ser el mismo día.
- **RN-01.3:** no puede haber dos proyectos con el mismo nombre, sin distinguir mayúsculas de minúsculas.
- **RN-01.4:** un proyecto nuevo queda en estado `Planificado`.

## Restricciones

- Las validaciones de RN-01.1, RN-01.2 y RN-01.4 están en el constructor `project.Nuevo(nombre, descripcion, inicio, fin) (Proyecto, error)`.
- La unicidad (RN-01.3) la controla el servicio consultando el repositorio, porque requiere conocer los otros proyectos.
- Los caracteres se cuentan como runas (letras), no bytes: "Año" tiene 3.

## Casos límite

- Nombre de exactamente 3 y de exactamente 100 caracteres: válidos. De 2 y de 101: inválidos.
- Nombre `"  Demo  "`: se guarda como `"Demo"`.
- Inicio y fin el mismo día: válido.
- `"demo"` cuando ya existe `"Demo"`: duplicado.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Nombre con menos de 3 o más de 100 caracteres (o vacío) | `project.ErrNombreInvalido`: "el nombre debe tener entre 3 y 100 caracteres" |
| Fin anterior al inicio | `project.ErrFechasInvalidas`: "la fecha de fin no puede ser anterior a la de inicio" |
| Ya existe un proyecto con ese nombre | `project.ErrNombreDuplicado`: "ya existe un proyecto con ese nombre" |
| Falta una de las fechas | `project.ErrFechasInvalidas` |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-01.1 | El nombre es obligatorio y tiene entre 3 y 100 caracteres | `Esquema del escenario: largo del nombre` |
| CA-01.2 | La fecha de fin no puede ser anterior a la de inicio | `Escenario: fin anterior al inicio`, `Escenario: inicio y fin el mismo día` |
| CA-01.3 | No puede haber dos proyectos con el mismo nombre | `Escenario: nombre duplicado` |
| CA-01.4 | Un proyecto nuevo queda en estado "Planificado" | `Escenario: crear un proyecto válido` |
