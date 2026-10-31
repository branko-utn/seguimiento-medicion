# Sprint 3 · Retrospectiva

- **Fecha:** sábado 31/10/2026, después de la review · **Duración:** 15'
- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Formato:** Seguir / Dejar de hacer / Probar

## Acciones de la retro anterior

| Acción | ¿Se hizo? |
|---|---|
| Escribir en `CONTRIBUTING.md` la resolución de conflictos de los archivos compartidos | No: tercer sprint seguido. Se resolvieron en las auditorías, PR por PR |
| Volcar la daily en `dailies.md` | No |
| En la daily del segundo día, reasignar o recortar los ítems sin PR abierto | Parcialmente: no se recortó nada, pero tres PRs llegaron recién el último día |
| Respetar la rotación de revisores | No: para llegar al cierre, el Agile Enabler aprobó #105, #109 y #111 |
| Planificar con menos SP firmes que la velocidad y los recortables marcados | Sí: 36 firmes + 5 recortables, y entró todo |

## Seguir haciendo

- **Una base compartida que entra primero** (TEC-20): Fede, Giorlando y Branko usaron las mismas funciones sin
  esperarse.
- **Specs sin dependencias entre personas:** nadie quedó bloqueado por la historia de otro.
- **Auditar cada PR contra su spec antes de mergear:** aparecieron 7 defectos, todos corregidos menos uno.
- **Corregir lo menor en la misma rama, con un comentario**, en lugar de devolver el PR: ahorró un día.

## Dejar de hacer

- **Llegar al cierre con los PRs más grandes abiertos:** HU-14 y HU-28 se abrieron el 29 y el 30, y se mergearon el
  31.
- **Copiar el número de issue sin verificarlo:** tres PRs cerraron el issue equivocado en dos sprints (#96, #109 y
  #111).
- **Apilar PRs** (#111 sobre #109) sin decirlo en la descripción.
- **Subir archivos sin `gofmt`:** el CI los rechaza.

## Probar

- **Abrir el PR como borrador el primer día**, aunque tenga solo los tests (RED).
- **Revisar `Closes #N` con `gh issue view N`** antes de abrir el PR.

## Acciones para el próximo sprint

| Acción | Responsable |
|---|---|
| Code freeze el sábado 31 a las 23:59: después solo se corrigen defectos críticos | Todo el equipo |
| Escribir la resolución de conflictos y el chequeo de `Closes` en la documentación técnica (TEC-07) | Branko Almeira |
| Corregir el formato de los números del dashboard (defecto abierto) | Agustín Giorlando |
| Ensayar la demo el domingo 1/11 con los datos de `make demo` | Todo el equipo |
