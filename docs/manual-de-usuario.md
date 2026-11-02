# Manual de usuario · Software Metrics & Estimation

> Plantilla de TEC-07. Cada responsable completa la sección de su épica, y Branko completa la introducción y la
> documentación técnica. Escribí para alguien que nunca usó la app: qué puede hacer, paso a paso y con las reglas
> que va a encontrar (mensajes de error incluidos). Una o dos capturas por sección como máximo, en
> `docs/manual/img/`.

## 1. Primeros pasos — Branko Almeira

**Software Metrics & Estimation** es una aplicación web para seguir un proyecto ágil de punta a punta: el Product
Backlog, la estimación con Planning Poker, los sprints, las horas trabajadas, los defectos y las métricas que salen
de todo eso (velocidad, porcentaje completado, desviación de horas). Al final, arma los reportes para la Sprint Review
y para el cierre del proyecto.

### Levantar la aplicación

Hace falta **Go 1.22 o posterior** y una terminal en la carpeta del repositorio.

1. **(Opcional) Cargar datos de ejemplo:** `make demo` carga nuestro propio proyecto (el backlog real, los Sprints 1
   a 3 cerrados y el 4 en curso, horas y defectos). Sirve para recorrer la aplicación con datos ya cargados.
2. **Levantar la aplicación:** `make run`.
3. Abrir **http://localhost:8080** en el navegador.

Los datos se guardan en archivos JSON en la carpeta `data/` y siguen ahí la próxima vez que se levanta la app.
Para trabajar con otra carpeta (por ejemplo, una limpia para una demo): `DATA_DIR=otra make run`. Para usar otro
puerto: `PORT=9090 make run`.

> **No hay usuarios ni contraseñas.** Cuando una pantalla necesita saber quién actúa (por ejemplo, para votar en
> el Planning Poker o registrar horas), elegís tu nombre de la lista de integrantes del proyecto. Por eso, antes que
> nada, cargá el proyecto y sus integrantes (sección 2).

### La portada y la navegación

La **portada** muestra una tarjeta por sección con su ícono y una descripción corta. Tocá cualquiera para entrar.
Arriba de todas las pantallas está la **barra de navegación**, con la sección en la que estás resaltada:

| Sección | Para qué sirve | Ver |
|---|---|---|
| **Proyectos** | Alta y modificación de proyectos, sus fechas, su estado y sus integrantes | Sección 2 |
| **Backlog** | Las historias de usuario con prioridad, estado, Story Points y criterios de aceptación | Sección 3 |
| **Sprints** | Sprint Goal, historias asignadas, inicio, cierre e historial de sprints | Sección 4 |
| **Planning Poker** | Estimar una historia con votos ocultos y llegar a un acuerdo | Sección 5 |
| **Esfuerzo** | Las horas trabajadas por integrante, historia y actividad | Sección 6 |
| **Defectos** | Los defectos con su severidad, su estado y en qué sprint se detectaron y resolvieron | Sección 7 |
| **Dashboard** | Las métricas del proyecto con gráficos | Sección 8 |
| **Reportes** | El reporte de un sprint o del proyecto, para imprimir o descargar en PDF | Sección 9 |

### El recorrido típico

```text
Proyecto e integrantes → Historias del backlog → Estimación (Planning Poker) → Sprint (asignar, iniciar)
→ Registrar horas y defectos → Completar historias → Cerrar el sprint → Dashboard y reportes
```

### Convenciones de las pantallas

- **Colores de los estados:** gris = todavía no empezó; azul = en curso; verde = terminado. En los defectos, rojo =
  abierto o crítico y ámbar = en progreso o de severidad alta.
- **Números:** los decimales se escriben con coma (por ejemplo, `25,0 %`).
- **Errores:** si un formulario no cumple una regla, vuelve con lo que escribiste y el motivo arriba, en rojo. Cada
  sección lista sus reglas y los mensajes exactos.
- **Listas vacías:** cuando todavía no hay nada cargado, la pantalla lo dice y ofrece el botón para crear el primero.

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

La pantalla **Backlog** (`/backlog`) agrupa las historias por proyecto y las lista en una tabla: el ID, el título
(que lleva al detalle), la prioridad y el estado como insignias de color, los Story Points (o "—" si no está
estimada) y la cantidad de criterios. Con **Nueva historia** entrás al formulario de alta.

![Listado del backlog con las insignias de prioridad y estado](manual/img/backlog.jpg)

### Crear una historia (HU-05) — Ignacio Bianchi

1. En Backlog, tocá **Nueva historia**.
2. Elegí el proyecto, escribí el título (obligatorio), la descripción (opcional) y la prioridad (Alta, Media o Baja).
3. Tocá **Crear historia** y volvés al listado, donde la nueva historia ya aparece.

La aplicación asigna el identificador automáticamente (`HU-1`, `HU-2`, …) y es **único en todo el sistema**: no se
reinicia por proyecto. La historia nace en estado *Pendiente*, sin Story Points y sin criterios.

Mensajes de error que puede mostrar: "el título es obligatorio" si el título está vacío, y "el proyecto no existe".

