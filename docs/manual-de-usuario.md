# Manual de usuario · Software Metrics & Estimation

> Plantilla de TEC-07. Cada responsable completa la sección de su épica, y Branko completa la introducción y la
> documentación técnica. Escribí para alguien que nunca usó la app: qué puede hacer, paso a paso y con las reglas
> que va a encontrar (mensajes de error incluidos). Una o dos capturas por sección como máximo, en
> `docs/manual/img/`.

## 1. Primeros pasos — Branko Almeira

<!-- Cómo levantar la app (make run / make demo), la navegación y qué hay en la portada. -->

## 2. Proyectos e integrantes (E1) — Federico Sosa

Todo empieza por un proyecto: el backlog, los sprints, el esfuerzo y los defectos son siempre de un proyecto. En la
navegación, **Proyectos** muestra la lista de proyectos con sus fechas y su estado.

### Crear un proyecto

1. En **Proyectos**, tocá **Nuevo proyecto**.
2. Completá el **nombre**, una **descripción** (opcional) y las **fechas de inicio y de fin**.
3. Tocá **Crear proyecto**. La app te lleva al detalle del proyecto nuevo.

Si algo no cumple las reglas, el formulario vuelve con lo que escribiste y el motivo:

| Regla | Mensaje |
|---|---|
| El nombre tiene entre 3 y 100 caracteres (los espacios del principio y del final no cuentan) | "el nombre debe tener entre 3 y 100 caracteres" |
| La fecha de fin no es anterior a la de inicio (pueden ser el mismo día) | "la fecha de fin no puede ser anterior a la de inicio" |
| No hay dos proyectos con el mismo nombre, sin importar mayúsculas ("Demo" y "demo" son el mismo) | "ya existe un proyecto con ese nombre" |

### El estado del proyecto

El estado **no se cambia a mano**: sale de las fechas del proyecto y del día de hoy.

| Estado | Cuándo |
|---|---|
| **Planificado** (gris) | Todavía no llegó la fecha de inicio |
| **En curso** (azul) | Entre el inicio y el fin, los dos días incluidos |
| **Finalizado** (verde) | Ya pasó la fecha de fin |

### Ver el estado del proyecto de un vistazo

Al abrir un proyecto, el detalle muestra arriba su estado y sus fechas, y tres tarjetas:

- **Días restantes:** los días que faltan hasta la fecha de fin, sin contar hoy. El último día y en un proyecto
  finalizado es 0.
- **Sprint activo:** el sprint que está en curso, con un enlace para abrirlo. Si no hay ninguno, dice "Ninguno".
- **Historias completadas:** el porcentaje de historias del backlog del proyecto que ya están completadas. Sin
  historias es 0,0 %.

### Modificar un proyecto

1. En el detalle del proyecto, tocá **Editar proyecto**. Se abre el mismo formulario del alta, con los datos cargados.
2. Cambiá el nombre, la descripción o las fechas y tocá **Guardar cambios**.

Valen las mismas reglas que al crearlo. Además:

- Podés dejar el mismo nombre o cambiarle solo las mayúsculas, pero no ponerle el nombre de **otro** proyecto
  ("ya existe un proyecto con ese nombre").
- **Un proyecto finalizado no se puede modificar** ("un proyecto finalizado no se puede modificar"). Tampoco se puede
  "reabrir" corriéndole la fecha de fin: se mira cómo está hoy, no cómo quedaría con los cambios.

### Registrar integrantes

En el detalle del proyecto, la sección **Integrantes** muestra quiénes forman el equipo y su rol. Para sumar a alguien:

1. Escribí su **nombre** y elegí el **rol**: Agile Enabler, Product Builder o Product Architect.
2. Tocá **Registrar integrante**.

| Regla | Mensaje |
|---|---|
| El nombre es obligatorio | "el nombre del integrante es obligatorio" |
| El rol es uno de los tres de la lista | "el rol debe ser Agile Enabler, Product Builder o Product Architect" |
| Una persona no se registra dos veces en el mismo proyecto, sin importar mayúsculas | "el integrante ya está registrado en el proyecto" |

Una misma persona puede ser integrante de varios proyectos. Los integrantes son los que después votan en el Planning
Poker y cargan horas.

Si abrís un proyecto que no existe (por ejemplo, un enlace viejo), la app muestra "el proyecto no existe".

## 3. Product Backlog (E2) — Ignacio Bianchi y Agustín Lara

<!-- Bianchi: crear historias, criterios, estados, estimar en SP (HU-05, 07, 08, 15).
     Lara: editar y eliminar, backlog ordenado y filtro por estado (HU-06, HU-09). -->

## 4. Sprints (E3) — Federico Sosa

Un sprint es una iteración del proyecto con fechas y un objetivo (el **Sprint Goal**). En la navegación, **Sprints**
muestra, por cada proyecto, sus sprints con su estado y, más abajo, los **sprints cerrados** con sus resultados.

El ciclo de un sprint es siempre el mismo:

| Estado | Qué significa | Qué se puede hacer |
|---|---|---|
| **Planificado** (gris) | Recién creado | Asignarle historias e **iniciarlo** |
| **Activo** (azul) | En curso | Asignar historias, **completarlas** y **cerrarlo** |
| **Cerrado** (verde) | Terminado: su resultado queda guardado | Solo consultarlo |

