# TEC-04 · Mecanismo común de persistencia a archivo JSON

- **Issue:** #5 · **Responsable:** Branko Almeira
- **Paquete:** `internal/store/memory` (tipo `Store[T]`)
- **Estado:** implementada

> Ítem técnico interno: no tiene escenarios BDD porque no hay comportamiento visible para el usuario.
> Se valida con tests unitarios (`internal/store/memory/store_test.go`).

## Objetivo

Dar a cada épica una forma única y probada de guardar sus datos en un archivo JSON, para que los repositorios
solo se ocupen de su lógica (buscar, filtrar, generar IDs) y no de leer y escribir archivos.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Ruta del archivo | `string` | No | Vacía = modo solo memoria (para tests) |
| Tipo de dato | parámetro genérico `T` | Sí | Debe poder serializarse a JSON |
| Función de lectura | `func(T)` | Sí | — |
| Función de actualización | `func(*T) error` | Sí | — |

## Salidas esperadas

- `Open` devuelve un `Store[T]` con el contenido del archivo, o el valor cero de `T` si el archivo no existe.
- `View` ejecuta la función con el estado actual, sin modificarlo.
- `Update` ejecuta la función sobre el estado y, si no devuelve error, lo guarda en el archivo.

## Reglas de negocio

- **RN-TEC04.1:** si la función de `Update` devuelve un error, **no se guarda nada** y el estado en memoria vuelve a ser el anterior.
- **RN-TEC04.2:** la escritura es atómica: se escribe un archivo temporal y se renombra, así nunca queda un JSON a medio escribir.
- **RN-TEC04.3:** el archivo no se crea hasta el primer `Update` exitoso; la carpeta se crea si no existe.

## Restricciones

- Seguro ante accesos concurrentes: varias lecturas a la vez, una escritura a la vez.
- Cada colección usa su propio archivo (por ejemplo `data/proyectos.json`, `data/sprints.json`), así las épicas no comparten estructuras.

## Casos límite

- Archivo inexistente: se arranca con el valor cero, sin error.
- Ruta vacía: funciona igual, sin tocar el disco.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| El archivo existe pero no es JSON válido | `Open` devuelve un error que menciona la ruta |
| No se puede escribir el archivo | `Update` devuelve el error y conserva el estado anterior |

## Criterios de aceptación

| ID | Criterio | Test |
|---|---|---|
| CA-TEC04.1 | Lo guardado se recupera al volver a abrir el archivo | `TestStore_persisteYRecupera` |
| CA-TEC04.2 | Archivo inexistente arranca con el valor cero y sin crear el archivo | `TestOpen_archivoInexistente` |
| CA-TEC04.3 | Un `Update` con error no guarda y restaura el estado anterior | `TestUpdate_conErrorNoGuarda` |
| CA-TEC04.4 | Un archivo corrupto devuelve error al abrir | `TestOpen_archivoCorrupto` |
| CA-TEC04.5 | Funciona en modo solo memoria | `TestStore_soloMemoria` |
| CA-TEC04.6 | Soporta actualizaciones concurrentes sin perder datos | `TestUpdate_concurrente` |