### Criterios de aceptación de una historia (HU-07) — Ignacio Bianchi

En el detalle de la historia, la sección **Criterios de aceptación**:

![Detalle de una historia: criterios, estimación y cambio de estado](manual/img/historia-detalle.jpg)

- **Agregar:** escribí el criterio en el campo "Nuevo criterio de aceptación" y tocá **Agregar**. Queda al final de
  la lista, que se numera sola empezando en 1.
- **Quitar:** tocá **Quitar** en la fila del criterio. Los que quedan se renumeran (si quitás el 2 de tres criterios,
  el 3 pasa a ser el 2).

El texto no puede estar vacío: si lo está, la pantalla avisa "el criterio de aceptación no puede estar vacío".
Una historia **sin criterios todavía no está lista para un sprint**: lo rechaza al querer asignarla (HU-11, Sprints).

### Cambiar el estado de una historia (HU-08) — Ignacio Bianchi

En el detalle, el desplegable **Estado** con el botón **Cambiar estado**. Los estados van Pendiente → En progreso →
Completada, con transiciones controladas:

| Pasa de | A | Condición |
|---|---|---|
| Pendiente | En progreso | La historia tiene que estar asignada a un sprint |
| En progreso | Pendiente | Libre: se puede volver atrás |
| Pendiente o En progreso | Completada | Asignada a un sprint |

*Completada* es final: la historia ya no cambia de estado. Para completarla desde la aplicación se usa el botón
**Completar** de la pantalla del sprint (HU-12), para que el sprint lo registre.

Errores que puede mostrar: "no se puede pasar de Pendiente a Completada" (sin sprint) o el par de estados
correspondiente, y "la historia debe estar asignada a un sprint".

### Estimar una historia en Story Points (HU-15) — Ignacio Bianchi

En el detalle, el desplegable **Story Points** con la escala 1, 2, 3, 5, 8, 13, 21 y el botón **Estimar**. Ese valor
es el tamaño de la historia y alimenta las métricas de los sprints (E7).

- Una historia estimada se puede **re-estimar** cuantas veces haga falta mientras no esté completada.
- Una completada no: si lo intentás, la pantalla avisa "no se puede re-estimar una historia completada".
- También se puede acordar el valor con el equipo en el Planning Poker (sección 5): es la misma operación de fondo.

### Editar y eliminar una historia (HU-06) — Agustín Lara

Al final del detalle de la historia, la sección **Editar** trae el formulario con el título, la descripción y la
prioridad que tiene hoy.

- **Editar:** cambiá lo que necesites y tocá **Guardar cambios**. Solo se editan esos tres datos: el estado, los
  Story Points y los criterios tienen su propia sección. Valen las reglas del alta: el título es obligatorio y la
  prioridad es Alta, Media o Baja. Una historia asignada a un sprint se puede editar mientras no esté completada.
- **Eliminar:** tocá **Eliminar** y volvés al backlog, donde la historia ya no aparece. La app **no pide
  confirmación**: revisá antes de tocarlo. El ID de una historia eliminada no se vuelve a usar (si borrás la HU-3,
  la próxima historia nueva no es la HU-3).

| Regla | Mensaje |
|---|---|
| El título es obligatorio | "el título es obligatorio" |
| La prioridad es Alta, Media o Baja | "la prioridad debe ser Alta, Media o Baja" |
| Una historia completada no se edita: ya es parte del resultado de un sprint | "una historia completada no se puede editar" |
| No se elimina una historia completada | "no se puede eliminar una historia completada" |
| No se elimina una historia asignada a un sprint (primero hay que quitarla del sprint) | "no se puede eliminar una historia asignada a un sprint" |
| No se elimina una historia con horas o defectos cargados: quedarían registros sin historia | "no se puede eliminar una historia con esfuerzo o defectos registrados" |

### Ver el backlog ordenado y filtrar por estado (HU-09) — Agustín Lara

La pantalla **Backlog** muestra el backlog de cada proyecto en el orden en que conviene trabajarlo, para saber qué
sigue:

1. Primero por **prioridad**: Alta, después Media y al final Baja.
2. A igual prioridad, por **número de historia**: HU-2 va antes que HU-10.

Arriba de la lista, el desplegable **Todos / Pendiente / En progreso / Completada** con el botón **Filtrar** deja
solo las historias de ese estado, en el mismo orden. Por ejemplo, **Pendiente** muestra lo que todavía no se empezó.
Para volver a ver todo, elegí **Todos**.

- Un proyecto sin historias, o sin ninguna en el estado elegido, no aparece en la lista. Si no queda ninguno, la
  pantalla dice "No hay historias en el backlog." con un enlace para crear la primera.
- Si la dirección trae un estado que no existe (por ejemplo, un enlace mal copiado), la app avisa "el estado debe ser
  Pendiente, En progreso o Completada".

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

