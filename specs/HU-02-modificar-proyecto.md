# HU-02 · Modificar proyecto

- **Issue:** #9 · **Épica:** E1 Proyectos · **Responsable:** Federico Sosa · **Story Points:** 2
- **Escenarios BDD:** [`features/hu-02-modificar-proyecto.feature`](../features/hu-02-modificar-proyecto.feature)
- **Depende de:** HU-01 (en `main`)
- **Redactada por:** Branko Almeira (Agile Enabler) después de la planning del Sprint 2 · **Validada por:** Federico Sosa (al aprobar el PR)
- **Estado:** aprobada para el Sprint 2

## Historia

**Como** líder del proyecto **quiero** modificar los datos de un proyecto **para** corregirlos o actualizarlos.

## Objetivo

Editar nombre, descripción y fechas de un proyecto existente con las mismas validaciones del alta (HU-01), incluida su
pantalla de edición.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir y no estar finalizado |
| Nombre | `string` | Sí | Entre 3 y 100 caracteres (RN-01.1); único sin distinguir mayúsculas, salvo el propio proyecto |
| Descripción | `string` | No | — |
| Fecha de inicio y de fin | `time.Time` (día) | Sí | Fin no anterior al inicio (RN-01.2) |

## Salidas esperadas

El proyecto queda guardado con los datos nuevos (nombre sin espacios sobrantes). El `ID` y el estado no cambian.

## Reglas de negocio

- **RN-02.1:** se aplican las mismas validaciones de nombre y fechas del alta (RN-01.1 y RN-01.2).
- **RN-02.2:** el nombre nuevo no puede coincidir con el de **otro** proyecto (sin distinguir mayúsculas). Se puede
  mantener el propio nombre o cambiarle solo las mayúsculas.
- **RN-02.3:** un proyecto en estado `Finalizado` no se puede modificar.
- **RN-02.4:** el ID y el estado del proyecto no se modifican con esta operación.

## Restricciones

- `Proyecto.Modificar(nombre, descripcion, inicio, fin) error` en `internal/domain/project` reutiliza las validaciones de
  `Nuevo` y aplica RN-02.3. La unicidad (RN-02.2) la controla el servicio.
- La pantalla de edición reutiliza el formulario del alta (TEC-13).

## Casos límite

- Guardar sin cambios: válido.
- Cambiar "Demo" por "DEMO" (el mismo proyecto): válido.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| Nombre inválido | `project.ErrNombreInvalido`: "el nombre debe tener entre 3 y 100 caracteres" |
| Fin anterior al inicio | `project.ErrFechasInvalidas`: "la fecha de fin no puede ser anterior a la de inicio" |
| Nombre de otro proyecto | `project.ErrNombreDuplicado`: "ya existe un proyecto con ese nombre" |
| Proyecto finalizado | `project.ErrProyectoFinalizado`: "un proyecto finalizado no se puede modificar" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-02.1 | Se pueden editar nombre, descripción y fechas, con las mismas validaciones del alta | `Escenario: modificar un proyecto`, `Escenario: nombre inválido`, `Escenario: fin anterior al inicio`, `Escenario: nombre de otro proyecto` |
| CA-02.2 | Modificar un proyecto inexistente informa un error | `Escenario: proyecto inexistente` |
| CA-02.3 | Un proyecto "Finalizado" no se puede modificar | `Escenario: proyecto finalizado` |
