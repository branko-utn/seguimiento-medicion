# ADR 0002 · Contratos de dominio entre épicas

- **Estado:** propuesta (23/10/2026), a aceptar con la revisión del PR
- **Responsable:** Branko Almeira
- **Issue:** #43 (TEC-09)

## Contexto

En el Sprint 1 cinco personas programan a la vez historias que dependen entre sí. Por ejemplo, asignar historias
a un sprint (HU-11, Sosa) y completarlas (HU-12, Branko) usan la `Historia` del backlog (HU-05, Bianchi), y las
métricas (HU-24, Giorlando) necesitan saber qué historias se completaron. Sin acuerdos previos, cada uno
inventaría su propia versión de "historia" y cada PR tendría conflictos.

## Decisión

### 1. Las épicas se referencian por ID

Los identificadores compartidos están en `internal/domain/ids.go` (`ProyectoID`, `IntegranteID`, `HistoriaID`,
`SprintID`, `RegistroID`, `DefectoID`). Una épica guarda el **ID** de lo que pertenece a otra, nunca una copia:
un defecto guarda `HistoriaID`, no una `Historia`.

### 2. Quién define qué

Cada responsable define estos tipos **en su primer PR**, con estos nombres y campos como mínimo. Puede agregar
campos y métodos; para cambiar o quitar algo de esta tabla, avisa en la Daily, porque otros dependen de eso.

| Paquete | Responsable | Tipos que usan otras épicas |
|---|---|---|
| `internal/domain/project` | Sosa | `Proyecto`, `Estado`, `Integrante`, `Rol` |
| `internal/domain/backlog` | Bianchi | `Historia`, `Prioridad`, `Estado`, `Criterio` |
| `internal/domain/sprint` | Sosa | `Sprint`, `Estado` |
| `internal/domain/estimation` | Bianchi | `Sesion` (Planning Poker), escala de Story Points |
| `internal/domain/effort` | Lara | `Registro`, `EstimacionHoras` |
| `internal/domain/defect` | Lara | `Defecto`, `Severidad`, `Estado` |
| `internal/metrics` | Giorlando | funciones puras con **sus propios tipos de entrada** |
| `internal/report` | Branko | `ReporteSprint`, `ReporteProyecto` |

```go
// project (Sosa)
type Proyecto struct {
	ID          domain.ProyectoID
	Nombre      string
	Descripcion string
	Inicio, Fin time.Time
	Estado      Estado // Planificado, EnCurso, Finalizado
}
type Integrante struct {
	ID         domain.IntegranteID
	ProyectoID domain.ProyectoID
	Nombre     string
	Rol        Rol // AgileEnabler, ProductBuilder, ProductArchitect
}

// backlog (Bianchi)
type Historia struct {
	ID          domain.HistoriaID
	ProyectoID  domain.ProyectoID
	Titulo      string
	Descripcion string
	Prioridad   Prioridad // Alta, Media, Baja
	Estado      Estado    // Pendiente, EnProgreso, Completada
	StoryPoints int       // 0 = sin estimar
	Criterios   []Criterio
	SprintID    domain.SprintID // "" = no está asignada a ningún sprint
}

// sprint (Sosa)
type Sprint struct {
	ID          domain.SprintID
	ProyectoID  domain.ProyectoID
	Goal        string
	Inicio, Fin time.Time
	Estado      Estado              // Planificado, Activo, Cerrado
	Historias   []domain.HistoriaID // asignadas (HU-11)
	Completadas []domain.HistoriaID // completadas en este sprint (HU-12)
	// Al cerrar (HU-13) quedan fijos:
	SPPlanificados, SPCompletados int
}

// effort (Lara)
type Registro struct {
	ID           domain.RegistroID
	IntegranteID domain.IntegranteID
	HistoriaID   domain.HistoriaID
	Fecha        time.Time
	Actividad    string
	Horas        float64
}
type EstimacionHoras struct {
	HistoriaID domain.HistoriaID
	Horas      float64
}

// defect (Lara)
type Defecto struct {
	ID               domain.DefectoID
	Descripcion      string
	Severidad        Severidad // Critica, Alta, Media, Baja
	Estado           Estado    // Abierto, EnProgreso, Resuelto
	HistoriaID       domain.HistoriaID
	SprintDeteccion  domain.SprintID
	SprintResolucion domain.SprintID // "" mientras no esté resuelto
}
```

