# HU-21 · Consultar el esfuerzo registrado

- **Issue:** #28 · **Épica:** E5 Esfuerzo · **Responsable:** Agustín Lara · **Story Points:** 3
- **Escenarios BDD:** [`features/hu-21-consultar-el-esfuerzo-registrado.feature`](../features/hu-21-consultar-el-esfuerzo-registrado.feature)
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo en la planning del Sprint 3
- **Depende de:** HU-20 (`effort.Registro`), HU-03 (integrantes), HU-05 (historias), HU-26 (`metrics.HorasReales`)
- **Estado:** borrador para la planning del Sprint 3

## Historia

**Como** integrante del equipo **quiero** ver el esfuerzo registrado por historia y por integrante **para** saber en qué
se fue el tiempo.

## Objetivo

Listar el esfuerzo registrado en un proyecto y resumirlo en totales por historia y por integrante. Las horas reales que
usan el dashboard (HU-28) y los reportes (E9) salen de estos mismos registros.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | `domain.ProyectoID` | Sí | Debe existir |

## Salidas esperadas

Un resumen del esfuerzo del proyecto con:

- los registros (fecha, integrante, historia, actividad y horas);
- el total de horas por historia;
- el total de horas por integrante;
- el total del proyecto.

La consulta no modifica nada.

## Reglas de negocio

- **RN-21.1:** se incluyen los registros de las historias del proyecto.
- **RN-21.2:** los registros se ordenan por fecha, del más reciente al más antiguo; en la misma fecha, por orden de carga.
- **RN-21.3:** los totales son la suma de las horas (`metrics.HorasReales`, HU-26). Los totales por historia se ordenan
  por identificador numérico (HU-2 antes que HU-10); los totales por integrante, por nombre.
- **RN-21.4:** las historias e integrantes sin registros no aparecen en los totales.
- **RN-21.5:** el total del proyecto es igual a la suma de los totales por historia y a la de los totales por integrante.

## Restricciones

- El servicio de registros suma `DelProyecto(proyectoID) (ResumenEsfuerzo, error)`. Busca las historias del proyecto y
  sus registros; no hace falta un método nuevo en el repositorio si alcanza con `DeHistoria`.
- `ResumenEsfuerzo` tiene `Registros`, `PorHistoria` y `PorIntegrante` (listas ordenadas, no mapas, para que el orden de
  RN-21.3 sea estable en la pantalla) y `Total`.

## Casos límite

- Proyecto sin registros: listas vacías y total 0; la pantalla muestra "Todavía no hay esfuerzo registrado".
- Un integrante con registros en varias historias: aparece una vez en los totales por integrante, con la suma.
- Horas con decimales (por ejemplo 2.5): se suman sin redondear; la pantalla las muestra con un decimal.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El proyecto no existe | `project.ErrNoExiste`: "el proyecto no existe" |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-21.1 | Se listan los registros con sus totales por historia y por integrante | `Escenario: totales por historia y por integrante`, `Escenario: los registros se listan del más reciente al más antiguo`, `Escenario: proyecto sin esfuerzo registrado` |
