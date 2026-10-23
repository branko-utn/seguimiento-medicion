# ADR 0003 · Enumeradores del dominio: int con iota y String()

- **Estado:** propuesta (23/10/2026), a aceptar con la revisión del PR
- **Responsable:** Ignacio Bianchi

## Contexto

Varios tipos del dominio tienen un conjunto cerrado de valores (ADR 0002, sección 2): `Prioridad` y `Estado` del
backlog, `Estado` de project y sprint, `Rol`, `Severidad` de defectos. Esos valores se comparan **por su texto** en
las specs y los escenarios BDD ("Alta", "En progreso"), se guardan en el JSON de `data/` y se interpolan en mensajes
de error ("no se puede pasar de Completada a En progreso", HU-08).

La decisión surgió al implementar el backlog (HU-05). Cinco personas programamos en paralelo sobre esos tipos: si
cada épica elige una representación distinta, cada PR rompe a los demás. La representación se define una sola vez, acá.

## Decisión

Cada enumerador es un **tipo definido sobre `int`**, con constantes `iota` que **arrancan en `iota + 1`** (el valor
cero queda reservado para "sin valor"), y cada tipo implementa **`String()`** con el texto que ve el usuario:

```go
// internal/domain/backlog
type Prioridad int

const (
	PrioridadAlta Prioridad = iota + 1 // 1
	PrioridadMedia                     // 2
	PrioridadBaja                      // 3
)

func (p Prioridad) String() string {
	switch p {
	case PrioridadAlta:
		return "Alta"
	case PrioridadMedia:
		return "Media"
	case PrioridadBaja:
		return "Baja"
	}
	return "" // valor cero: no es válido
}
```

Reglas:

1. **El cero nunca es un valor válido.** Las constantes arrancan en `iota + 1`. Así el struct cero —y un JSON al que
   le falta un campo— queda con un valor inválido que los constructores detectan y rechazan con el error de la spec
   (`ErrPrioridadInvalida`, `ErrEstadoInvalido`), en vez de aceptarlo por descuido.
2. **Nombres de constantes:** `<Tipo><Valor>` cuando agrega contexto (`PrioridadAlta`); el valor a secas cuando el
   tipo ya lo dice (`backlog.Estado`: `Pendiente`, `EnProgreso`, `Completada`).
3. **`String()` devuelve el texto de las specs**, en español, tal cual lo escriben los escenarios: `"En progreso"`,
   con espacio. Para el valor cero devuelve cadena vacía: un inválido no debería mostrarse nunca como número.
4. **Los mensajes de error se arman con `String()`**: `fmt.Errorf("no se puede pasar de %s a %s", actual, nuevo)`
   produce exactamente `"no se puede pasar de Completada a En progreso"` (HU-08). Cambiar un texto rompe los tests y
   desalinea la spec: se detecta en el CI.
5. **El texto del usuario se convierte en el borde** (pasos BDD, formularios HTTP): ahí se mapea `"Alta"` →
   `PrioridadAlta`. El dominio recibe y devuelve constantes, nunca cadenas. Si la conversión se repite, la función
   auxiliar va junto a ese borde, no en el dominio.
6. **Persistencia:** se serializan como número en el JSON de `data/`. Si más adelante el archivo necesita ser legible,
   se agrega `MarshalJSON`/`UnmarshalJSON` y se documenta en un ADR nuevo.

`estimation.Escala` no es un enumerador: es la lista de valores válidos de Story Points (HU-15).

## Alternativas descartadas

- **Tipos string** (`type Prioridad string`): el JSON queda legible, pero el cero `""` confunde "sin valor" con
  "inválido", y los literales (`"En progreso"`) quedan duplicados en cada comparación.
- **`iota` sin `+1`** (cero como valor válido): el struct cero pasaría por válido y las validaciones no detectarían
  datos faltantes.

## Consecuencias

- Todas las épicas aplican esta convención a sus enumeradores desde su primer PR; se revisa igual que los nombres de
  tipos del ADR 0002.
- Agregar un valor nuevo son dos cambios: la constante y su caso en `String()`.
- Las pantallas muestran `%s` (o `.String()`), nunca el número.
