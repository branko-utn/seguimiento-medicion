# seguimiento-medicion

**Software Metrics & Estimation** — Trabajo Práctico Integrador de *Ingeniería y Calidad de Software 2026*, UTN Facultad Regional San Rafael.

Aplicación para administrar proyectos de software y obtener información sobre su **estimación, planificación, seguimiento y calidad**.

## Funcionalidades

- **Gestión de proyectos:** alta y modificación, integrantes, fechas y estado.
- **Product Backlog:** historias con prioridad, estado, Story Points y criterios de aceptación.
- **Sprints:** Sprint Goal, asignación de historias, cierre e historial.
- **Estimación:** Story Points mediante Planning Poker (votación oculta, rondas y estimación acordada).
- **Registro de esfuerzo:** horas trabajadas por integrante y actividad, para comparar lo estimado con lo real.
- **Defectos:** severidad, estado, historia relacionada y sprints de detección y resolución.
- **Métricas y dashboard:** velocidad, Story Points planificados y completados, desviación de esfuerzo, defectos, entre otras.
- **Reportes:** reporte de proyecto o Sprint, exportable a PDF.

## Tecnologías y prácticas

- **Go** para el núcleo y las reglas de negocio.
- **Scrum** para la gestión del proyecto.
- **SDD**, **BDD** (Given–When–Then) y **TDD** (Red → Green → Refactor).
- **Git** y **GitHub Projects** para el control de versiones y el tablero.

## Cómo empezar

Requisitos: [Go](https://go.dev/dl/) 1.22 o superior y Git.

```bash
git clone https://github.com/branko-utn/seguimiento-medicion.git
cd seguimiento-medicion
make test   # tests unitarios + escenarios BDD
make run    # http://localhost:8080
```

Sin `make` (por ejemplo en Windows): `go test ./...` y `go run ./cmd/server`.

## Documentación

| Documento | Qué tiene |
|---|---|
| [CONTRIBUTING.md](CONTRIBUTING.md) | Ramas, commits de TDD, Pull Requests, revisión cruzada y estructura del código |
| [docs/plan/cronograma-sprints.html](docs/plan/cronograma-sprints.html) | Cronograma de sprints, ceremonias y reparto de épicas |
| [docs/scrum/](docs/scrum/) | Ceremonias, DoR/DoD, guía de Planning Poker y actas de cada sprint |
| [docs/adr/](docs/adr/) | Decisiones de arquitectura |
| [specs/](specs/) | Especificaciones SDD |
| [features/](features/) | Escenarios BDD (Given–When–Then) |
| [docs/ia/registro.md](docs/ia/registro.md) | Registro de uso de IA |

## Equipo

| Apellido y Nombre | Rol |
|---|---|
| Almeira, Branko | Agile Enabler y Product Builder (E9 Reportes) |
| Sosa, Federico | Product Builder (E1 Proyectos, E3 Sprints) |
| Bianchi, Ignacio | Product Builder (E2 Backlog, E4 Estimación) |
| Lara, Agustín | Product Builder (E5 Esfuerzo, E6 Defectos) |
| Giorlando, Agustín | Product Builder (E7 Métricas, E8 Dashboard) |

**Product Architects:** profesores de la cátedra.

## Enlaces

- [Tablero del proyecto](https://github.com/users/branko-utn/projects/1)
- [Repositorio del proyecto](https://github.com/branko-utn/seguimiento-medicion)
