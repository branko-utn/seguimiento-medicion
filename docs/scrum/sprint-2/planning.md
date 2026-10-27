# Sprint 2 · Planning

- **Asistentes:** Branko Almeira, Federico Sosa, Ignacio Bianchi, Agustín Lara, Agustín Giorlando
- **Sprint:** La Interfaz

## Sprint Goal

> Usar la aplicación desde el navegador, recorriendo proyecto → backlog → estimación → sprint → métricas, y estimar
> con Planning Poker.

## Capacidad

Velocidad del Sprint 1: **44 SP** en 3 días, lo que da unos 30 SP para un sprint de 2 días. Se comprometen **38 SP**:
el sprint se apoya en código ya probado (las pantallas usan servicios que están en `main`). Si en la daily del segundo
día no se llega, se recorta el último ítem de cada persona.

## Sprint Backlog

| Responsable | Ítem | SP |
|---|---|---|
| Federico Sosa | TEC-13 · Pantallas de proyectos e integrantes | 3 |
| | TEC-14 · Pantalla de sprints (crear, iniciar, asignar) | 3 |
| | HU-02 · Modificar proyecto (#9) | 2 |
| Ignacio Bianchi | HU-16 · Votar en Planning Poker con votos ocultos (#23), con pantalla | 5 |
| | HU-17 · Revelar votos y detectar diferencias (#24), con pantalla | 3 |
| Agustín Lara | TEC-15 · Pantallas de esfuerzo y defectos | 3 |
| | HU-23 · Cambiar estado de un defecto (#30), con pantalla | 3 |
| | HU-27 · Defectos detectados y resueltos (#34) | 2 |
| Agustín Giorlando | TEC-16 · Pantalla del backlog | 3 |
| | TEC-17 · Vista de métricas del sprint | 2 |
| | HU-26 · Horas estimadas, reales y desviación (#33) | 3 |
| Branko Almeira | TEC-18 · Conectar los servicios a la app web | 1 |
| | HU-13 · Cerrar sprint (#20), con pantalla | 5 |
| | **Total** | **38** |

Branko tiene 2 SP menos porque además facilita las ceremonias y arma las actas.

**Pasan al Sprint 3:** HU-29 (reporte de sprint), HU-04 (estado del proyecto), HU-06 (editar y eliminar del backlog) y
HU-09 (ordenar el backlog, recortable).

## Decisiones

- **El incremento del Sprint 1 fue el núcleo funcional verificado con BDD, sin interfaz.** El MVP usable es el objetivo
  de este sprint.
- **Identidad sin login (ADR 0004):** quien vota en el Planning Poker se elige de la lista de integrantes del proyecto.
- **HU-13:** un sprint que nunca se inició no se puede cerrar.
- **HU-16:** sin cartas "?" ni "☕", y no se cambia el voto antes de revelar.
- **Carga pareja entre los cinco**, incluido el Agile Enabler.
- **Specs sin esperas entre personas:** cada uno hace lo suyo de punta a punta (dominio, servicio, escenarios y
  pantalla), y una spec solo depende de código que ya está en `main`. HU-27 pasa de Giorlando a Lara por ese motivo.
- **Desde este sprint, "probado en la aplicación" deja de ser N/A en la Definition of Done.**

## Riesgos

- La capacidad comprometida supera la estimada en unos 8 SP.
- Las pantallas registran rutas en `internal/server`: TEC-18 entra primero para conectar todos los servicios y evitar conflictos.
