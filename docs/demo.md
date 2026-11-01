# Guía de la demo · Presentación final (TEC-22)

Cómo preparar el entorno y qué mostrar en la demo del sistema. Complementa las diapositivas y el guion general de la
presentación, que arma el Agile Enabler ([`docs/scrum/sprint-4/guion-presentacion.md`](scrum/sprint-4/guion-presentacion.md)). Cada integrante presenta su épica; el recorrido
completo de una funcionalidad (historia → spec → criterios → BDD → tests → código) lo hace Ignacio Bianchi con HU-16.

## 1. Preparar el entorno

Requisitos: Go 1.22 o más nuevo y `make` (en Windows, desde WSL). No hace falta internet: Pico CSS, htmx y Chart.js
están dentro del repositorio (`web/static/vendor/`).

```bash
git checkout main && git pull        # la demo se hace con lo que está en main
make test                            # todo en verde antes de presentar
rm -rf data/                         # arrancar de cero (data/ no se versiona)
make demo                            # carga nuestro proyecto en data/
make run                             # http://localhost:8080
```

La aplicación no tiene un `main.go` en la raíz: el servidor es `cmd/server` (`make run` = `go run ./cmd/server`) y
la carga de datos es `cmd/demo` (`make demo` = `go run ./cmd/demo`).

`make demo` termina con:

```
Listo: "Software Metrics & Estimation" (P1) cargado en data. Levantá la app con `make run` y abrí /dashboard?proyecto=P1
```

### Qué carga la demo

Nuestro propio proyecto, con los mismos servicios que usa la aplicación (los datos pasan por todas las validaciones):

- el proyecto **Software Metrics & Estimation** (21/10 al 02/11/2026) con los cinco integrantes y sus roles;
- el Product Backlog HU-01 a HU-31, cada historia con su criterio, sus Story Points y sus horas estimadas;
- los Sprints 1 a 3 **cerrados** y el Sprint 4 **en curso** (S4, con HU-31 sin completar);
- horas trabajadas por día en cada sprint y 9 defectos (8 resueltos y 1 abierto).

Las horas se estimaron a 2 horas por Story Point, porque no se registraron durante el proyecto; todo lo demás es tal
cual pasó (ver `cmd/demo/datos.go`).

### Problemas conocidos

| Síntoma | Causa | Qué hacer |
|---|---|---|
| `cargar la demo: ya existe un proyecto con ese nombre` | La demo ya estaba cargada en `data/` | Nada, si los datos están bien. Para empezar de cero: `rm -rf data/` y `make demo` |
| `listen tcp :8080: bind: address already in use` | Quedó otra instancia de la app abierta | Cerrarla, o levantar en otro puerto: `PORT=8081 make run` |

## 2. Verificación antes de presentar

Con la app levantada sobre los datos de demo, revisar que todas estas pantallas abran (verificado el 01/11/2026 sobre
`main`, todas con código 200):

- [ ] Portada `/` y `/health`
- [ ] `/proyectos` y `/proyectos/P1`
- [ ] `/backlog` (con el filtro por estado)
- [ ] `/sprints`, `/sprints/S3` (cerrado) y `/sprints/S4` (activo)
- [ ] `/planning-poker`
- [ ] `/esfuerzo` y `/defectos`
- [ ] `/dashboard` y `/dashboard?proyecto=P1`, con los cuatro gráficos dibujados
- [ ] `/reportes`, `/reportes/proyecto?proyecto=P1`, `/reportes/sprints/S3` y `/reportes/sprints/S4` ("Reporte parcial")
- [ ] La descarga en PDF de los dos reportes

Antes de empezar: navegador en pantalla completa, zoom al 125 % para que se lea desde lejos y una pestaña ya abierta
en `/dashboard?proyecto=P1`.

## 3. Guion de la demo

Orden sugerido: del estado general al detalle, para que el público vea primero el resultado y después cómo se llega
a él. Cada bloque lo presenta el responsable de la épica (ver el manual de usuario).

### 3.1 Portada y proyecto (E1)