### Crear un sprint

1. En **Sprints**, tocá **Nuevo sprint** en el proyecto que corresponda (ya viene elegido en el formulario).
2. Escribí el **Sprint Goal** y las **fechas de inicio y de fin**.
3. Tocá **Crear sprint**. La app te lleva al detalle del sprint, en estado Planificado.

| Regla | Mensaje |
|---|---|
| El Sprint Goal es obligatorio | "el Sprint Goal es obligatorio" |
| La fecha de fin no es anterior a la de inicio | "la fecha de fin no puede ser anterior a la de inicio" |
| Las fechas están dentro de las del proyecto (puede empezar el día de inicio y terminar el día de fin) | "las fechas del sprint deben estar dentro de las del proyecto" |
| No comparte ningún día con otro sprint del mismo proyecto (sí puede empezar el día siguiente al fin del otro) | "las fechas se superponen con otro sprint del proyecto" |

### Asignar historias al sprint (armar el Sprint Backlog)

En el detalle del sprint, **Asignar una historia del backlog** ofrece las historias del proyecto que todavía no están
en ningún sprint. Elegí una y tocá **Asignar al sprint**. La historia aparece en la tabla **Historias del sprint** con
sus Story Points.

Solo entran historias **listas para trabajar**. Si no lo están, la app lo dice:

| Regla | Mensaje |
|---|---|
| La historia está estimada (tiene Story Points) | "la historia no está estimada" |
| Tiene al menos un criterio de aceptación | "la historia no tiene criterios de aceptación" |
| No está en otro sprint (ni dos veces en el mismo) | "la historia ya está asignada a un sprint" |
| No está completada | "la historia ya está completada" |
| El sprint no está cerrado | "el sprint está cerrado y no se puede modificar" |

Para estimar una historia o agregarle criterios, andá al **Backlog** (sección 3) o al **Planning Poker** (sección 5).

### Iniciar el sprint

En un sprint planificado, tocá **Iniciar sprint**. Solo puede haber **un sprint activo por proyecto**: si ya hay otro,
la app responde "ya hay un sprint activo en el proyecto". Un sprint que ya está activo o cerrado no se vuelve a
iniciar ("solo se puede iniciar un sprint planificado").

### Completar historias

Mientras el sprint está activo, cada historia de la tabla tiene el botón **Completar**. Al tocarlo, la historia queda
como **Completada** (verde) y sus Story Points pasan a contar como completados en el sprint. Una historia
completada no se vuelve a completar.

### Métricas del sprint

Arriba del detalle, cuatro tarjetas resumen cómo va el sprint:

- **Story Points planificados:** la suma de los SP de las historias del sprint.
- **Story Points completados:** la suma de los SP de las historias completadas.
- **Historias completadas:** el porcentaje, con cuántas son de cuántas.
- **Velocidad del equipo:** el promedio de SP completados de los sprints **cerrados** del proyecto (el sprint en curso
  todavía no cuenta).

Qué significa cada métrica en detalle está en la sección 8.

### Cerrar el sprint

En un sprint activo, tocá **Cerrar sprint**. Al cerrarlo:

- quedan fijos sus Story Points planificados y completados, que pasan a contar para la velocidad del equipo;
- las historias que **no** se completaron vuelven al backlog como Pendientes y sin sprint, listas para planificarse en
  el próximo;
- el sprint ya no se puede modificar ("el sprint está cerrado y no se puede modificar").

Solo se cierra un sprint activo: uno que nunca se inició no se puede cerrar ("solo se puede cerrar un sprint activo").

### Consultar los sprints anteriores

En **Sprints**, debajo de los sprints de cada proyecto, **Sprints cerrados** lista los que ya terminaron, del más
antiguo al más reciente, con sus fechas, su Sprint Goal y los Story Points planificados y completados. Al abrir uno
se ven todas las historias que se planificaron: las completadas y las que volvieron al backlog
(**No completada**, en gris).

Si abrís un sprint que no existe, la app muestra "el sprint no existe".

## 5. Planning Poker (E4) — Ignacio Bianchi

<!-- Abrir una votación, votar, revelar, nueva ronda y estimación acordada (HU-16 a HU-18). -->

## 6. Esfuerzo (E5) — Agustín Lara

<!-- Horas estimadas, registro de horas trabajadas y consulta del esfuerzo (HU-19 a HU-21). -->

## 7. Defectos (E6) — Agustín Lara

<!-- Registrar un defecto, cambiar su estado y resolverlo (HU-22, HU-23). -->

## 8. Métricas y dashboard (E7 y E8) — Agustín Giorlando

<!-- Qué significa cada métrica (SP, velocidad, %, horas, desviación, defectos) y cómo leer el dashboard
     (HU-24 a HU-28). -->

## 9. Reportes y exportación a PDF (E9) — Agustín Giorlando

<!-- Reporte de sprint y de proyecto, y cómo guardarlos como PDF (HU-29 a HU-31). -->

## Anexo · Documentación técnica — Branko Almeira

<!-- Arquitectura (ADR 0001 a 0004), estructura de paquetes, cómo correr, cómo probar (make test, make cover,
     BDD), CI, cómo contribuir: SDD → BDD → TDD, resolución de conflictos en los archivos compartidos y chequeo de
     `Closes #N` antes de abrir un PR (acciones pendientes de las retros). -->
