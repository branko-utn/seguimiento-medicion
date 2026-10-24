package service

import (
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

// RepoEstimaciones es lo que el servicio de estimaciones necesita del almacenamiento.
type RepoEstimaciones interface {
	Guardar(e effort.EstimacionHoras) error
	Buscar(id domain.HistoriaID) (effort.EstimacionHoras, bool)
}

// Estimaciones son los casos de uso de estimacion de horas (Epica E5).
type Estimaciones struct {
	estimaciones RepoEstimaciones
	historias    RepoHistorias
}

func NuevoEstimaciones(estimaciones RepoEstimaciones, historias RepoHistorias) *Estimaciones {
	return &Estimaciones{estimaciones: estimaciones, historias: historias}
}

// Estimar registra o reemplaza las horas estimadas de una historia (HU-19).
func (s *Estimaciones) Estimar(historiaID domain.HistoriaID, horas float64) (effort.EstimacionHoras, error) {
	if _, err := s.historias.Buscar(historiaID); err != nil {
		return effort.EstimacionHoras{}, err
	}
	e, err := effort.NuevaEstimacion(historiaID, horas)
	if err != nil {
		return effort.EstimacionHoras{}, err
	}
	if err := s.estimaciones.Guardar(e); err != nil {
		return effort.EstimacionHoras{}, err
	}
	return e, nil
}

// Buscar devuelve las horas estimadas de una historia.
func (s *Estimaciones) Buscar(historiaID domain.HistoriaID) (effort.EstimacionHoras, bool, error) {
	if _, err := s.historias.Buscar(historiaID); err != nil {
		return effort.EstimacionHoras{}, false, err
	}
	e, ok := s.estimaciones.Buscar(historiaID)
	return e, ok, nil
}