1. Portada `/`: las secciones como tarjetas, una por épica.
2. **Proyectos** → *Software Metrics & Estimation*: estado **En curso** con su insignia azul, días restantes, sprint
   activo S4, porcentaje de historias completadas e integrantes con sus roles.

### 3.2 Dashboard: el estado del proyecto de un vistazo (E7 y E8)

1. **Dashboard** → elegir el proyecto en **Proyecto** → **Ver dashboard** (filtro por proyecto). Contar que con
   **Cambiar de proyecto** se vuelve al selector.
2. **Resumen**, con los datos de demo:

   | Métrica | Valor | Qué decir |
   |---|---|---|
   | Velocidad del proyecto | 34,0 | Promedio de SP completados de los tres sprints cerrados (44, 23 y 35) |
   | Historias completadas | 96,8 % | 30 de 31: solo falta HU-31, que está en el sprint en curso |
   | Horas estimadas / reales | 214,0 / 227,7 | Trabajamos un poco más de lo estimado |
   | Desviación | 6,4 % | Positiva: más horas que las estimadas |
   | Defectos | 9 detectados · 8 resueltos · 1 abierto | El abierto es "el dashboard muestra los números con punto decimal", del Sprint 3: lo resuelve HU-32.4, y por eso esta pantalla ya muestra coma |

3. **Gráficos** (solo sprints cerrados: S1, S2 y S3):
   - *Velocidad por sprint*: 44, 23 y 35 SP, con la línea punteada del promedio (34).
   - *SP planificados vs. completados*: en los tres sprints se completó todo lo planificado.
   - *Horas estimadas vs. reales*: 88/92,1, 46/50,8 y 70/82,8; el Sprint 3 es el de mayor desviación.
   - *Defectos detectados vs. resueltos*: 3/3, 2/2 y 4/3.
   - Pasar el mouse por una barra para ver el valor exacto y tocar una serie de la leyenda para ocultarla.
4. Por qué el sprint en curso no está en los gráficos: sus valores todavía cambian; se sigue desde el resumen
   (sprint activo S4).
5. Si el tiempo alcanza: achicar la ventana para mostrar que tarjetas y gráficos pasan a una columna.

### 3.3 Backlog y Planning Poker (E2 y E4)

1. **Backlog**: historias ordenadas por prioridad, con insignias de prioridad y estado. Usar el filtro **Filtrar por
   estado** → *Pendiente*: queda solo HU-31.
2. Detalle de una historia: criterios de aceptación y Story Points.
3. **Planning Poker**: abrir una votación, votar con dos integrantes, revelar y registrar la estimación acordada.
   Hacerlo sobre una historia que no se use después, porque cambia sus Story Points.

### 3.4 Sprints (E3)

1. **Sprints** → S4 (activo): Sprint Goal, historias y sus cuatro tarjetas de métricas.
2. **Sprints cerrados** → S3: historias planificadas y completadas.
3. No cerrar S4 durante la demo: cambiaría el dashboard que ya se mostró.

### 3.5 Esfuerzo y defectos (E5 y E6)

1. **Esfuerzo**: registros por integrante e historia, con los totales.
2. **Defectos**: severidad y estado con insignias. El único abierto es el del punto decimal del dashboard: pasarlo a
   *En progreso* y explicar el flujo Abierto → En progreso → Resuelto (al resolverlo se elige el sprint de
   resolución, S4).

### 3.6 Reportes y exportación (E9)

1. **Reportes** → **Reporte del proyecto**: resumen de los sprints, métricas acumuladas y defectos abiertos.
2. Reporte del Sprint 4: lleva la marca **Reporte parcial** porque el sprint está activo.
3. **Descargar PDF** (HU-31) y abrir el archivo descargado: mismo contenido que en pantalla y con coma decimal.
   Plan B si la descarga falla: **Imprimir o guardar como PDF**.

### 3.7 Cierre

Volver al dashboard: es la pantalla que resume todo lo que se mostró.

## 4. Después de la demo

Votar en Planning Poker y pasar el defecto a *En progreso* cambian los datos. Para dejarlos como al principio:
cerrar la app, `rm -rf data/`, `make demo` y `make run`.