La pantalla **Planning Poker** (`/planning-poker`) muestra todas las historias del backlog en dos tablas:
**Pendientes** (las que se pueden votar) y **Completadas** (no se votan, y si tenían estimación la muestran).
La regla de la votación es que las cartas quedan **ocultas** hasta que alguien las revela: mientras la votación está
abierta, cada fila solo muestra "Votaron: …" con los nombres de quienes ya votaron.

La aplicación no tiene login: el equipo se entiende eligiendo el nombre en la lista ("Votás como…"). Es una
limitación aceptada para una herramienta interna, y solo votan los integrantes del proyecto de la historia.

### Abrir una votación

En la fila de la historia, tocá **Abrir votación** (o **Abrir nueva votación** si ya había una estimación de antes).
La votación arranca en la ronda 1, sin votos. Una historia solo puede tener una votación abierta a la vez; si la
historia está completada, no se puede abrir ("no se puede votar una historia completada").

### Votar

Elegí tu nombre en "Votás como…" y tocá una carta: 1, 2, 3, 5, 8, 13 o 21, los botones de la fila. El valor queda
registrado sin mostrarse, y tu nombre pasa a la lista "Votaron: …" (mientras nadie votó, muestra "Votaron: todavía
nadie"). Cada integrante vota **una vez por ronda**: si lo intentás de nuevo, la pantalla avisa "el integrante ya
votó en esta ronda".

![Votación abierta: cartas de la escala y lista de votantes oculta](manual/img/planning-poker-votando.jpg)

### Revelar los votos

Cuando ya votaron, en la fila aparecen "Revelo como…" y el botón **Revelar votos**. La revelación es tarea del
Agile Enabler: si el proyecto no tiene uno, la fila lo avisa con "Se necesita Agile Enabler". Al revelar se muestran
todos los votos de la ronda (integrante y carta) y:

- si coinciden, una insignia de éxito: "Hay consenso en 5." (un único voto también es consenso);
- si no, la diferencia: "No hay consenso: el voto más alto es 8 (…) y el voto más bajo es 2 (…)", con quiénes votaron
  cada extremo. Esos dos grupos son los que explican su voto antes de la próxima ronda.

![Ronda revelada con diferencias: votos de la tabla, más alto y más bajo, y acciones al pie](manual/img/planning-poker-revelado.jpg)

Después de revelar, la ronda cierra: no acepta votos ("la ronda ya fue revelada").

### Nueva ronda (si falta acuerdo)

Si no hubo consenso, se discute y tocá **Nueva ronda**: la ronda actual pasa al historial (el desplegable "Historial
de rondas" muestra cada una con sus votos) y arranca la siguiente, sin votos, para volver a votar. Solo se puede abrir
una ronda nueva después de revelar la actual: si no, avisa "primero hay que revelar los votos de la ronda".

### Registrar la estimación acordada

Cuando el equipo se pone de acuerdo, elegí el valor en el desplegable **Estimar** y tocá **Registrar**. Puede ser
cualquier valor de la escala, aunque nadie lo haya votado. La historia queda con esos Story Points y la votación se
**cierra**: la fila pasa a mostrar "Estimación acordada: 5 SP" y aparece **Abrir nueva votación** por si más
adelante hay que re-estimarla.

### Las reglas, en corto

- Una sola votación abierta por historia; no se ve el valor de ningún voto antes de revelar.
- Una carta por integrante por ronda; la revelación es del Agile Enabler.
- Las cartas y la estimación acordada siempre son valores de la escala: 1, 2, 3, 5, 8, 13 o 21
  ("los Story Points deben ser 1, 2, 3, 5, 8, 13 o 21").
- "Primero hay que revelar los votos de la ronda" aparece si se quiere acordar o abrir ronda antes de tiempo;
  "no hay votos para revelar", si nadie votó todavía.

## 6. Esfuerzo (E5) — Agustín Lara

La pantalla **Esfuerzo** (`/esfuerzo`) junta todo lo que tiene que ver con horas: el formulario para cargar las
horas trabajadas, la lista de lo cargado, el esfuerzo de cada proyecto y las horas estimadas de cada historia. Las
horas se escriben y se muestran con coma decimal (1,5 es una hora y media), aunque los formularios también aceptan
punto.

### Registrar horas trabajadas (HU-20)

1. En **Registrar horas trabajadas**, elegí el **integrante** (al lado del nombre figura su proyecto) y la
   **historia** en la que trabajó.
2. Completá la **fecha** (viene cargada la de hoy), la **actividad** (por ejemplo, Desarrollo, Tests o Revisión de
   PR) y las **horas**. Se aceptan fracciones, como 0,5 o 2,5.
3. Tocá **Registrar horas**. El registro aparece en la lista **Horas registradas**, agrupado por historia.

Si algo no cumple las reglas, el formulario vuelve con lo que escribiste y el motivo:

| Regla | Mensaje |
|---|---|
| La actividad es obligatoria | "la actividad es obligatoria" |
| La fecha es obligatoria | "la fecha es obligatoria" |
| Las horas son mayores que cero | "las horas deben ser mayores que cero" |
| La fecha no puede ser futura (hoy sí vale) | "la fecha no puede ser futura" |
| Una persona no carga más de 24 horas en un mismo día, sumando todos sus registros de ese día | "el integrante superaría las 24 horas en ese día" |
| El integrante tiene que ser del proyecto de la historia | "el integrante no pertenece al proyecto de la historia" |

### Consultar el esfuerzo de un proyecto (HU-21)

En **Esfuerzo por proyecto**, cada proyecto tiene su bloque con tres tarjetas:

- **Horas registradas:** el total de horas cargadas en las historias del proyecto.
- **Historias con esfuerzo:** cuántas historias tienen al menos un registro.
- **Integrantes con esfuerzo:** cuántas personas cargaron horas.

Debajo, dos tablas con el detalle: las horas **por historia** (en orden de número, con un enlace a cada una) y las
horas **por integrante** (en orden alfabético). Las dos suman lo mismo que la tarjeta de horas registradas. Un
proyecto sin horas cargadas muestra "Todavía no hay esfuerzo registrado." con un enlace al formulario.

![Pantalla de esfuerzo con las tarjetas y los totales de un proyecto](manual/img/esfuerzo.jpg)

### Estimar las horas de una historia (HU-19)

En **Horas estimadas por historia**, cada historia del backlog tiene un campo para escribir cuántas horas se calcula
que va a llevar y el botón **Guardar**. Las que todavía no tienen estimación dicen "Sin estimar".

- Una historia tiene **una sola** estimación de horas: si guardás otra, reemplaza a la anterior.
- Las horas tienen que ser mayores que cero ("las horas deben ser mayores que cero"). El error aparece en la fila de
  esa historia.
- Es distinta de los Story Points: los Story Points miden el tamaño relativo (sección 3) y las horas estimadas se
  comparan con las horas reales en las métricas (sección 8).

## 7. Defectos (E6) — Agustín Lara

La pantalla **Defectos** (`/defectos`) lista los defectos encontrados, con su severidad y su estado como insignias de
color, y abajo tiene el formulario para registrar uno nuevo.

![Lista de defectos con las insignias de severidad y estado](manual/img/defectos.jpg)

### Registrar un defecto (HU-22)

1. En **Registrar defecto**, escribí la **descripción** del problema.
2. Elegí la **severidad** (Crítica, Alta, Media o Baja), la **historia** en la que apareció y el **sprint de
   detección**, es decir, el sprint en el que se encontró.
3. Tocá **Registrar defecto**. Aparece en la lista con un ID (`D1`, `D2`, …) y en estado **Abierto**.

Se puede registrar un defecto en una historia ya completada (suele ser justo cuando aparecen) y en un sprint ya
cerrado, si se encontró en ese sprint y se carga tarde.

| Regla | Mensaje |
|---|---|
| La descripción es obligatoria | "la descripción es obligatoria" |
| La severidad es una de las cuatro de la lista | "la severidad debe ser Crítica, Alta, Media o Baja" |
| El sprint de detección es del mismo proyecto que la historia | "el sprint de detección no pertenece al proyecto de la historia" |

Las insignias de severidad van de más grave a menos grave: **Crítica** en rojo, **Alta** en ámbar, **Media** en
azul y **Baja** en gris.

### Seguir y resolver un defecto (HU-23)

Un defecto pasa por tres estados:

| Estado | Qué significa | Insignia |
|---|---|---|
| **Abierto** | Se encontró y nadie lo está corrigiendo | Rojo |
| **En progreso** | Alguien lo está corrigiendo | Ámbar |
| **Resuelto** | Ya se corrigió, en un sprint de resolución | Verde |

En la columna **Acciones** de cada defecto:

- **Pasar a En progreso / Pasar a Abierto:** mueve el defecto entre esos dos estados, para cualquiera de los dos
  lados.
- **Resolver:** elegí el **sprint de resolución** en el desplegable y tocá **Resolver**. Se puede resolver desde
  Abierto o desde En progreso. El desplegable solo ofrece los sprints del mismo proyecto que no empiezan antes que el
  de detección, y puede ser el mismo sprint.

**Resuelto es final:** un defecto resuelto ya no ofrece acciones y no vuelve a cambiar de estado. Si la app rechaza
un cambio, avisa con alguno de estos mensajes:

| Situación | Mensaje |
|---|---|
| Cambiar un defecto resuelto o pasarlo al estado que ya tiene | "no se puede pasar de Resuelto a Abierto" (con los estados que correspondan) |
| El sprint de resolución es de otro proyecto | "el sprint de resolución no pertenece al proyecto del defecto" |
| El sprint de resolución empieza antes que el de detección | "el sprint de resolución no puede ser anterior al de detección" |

## 8. Métricas y dashboard (E7 y E8) — Agustín Giorlando

El **Dashboard** (`/dashboard`) junta en una sola pantalla cómo viene un proyecto: un resumen con las métricas en
tarjetas y cuatro gráficos con la evolución sprint a sprint. Sirve para la Sprint Review o para mostrar el estado del
proyecto sin tener que recorrer el backlog, los sprints, el esfuerzo y los defectos por separado. La pantalla solo
muestra datos: para cambiarlos se usan las otras secciones.

![Dashboard con las tarjetas del resumen y los cuatro gráficos, con los datos de demo](manual/img/dashboard.jpg)

### Elegir el proyecto

En la navegación, tocá **Dashboard**. Elegí el proyecto en el desplegable **Proyecto** y tocá **Ver dashboard**.
Para mirar otro proyecto, usá **Cambiar de proyecto**, arriba a la derecha.

Si todavía no hay proyectos, en lugar del desplegable aparece "Todavía no hay proyectos." con el enlace **Crear el
primer proyecto** (sección 2). Si se abre la dirección de un proyecto que no existe, la app muestra "el proyecto no
existe".

### El resumen

Debajo del título está el nombre del proyecto, su estado con una insignia de color (gris *Planificado*, azul
*En curso*, verde *Finalizado*; ver sección 2) y sus fechas. Las tarjetas del resumen muestran:

| Tarjeta | Qué muestra |
|---|---|
| Días restantes | Los días que faltan para la fecha de fin del proyecto |
| Sprint activo | El sprint en curso, o "Ninguno" |
| Velocidad del proyecto | El promedio de Story Points completados de los sprints **cerrados** |
| Historias completadas | El porcentaje de historias del proyecto que están completadas |
| Horas estimadas | La suma de las horas estimadas de todas las historias del proyecto |
| Horas reales | La suma de todas las horas registradas en el proyecto |
| Desviación | Cuánto se apartaron las horas reales de las estimadas, en porcentaje, o "no calculable" |
| Defectos detectados, resueltos y abiertos | Los defectos del proyecto: los abiertos son los detectados que todavía no se resolvieron |

Los números con decimales se muestran con coma y un decimal, igual que en el resto de la app (por ejemplo, 10,3).

### Qué significa cada métrica

Las mismas métricas aparecen en el detalle del sprint (sección 4) y en los reportes (sección 9), y se calculan igual
en todas las pantallas:

- **Story Points planificados:** la suma de los SP de las historias asignadas al sprint.
- **Story Points completados:** la suma de los SP de las historias completadas en el sprint. Al cerrar el sprint, los
  dos valores quedan fijos.
- **Velocidad:** el promedio de SP completados de los sprints cerrados. Sin sprints cerrados es 0, y un sprint cerrado
  sin historias completadas también cuenta (con 10 y 0 SP completados, la velocidad es 5). El sprint en curso no
  cuenta hasta que se cierra.
- **Porcentaje de historias completadas:** completadas dividido el total de historias, por 100. Cuenta historias, no
  Story Points: 1 de 3 da 33,3 %. Sin historias es 0 %.
- **Horas reales:** la suma de las horas registradas (sección 6).
- **Desviación:** (horas reales − horas estimadas) / horas estimadas × 100. Si es **positiva**, se trabajó más de lo
  estimado; si es **negativa**, menos. Si no hay horas estimadas no se puede calcular, y la tarjeta dice
  "no calculable". Con horas estimadas y ninguna hora registrada, da −100 %.
- **Defectos detectados y resueltos:** un defecto cuenta como detectado en su sprint de detección y como resuelto en
  su sprint de resolución (sección 7). Uno que todavía no tiene sprint de resolución no está resuelto.

### Los gráficos

Debajo del resumen hay cuatro gráficos, en dos columnas (en un celular o una ventana angosta pasan a una sola
columna). Cada uno tiene un punto o una barra por **sprint cerrado**, ordenados por fecha de inicio. El sprint en
curso no aparece porque sus valores todavía cambian: se sigue desde el resumen y desde su detalle.

- **Velocidad por sprint:** los SP completados de cada sprint (línea) y, en línea punteada, la velocidad promedio del
  proyecto. Muestra si el equipo acelera o se frena respecto de su promedio.
- **Story Points planificados vs. completados:** una barra de cada uno por sprint. Si los completados quedan muy por
  debajo de los planificados, el equipo se está comprometiendo de más.
- **Horas estimadas vs. reales:** las horas estimadas de las historias planificadas en el sprint y las horas
  registradas entre su fecha de inicio y de fin.
- **Defectos detectados vs. resueltos:** cuántos defectos se encontraron y cuántos se resolvieron en cada sprint.

Al pasar el mouse sobre un punto o una barra se ve el valor exacto. Tocá un elemento de la leyenda para ocultar o
mostrar esa serie.

### Cuando todavía no hay datos

Un gráfico sin datos no se dibuja vacío: en su lugar aparece un mensaje con un enlace a la sección que lo completa.

| Mensaje | Cuándo aparece | Enlace |
|---|---|---|
| "Todavía no hay sprints cerrados para graficar" | En los cuatro gráficos, mientras el proyecto no tenga ningún sprint cerrado | **Ir a Sprints** |
| "Todavía no hay horas registradas" | En el de horas, si los sprints cerrados no tienen horas estimadas ni registradas | **Registrar esfuerzo** |
| "No hay defectos registrados" | En el de defectos, si los sprints cerrados no tienen defectos | **Ir a Defectos** |

En un proyecto recién creado el resumen igual se ve completo: las métricas en 0, la desviación "no calculable" y los
cuatro mensajes en lugar de los gráficos. A medida que se cierran sprints, los gráficos se van llenando.

## 9. Reportes y exportación a PDF (E9) — Agustín Giorlando

La sección **Reportes** (`/reportes`) arma dos documentos listos para compartir: el **reporte de un sprint**, para la
Sprint Review, y el **reporte del proyecto**, para el cierre o para mostrar cómo viene. Los dos usan los mismos
cálculos que el dashboard y el detalle del sprint (ver "Qué significa cada métrica" en la sección 8), así que los
números coinciden en todas las pantallas. Igual que el dashboard, los reportes solo muestran datos.

### Elegir un reporte

En la navegación, tocá **Reportes**. La pantalla muestra un bloque por proyecto con:

- el botón **Reporte del proyecto**, y
- una tabla con los sprints que tienen reporte (Sprint Goal, inicio, fin y estado). Tocá el Sprint Goal para abrir el
  reporte de ese sprint.

Solo aparecen los sprints **activos** y **cerrados**. Si un proyecto todavía no inició ninguno, el bloque dice
"Todavía no hay sprints iniciados: un sprint planificado no tiene resultado para reportar." Si no hay proyectos, la
pantalla ofrece el enlace **Creá un proyecto** (sección 2).

### El reporte de un sprint (HU-29)

Dirección: `/reportes/sprints/{sprint}`, por ejemplo `/reportes/sprints/S3`. Debajo del título están el Sprint Goal,
el estado y las fechas del sprint, y después cuatro partes:

| Parte | Qué muestra |
|---|---|
| Historias planificadas | Cada historia asignada al sprint, en el orden en que se asignó: ID, título, Story Points, horas estimadas, horas reales y si se completó (✔ o "No") |
| Métricas | SP planificados y completados, porcentaje de historias completadas, horas estimadas, horas reales y la desviación (o "desviación no calculable (sin horas estimadas)") |
| Esfuerzo registrado | Cada registro de horas con fecha **dentro** del sprint (fecha, integrante, historia, actividad y horas), ordenados por fecha, y abajo el **total por integrante** |
| Defectos | Cuántos defectos se detectaron y cuántos se resolvieron en el sprint, y la lista de los detectados en él con su severidad y su estado |

Las horas reales de cada historia son solo las registradas entre la fecha de inicio y la de fin del sprint, ambas
incluidas. En un sprint cerrado, los Story Points planificados y completados son los que quedaron fijos al cerrarlo
(sección 4).

![Reporte de un sprint activo, con la marca de reporte parcial y los botones para descargar o imprimir](manual/img/reporte-sprint.jpg)

### El reporte del proyecto (HU-30)

Dirección: `/reportes/proyecto?proyecto=P1` (cambiá `P1` por el ID del proyecto). Debajo del título están el nombre,
el estado, las fechas y los días restantes del proyecto, y después tres partes:

| Parte | Qué muestra |
|---|---|
| Resumen de sprints | Una fila por sprint activo o cerrado, por fecha de inicio: estado, SP planificados y completados, % de historias completadas, horas estimadas y reales, y defectos detectados y resueltos. Tocá el sprint para abrir su reporte |
| Métricas acumuladas | La velocidad (con cuántos sprints cerrados se calculó), los SP completados, el porcentaje de historias completadas del proyecto ("12 de 15"), las horas estimadas y reales y la desviación |
| Defectos | Cuántos se detectaron, cuántos se resolvieron y cuántos siguen abiertos, y la lista de los **abiertos**, de más grave a menos grave |

Cada fila del resumen tiene los mismos valores que el reporte de ese sprint. En cambio, las horas reales de las
métricas acumuladas suman **todo** el esfuerzo del proyecto, también el registrado fuera de las fechas de un sprint;
el reporte lo aclara debajo de la lista.

### Guardar un reporte en PDF (HU-31)

Arriba de cada reporte hay dos botones. Ninguno de los dos aparece en el papel ni en el PDF.

- **Descargar PDF:** la app genera el archivo y el navegador lo descarga. Se llama `reporte-sprint-S3.pdf` o
  `reporte-proyecto-P1.pdf`, es una hoja A4 con las mismas partes y los mismos números que la pantalla, y debajo del
  título dice "Generado el" con la fecha del día. Es la opción recomendada para mandar el reporte por mail o subirlo.
- **Imprimir o guardar como PDF:** abre el diálogo de impresión del navegador. Elegí una impresora o, como destino,
  **Guardar como PDF** (en Chrome y Edge) o **Microsoft Print to PDF**. Para imprimir, la app esconde la barra de
  navegación, el pie y los botones, y evita cortar una parte entre dos hojas. Sirve como alternativa si la descarga
  falla, o para imprimir en papel.

Las direcciones de descarga también se pueden abrir directo: `/reportes/sprints/{sprint}/pdf` y
`/reportes/proyecto.pdf?proyecto=P1`.

### Sprints planificados y sprints activos

- **Un sprint planificado no tiene reporte:** todavía no empezó, así que no hay resultado. No aparece en la lista de
  Reportes ni en el resumen del reporte del proyecto, y si se abre su dirección (o la de su PDF) la app muestra "el
  sprint todavía no empezó".
- **Un sprint activo tiene un reporte parcial:** arriba aparece, resaltado, "Reporte parcial: el sprint sigue activo
  y el resultado todavía puede cambiar." Los SP son los del momento y cambian a medida que se completan historias. El
  PDF lleva la misma marca. En el reporte del proyecto, el sprint activo figura en el resumen con su estado, pero no
  cuenta para la velocidad hasta que se cierra.

| Situación | Mensaje |
|---|---|
| El sprint todavía está planificado | "el sprint todavía no empezó" |
| El sprint no existe | "el sprint no existe" |
| El proyecto no existe (o falta `?proyecto=`) | "el proyecto no existe" |

## Anexo · Documentación técnica — Branko Almeira

Esta parte es para quien tenga que mantener o extender la aplicación. Las reglas completas de trabajo están en
[`CONTRIBUTING.md`](../CONTRIBUTING.md) y las decisiones de diseño en [`docs/adr/`](adr/).

### A.1 Arquitectura

Una aplicación web **en Go, sin un segundo stack**: el servidor arma el HTML y el navegador solo lo muestra.

| Pieza | Tecnología | Decisión |
|---|---|---|
| Lenguaje y servidor | Go 1.22, `net/http` de la biblioteca estándar (rutas `GET /sprints/{sprint}`) | [ADR 0001](adr/0001-arquitectura.md) |
| Pantallas | `html/template` + Pico CSS + HTMX; gráficos con Chart.js | ADR 0001 |
| Persistencia | Repositorios en memoria que se guardan en archivos JSON (`data/`) | ADR 0001 |
| BDD | Gherkin en español con [godog](https://github.com/cucumber/godog) | ADR 0001 |
| Exportación a PDF | `github.com/go-pdf/fpdf` (HU-31) | ADR 0001 |
| Contratos entre épicas | Cada épica guarda **IDs** de las demás, nunca copias; las reglas que cruzan épicas las coordina un servicio | [ADR 0002](adr/0002-contratos-de-dominio.md) |
| Enumeradores | `int` con `iota + 1` y `String()`; el cero es "sin valor". Las plantillas no comparan textos de estado: el handler decide (por ejemplo, el color de una insignia con `tonoEstado`) | [ADR 0003](adr/0003-enumeradores-del-dominio.md) |
| Identidad | Sin autenticación: quien actúa elige su nombre de la lista de integrantes | [ADR 0004](adr/0004-identidad-sin-autenticacion.md) |

Pico CSS, HTMX y Chart.js están **incluidos en el repositorio** (`web/static/vendor/`), así la app funciona sin
internet, por ejemplo durante una presentación.

### A.2 Estructura de paquetes

Las dependencias van en un solo sentido: de afuera (HTTP) hacia adentro (dominio). El dominio no depende de
ninguna otra capa.

```text
cmd/server ─► internal/app ─► internal/server ─► internal/service ─► internal/domain/<épica>
                   │                                   │
                   └─► internal/store/memory           ├─► internal/report ─► internal/metrics
                         (repositorios JSON)           └─► internal/dashboard
```

`internal/metrics` no importa a ningún otro paquete del proyecto: recibe números y fechas y devuelve resultados, así
se prueba sin armar entidades.

| Paquete | Responsabilidad |
|---|---|
| `internal/domain/{project,backlog,sprint,estimation,effort,defect}` | Entidades y reglas de negocio de cada épica, con sus errores. Acá se hace TDD |
| `internal/metrics` | Cálculo de métricas como funciones puras: SP planificados y completados, velocidad, porcentaje, desviación, defectos, horas en un período |
| `internal/report` | Arma el reporte de sprint (HU-29) y el de proyecto (HU-30) a partir de los datos que recibe; no accede a repositorios |
| `internal/dashboard` | Arma los datos del dashboard (HU-28) |
| `internal/service` | Casos de uso: buscan en los repositorios, aplican las reglas del dominio y guardan |
| `internal/store/memory` | Repositorios en memoria con persistencia a JSON (una ruta vacía = sin archivo, para los tests) |
| `internal/server` | Handlers HTTP, rutas y plantillas. Error de formulario → 400 con el formulario; recurso inexistente → 404; éxito → 303 |
| `internal/app` | Arma la aplicación: crea los repositorios y los servicios y se los pasa al servidor |
| `cmd/server`, `cmd/demo` | La aplicación (`make run`) y la carga de datos de ejemplo (`make demo`) |
| `features` | Escenarios BDD (`*.feature`) y sus pasos |
| `tools/checkcoverage` | Verifica el umbral de cobertura en el CI |

### A.3 Cómo correr y configurar

| Quiero… | Comando |
|---|---|
| Levantar la app | `make run` → http://localhost:8080 |
| Cargar los datos de ejemplo | `make demo` |
| Usar otra carpeta de datos / otro puerto | `DATA_DIR=otra make run` / `PORT=9090 make run` |
| Empezar de cero | Borrar la carpeta de datos (o usar una nueva con `DATA_DIR`) |

### A.4 Cómo probar

| Quiero… | Comando |
|---|---|
| Todo: `go vet`, tests unitarios con detector de carreras y BDD | `make test` |
| Solo los tests unitarios | `make unit` |
| Solo los escenarios BDD | `make bdd` |
| Cobertura total y reporte HTML (`coverage.html`) | `make cover` |
| Formatear | `make fmt` |

- **Tests unitarios:** en cada paquete (`*_test.go`), con tablas de casos. Los handlers se prueban con `httptest`.
- **BDD:** cada historia tiene su `features/hu-XX-*.feature`, en español, y sus pasos en `features/huXX_steps_test.go`.
  Los pasos comunes (por ejemplo, `la página contiene "…"` o `la página no contiene "…"`) están en
  `features/http_steps_test.go`. Un escenario con la etiqueta `@pendiente` no se corre.

### A.5 Integración continua

Cada PR y cada push a `main` corren [`.github/workflows/ci.yml`](../.github/workflows/ci.yml):

1. Formato (`gofmt`): falla si algún archivo no está formateado.
2. Análisis estático (`go vet`).
3. Tests unitarios con `-race` y cobertura.
4. **Umbral de cobertura:** `internal/domain` e `internal/metrics` en **80 % o más**, ponderado por sentencias
   (`tools/checkcoverage`).
5. Escenarios BDD.

Además, [`revisor.yml`](../.github/workflows/revisor.yml) le pide la revisión del PR a quien le toca según la
rotación (Sosa → Bianchi → Lara → Giorlando → Branko → Sosa). `main` está protegida: solo entra un PR con el CI en
verde y una aprobación.

**Despliegue continuo:** cada merge a `main` que cambia código corre
[`cd-azure.yml`](../.github/workflows/cd-azure.yml), que publica la app en Azure Container Apps con HTTPS y los datos
de la demo. Cómo funciona, la URL, los secretos y cómo volver atrás: [despliegue.md](despliegue.md).

### A.6 Cómo contribuir: SDD → BDD → TDD

Cada historia recorre la misma cadena, con el mismo ID en todos los artefactos:

```text
Issue HU-XX → specs/HU-XX-*.md → criterios CA-XX.n → features/hu-xx-*.feature → *_test.go → código Go
```

1. **SDD:** la spec (`specs/_plantilla.md`) define entradas, salidas, reglas de negocio (RN), casos límite, errores
   con su mensaje exacto y criterios de aceptación (CA). Es la fuente de verdad.
2. **BDD:** cada criterio tiene al menos un escenario en el `.feature`.
3. **TDD:** un commit por fase: `test(HU-XX): RED - …` (el test falla), `feat(HU-XX): GREEN - …` (el código mínimo
   que lo hace pasar) y, si hace falta, `refactor(HU-XX): …`. Por eso los PRs se mergean sin *squash*.
4. **PR:** rama `feature/HU-XX-…`, plantilla con la tabla de trazabilidad y una fila en
   [`docs/ia/registro.md`](ia/registro.md) si se usó IA.

Para agregar una pantalla, seguí la sección 6 de `CONTRIBUTING.md`.

### A.7 Antes de abrir un PR (acciones de las retrospectivas)

**`Closes #N` solo si el PR termina el ítem.** `Closes` cierra el issue al mergear y lo mueve a *Done* en el
tablero. Para un PR que hace una parte, usá `Parte de #N`. Antes de abrirlo, confirmá el número:

```bash
gh issue view N    # ¿es el ítem correcto? ¿este PR lo termina?
```

**Conflictos en los archivos compartidos: quedarse con los dos lados.** Varios PRs agregan líneas a los mismos
archivos al mismo tiempo. En un conflicto, **nunca se elige un lado**: se conservan las líneas de `main` y después
las de la rama.

| Archivo | Qué agrega cada PR | Cómo resolver |
|---|---|---|
| `docs/ia/registro.md` | Una fila de la tabla | Las dos filas, primero la de `main` |
| `features/features_test.go`, `features/mundo_test.go` | El registro de sus pasos y sus campos del mundo | Los dos bloques |
| `internal/app/app.go`, `internal/server/server.go` | Su servicio, sus repositorios y sus rutas | Los dos bloques; si cambió la firma de un constructor en `main`, usar la nueva |
| `web/static/app.css` | Sus clases, al final y bajo un comentario con su ítem (`/* HU-32.3 */`) | Los dos bloques |

```bash
git switch <rama> && git fetch origin && git merge origin/main
# resolver cada conflicto conservando los dos lados, sin dejar marcas <<<<<<< ======= >>>>>>>
grep -rn '^<<<<<<<\|^>>>>>>>' . && echo "quedan marcas"
make test && git commit && git push
```
