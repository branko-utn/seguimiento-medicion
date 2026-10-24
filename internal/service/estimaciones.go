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
	return effort.EstimacionHoras{}, effort.ErrHorasInvalidas
}

// Buscar devuelve las horas estimadas de una historia.
func (s *Estimaciones) Buscar(historiaID domain.HistoriaID) (effort.EstimacionHoras, bool, error) {
	return effort.EstimacionHoras{}, false, nil
}
