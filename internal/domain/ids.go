// Package domain define los identificadores que comparten las épicas (ADR 0002).
//
// Las épicas se referencian entre sí por estos IDs y no copiando las estructuras de otra épica.
// Los genera el repositorio de cada entidad al guardarla por primera vez.
package domain

// ProyectoID identifica un proyecto (épica E1). Formato: "P1", "P2", …
type ProyectoID string

// IntegranteID identifica a un integrante de un proyecto (épica E1). Formato: "I1", "I2", …
type IntegranteID string

// HistoriaID identifica una historia del Product Backlog (épica E2). Formato: "HU-1", "HU-2", … único en todo el sistema.
type HistoriaID string

// SprintID identifica un sprint (épica E3). Formato: "S1", "S2", …
type SprintID string

// RegistroID identifica un registro de esfuerzo (épica E5). Formato: "E1", "E2", …
type RegistroID string

// DefectoID identifica un defecto (épica E6). Formato: "D1", "D2", …
type DefectoID string