### 3. Dependencias permitidas entre paquetes

```text
domain (IDs) ◄── project ◄── backlog ◄── sprint
            ◄── estimation, effort, defect   (solo importan domain)
metrics: no importa paquetes de dominio
report:  importa metrics y los tipos que necesita mostrar
service: importa todo lo anterior; server importa service
```

- `sprint` puede importar `backlog` (para validar que una historia esté estimada al asignarla). `backlog` **no** importa
  `sprint`: guarda solo el `SprintID`. Así no hay ciclos.
- **`metrics` no depende de nadie.** Recibe datos simples, por ejemplo
  `metrics.SPCompletados([]metrics.HistoriaSprint{{StoryPoints: 3, Completada: true}})`. El servicio traduce de los tipos
  de dominio a esos datos. Giorlando puede trabajar desde el primer minuto, sin esperar a nadie.

### 4. Las operaciones que tocan dos épicas van en el servicio

Si una acción cambia entidades de dos épicas, cada entidad valida sus reglas en su paquete de dominio y el servicio
coordina. Ejemplo, completar una historia (HU-12):

```go
// internal/service: carga las dos entidades, aplica las reglas de cada una y guarda ambas.
func (s *Sprints) CompletarHistoria(sprintID domain.SprintID, historiaID domain.HistoriaID) error {
	sp, err := s.sprints.Buscar(sprintID)      // error si no existe
	h, err := s.historias.Buscar(historiaID)   // error si no existe
	if err := sp.Completar(h.ID); err != nil { … }   // regla del sprint
	if err := h.Completar(); err != nil { … }        // regla de la historia
	// guardar sp y h
}
```

### 5. Convenciones comunes

- **Errores:** cada paquete de dominio exporta errores con mensaje para el usuario, en español:
  `var ErrNombreInvalido = errors.New("el nombre debe tener entre 3 y 100 caracteres")`.
  Los tests los comparan con `errors.Is`; los handlers muestran `err.Error()`.
- **Constructores que validan:** `project.Nuevo(nombre, descripcion, inicio, fin) (Proyecto, error)`. Una entidad
  inválida no se puede crear.
- **Fechas:** `time.Time` sin hora (se usa solo el día). En formularios, formato `2006-01-02`.
- **Repositorios:** cada servicio declara la interfaz del repositorio que necesita y `internal/store/memory` la implementa
  con `Store[T]` (ADR 0001 y TEC-04). Los IDs los genera el repositorio al guardar por primera vez.

### 6. Orden de merge del Sprint 1

Para que nadie quede bloqueado, cada responsable abre primero un PR chico **solo con sus tipos y constructores**
(spec, tests y código de la primera historia), y lo demás va en PRs siguientes.

| Prioridad | PR | Desbloquea a |
|---|---|---|
| 1 | Sosa · HU-01/HU-03 (`Proyecto`, `Integrante`) | Sprints, esfuerzo, Planning Poker |
| 1 | Bianchi · HU-05 (`Historia`) | Sosa (HU-11), Branko (HU-12), reportes |
| 2 | Sosa · HU-10 (`Sprint`) | Branko (HU-12), defectos, métricas en el servicio |
| — | Giorlando, Lara | No dependen de nadie: arrancan ya |

Meta: los PR de prioridad 1 mergeados el **sábado 24 al mediodía**.

## Consecuencias

- Cada integrante trabaja en su paquete y los conflictos se limitan a `server.go` (rutas) y `Servicios`.
- Mientras un tipo de otra épica no esté mergeado, se puede avanzar con la spec, los escenarios BDD y los tests del propio dominio.
- Un cambio a un tipo de esta tabla se avisa en la Daily y se actualiza este ADR.
