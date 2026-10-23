# HU-03 · Registrar integrantes

- **Issue:** #10 · **Épica:** E1 Proyectos · **Responsable:** Federico Sosa · **Story Points:** 3 (historia de referencia)
- **Escenarios BDD:** [`features/hu-03-registrar-integrantes.feature`](../features/hu-03-registrar-integrantes.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** Federico Sosa (al aprobar el PR)
- **Depende de:** HU-01 (`Proyecto`)
- **Estado:** revisada

## Historia

**Como** líder del proyecto **quiero** registrar los integrantes del proyecto con su rol **para** asignarles trabajo, votos y horas.

## Objetivo

Registrar quiénes forman parte de un proyecto y con qué rol. Los integrantes votan en Planning Poker (E4) y cargan horas (E5).

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |
| Nombre | `string` | Sí | No vacío (sin contar espacios al inicio y al final) |
| Rol | `project.Rol` | Sí | `AgileEnabler`, `ProductBuilder` o `ProductArchitect` |

## Salidas esperadas

- Un `project.Integrante` guardado, con `ID` generado (`I1`, `I2`, …) y el nombre sin espacios sobrantes.
- La consulta de integrantes de un proyecto devuelve todos sus integrantes en el orden en que se registraron.

## Reglas de negocio

- **RN-03.1:** el nombre es obligatorio.
- **RN-03.2:** el rol es uno de: Agile Enabler, Product Builder, Product Architect.
- **RN-03.3:** no se puede registrar dos veces el mismo integrante (mismo nombre, sin distinguir mayúsculas) en un proyecto.
- **RN-03.4:** la misma persona puede ser integrante de proyectos distintos.

## Restricciones

- RN-03.1 y RN-03.2 se validan en `project.NuevoIntegrante(proyectoID, nombre, rol) (Integrante, error)`.
- RN-03.3 la controla el servicio con los integrantes ya registrados del proyecto.

## Casos límite

- Proyecto sin integrantes: la consulta devuelve una lista vacía, sin error.
- "ana pérez" cuando ya existe "Ana Pérez" en el mismo proyecto: duplicado.
- "Ana Pérez" en otro proyecto: válido.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |
| Nombre vacío | `project.ErrIntegranteSinNombre`: "el nombre del integrante es obligatorio" |
| Rol fuera de la lista | `project.ErrRolInvalido`: "el rol debe ser Agile Enabler, Product Builder o Product Architect" |
| Integrante ya registrado en el proyecto | `project.ErrIntegranteDuplicado`: "el integrante ya está registrado en el proyecto" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-03.1 | El nombre es obligatorio y el rol es uno de los tres de la consigna | `Escenario: registrar un integrante`, `Escenario: nombre vacío`, `Escenario: rol inválido` |
| CA-03.2 | No se registra dos veces el mismo integrante en un proyecto | `Escenario: integrante duplicado`, `Escenario: mismo integrante en otro proyecto` |
| CA-03.3 | Se puede consultar la lista de integrantes del proyecto | `Escenario: consultar los integrantes` |
