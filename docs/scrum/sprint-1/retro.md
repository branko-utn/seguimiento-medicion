# Sprint 1 · Retrospectiva

- **Fecha:** domingo 25/10/2026 por la noche, después de la review · **Duración:** 15'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Formato:** Seguir / Dejar de hacer / Probar

## Acciones de la retro anterior

| Acción | ¿Se hizo? |
|---|---|
| Automatizar la asignación del revisor según la rotación | Sí: el CI asigna el revisor en cada PR |
| PRs separados por capa | Sí: HU-05, HU-07, HU-08 y HU-15 entraron en dos PRs cada una |
| Volcar la daily en `dailies.md` | No: la daily se hizo por el grupo, pero no se pasó al repo |

## Seguir haciendo

- **Specs y escenarios BDD redactados por una sola persona:** los mensajes de error y los nombres quedaron iguales en
  todas las épicas, y nadie tuvo que adivinar el alcance.
- **Contratos de dominio y orden de merge (ADR 0002):** los tipos de Sosa entraron primero y los demás pudieron
  conectar sus servicios sin reescribir nada.
- **Revisiones con el detalle en el PR:** los pedidos de cambio traían el código sugerido, y se corrigieron en el día.
- **TDD visible en los commits** (`test: RED` → `feat: GREEN`): facilita revisar.

## Dejar de hacer

- **Mergear de madrugada:** muchos PRs entraron después de las 23:00, cuando ya no había nadie para revisar con
  calma.
- **Resolver los conflictos quedándose con un solo lado:** en `features/features_test.go` y `docs/ia/registro.md` se
  pisaron registros de otros (#53 y #54).
- **Usar `Closes #N` en PRs parciales:** cierra historias que no están terminadas y desordena el tablero.
- **Llegar tarde con los tipos compartidos:** el corte era el sábado 24 a las 12:00 y los tipos entraron a las 16:30.

## Probar

- Hacer `git pull origin main` antes de abrir el PR y, en un conflicto, **quedarse con los dos lados**.
- Que "probado en la aplicación" sea obligatorio desde que haya pantallas.
- Hacer primero el PR que conecta todos los servicios a la app (TEC-18), para que las pantallas no choquen en
  `internal/server`.

## Acciones para el próximo sprint

| Acción | Responsable |
|---|---|
| Agregar a `CONTRIBUTING.md` cómo resolver los conflictos de los archivos compartidos (quedarse con los dos lados) | Branko Almeira |
| Revisar `Closes` / `Parte de` en cada revisión de PR | Revisor de turno |
| Mergear TEC-18 antes de cualquier pantalla | Branko Almeira |
| Volcar la daily en `docs/scrum/sprint-2/dailies.md` | Branko Almeira |
| No mergear después de las 23:00 salvo para cerrar el sprint | Todo el equipo |
