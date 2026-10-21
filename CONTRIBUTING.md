# Cómo trabajamos en este repositorio

Resumen operativo del plan aprobado ([`docs/plan/cronograma-sprints.html`](docs/plan/cronograma-sprints.html)).
Las reglas de Scrum (ceremonias, DoR y DoD) están en [`docs/scrum/`](docs/scrum/).

## 1. Ramas

`main` está protegida: solo se modifica mediante un Pull Request con el CI en verde y **una aprobación** de otro integrante.

| Tipo | Nombre de la rama |
|---|---|
| Historia de usuario | `feature/HU-01-crear-proyecto` |
| Ítem técnico | `chore/TEC-02-persistencia` |
| Defecto | `fix/DEF-3-descripcion-corta` |
| Documentación / actas | `docs/sprint-1-retro` |

```bash
git switch main && git pull
git switch -c feature/HU-01-crear-proyecto
```

## 2. Commits: evidencia de TDD

Usamos [Conventional Commits](https://www.conventionalcommits.org/es/) con el ID de la historia y la fase de TDD.
**Un commit por fase**; los profesores evalúan TDD mirando este historial.

```text
docs(HU-01): spec SDD y escenarios BDD
test(HU-01): RED - nombre vacío es inválido
feat(HU-01): GREEN - valida nombre obligatorio
refactor(HU-01): extrae normalización del nombre
```

1. **RED**: el test se escribe primero y se commitea **fallando** (`go test` en rojo).
2. **GREEN**: el código mínimo que lo hace pasar.
3. **REFACTOR**: se mejora el código sin romper los tests.

El commit RED puede romper el CI de la rama; no importa, el PR se mergea cuando todo está en verde.
El ejemplo completo está en el historial de `chore/TEC-01-esqueleto` (endpoint `/health`).

## 3. Pull Requests

1. `git push -u origin <rama>` y `gh pr create` (o desde la web). Se puede abrir como *draft* desde el primer commit.
2. Completar la plantilla del PR, **incluida la tabla de trazabilidad**, y escribir `Closes #N` para vincular el issue.
3. **Revisión cruzada**: cada uno revisa los PRs de quien está a su izquierda en la rotación

   `Sosa → Bianchi → Lara → Giorlando → Branko → Sosa`

   | Autor del PR | Revisa |
   |---|---|
   | Federico Sosa | Ignacio Bianchi |
   | Ignacio Bianchi | Agustín Lara |
   | Agustín Lara | Agustín Giorlando |
   | Agustín Giorlando | Branko Almeira |
   | Branko Almeira | Federico Sosa |

   El workflow [`revisor.yml`](.github/workflows/revisor.yml) pide la revisión automáticamente al abrir el PR.
   Si el revisor no responde en el día, se avisa en la Daily y revisa cualquier otro integrante.
4. Se mergea con **"Create a merge commit"** o **"Rebase and merge"**. El *squash* está deshabilitado porque borraría los commits RED/GREEN/REFACTOR.
5. Mover la tarjeta a **Done** en el tablero.

## 4. Trazabilidad

Cada historia recorre la cadena completa con el mismo ID:

```text
Issue HU-XX → specs/HU-XX-*.md → criterios CA-XX.n → features/hu-xx-*.feature → *_test.go → código Go
```

Plantillas: [`specs/_plantilla.md`](specs/_plantilla.md) y [`features/_plantilla.feature`](features/_plantilla.feature).

## 5. Estructura del código

| Carpeta | Qué va | Depende de |
|---|---|---|
| `internal/domain/<épica>` | Entidades, reglas de negocio y validaciones. **Acá se hace TDD.** | nada |
| `internal/metrics` | Cálculo de métricas (funciones puras) | `domain` |
| `internal/report` | Armado de reportes y exportación a PDF | `domain`, `metrics` |
| `internal/service` | Casos de uso: coordinan dominio y repositorios | `domain`, `metrics` |
| `internal/store/memory` | Repositorios (en memoria con persistencia a archivo) | `domain` |
| `internal/server` | Handlers HTTP y rutas | `service` |
| `web/templates` | Pantallas HTML (`html/template` + HTMX) | — |
| `features` | Escenarios BDD y sus *steps* (godog) | `service`, `server` |
| `cmd/server` | Punto de entrada de la aplicación | todo |

Ver la decisión completa en [`docs/adr/0001-arquitectura.md`](docs/adr/0001-arquitectura.md).

## 6. Comandos

| Quiero… | Comando |
|---|---|
| Correr la aplicación | `make run` → http://localhost:8080 |
| Correr todos los tests (unitarios + BDD) | `make test` |
| Solo los escenarios BDD | `make bdd` |
| Ver la cobertura | `make cover` |
| Formatear el código | `make fmt` |

## 7. Uso de IA

Está permitido, pero **cada integrante debe entender y poder explicar todo lo que sube**.
Registrar qué se pidió, qué se usó y cómo se validó en [`docs/ia/registro.md`](docs/ia/registro.md).
