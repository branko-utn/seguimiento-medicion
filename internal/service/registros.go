package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoRegistros es lo que el servicio de esfuerzo necesita del almacenamiento.
type RepoRegistros interface {
	Agregar(r effort.Registro) (effort.Registro, error)
	DelIntegrante(id domain.IntegranteID) []effort.Registro
	DeHistoria(id domain.HistoriaID) []effort.Registro
}

// BuscadorIntegrantes es lo que el servicio de esfuerzo necesita del repositorio de integrantes.
type BuscadorIntegrantes interface {
	Buscar(id domain.IntegranteID) (project.Integrante, error)
}

// Registros son los casos de uso del registro de esfuerzo (épica E5).
type Registros struct {
	registros   RepoRegistros
	integrantes BuscadorIntegrantes
	historias   RepoHistorias
}

func NuevoRegistros(registros RepoRegistros, integrantes BuscadorIntegrantes, historias RepoHistorias) *Registros {
	return &Registros{registros: registros, integrantes: integrantes, historias: historias}
}

func (s *Registros) Registrar(integranteID domain.IntegranteID, historiaID domain.HistoriaID, fecha time.Time, actividad string, horas float64, hoy time.Time) (effort.Registro, error) {
	return effort.Registro{}, nil
}

func (s *Registros) DelIntegrante(id domain.IntegranteID) []effort.Registro {
	return nil
}

func (s *Registros) DeHistoria(id domain.HistoriaID) []effort.Registro {
	return nil
}
