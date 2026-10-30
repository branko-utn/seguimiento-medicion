package service

import (
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

// RepoRegistros es lo que el servicio de esfuerzo necesita del almacenamiento.
type RepoRegistros interface {
	// Agregar guarda un registro nuevo y lo devuelve con el ID generado (E1, E2, …).
	Agregar(r effort.Registro) (effort.Registro, error)
	// DelIntegrante devuelve los registros de un integrante, en el orden en que se registraron.
	DelIntegrante(id domain.IntegranteID) []effort.Registro
	// DeHistoria devuelve los registros de una historia, en el orden en que se registraron.
	DeHistoria(id domain.HistoriaID) []effort.Registro
	// Listar devuelve todos los registros.
	Listar() []effort.Registro
}

// BuscadorIntegrantes es lo que el servicio de esfuerzo necesita del repositorio de integrantes (épica E1).
type BuscadorIntegrantes interface {
	// Buscar devuelve el integrante, o project.ErrIntegranteNoExiste.
	Buscar(id domain.IntegranteID) (project.Integrante, error)
}

// Registros son los casos de uso del registro de esfuerzo (épica E5).
type Registros struct {
	registros   RepoRegistros
	integrantes BuscadorIntegrantes
	historias   RepoHistorias
}

// NuevoRegistros crea el servicio. Usa los integrantes (épica E1) y las historias (épica E2) para validar el registro.
func NuevoRegistros(registros RepoRegistros, integrantes BuscadorIntegrantes, historias RepoHistorias) *Registros {
	return &Registros{registros: registros, integrantes: integrantes, historias: historias}
}

// Registrar guarda las horas que un integrante trabajó en una historia (HU-20). Verifica que el integrante y la
// historia existan y sean del mismo proyecto (RN-20.5); el dominio valida actividad, horas y fecha (RN-20.1, RN-20.2,
// RN-20.4) y el tope de 24 horas por día con los registros existentes del integrante (RN-20.3).
func (s *Registros) Registrar(integranteID domain.IntegranteID, historiaID domain.HistoriaID, fecha time.Time, actividad string, horas float64, hoy time.Time) (effort.Registro, error) {
	integrante, err := s.integrantes.Buscar(integranteID)
	if err != nil {
		return effort.Registro{}, err
	}
	historia, err := s.historias.Buscar(historiaID)
	if err != nil {
		return effort.Registro{}, err
	}
	if integrante.ProyectoID != historia.ProyectoID {
		return effort.Registro{}, effort.ErrIntegranteDeOtroProyecto
	}
	r, err := effort.NuevoRegistro(integranteID, historiaID, fecha, actividad, horas, hoy)
	if err != nil {
		return effort.Registro{}, err
	}
	if err := effort.ValidarHorasDiarias(s.registros.DelIntegrante(integranteID), r); err != nil {
		return effort.Registro{}, err
	}
	return s.registros.Agregar(r)
}

// DelIntegrante devuelve los registros de un integrante, en el orden en que se registraron.
func (s *Registros) DelIntegrante(id domain.IntegranteID) []effort.Registro {
	return s.registros.DelIntegrante(id)
}

// DeHistoria devuelve los registros de una historia, en el orden en que se registraron.
func (s *Registros) DeHistoria(id domain.HistoriaID) []effort.Registro {
	return s.registros.DeHistoria(id)
}
